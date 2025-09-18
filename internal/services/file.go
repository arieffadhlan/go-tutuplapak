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
	minioUploader "tutuplapak-user/internal/utils"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type FileService struct {
	config         config.Config
	minio          *minio.Client
	fileRepository repository.FileRepositoryInterface
}

func NewFileService(config config.Config, fileRepository repository.FileRepositoryInterface) *FileService {
	minioConfig := &minioUploader.MinioConfig{
		AccessKeyID:     config.Minio.AccessKeyID,
		SecretAccessKey: config.Minio.SecretAccessKey,
		UseSSL:          config.Minio.UseSSL,
		Endpoint:        config.Minio.Endpoint,
	}

	minioClient, _ := minioUploader.NewUploader(minioConfig)

	return &FileService{
		config:         config,
		minio:          minioClient,
		fileRepository: fileRepository,
	}
}

func (uc *FileService) GetFileById(ctx context.Context, id uuid.UUID) (dto.FileResponse, error) {
	res, err := uc.fileRepository.GetFileById(ctx, id)
	if err != nil {
		return dto.FileResponse{}, err
	} else {
		return res, nil
	}
}

func (uc *FileService) UploadFile(ctx context.Context, file *multipart.FileHeader, src multipart.File, fileName string) (dto.FileResponse, error) {
	bucketName := uc.config.Minio.BucketName
	fileNameCompress := minioUploader.AddFileNameSuffix(fileName)
	// compress file
	srcForCompress, err := file.Open()
	if err != nil {
		return dto.FileResponse{}, err
	}
	defer srcForCompress.Close()

	img, _, _ := image.Decode(srcForCompress)

	compressedBytes, contentType, _ := minioUploader.CompressThumbnail(file, img, 10)

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
		return dto.FileResponse{}, err
	}

	// upload for the thumbnail
	_, err = uc.minio.PutObject(
		ctx,
		bucketName,
		fileNameCompress,
		bytes.NewReader(compressedBytes),
		int64(len(compressedBytes)),
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		return dto.FileResponse{}, err
	}

	// get url original
	url, err := uc.minio.PresignedGetObject(ctx, bucketName, fileName, time.Hour*24*7, nil)
	if err != nil {
		return dto.FileResponse{}, err
	}
	// get url thumbnail
	thumbUrl, err := uc.minio.PresignedGetObject(ctx, bucketName, fileNameCompress, time.Hour*24*7, nil)
	if err != nil {
		return dto.FileResponse{}, err
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
		return dto.FileResponse{}, err
	}

	return dto.FileResponse{
		ID:           res.ID,
		Url:          res.Url,
		ThumbnailUrl: res.ThumbnailUrl,
	}, nil
}
