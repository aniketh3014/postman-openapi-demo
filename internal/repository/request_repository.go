package repository

import (
	"context"
	"fmt"
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
	"time"

	"github.com/uptrace/bun"
)

// RequestRepository handles database operations for requests
type RequestRepository struct {
	db *bun.DB
}

// NewRequestRepository creates a new request repository
func NewRequestRepository(db *bun.DB) interfaces.RequestRepository {
	return &RequestRepository{db: db}
}

// Create adds a new request to the database
func (r *RequestRepository) Create(ctx context.Context, request *models.Request) error {
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	_, err := r.db.NewInsert().
		Model(request).
		Returning("id").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return nil
}

// GetByID retrieves a request by its ID
func (r *RequestRepository) GetByID(ctx context.Context, id int64) (*models.Request, error) {
	request := &models.Request{}
	err := r.db.NewSelect().
		Model(request).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get request by ID: %w", err)
	}

	return request, nil
}

// GetByIDWithCollection retrieves a request by its ID with collection data
func (r *RequestRepository) GetByIDWithCollection(ctx context.Context, id int64) (*models.Request, error) {
	request := &models.Request{}
	err := r.db.NewSelect().
		Model(request).
		Where("id = ?", id).
		Relation("Collection").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get request with collection: %w", err)
	}

	return request, nil
}

// List returns all requests with pagination
func (r *RequestRepository) List(ctx context.Context, offset, limit int) ([]*models.Request, error) {
	var requests []*models.Request
	err := r.db.NewSelect().
		Model(&requests).
		OrderExpr("created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list requests: %w", err)
	}

	return requests, nil
}

// ListByCollectionID returns all requests for a specific collection
func (r *RequestRepository) ListByCollectionID(ctx context.Context, collectionID int64, offset, limit int) ([]*models.Request, error) {
	var requests []*models.Request
	err := r.db.NewSelect().
		Model(&requests).
		Where("collection_id = ?", collectionID).
		OrderExpr("created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list requests by collection ID: %w", err)
	}

	return requests, nil
}

// Update modifies an existing request
func (r *RequestRepository) Update(ctx context.Context, request *models.Request) error {
	request.UpdatedAt = time.Now()

	_, err := r.db.NewUpdate().
		Model(request).
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	return nil
}

// Delete removes a request from the database
func (r *RequestRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*models.Request)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}

	return nil
}

// DeleteByCollectionID removes all requests associated with a collection
func (r *RequestRepository) DeleteByCollectionID(ctx context.Context, collectionID int64) error {
	_, err := r.db.NewDelete().
		Model((*models.Request)(nil)).
		Where("collection_id = ?", collectionID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete requests by collection ID: %w", err)
	}

	return nil
}

// Count returns the total number of requests
func (r *RequestRepository) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Request)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count requests: %w", err)
	}

	return count, nil
}

// CountByCollectionID returns the number of requests in a collection
func (r *RequestRepository) CountByCollectionID(ctx context.Context, collectionID int64) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Request)(nil)).
		Where("collection_id = ?", collectionID).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count requests by collection ID: %w", err)
	}

	return count, nil
}
