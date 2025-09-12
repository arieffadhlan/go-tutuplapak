package handlers

// import (
// 	"context"
// 	"mime/multipart"
// 	"tutuplapak-user/internal/config"
// 	"tutuplapak-user/internal/dto"

// 	minioUploader "github.com/arieffadhlan/go-fitbyte/internal/pkg/minio"
// 	"github.com/minio/minio-go/v7"
// )

// type UseCase interface {
// 	UploadFile(context.Context, *multipart.FileHeader, multipart.File, string)(dto.File, error)
// }

// type useCase struct{
// 	config config.Config
// 	minio *minio.Client
// }

// func NewUseCase(config config.Config) UseCase {
// 	minioConfig := &minioUploader.MinioConfig{
// 		AccessKeyID:     config.,
// 		SecretAccessKey: config.Minio.SecretAccessKey,
// 		UseSSL:          config.Minio.UseSSL,
// 		Endpoint:        config.Minio.Endpoint,
// 	}

// 	minioClient, _ := minioUploader.NewUploader(minioConfig)

// 	return &useCase{
// 		config: config,
// 		minio:  minioClient,
// 	}
// }
