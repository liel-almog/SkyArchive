package controllers

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/be-file-streaming/configs"
	"github.com/lielalmog/be-file-streaming/models"
	"github.com/lielalmog/be-file-streaming/services"
)

type UploadController interface {
	StartUpload(c *fiber.Ctx) error
	UploadChunk(c *fiber.Ctx) error
	CompleteUpload(c *fiber.Ctx) error
}

type uploadControllerImpl struct {
	uploadService services.UploadService
}

var (
	initUploadControllerOnce sync.Once
	uploadController         *uploadControllerImpl
)

func (u *uploadControllerImpl) StartUpload(c *fiber.Ctx) error {
	fileMetadata := new(models.FileMetadateDTO)
	if err := c.BodyParser(fileMetadata); err != nil {
		return err
	}

	if err := configs.GetValidator().Struct(fileMetadata); err != nil {
		return err
	}

	id, err := u.uploadService.StartUpload(fileMetadata)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id": id,
	})
}

func (u *uploadControllerImpl) UploadChunk(c *fiber.Ctx) error {
	// Convert to types
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	chunkIndex, err := strconv.Atoi(c.Params("chunkIndex"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	// Validate id and chunkIndex
	if err := configs.GetValidator().Var(id, "min=1"); err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Var(chunkIndex, "min=0"); err != nil {
		return fiber.ErrBadRequest
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	u.uploadService.UploadChunk(fileHeader, id, chunkIndex)

	return c.SendStatus(fiber.StatusOK)
}

func (u *uploadControllerImpl) CompleteUpload(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Var(id, "min=1"); err != nil {
		return fiber.ErrBadRequest
	}

	if err = u.uploadService.CombineChunksAndUploadToPermanent(id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func newUploadController() *uploadControllerImpl {
	return &uploadControllerImpl{
		uploadService: services.GetUploadService(),
	}
}

func GetUploadController() UploadController {
	initUploadControllerOnce.Do(func() {
		uploadController = newUploadController()
	})

	return uploadController
}
