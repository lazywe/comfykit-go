package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func ResizeImage(inputPath, outputPath string, width, height int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	resized := resize(img, width, height)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	ext := filepath.Ext(outputPath)
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(outFile, resized, &jpeg.Options{Quality: 90})
	case ".png":
		return png.Encode(outFile, resized)
	default:
		return jpeg.Encode(outFile, resized, &jpeg.Options{Quality: 90})
	}
}

func resize(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	scale := float64(width) / float64(srcWidth)
	if float64(height)/float64(srcHeight) < scale {
		scale = float64(height) / float64(srcHeight)
	}

	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)

	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / scale)
			srcY := int(float64(y) / scale)
			result.Set(x, y, img.At(srcX, srcY))
		}
	}

	return result
}

func ImageToBase64(imagePath string) (string, error) {
	data, err := ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", GetContentType(imagePath), Base64Encode(data)), nil
}

func Base64Encode(data []byte) string {
	return fmt.Sprintf("%x", data)
}

func GetImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

func ImageToBytes(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return nil, err
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
	default:
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
