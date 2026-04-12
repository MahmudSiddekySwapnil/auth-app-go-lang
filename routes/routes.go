package routes

import (
	"auth-app/handlers"   // Import handlers package containing Register, Login, GetUsers functions
	"auth-app/middleware" // Import middleware package for authentication/authorization logic

	"github.com/gin-gonic/gin" // Import Gin web framework for HTTP routing
)

// SetupRoutes configures all HTTP routes and their handlers for the application
// Parameter r *gin.Engine: The Gin router instance that will handle all HTTP requests
func SetupRoutes(r *gin.Engine) {
	// Create a route group with prefix "/api" - all routes defined under this group will start with /api
	// This helps organize routes and apply group-level configurations
	api := r.Group("/api")

	// =====================
	// PUBLIC ROUTES (No Authentication Required)
	// =====================
	// These routes are accessible to anyone without needing a JWT token

	// POST /api/register - Handles user registration/signup
	// Registers a new user by accepting email and password, then stores in database with bcrypt hashing
	api.POST("/register", handlers.Register)

	// POST /api/login - Handles user login/authentication
	// Validates user credentials and returns a JWT token if authentication is successful
	api.POST("/login", handlers.Login)

	// GET /api/users - COMMENTED OUT (previously used for public access)
	// This was disabled because user listing should be protected to maintain privacy
	//api.GET("users", handlers.GetUsers)

	// =====================
	// PROTECTED ROUTES (Authentication Required)
	// =====================
	// Create a sub-group under /api with empty path prefix
	// This group is used to apply middleware to specific routes
	protected := api.Group("/")

	// Apply AuthMiddleware to all routes in this protected group
	// AuthMiddleware checks if request has valid JWT token before allowing access
	// If token is missing or invalid, request is rejected with 401 Unauthorized
	protected.Use(middleware.AuthMiddleware())

	// GET /api/users - Handles fetching list of all users (Protected)
	// Only accessible with valid JWT token. Returns id and email of all registered users
	protected.GET("/users", handlers.GetUsers)
}
