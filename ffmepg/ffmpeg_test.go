//go:build windows && amd64

package ffmepg

import (
	"os/exec"
	"testing"
)

func Test_Exec(t *testing.T) {
	if err := CheckFfmpeg(); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("./ffmpeg.exe", "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
}
