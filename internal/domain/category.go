package domain

import (
	"context"
	"mime/multipart"
	"time"
)

type Category struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	CoverImage  string    `gorm:"column:cover_image;type:varchar(255)" json:"cover_image"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

type CategoryCreateRequest struct {
	Name        string                `form:"name" binding:"required,max=255"`
	Description string                `form:"description"`
	CoverImage  *multipart.FileHeader `form:"cover_image" binding:"required"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
}

type CategoryService interface {
	Create(ctx context.Context, req *CategoryCreateRequest, coverImageURL string) (*Category, error)
}
