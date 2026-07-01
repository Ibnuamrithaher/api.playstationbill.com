package http

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"api.poster.com/internal/domain"
	"api.poster.com/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService domain.CategoryService
}

func NewCategoryHandler(categoryService domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

type createCategoryRequest struct {
	Name        string                `form:"name" binding:"required,max=255"`
	Description string                `form:"description"`
	CoverImage  *multipart.FileHeader `form:"cover_image" binding:"required"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req createCategoryRequest

	// Bind form data (multipart form-data)
	if err := c.ShouldBind(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Validate cover_image extension
	ext := filepath.Ext(req.CoverImage.Filename)
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !validExtensions[strings.ToLower(ext)] {
		response.SendError(c, http.StatusBadRequest, "Validation failed", "Invalid cover_image extension. Only jpg, jpeg, png are allowed.")
		return
	}

	// Validate MIME type (Content-Type)
	file, err := req.CoverImage.Open()
	if err != nil {
		response.SendError(c, http.StatusBadRequest, "Validation failed", "Failed to open cover_image.")
		return
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		response.SendError(c, http.StatusInternalServerError, "Failed to read cover_image", err.Error())
		return
	}

	contentType := http.DetectContentType(buffer[:n])
	validContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !validContentTypes[contentType] {
		response.SendError(c, http.StatusBadRequest, "Validation failed", "Invalid cover_image content type. Only jpg, jpeg, png are allowed.")
		return
	}

	// Validate cover_image size (Max 2MB)
	if req.CoverImage.Size > 2*1024*1024 {
		response.SendError(c, http.StatusBadRequest, "Validation failed", "cover_image size exceeds maximum limit of 2MB.")
		return
	}

	// Ensure upload directory exists (0755 permission instead of os.ModePerm)
	uploadDir := "public/uploads/categories"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.SendError(c, http.StatusInternalServerError, "Failed to create upload directory", err.Error())
		return
	}

	// Generate unique filename for the cover image
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Save file locally
	if err := c.SaveUploadedFile(req.CoverImage, filePath); err != nil {
		response.SendError(c, http.StatusInternalServerError, "Failed to save cover_image", err.Error())
		return
	}

	// Create relative URL path
	// Replace any backslashes (Windows) with forward slashes for URLs
	coverImageURL := fmt.Sprintf("/public/uploads/categories/%s", fileName)

	// Map to service DTO
	input := &domain.CategoryCreateInput{
		Name:        req.Name,
		Description: req.Description,
	}

	// Call service layer to save category
	category, err := h.categoryService.Create(c.Request.Context(), input, coverImageURL)
	if err != nil {
		// Clean up the uploaded file if DB save fails
		_ = os.Remove(filePath)
		response.SendError(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Category created successfully", category)
}
