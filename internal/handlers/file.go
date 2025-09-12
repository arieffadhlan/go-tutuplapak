package handlers

import (
	"fmt"
	"time"
	"tutuplapak-user/internal/services"
	"tutuplapak-user/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type FileHandler interface {
	Post(ctx *fiber.Ctx) error
}

type fileHandler struct {
	fileUseCase services.UseCase
}

func NewFileHandler(fileUseCase services.UseCase) FileHandler {
	return &fileHandler{
		fileUseCase: fileUseCase,
	}
}

func (uc *fileHandler) Post(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	defer src.Close()

	if file.Size > (100 * 1024) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file exceeds the maximum limit of 100KiB",
		})
	}

	fileName := file.Filename
	fileType := file.Header.Get("Content-Type")

	if !utils.IsAllowedFileType(fileName, fileType) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file type is not allowed",
		})
	}

	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)

	files, err := uc.fileUseCase.UploadFile(ctx.Context(), file, src, filename)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(files)
}

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
