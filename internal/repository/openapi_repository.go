package repository

import (
	"context"
	"fmt"
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
	"time"

	"github.com/uptrace/bun"
)

// OpenAPIRepository handles database operations for OpenAPI specifications
type OpenAPIRepository struct {
	db *bun.DB
}

func NewOpenAPIRepository(db *bun.DB) interfaces.OpenAPIRepository {
	return &OpenAPIRepository{db: db}
}

// Create adds a new OpenAPI specification to the database
func (r *OpenAPIRepository) Create(ctx context.Context, spec *models.OpenAPISpec) error {
	spec.CreatedAt = time.Now()
	spec.UpdatedAt = time.Now()

	_, err := r.db.NewInsert().
		Model(spec).
		Returning("id").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create OpenAPI spec: %w", err)
	}

	return nil
}

// GetByID retrieves an OpenAPI specification by its ID
func (r *OpenAPIRepository) GetByID(ctx context.Context, id int64) (*models.OpenAPISpec, error) {
	spec := &models.OpenAPISpec{}
	err := r.db.NewSelect().
		Model(spec).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI spec by ID: %w", err)
	}

	return spec, nil
}

// GetByTitle retrieves an OpenAPI specification by its title
func (r *OpenAPIRepository) GetByTitle(ctx context.Context, title string) (*models.OpenAPISpec, error) {
	spec := &models.OpenAPISpec{}
	err := r.db.NewSelect().
		Model(spec).
		Where("title = ?", title).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI spec by title: %w", err)
	}

	return spec, nil
}

// List returns all OpenAPI specifications with pagination
func (r *OpenAPIRepository) List(ctx context.Context, offset, limit int) ([]*models.OpenAPISpec, error) {
	var specs []*models.OpenAPISpec
	err := r.db.NewSelect().
		Model(&specs).
		OrderExpr("created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list OpenAPI specs: %w", err)
	}

	return specs, nil
}

// Update modifies an existing OpenAPI specification
func (r *OpenAPIRepository) Update(ctx context.Context, spec *models.OpenAPISpec) error {
	spec.UpdatedAt = time.Now()

	_, err := r.db.NewUpdate().
		Model(spec).
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update OpenAPI spec: %w", err)
	}

	return nil
}

// Delete removes an OpenAPI specification from the database
func (r *OpenAPIRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*models.OpenAPISpec)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete OpenAPI spec: %w", err)
	}

	return nil
}

// Count returns the total number of OpenAPI specifications
func (r *OpenAPIRepository) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.OpenAPISpec)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count OpenAPI specs: %w", err)
	}

	return count, nil
}

// Search searches OpenAPI specifications by title or description
func (r *OpenAPIRepository) Search(ctx context.Context, query string, offset, limit int) ([]*models.OpenAPISpec, error) {
	var specs []*models.OpenAPISpec
	err := r.db.NewSelect().
		Model(&specs).
		Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		OrderExpr("created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to search OpenAPI specs: %w", err)
	}

	return specs, nil
}
