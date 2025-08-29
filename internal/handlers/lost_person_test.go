package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLostPersonHandler_CreateLostPerson(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	handler := NewLostPersonHandler()

	// Setup route
	router.POST("/lost-persons", handler.CreateLostPerson)

	// Create a test image file
	testImagePath := "test_image.jpg"
	testImageContent := []byte("fake image content")
	err := os.WriteFile(testImagePath, testImageContent, 0644)
	assert.NoError(t, err)
	defer os.Remove(testImagePath)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	writer.WriteField("name", "John Doe")
	writer.WriteField("aadhar_number", "123456789012")
	writer.WriteField("contact_number", "9876543210")
	writer.WriteField("place_lost", "Mumbai Central")
	writer.WriteField("permanent_address", "123 Main St, Mumbai")

	// Add file with proper content type
	part, err := writer.CreateFormFile("image", "test_image.jpeg")
	assert.NoError(t, err)
	part.Write(testImageContent)

	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", "/lost-persons", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Print response for debugging
	t.Logf("Response Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert response structure
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Lost person report created successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Clean up
	os.Remove("database.json")
	os.RemoveAll("uploads")
}

func TestLostPersonHandler_GetAllLostPersons(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	handler := NewLostPersonHandler()

	// Setup route
	router.GET("/lost-persons", handler.GetAllLostPersons)

	// Create request
	req, err := http.NewRequest("GET", "/lost-persons", nil)
	assert.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert response structure
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Lost person reports retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Clean up
	os.Remove("database.json")
	os.RemoveAll("uploads")
}
