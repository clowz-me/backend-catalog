package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/clowz-me/catalog/internal/domain"
)

type postgresCatalogRepository struct {
	db *sql.DB
}

func NewPostgresCatalogRepository(db *sql.DB) domain.CatalogRepository {
	return &postgresCatalogRepository{db: db}
}

func (r *postgresCatalogRepository) CreateCategory(c *domain.Category) error {
	query := `
		INSERT INTO categories (id, establishment_id, name, description, "order", created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, c.ID, c.EstablishmentID, c.Name, c.Description, c.Order, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *postgresCatalogRepository) GetCategoriesByEstablishment(estID uuid.UUID) ([]*domain.Category, error) {
	query := `SELECT id, establishment_id, name, description, "order", created_at, updated_at FROM categories WHERE establishment_id = $1 ORDER BY "order" ASC`
	rows, err := r.db.Query(query, estID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.EstablishmentID, &c.Name, &c.Description, &c.Order, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, &c)
	}
	return results, nil
}

func (r *postgresCatalogRepository) CreateProduct(p *domain.Product) error {
	query := `
		INSERT INTO products (id, establishment_id, category_id, name, description, price, is_active, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, p.ID, p.EstablishmentID, p.CategoryID, p.Name, p.Description, p.Price, p.IsActive, p.ImageURL, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *postgresCatalogRepository) GetProductsByCategory(categoryID uuid.UUID) ([]*domain.Product, error) {
	query := `
		SELECT id, establishment_id, category_id, name, description, price, is_active, image_url, created_at, updated_at 
		FROM products WHERE category_id = $1
	`
	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.EstablishmentID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.IsActive, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, &p)
	}
	return results, nil
}

func (r *postgresCatalogRepository) GetProductsByEstablishment(estID uuid.UUID) ([]*domain.Product, error) {
	query := `
		SELECT id, establishment_id, category_id, name, description, price, is_active, image_url, created_at, updated_at 
		FROM products WHERE establishment_id = $1
	`
	rows, err := r.db.Query(query, estID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*domain.Product{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var results []*domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.EstablishmentID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.IsActive, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, &p)
	}
	return results, nil
}

func (r *postgresCatalogRepository) GetCategoryByID(id uuid.UUID) (*domain.Category, error) {
	query := `SELECT id, establishment_id, name, description, "order", created_at, updated_at FROM categories WHERE id = $1`
	var c domain.Category
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.EstablishmentID, &c.Name, &c.Description, &c.Order, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresCatalogRepository) UpdateCategory(c *domain.Category) error {
	query := `UPDATE categories SET name = $1, description = $2, "order" = $3, updated_at = $4 WHERE id = $5`
	res, err := r.db.Exec(query, c.Name, c.Description, c.Order, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *postgresCatalogRepository) DeleteCategory(id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *postgresCatalogRepository) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, establishment_id, category_id, name, description, price, is_active, image_url, created_at, updated_at 
		FROM products WHERE id = $1
	`
	var p domain.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.EstablishmentID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.IsActive, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresCatalogRepository) UpdateProduct(p *domain.Product) error {
	query := `
		UPDATE products 
		SET category_id = $1, name = $2, description = $3, price = $4, is_active = $5, image_url = $6, updated_at = $7
		WHERE id = $8
	`
	res, err := r.db.Exec(query, p.CategoryID, p.Name, p.Description, p.Price, p.IsActive, p.ImageURL, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *postgresCatalogRepository) DeleteProduct(id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("product not found")
	}
	return nil
}
