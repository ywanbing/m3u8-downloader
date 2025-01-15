package ffmepg

import (
	"log"
	"os/exec"
	"path/filepath"
)

// FfmpegCmd 调用ffmpeg
func FfmpegCmd(cmd string, args ...string) error {
	// 1. 检查ffmpeg是否可用
	if err := CheckFfmpeg(); err != nil {
		return err
	}

	cmdPath := filepath.Clean(filepath.Join(filepath.Join(executablePath, Tools), cmd))
	// 2. 调用ffmpeg
	command := exec.Command(cmdPath, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		return err
	}

	// 3. 打印日志
	log.Println(string(output))

	return nil
}
