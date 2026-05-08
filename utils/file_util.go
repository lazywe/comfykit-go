package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func ReadJSONFile(path string, v interface{}) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func WriteJSONFile(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return WriteFile(path, data, 0644)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func CreateDirIfNotExists(path string) error {
	if !DirExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

func GetFileName(path string) string {
	return filepath.Base(path)
}

func GetFileNameWithoutExt(path string) string {
	base := filepath.Base(path)
	return base[:len(base)-len(filepath.Ext(base))]
}
