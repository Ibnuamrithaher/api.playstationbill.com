package domain

import (
	"context"
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

type CategoryCreateInput struct {
	Name        string
	Description string
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
}

type CategoryService interface {
	Create(ctx context.Context, input *CategoryCreateInput, coverImageURL string) (*Category, error)
}
