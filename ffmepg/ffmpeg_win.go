//go:build windows && amd64

package ffmepg

import (
	"embed"
)

//go:embed win
var ffmepgFs embed.FS

func CheckFfmpeg() error {
	ok, err := checkFfmpeg()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	file, err := ffmepgFs.ReadFile("win/ffmpeg.exe")
	if err != nil {
		return err
	}

	// 检查目录
	err = checkDir(ToolDir)
	if err != nil {
		return err
	}

	// 写入本地文件系统，进行命令调用
	err = writeFile("./tools/ffmpeg.exe", file)
	if err != nil {
		return err
	}

	return nil
}
