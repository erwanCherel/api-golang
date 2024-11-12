package userHandler

import (
	"fmt"
	"myapi/database"
	"myapi/models"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []models.User

	pageStr := c.Query("page", "1")
	perPageStr := c.Query("per_page", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid page number",
			"data":    nil,
		})
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid per_page number",
			"data":    nil,
		})
	}

	var totalUsers int64
	db.Model(&models.User{}).Count(&totalUsers)

	totalPages := (int(totalUsers) + perPage - 1) / perPage

	if page > totalPages {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Page number exceeds total pages",
			"data":    nil,
		})
	}

	offset := (page - 1) * perPage
	result := db.Order("id DESC").Limit(perPage).Offset(offset).Find(&users)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Not found", "data": nil})
	}

	if len(users) == 0 {
		return c.JSON(fiber.Map{
			"message": "Ok",
			"data":    []models.User{},
			"pager": fiber.Map{
				"current": page,
				"total":   totalPages,
			},
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ok",
		"data":    users,
		"pager": fiber.Map{
			"current": page,
			"total":   totalPages,
		},
	})
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(models.User)

	err := c.BodyParser(user)
	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"code":    fiber.StatusBadRequest,
			"data":    err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not hash password",
			"code":    fiber.StatusInternalServerError,
			"data":    err.Error(),
		})
	}

	user.Password = string(hashedPassword)

	err = db.Create(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create user",
			"code":    fiber.StatusInternalServerError,
			"data":    err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Created user",
		"data":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	db := database.DB
	var user models.User

	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {

		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid user ID", "data": nil})
	}

	result := db.First(&user, "id = ?", userID)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID",
			"data":    nil,
		})
	}

	var user models.User
	result := database.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found",
			"data":    nil,
		})
	}

	if err := database.DB.Select(clause.Associations).Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete user",
			"data":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
		"data":    nil,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID",
			"data":    nil,
		})
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized, token is missing",
			"data":    nil,
		})
	}

	var input models.User
	err = c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	var user models.User
	result := database.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    nil,
		})
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to hash password",
				"data":    err.Error(),
			})
		}
		input.Password = string(hashedPassword)
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Pseudo != "" {
		user.Pseudo = input.Pseudo
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		user.Password = input.Password
	}

	err = database.DB.Save(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update user",
			"data":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ok",
		"data":    user,
	})
}

func CreateVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
			"data":    nil,
		})
	}

	file, err := c.FormFile("source")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No video file uploaded",
			"data":    err.Error(),
		})
	}

	videoName := c.FormValue("name")
	if videoName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Video name is required",
			"data":    nil,
		})
	}

	publicPath := "./public/videos/"
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(publicPath, fileName)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save video",
			"data":    err.Error(),
		})
	}

	video := models.Video{
		Name:      videoName,
		Source:    filePath,
		UserID:    userID,
		CreatedAt: time.Now(),
		Enabled:   true,
		Views:     0,
	}

	if err := database.DB.Create(&video).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save video in database",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "OK",
		"data":    video,
	})
}

func GetVideosByUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	db := database.DB

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	perPage, err := strconv.Atoi(c.Query("perPage", "10"))
	if err != nil || perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	var videos []models.Video
	var totalVideos int64

	err = db.Model(&models.Video{}).Where("user_id = ?", userID).Count(&totalVideos).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching video count",
		})
	}

	err = db.Preload("Formats").Preload("Comments").Where("user_id = ?", userID).
		Limit(perPage).Offset(offset).Find(&videos).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching videos",
		})
	}

	totalPages := int(totalVideos / int64(perPage))
	if totalVideos%int64(perPage) > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"message": "OK",
		"data":    videos,
		"pager": fiber.Map{
			"current": page,
			"total":   totalPages,
		},
	})
}
