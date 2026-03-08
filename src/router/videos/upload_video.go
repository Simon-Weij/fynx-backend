package videos

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"simon-weij/wayland-recorder-backend/src/database"
	"simon-weij/wayland-recorder-backend/src/router/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func UploadVideo(ctx fiber.Ctx) error {
	var req struct {
		Title     string `form:"title"`
		IsPrivate *bool  `form:"is_private"`
	}

	if err := ctx.Bind().Form(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Title == "" {
		return fiber.ErrBadRequest
	}

	isPrivate := true
	if req.IsPrivate != nil {
		isPrivate = *req.IsPrivate
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		log.Error("upload FormFile:", err)
		return fiber.ErrInternalServerError
	}

	hashSum, err := calculateHash(fileHeader)
	if err != nil {
		return err
	}

	extension := filepath.Ext(fileHeader.Filename)

	fullLocation := getStoragePath(hashSum, extension)

	if err := saveToDisk(ctx, fileHeader, fullLocation); err != nil {
		return err
	}

	uid, err := auth.GetUID(ctx)
	if err != nil {
		return err
	}

	database.InsertVideo(uid, req.Title, hashSum, extension, isPrivate)

	return ctx.SendString("File uploaded successfully")
}

func calculateHash(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		log.Error("upload open:", err)
		return "", fiber.ErrInternalServerError
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		log.Error("upload hash:", err)
		return "", fiber.ErrInternalServerError
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func getStoragePath(hashSum string, extension string) string {
	uploadLocation := os.Getenv("UPLOAD_DIR")
	firstFolder := filepath.Join(uploadLocation, hashSum[0:1], hashSum[1:2])
	return filepath.Join(firstFolder, hashSum[2:3], hashSum[3:4], hashSum+extension)
}

func saveToDisk(ctx fiber.Ctx, fileHeader *multipart.FileHeader, fullLocation string) error {
	if err := os.MkdirAll(filepath.Dir(fullLocation), 0750); err != nil {
		log.Error("upload mkdir:", err)
		return fiber.ErrInternalServerError
	}
	if err := ctx.SaveFile(fileHeader, fullLocation); err != nil {
		log.Error("upload save:", err)
		return fiber.ErrInternalServerError
	}
	return nil
}
