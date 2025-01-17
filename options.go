package m3u8

import (
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
