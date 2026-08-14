package dtos

import (
	"time"

	"github.com/reshap0318/serv-master/internal/models"
)

// CategoryCreateRequest represents the request to create a category.
type CategoryCreateRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

// CategoryUpdateRequest represents the request to update a category.
type CategoryUpdateRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

// CategoryDTO represents category data transfer object.
type CategoryDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToCategoryDTO converts a Category model to CategoryDTO.
func ToCategoryDTO(c *models.Category) CategoryDTO {
	return CategoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// ToCategoryDTOList converts a slice of Category models to CategoryDTOs.
func ToCategoryDTOList(categories []models.Category) []CategoryDTO {
	result := make([]CategoryDTO, len(categories))
	for i, c := range categories {
		result[i] = ToCategoryDTO(&c)
	}
	return result
}
