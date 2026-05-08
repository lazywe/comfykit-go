package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 300 * time.Second,
}

func DownloadFile(url, destPath string) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadFileToTemp(url string) (string, error) {
	tempFile, err := os.CreateTemp("", "download_*")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	resp, err := httpClient.Get(url)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

func IsURL(str string) bool {
	return strings.HasPrefix(strings.ToLower(str), "http://") ||
		strings.HasPrefix(strings.ToLower(str), "https://")
}

func GetBaseURL(url string) string {
	if strings.HasPrefix(url, "https://") {
		url = url[8:]
	} else if strings.HasPrefix(url, "http://") {
		url = url[7:]
	}

	idx := strings.Index(url, "/")
	if idx != -1 {
		return url[:idx]
	}
	return url
}

func GetContentType(url string) string {
	ext := strings.ToLower(GetFileExtension(url))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".ogg":
		return "audio/ogg"
	default:
		return "application/octet-stream"
	}
}
