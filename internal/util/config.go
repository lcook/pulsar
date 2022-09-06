package util

import (
	"encoding/json"
	"io"
	"os"
)

func GetConfig[T any](config string) (T, error) {
	var cfg T

	file, err := os.Open(config)
	if err != nil {
		return cfg, err
	}
	//nolint
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
