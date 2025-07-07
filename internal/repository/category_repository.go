package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nordicmanx/videos-api/internal/models"
)

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	category := &models.Category{
		Name: name,
	}

	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, created_at, updated_at`
	err := r.DB.QueryRow(ctx, query, name).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	query := `SELECT id, name, created_at, updated_at FROM categories ORDER BY name ASC`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`
	err := r.DB.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, id uuid.UUID, name string) error {
	query := `UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(ctx, query, name, time.Now(), id)
	return err
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, id)
	return err
}
