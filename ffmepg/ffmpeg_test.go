//go:build windows && amd64

package ffmepg

import (
	"os"
	"os/exec"
	"testing"
)

func Test_Exec(t *testing.T) {
	file, err := ffmepgFs.ReadFile("win/ffmpeg.exe")
	if err != nil {
		t.Fatal(err)
	}

	// 写入本地文件系统，进行命令调用
	err = os.WriteFile("./ffmpeg.exe", file, 0777)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("./ffmpeg.exe", "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
}
