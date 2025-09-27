package services

import (
	"bytes"
	"context"
	"image"
	"mime/multipart"
	"time"
	"tutuplapak-user/internal/config"
	"tutuplapak-user/internal/dto"
	minioUploader "tutuplapak-user/internal/utils"

	"github.com/minio/minio-go/v7"
)

type FileBeliService struct {
	config config.Config
	minio  *minio.Client
}

func NewFileBeliService(config config.Config) *FileService {
	minioConfig := &minioUploader.MinioConfig{
		AccessKeyID:     config.Minio.AccessKeyID,
		SecretAccessKey: config.Minio.SecretAccessKey,
		UseSSL:          config.Minio.UseSSL,
		Endpoint:        config.Minio.Endpoint,
	}

	minioClient, _ := minioUploader.NewUploader(minioConfig)

	return &FileService{
		config: config,
		minio:  minioClient,
	}
}

func (uc *FileBeliService) UploadFileBeli(ctx context.Context, file *multipart.FileHeader, src multipart.File, fileName string) (dto.UploadBeliResponse, error) {
	bucketName := uc.config.Minio.BucketName

	// compress file
	srcForCompress, err := file.Open()
	if err != nil {
		return dto.UploadBeliResponse{}, err
	}
	defer srcForCompress.Close()

	img, _, _ := image.Decode(srcForCompress)

	compressedBytes, contentType, _ := minioUploader.CompressThumbnail(file, img, 10)

	// upload file compress
	_, err = uc.minio.PutObject(
		ctx,
		bucketName,
		fileName,
		bytes.NewReader(compressedBytes),
		int64(len(compressedBytes)),
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		return dto.UploadBeliResponse{}, err
	}

	url, err := uc.minio.PresignedGetObject(ctx, bucketName, fileName, time.Hour*24*7, nil)
	if err != nil {
		return dto.UploadBeliResponse{}, err
	}

	return dto.UploadBeliResponse{
		Message: "File uploaded sucessfully",
		Data: dto.UploadBeliData{
			ImageUrl: url.String(),
		},
	}, nil
}
