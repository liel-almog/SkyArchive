package controllers

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/lielalmog/SkyArchive/backend/configs"
	"github.com/lielalmog/SkyArchive/backend/models"
	"github.com/lielalmog/SkyArchive/backend/services"
)

type FileController interface {
	StartFileUpload(c *fiber.Ctx) error
	CompleteFileUpload(c *fiber.Ctx) error
	GetUserFiles(c *fiber.Ctx) error
	UpdateFavorite(c *fiber.Ctx) error
	UpdateDisplayName(c *fiber.Ctx) error
}

type fileControllerImpl struct {
	fileService services.FileService
}

var (
	initFileControllerOnce sync.Once
	fileController         *fileControllerImpl
)

func (u *fileControllerImpl) StartFileUpload(c *fiber.Ctx) error {
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

	id, err := u.fileService.SaveFileMetadata(c.Context(), fileMetadata)
	if err != nil {
		return err
	}

	sasToken, err := u.fileService.GenerateSasToken(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id":        id,
		"signedUrl": sasToken,
	})
}

func (u *fileControllerImpl) CompleteFileUpload(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Var(id, "min=1"); err != nil {
		return fiber.ErrBadRequest
	}

	if err := u.fileService.CompleteFileUploadEvent(c.Context(), &id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (u *fileControllerImpl) GetUserFiles(c *fiber.Ctx) error {
	claims, ok := c.Locals("userClaims").(*configs.CustomJwtClaims)
	if !ok {
		return fiber.ErrUnauthorized
	}

	files, err := u.fileService.GetUserFiles(c.Context(), &claims.Id)
	if err != nil {
		return err
	}

	return c.JSON(files)
}

func (u *fileControllerImpl) UpdateFavorite(c *fiber.Ctx) error {
	updateFavoriteDTO := new(models.UpdateFavoriteDTO)
	if err := c.BodyParser(updateFavoriteDTO); err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Struct(updateFavoriteDTO); err != nil {
		return fiber.ErrBadRequest
	}

	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Var(id, "min=1"); err != nil {
		return fiber.ErrBadRequest
	}

	claims := c.Locals("userClaims").(*configs.CustomJwtClaims)

	if err := u.fileService.UpdateFavorite(c.Context(), &id, &claims.Id, updateFavoriteDTO); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (u *fileControllerImpl) UpdateDisplayName(c *fiber.Ctx) error {
	updateDisplayNameDTO := new(models.UpdateDisplayNameDTO)
	if err := c.BodyParser(updateDisplayNameDTO); err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Struct(updateDisplayNameDTO); err != nil {
		return fiber.ErrBadRequest
	}

	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := configs.GetValidator().Var(id, "min=1"); err != nil {
		return fiber.ErrBadRequest
	}

	claims := c.Locals("userClaims").(*configs.CustomJwtClaims)

	if err := u.fileService.UpdateDisplayName(c.Context(), &id, &claims.Id, updateDisplayNameDTO); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func newFileController() *fileControllerImpl {
	return &fileControllerImpl{
		fileService: services.GetFileService(),
	}
}

func GetFileController() FileController {
	initFileControllerOnce.Do(func() {
		fileController = newFileController()
	})

	return fileController
}
