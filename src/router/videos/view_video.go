package videos

import (
	"os"
	"path/filepath"
	"simon-weij/wayland-recorder-backend/src/database"
	"simon-weij/wayland-recorder-backend/src/router/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func ServeVideoById(ctx fiber.Ctx) error {
	var req struct {
		Id int `params:"id"`
	}

	if err := ctx.Bind().URI(&req); err != nil {
		return fiber.ErrBadRequest
	}

	uid, err := auth.GetUID(ctx)
	if err != nil {
		return err
	}

	video, err := database.GetVideoByID(req.Id, uid)
	if video == nil {
		log.Warn("Video was nil!")
		return fiber.ErrInternalServerError
	}
	if err != nil {
		log.Warn(err)
		return fiber.ErrInternalServerError
	}

	videoPath := getVideoFromHash(video.VideoHash, video.Extension)

	return ctx.SendFile(videoPath)
}

func getVideoFromHash(hash string, extension string) string {
	uploadDir := os.Getenv("UPLOAD_DIR")

	filePath := filepath.Join(
		uploadDir,
		string(hash[0]),
		string(hash[1]),
		string(hash[2]),
		string(hash[3]),
		string(hash)+string(extension),
	)
	return filePath
}
