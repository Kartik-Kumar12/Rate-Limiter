package utils

import (
	"fmt"
	"os"
)

func ReadFileContent(filepath string) ([]byte, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %v", err.Error())
	}
	return bytes, nil
}
