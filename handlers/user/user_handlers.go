package userHandler

import (
	"myapi/database"
	"myapi/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetUsers handles the endpoint to get all users
func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []models.User

	// Query the database to find all users
	db.Find(&users)

	// If no users are found, return a 404 error with a message
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No users present", "data": nil})
	}

	// Return the list of users if found
	return c.JSON(fiber.Map{"status": "success", "message": "Users Found", "data": users})
}

// CreateUser handles the endpoint to create a new user
func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(models.User)

	// Parse the incoming request body to populate the user model
	err := c.BodyParser(user)
	if err != nil {
		// If parsing fails, return a 500 error with a message
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Attempt to create the user in the database
	err = db.Create(&user).Error
	if err != nil {
		// If creating the user fails, return a 500 error with a message
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}

	// Return the created user with a success message
	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": user})
}

// GetUser handles the endpoint to get a user by their ID
func GetUser(c *fiber.Ctx) error {
	db := database.DB
	var user models.User

	// Read the param 'id' from the URL and convert it to an int
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		// If the ID is invalid, return a 400 error with a message
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid user ID", "data": nil})
	}

	// Find the user in the database with the given userID
	result := db.First(&user, "id = ?", userID)

	// If no such user is found, return a 404 error with a message
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found", "data": nil})
	}

	// Return the found user with a success message
	return c.JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}
