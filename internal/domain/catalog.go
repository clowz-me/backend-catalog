package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID              uuid.UUID `json:"id"`
	EstablishmentID uuid.UUID `json:"establishment_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	ImageURL        string    `json:"image_url,omitempty"`
	Order           int       `json:"order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Product struct {
	ID              uuid.UUID `json:"id"`
	EstablishmentID uuid.UUID `json:"establishment_id"`
	CategoryID      uuid.UUID `json:"category_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           int64     `json:"price"` // Price in cents to avoid floating point issues
	IsActive        bool      `json:"is_active"`
	ImageURL        string    `json:"image_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CatalogRepository interface {
	CreateCategory(c *Category) error
	GetCategoriesByEstablishment(estID uuid.UUID) ([]*Category, error)
	GetCategoryByID(id uuid.UUID) (*Category, error)
	UpdateCategory(c *Category) error
	DeleteCategory(id uuid.UUID) error
	CreateProduct(p *Product) error
	GetProductByID(id uuid.UUID) (*Product, error)
	UpdateProduct(p *Product) error
	DeleteProduct(id uuid.UUID) error
	GetProductsByCategory(categoryID uuid.UUID) ([]*Product, error)
	GetProductsByEstablishment(estID uuid.UUID) ([]*Product, error)
}
