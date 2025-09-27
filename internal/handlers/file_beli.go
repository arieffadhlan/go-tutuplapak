package handlers

import (
	"fmt"
	"time"
	"tutuplapak-user/internal/services"
	"tutuplapak-user/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type FileBeliHandler interface {
	Post(ctx *fiber.Ctx) error
}

type fileBeliHandler struct {
	fileUseCase *services.FileBeliService
}

func NewFileBeliHandler(fileUseCase *services.FileBeliService) FileBeliHandler {
	return &fileBeliHandler{
		fileUseCase: fileUseCase,
	}
}

func (uc *fileBeliHandler) Post(ctx *fiber.Ctx) error {
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

	// size limits
	const (
		minFileSize = 10 * 1024       // 10 KB
		maxFileSize = 2 * 1024 * 1024 // 2 MB
	)

	if file.Size < minFileSize {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file size is too small (minimum 10KB required)",
		})
	}

	if file.Size > maxFileSize {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file size exceeds the maximum limit of 2MB",
		})
	}

	fileName := file.Filename
	fileType := file.Header.Get("Content-Type")

	if !utils.IsAllowedFileBeliType(fileName, fileType) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file type is not allowed",
		})
	}

	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)

	files, err := uc.fileUseCase.UploadFileBeli(ctx.Context(), file, src, filename)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(files)
}
