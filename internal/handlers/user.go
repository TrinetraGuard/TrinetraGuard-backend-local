package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct{}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Mock data for demonstration
	users := []gin.H{
		{"id": 1, "name": "John Doe", "email": "john@example.com"},
		{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

// GetUser returns a specific user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	// Mock user data
	user := gin.H{
		"id":    id,
		"name":  "John Doe",
		"email": "john@example.com",
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock created user
	user := gin.H{
		"id":    3,
		"name":  req.Name,
		"email": req.Email,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email" binding:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock updated user
	user := gin.H{
		"id":    id,
		"name":  req.Name,
		"email": req.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      id,
	})
}

// GetExamples returns example data
func (h *UserHandler) GetExamples(c *gin.Context) {
	examples := []gin.H{
		{"id": 1, "title": "Example 1", "description": "This is example 1"},
		{"id": 2, "title": "Example 2", "description": "This is example 2"},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Examples retrieved successfully",
		"data":    examples,
	})
}

// CreateExample creates a new example
func (h *UserHandler) CreateExample(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	example := gin.H{
		"id":          3,
		"title":       req.Title,
		"description": req.Description,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Example created successfully",
		"data":    example,
	})
}
