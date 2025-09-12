package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nfnt/resize"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

func NewUploader(cfg *MinioConfig) (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func CompressThumbnail(fileHeader *multipart.FileHeader, img image.Image, maxKB int) ([]byte, string, error) {
	var buf bytes.Buffer
	ext := strings.ToLower(fileHeader.Filename)

	// Resize smaller for thumbnail
	thumb := resize.Resize(0, 256, img, resize.Lanczos3)

	// --- If JPEG/JPG ---
	if strings.HasSuffix(ext, ".jpg") || strings.HasSuffix(ext, ".jpeg") {
		quality := 85
		for quality > 20 {
			buf.Reset()
			_ = jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: quality})
			if buf.Len() <= maxKB*1024 {
				break
			}
			quality -= 5
		}
		return buf.Bytes(), "image/jpeg", nil
	}

	// --- If PNG ---
	// PNG compression is lossless and won’t easily shrink to 10KB
	// So either accept bigger size or convert to JPEG
	buf.Reset()
	err := png.Encode(&buf, thumb)
	if err != nil {
		return nil, "", err
	}

	if buf.Len() > maxKB*1024 {
		// fallback: convert PNG to JPEG to respect max size
		buf.Reset()
		_ = jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80})
		return buf.Bytes(), "image/jpeg", nil
	}

	return buf.Bytes(), "image/png", nil
}

func IsAllowedFileType(fileName, fileType string) bool {
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		// "application/octet-stream": true,
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedMimeTypes[fileType] {
		return false
	}

	if fileType == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(fileName))
		if !allowedExtensions[ext] {
			return false
		}
	}

	return true
}

func AddFileNameSuffix(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s_thumbnail%s", name, ext)
}
