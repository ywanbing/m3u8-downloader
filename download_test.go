package m3u8

import (
	"net/url"
	"testing"
)

const TempUrl = "https://hn.bfvvs.com/play/nelVLW5b/index.m3u8"

func ExampleDownloader_Download() {
	baseUrl, _ := url.Parse(TempUrl)
	downloader := NewDownloader(baseUrl, "./temp", "动画", nil)
	_ = downloader.Download()
}

func TestDownloader_Download(t *testing.T) {
	baseUrl, _ := url.Parse(TempUrl)
	downloader := NewDownloader(baseUrl, "./temp", "动画", nil)
	_ = downloader.Download()
}
