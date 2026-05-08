package utils

import (
	"encoding/base64"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SaveBase64Image(base64Data string, fileName string) error {
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(imageData)
	if err != nil {
		return err
	}
	return nil
}
