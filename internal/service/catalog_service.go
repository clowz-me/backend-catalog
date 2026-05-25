package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/clowz-me/catalog/internal/domain"
)

type CatalogService struct {
	repo domain.CatalogRepository
}

func NewCatalogService(repo domain.CatalogRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) AddCategory(estID uuid.UUID, name, description string, order int) (*domain.Category, error) {
	c := &domain.Category{
		ID:              uuid.New(),
		EstablishmentID: estID,
		Name:            name,
		Description:     description,
		Order:           order,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.CreateCategory(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CatalogService) GetCategories(estID uuid.UUID) ([]*domain.Category, error) {
	return s.repo.GetCategoriesByEstablishment(estID)
}

func (s *CatalogService) AddProduct(estID, catID uuid.UUID, name, description string, price int64, imgURL string) (*domain.Product, error) {
	now := time.Now()
	p := &domain.Product{
		ID:              uuid.New(),
		EstablishmentID: estID,
		CategoryID:      catID,
		Name:            name,
		Description:     description,
		Price:           price,
		IsActive:        true,
		ImageURL:        imgURL,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.CreateProduct(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *CatalogService) GetProductsForEstablishment(estID uuid.UUID) ([]*domain.Product, error) {
	return s.repo.GetProductsByEstablishment(estID)
}

func (s *CatalogService) UpdateCategory(id uuid.UUID, name, description string, order int) error {
	cat, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return err
	}
	cat.Name = name
	cat.Description = description
	cat.Order = order
	return s.repo.UpdateCategory(cat)
}

func (s *CatalogService) DeleteCategory(id uuid.UUID) error {
	products, err := s.repo.GetProductsByCategory(id)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return errors.New("não é possível excluir a categoria pois existem produtos associados a ela")
	}
	return s.repo.DeleteCategory(id)
}

func (s *CatalogService) UpdateProduct(id uuid.UUID, catID uuid.UUID, name, description string, price int64, imgURL string, isActive bool) error {
	prod, err := s.repo.GetProductByID(id)
	if err != nil {
		return err
	}
	prod.CategoryID = catID
	prod.Name = name
	prod.Description = description
	prod.Price = price
	prod.ImageURL = imgURL
	prod.IsActive = isActive
	prod.UpdatedAt = time.Now()
	
	return s.repo.UpdateProduct(prod)
}

func (s *CatalogService) DeleteProduct(id uuid.UUID) error {
	return s.repo.DeleteProduct(id)
}
