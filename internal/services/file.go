package services

import (
	"bytes"
	"context"
	"image"
	"mime/multipart"
	"time"
	"tutuplapak-user/internal/config"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/entities"
	"tutuplapak-user/internal/repository"
	"tutuplapak-user/internal/utils"
	minioUploader "tutuplapak-user/internal/utils"

	"github.com/minio/minio-go/v7"
)

type UseCase interface {
	UploadFile(context.Context, *multipart.FileHeader, multipart.File, string) (*dto.File, error)
}

type useCase struct {
	config         config.Config
	minio          *minio.Client
	fileRepository repository.FileRepositoryInterface
}

func NewUseCase(config config.Config, fileRepository repository.FileRepositoryInterface) UseCase {
	minioConfig := &minioUploader.MinioConfig{
		AccessKeyID:     config.Minio.AccessKeyID,
		SecretAccessKey: config.Minio.SecretAccessKey,
		UseSSL:          config.Minio.UseSSL,
		Endpoint:        config.Minio.Endpoint,
	}

	minioClient, _ := minioUploader.NewUploader(minioConfig)

	return &useCase{
		config:         config,
		minio:          minioClient,
		fileRepository: fileRepository,
	}
}

func (uc *useCase) UploadFile(ctx context.Context, file *multipart.FileHeader, src multipart.File, fileName string) (*dto.File, error) {
	bucketName := uc.config.Minio.BucketName

	// compress file
	srcForCompress, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer srcForCompress.Close()

	img, _, _ := image.Decode(srcForCompress)

	compressedBytes, contentType, _ := utils.CompressThumbnail(file, img, 10)

	// upload the original
	_, err = uc.minio.PutObject(
		ctx,
		bucketName,
		fileName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)

	if err != nil {
		return nil, err
	}

	_, err = uc.minio.PutObject(
		ctx,
		bucketName,
		"compressed_"+fileName,
		bytes.NewReader(compressedBytes),
		int64(len(compressedBytes)),
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		return nil, err
	}

	// get url original
	url, err := uc.minio.PresignedGetObject(ctx, bucketName, fileName, time.Hour*24*7, nil)
	if err != nil {
		return nil, err
	}
	// get url thumbnail
	thumbUrl, err := uc.minio.PresignedGetObject(ctx, bucketName, "compressed_"+fileName, time.Hour*24*7, nil)
	if err != nil {
		return nil, err
	}

	// build entity
	filePayload := &entities.File{
		Url:          url.String(),
		ThumbnailUrl: thumbUrl.String(),
	}

	// repository expects dto.File → map entity to dto
	dtoPayload := dto.File{
		Url:          filePayload.Url,
		ThumbnailUrl: filePayload.ThumbnailUrl,
	}
	res, err := uc.fileRepository.Post(ctx, dtoPayload)
	if err != nil {
		return nil, err
	}

	return &dto.File{
		ID:           res.ID,
		Url:          res.Url,
		ThumbnailUrl: res.ThumbnailUrl,
		CreateAt:     res.CreateAt,
	}, nil
}
