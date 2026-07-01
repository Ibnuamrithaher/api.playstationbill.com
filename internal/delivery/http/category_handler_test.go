package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	deliveryHTTP "api.poster.com/internal/delivery/http"
	"api.poster.com/internal/domain"
	"github.com/gin-gonic/gin"
)

type mockCategoryService struct {
	CreateFunc func(ctx context.Context, input *domain.CategoryCreateInput, coverImageURL string) (*domain.Category, error)
}

func (m *mockCategoryService) Create(ctx context.Context, input *domain.CategoryCreateInput, coverImageURL string) (*domain.Category, error) {
	return m.CreateFunc(ctx, input, coverImageURL)
}

func TestCategoryHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success create category", func(t *testing.T) {
		mockService := &mockCategoryService{
			CreateFunc: func(ctx context.Context, input *domain.CategoryCreateInput, coverImageURL string) (*domain.Category, error) {
				return &domain.Category{
					ID:          "mock-uuid",
					Name:        input.Name,
					Description: input.Description,
					CoverImage:  coverImageURL,
				}, nil
			},
		}

		handler := deliveryHTTP.NewCategoryHandler(mockService)
		router := gin.Default()
		router.POST("/api/category", handler.Create)

		// Create multipart form request
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add Name field
		_ = writer.WriteField("name", "Action Games")
		// Add Description field
		_ = writer.WriteField("description", "A category for action games")

		// Add CoverImage file
		part, _ := writer.CreateFormFile("cover_image", "test_cover.png")
		_, _ = part.Write([]byte("\x89PNG\r\n\x1a\nfake image data"))

		_ = writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/api/category", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert response
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["success"] != true {
			t.Errorf("Expected success true, got %v", resp["success"])
		}

		data := resp["data"].(map[string]interface{})
		if data["name"] != "Action Games" {
			t.Errorf("Expected category name 'Action Games', got '%s'", data["name"])
		}

		coverURL := data["cover_image"].(string)
		if !strings.HasPrefix(coverURL, "/public/uploads/categories/") {
			t.Errorf("Expected cover_image url prefix, got '%s'", coverURL)
		}

		// Clean up file in local upload folder generated during test
		fileName := filepath.Base(coverURL)
		filePath := filepath.Join("public/uploads/categories", fileName)
		_ = os.Remove(filePath)
	})

	t.Run("validation failure - missing name", func(t *testing.T) {
		mockService := &mockCategoryService{}
		handler := deliveryHTTP.NewCategoryHandler(mockService)
		router := gin.Default()
		router.POST("/api/category", handler.Create)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Name is missing
		_ = writer.WriteField("description", "A category for action games")
		part, _ := writer.CreateFormFile("cover_image", "test_cover.png")
		_, _ = part.Write([]byte("fake image data"))
		_ = writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/api/category", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("validation failure - invalid file extension", func(t *testing.T) {
		mockService := &mockCategoryService{}
		handler := deliveryHTTP.NewCategoryHandler(mockService)
		router := gin.Default()
		router.POST("/api/category", handler.Create)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("name", "Action Games")
		// Using .txt file extension which is not allowed
		part, _ := writer.CreateFormFile("cover_image", "test_cover.txt")
		_, _ = part.Write([]byte("fake file data"))
		_ = writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/api/category", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("validation failure - invalid MIME type (spoofed extension)", func(t *testing.T) {
		mockService := &mockCategoryService{}
		handler := deliveryHTTP.NewCategoryHandler(mockService)
		router := gin.Default()
		router.POST("/api/category", handler.Create)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("name", "Action Games")
		// Using .png extension but sending plain text content
		part, _ := writer.CreateFormFile("cover_image", "test_cover.png")
		_, _ = part.Write([]byte("fake file data which is plain text"))
		_ = writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/api/category", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}
