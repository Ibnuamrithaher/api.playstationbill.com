package service

import (
	"context"
	"time"

	"api.poster.com/internal/domain"
	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryService(categoryRepo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) Create(ctx context.Context, input *domain.CategoryCreateInput, coverImageURL string) (*domain.Category, error) {
	now := time.Now()
	category := &domain.Category{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		CoverImage:  coverImageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
