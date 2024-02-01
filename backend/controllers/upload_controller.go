package controllers

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/file-uploader/backend/configs"
	"github.com/lielalmog/file-uploader/backend/models"
	"github.com/lielalmog/file-uploader/backend/services"
)

type UploadController interface {
	StartUpload(c *fiber.Ctx) error
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
	fileMetadata := new(models.UploadFileMetadateDTO)
	if err := c.BodyParser(fileMetadata); err != nil {
		return err
	}

	if err := configs.GetValidator().Struct(fileMetadata); err != nil {
		return err
	}

	id, err := u.uploadService.SaveFileMetadata(c.Context(), fileMetadata)
	if err != nil {
		return err
	}

	sasToken, err := u.uploadService.GenerateSasToken(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id":        id,
		"signedUrl": sasToken,
	})
}

func (u *uploadControllerImpl) CompleteUpload(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Var(id, "min=1"); err != nil {
		return fiber.ErrBadRequest
	}

	if err := u.uploadService.CompleteUploadEvent(c.Context(), &id); err != nil {
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
