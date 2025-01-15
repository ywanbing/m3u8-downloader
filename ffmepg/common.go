package ffmepg

import (
	"os"
)

func checkDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkFfmpeg() (bool, error) {
	_, err := os.Stat("./tools/ffmpeg")
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
