package videoRoutes

import (
	videoHandler "myapi/handlers/video"
	"myapi/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpVideoRoutes(router fiber.Router) {
	public := router.Group("/")

	public.Get("/videos", videoHandler.GetVideos)
	public.Patch("/video/:id", videoHandler.EncodeVideoByID)

	private := router.Group("/private", middleware.AuthMiddleware)

	private.Put("/video/:id", videoHandler.UpdateVideoByID)
	private.Delete("/video/:id", videoHandler.DeleteVideoByID)
	private.Post("/video/:id/comment", videoHandler.PostComment)
	private.Get("/video/:id/comments", videoHandler.GetCommentsByVideoID)
}
