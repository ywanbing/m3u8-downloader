//go:build darwin && amd64

package ffmepg

import (
	"embed"
)

//go:embed mac
var ffmepgFs embed.FS

func CheckFfmpeg() error {
	ok, err := checkFfmpeg()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	file, err := ffmepgFs.ReadFile("mac/ffmpeg")
	if err != nil {
		return err
	}

	// 检查目录
	err = checkDir("./tools")
	if err != nil {
		return err
	}

	// 写入本地文件系统，进行命令调用
	err = os.WriteFile("./tools/ffmpeg", file, 0777)
	if err != nil {
		return err
	}

	return nil
}
