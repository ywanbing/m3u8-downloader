package m3u8

import (
	"context"
	"log/slog"

	"github.com/imroc/req/v3"
)

type Options func(d *Downloader)

// WithLogger 设置日志
func WithLogger(log *slog.Logger) Options {
	return func(d *Downloader) {
		d.log = log
	}
}

// WithClient 设置客户端
func WithClient(client *req.Client) Options {
	return func(d *Downloader) {
		d.client = client
	}
}

// WithMaxGoroutines 设置最大协程数,default: Cpus
func WithMaxGoroutines(maxGoroutines int) Options {
	return func(d *Downloader) {
		d.maxGoroutines = maxGoroutines
	}
}

// WithContext 设置上下文,可以用于取消下载
func WithContext(ctx context.Context) Options {
	return func(d *Downloader) {
		d.ctx = ctx
	}
}
