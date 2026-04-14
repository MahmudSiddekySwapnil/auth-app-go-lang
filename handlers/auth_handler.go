package handlers

import (
	"auth-app/models"
	"auth-app/services"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(userService services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// =====================
// REGISTER
// =====================
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User

	// Bind JSON request
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// LOG: Print incoming request data
	log.Printf("Register Request Received:\n  Email: %s\n  Name: %s\n  Password: %s\n", user.Email, user.Name, user.Password)

	// Basic validation
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password required",
		})
		return
	}

	// Email format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		log.Printf("⚠️  Invalid email format: %s\n", user.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	err := h.userService.Register(&user)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Email already registered" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}

// =====================
// LOGIN
// =====================
func (h *AuthHandler) Login(c *gin.Context) {
	var user models.User

	// Bind JSON
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate input
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password required",
		})
		return
	}

	token, err := h.userService.Login(user.Email, user.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Invalid email" || err.Error() == "Wrong password" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// =====================
// GET USERS
// =====================
func (h *AuthHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if users == nil {
		users = []models.User{}
	}

	c.JSON(http.StatusOK, users)
}

// =====================
// GET USER BY ID
// =====================
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "User not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// =====================
// DELETE USER
// =====================
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "User not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// =====================
// UPDATE USER
// =====================
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.userService.UpdateUser(id, &user)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "User not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}
