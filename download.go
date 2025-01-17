package m3u8

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/grafov/m3u8"
	"github.com/imroc/req/v3"
)

type Downloader struct {
	baseUrl       *url.URL
	dir           string // 保存文件的绝对路径
	fileName      string // 文件名,不带扩展名
	client        *req.Client
	log           *slog.Logger
	maxGoroutines int // 并发下载最大数
	ctx           context.Context
}

func NewDownloader(baseUrl *url.URL, dir string, fileName string, ops ...Options) *Downloader {
	if strings.Contains(fileName, ".mp4") {
		fileName, _ = strings.CutSuffix(fileName, ".mp4")
	}

	d := &Downloader{
		baseUrl:  baseUrl,
		dir:      dir,
		fileName: fileName,
	}
	for _, op := range ops {
		op(d)
	}

	// check param
	if d.client == nil {
		d.client = req.C().ImpersonateChrome()
	}
	if d.log == nil {
		d.log = slog.Default()
	}
	if d.maxGoroutines <= 0 {
		d.maxGoroutines = runtime.NumCPU()
	}
	if d.ctx == nil {
		d.ctx = context.Background()
	}

	return d
}

func (d *Downloader) Download() error {
	downloadDir := filepath.Join(d.dir, d.fileName)
	for _, dir := range []string{d.dir, downloadDir} {
		if !isDirExists(dir) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				d.log.Warn("os.MkdirAll error: " + err.Error())
				return err
			}
		}
	}

	baseUrl := d.baseUrl.String()
	res := d.client.Get(baseUrl).Do()
	if res.Err != nil || !res.IsSuccessState() {
		d.log.Warn("req %s err :%v", baseUrl, res.IsSuccessState())
		return res.Err
	}

	// 获取m3u8文件
	playlist, listType, err := m3u8.DecodeFrom(res.Body, true)
	if err != nil {
		d.log.Warn("m3u8.DecodeFrom error: " + err.Error())
		return err
	}

	// 类型不合理
	if listType != m3u8.MEDIA && listType != m3u8.MASTER {
		d.log.Error("listType error: %v ", slog.Any("listType", listType))
		return errors.New("listType error")
	}

	// TODO
	if listType == m3u8.MASTER {
		// 暂时不支持
		d.log.Error("listType is MASTER not at this time")
		return nil
	}

	return d.downloadTs(playlist)
}

func (d *Downloader) downloadTs(playlist m3u8.Playlist) error {
	dlNum := make(chan struct{}, d.maxGoroutines)
	ctx, cancel := context.WithCancelCause(d.ctx)
	defer cancel(nil)

	decryptKey := ""
	var err error
	// 获取ts文件
	tsList := playlist.(*m3u8.MediaPlaylist).Segments
	wg := new(sync.WaitGroup)
	for _, ts := range tsList {
		select {
		case <-d.ctx.Done():
			return d.ctx.Err()
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if ts == nil {
			continue
		}

		if ts.Key != nil {
			// 获取 enc key
			decryptKey, err = d.getAndSetEncKey(ts.Key)
			if err != nil {
				d.log.Warn("get enc key", slog.String("err", err.Error()))
				return err
			}
		}

		select {
		case <-d.ctx.Done():
			return d.ctx.Err()
		case <-ctx.Done():
			return ctx.Err()
		case dlNum <- struct{}{}:
			wg.Add(1)
			go func(ts *m3u8.MediaSegment, decryptKey string) {
				err := d.downloadTsFile(ctx, ts, decryptKey)
				if err != nil {
					cancel(err)
				}
				<-dlNum
				wg.Done()
			}(ts, decryptKey)
		}
	}

	// 等待下载完成
	wg.Wait()
	return nil
}

// getAndSetEncKey 获取 enc key
func (d *Downloader) getAndSetEncKey(key *m3u8.Key) (string, error) {
	if key == nil {
		return "", nil
	}
	uri, err := absolutist(key.URI, d.baseUrl)
	if err != nil {
		return "", err
	}

	if key.Method == Encrypt_AES128 {
		// 获取 enc key
		res := d.client.Get(uri.String()).Do()
		if res.Err != nil || !res.IsSuccessState() {
			d.log.Warn("req %s err :%v", key.URI, res.IsSuccessState())
			return "", res.Err
		}
		return res.String(), nil
	}

	return "", errors.New("unsupported encrypt method")
}

// 下载ts文件
// @modify: 2020-08-13 修复ts格式SyncByte合并不能播放问题
func (d *Downloader) downloadTsFile(ctx context.Context, ts *m3u8.MediaSegment, decryptKey string) error {
	downloadDir := filepath.Join(d.dir, d.fileName)
	tsUrl, err := absolutist(ts.URI, d.baseUrl)
	if err != nil {
		d.log.Warn("absolutist error: " + err.Error())
		return err
	}

	tsUrlStr := tsUrl.String()
	fileName := strconv.Itoa(int(ts.SeqId)) + ".ts"
	d.log.Info("parse ts", slog.String("tsUrl", tsUrlStr), slog.String("fileName", fileName))

	// TODO 检测文件是否存在,进行断点续传
	res := d.client.Get(tsUrlStr).Do(ctx)
	if res.Err != nil || !res.IsSuccessState() {
		d.log.Warn("download ts error", slog.String("tsUrl", tsUrlStr))
		return res.Err
	}

	origData := res.Bytes()
	// 解密出视频 ts 源文件
	if decryptKey != "" {
		// 解密 ts 文件，算法：aes 128 cbc pack5
		origData, err = AesDecrypt(origData, []byte(decryptKey))
		if err != nil {
			return err
		}
	}

	// Detect Fake png file
	if bytes.HasPrefix(origData, PNG_SIGN) {
		origData = origData[8:]
	}

	// https://en.wikipedia.org/wiki/MPEG_transport_stream
	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
	// Need to remove the bytes before the SyncByte 0x47(71).
	syncByte := uint8(71) // 0x47
	bLen := len(origData)
	for j := 0; j < bLen; j++ {
		if origData[j] == syncByte {
			origData = origData[j:]
			break
		}
	}

	_ = os.WriteFile(path.Join(downloadDir, fileName), origData, 0666)
	return nil
}
