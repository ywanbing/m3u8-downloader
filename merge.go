package m3u8

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"video/ffmepg"

	"github.com/samber/lo"
)

type MergeTsFileReq struct {
	Dir        string
	TsFileList []string
	OutputMp4  string
	AfterDel   bool
}

func MergeTsFileListToSingleMp4(req *MergeTsFileReq) (err error) {
	if !isDirExists(req.Dir) {
		return errors.New("dir not exists")
	}

	// 组装ts文件列表
	tsFiles := lo.Map(req.TsFileList, func(item string, _ int) string {
		return filepath.Join(req.Dir, item)
	})
	outPutFileName := filepath.Clean(filepath.Join(req.Dir, "../"+req.OutputMp4))

	// ffmpeg -i "concat:1.ts|2.ts|3.ts" -c copy output.mp4
	err = ffmepg.FfmpegCmd("ffmpeg", "-i", "concat:"+strings.Join(tsFiles, "|"), "-c", "copy", outPutFileName)
	if err != nil {
		return err
	}

	if req.AfterDel {
		// 删除ts文件
		err = os.RemoveAll(req.Dir)
		if err != nil {
			return err
		}
	}

	return nil
}
