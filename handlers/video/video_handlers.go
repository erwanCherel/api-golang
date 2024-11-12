package videoHandler

import (
	"fmt"
	"myapi/database"
	"myapi/models"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetVideos(c *fiber.Ctx) error {
	name := c.Query("name", "")
	user := c.Query("user", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("perPage", "10"))

	db := database.DB
	query := db.Model(&models.Video{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if user != "" {
		if userID, err := strconv.Atoi(user); err == nil {
			query = query.Where("user_id = ?", userID)
		} else {
			query = query.Joins("JOIN users ON users.id = videos.user_id").Where("users.name LIKE ?", "%"+user+"%")
		}
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * perPage
	var videos []models.Video
	query.Limit(perPage).Offset(offset).Find(&videos)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OK",
		"data":    videos,
		"pager": fiber.Map{
			"current": page,
			"total":   int(total),
		},
	})
}

func EncodeVideoByID(c *fiber.Ctx) error {
	videoID := c.Params("id")
	db := database.DB

	var input struct {
		Format string `json:"format"`
		File   string `json:"file"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	var video models.Video
	err := db.First(&video, videoID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Video not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving video",
		})
	}

	var existingFormat models.VideoFormat
	err = db.Where("video_id = ? AND code = ?", videoID, input.Format).First(&existingFormat).Error
	if err == nil {
		return c.JSON(fiber.Map{
			"message": "OK",
			"data":    video,
		})
	}

	newFormat := models.VideoFormat{
		Code:    input.Format,
		URI:     input.File,
		VideoID: video.ID,
	}

	err = db.Create(&newFormat).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving encoded format",
		})
	}

	err = db.Preload("Formats").First(&video, videoID).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error reloading video",
		})
	}

	return c.JSON(fiber.Map{
		"message": "OK",
		"data":    video,
	})
}

func UpdateVideoByID(c *fiber.Ctx) error {
	videoID := c.Params("id")
	user := c.Locals("user").(models.User)
	db := database.DB

	var input struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	var video models.Video
	err := db.First(&video, videoID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Video not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving video",
		})
	}

	if video.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You are not authorized to update this video",
		})
	}

	video.Name = input.Name

	err = db.Save(&video).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving video",
		})
	}

	return c.JSON(fiber.Map{
		"message": "OK",
		"data": fiber.Map{
			"id":         video.ID,
			"name":       video.Name,
			"source":     video.Source,
			"created_at": video.CreatedAt,
			"views":      video.Views,
			"enabled":    video.Enabled,
			"user_id":    video.UserID,
			"format":     video.FormatMap(),
		},
	})
}

func DeleteVideoByID(c *fiber.Ctx) error {
	videoID := c.Params("id")
	user := c.Locals("user").(models.User)
	db := database.DB

	var video models.Video
	err := db.First(&video, videoID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Video not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving video",
		})
	}

	if video.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You are not authorized to delete this video",
		})
	}

	if err := os.Remove(video.Source); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Error deleting video file: %v", err),
		})
	}

	err = database.DB.Select(clause.Associations).Delete(&video).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting video",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func PostComment(c *fiber.Ctx) error {
	videoID := c.Params("id")
	user := c.Locals("user").(models.User)
	db := database.DB

	var video models.Video
	err := db.First(&video, videoID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Video not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving video",
		})
	}

	var input struct {
		Body string `json:"body"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	comment := models.Comment{
		Body:    input.Body,
		UserID:  user.ID,
		VideoID: video.ID,
	}

	err = db.Create(&comment).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving comment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "OK",
		"data": fiber.Map{
			"id":       comment.ID,
			"body":     comment.Body,
			"user_id":  comment.UserID,
			"video_id": comment.VideoID,
		},
	})
}

func GetCommentsByVideoID(c *fiber.Ctx) error {
	videoID := c.Params("id")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("perPage", "10"))
	db := database.DB

	var comments []models.Comment
	var total int64

	db.Model(&models.Comment{}).Where("video_id = ?", videoID).Count(&total)

	offset := (page - 1) * perPage
	err := db.Where("video_id = ?", videoID).Limit(perPage).Offset(offset).Find(&comments).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving comments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OK",
		"data":    comments,
		"pager": fiber.Map{
			"current": page,
			"total":   int(total),
		},
	})
}
