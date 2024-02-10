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
	jwtService    services.JWTService
}

var (
	initUploadControllerOnce sync.Once
	uploadController         *uploadControllerImpl
)

func (u *uploadControllerImpl) StartUpload(c *fiber.Ctx) error {
	fileMetadataDTO := new(models.UploadFileMetadateDTO)
	if err := c.BodyParser(fileMetadataDTO); err != nil {
		return err
	}

	if err := configs.GetValidator().Struct(fileMetadataDTO); err != nil {
		return err
	}

	claims := c.Locals("userClaims").(*configs.CustomJwtClaims)

	fileMetadata := &models.FileMetadata{
		UserID:                claims.Id,
		UploadFileMetadateDTO: *fileMetadataDTO,
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
		jwtService:    services.GetJWTService(),
	}
}

func GetUploadController() UploadController {
	initUploadControllerOnce.Do(func() {
		uploadController = newUploadController()
	})

	return uploadController
}
