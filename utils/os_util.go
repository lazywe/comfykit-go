package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func GetTempDir() string {
	return os.TempDir()
}

func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func SetCurrentDir(path string) error {
	return os.Chdir(path)
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

func RenameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func CopyFile(src, dst string) error {
	data, err := ReadFile(src)
	if err != nil {
		return err
	}
	return WriteFile(dst, data, 0644)
}

func ListDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func GetFileModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, string(output))
	}
	return string(output), nil
}

func ExecuteCommandInDir(dir, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, string(output))
	}
	return string(output), nil
}

func ExpandHomePath(path string) string {
	if strings.HasPrefix(path, "~") {
		return filepath.Join(GetHomeDir(), path[1:])
	}
	return path
}
