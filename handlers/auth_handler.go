package handlers

import (
	"auth-app/config"
	"auth-app/models"
	"auth-app/utils"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// =====================
// REGISTER
// =====================
func Register(c *gin.Context) {
	var user models.User

	// Bind JSON request
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Basic validation
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password required",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Insert user into DB
	err = config.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email,
		string(hashedPassword),
	).Scan(&user.ID)

	// Handle DB errors
	if err != nil {

		// Duplicate email
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already registered",
			})
			return
		}

		// Other DB errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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
func Login(c *gin.Context) {
	var user models.User
	var storedPassword string

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

	// Get user password from DB
	err := config.DB.QueryRow(
		"SELECT password FROM users WHERE email=$1",
		user.Email,
	).Scan(&storedPassword)

	// Handle errors
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(
		[]byte(storedPassword),
		[]byte(user.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong password",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
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
func GetUsers(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}