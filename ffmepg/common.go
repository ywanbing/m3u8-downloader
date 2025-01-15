package ffmepg

import (
	"os"
	"path/filepath"
)

const (
	Tools   = "tools"
	ToolDir = "./tools"
)

var executablePath string

func init() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	executablePath = filepath.Dir(executable)
	return
}

func checkDir(dir string) error {
	rawDir := filepath.Clean(filepath.Join(executablePath, dir))
	_, err := os.Stat(rawDir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(rawDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkFfmpeg() (bool, error) {
	ffmpegPath := filepath.Join(executablePath, "tools/ffmpeg")
	_, err := os.Stat(ffmpegPath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func writeFile(name string, data []byte) error {
	name = filepath.Clean(filepath.Join(executablePath, name))
	return os.WriteFile(name, data, os.ModePerm)
}
