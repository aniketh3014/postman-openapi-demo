package repository

import (
	"context"
	"fmt"
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
	"time"

	"github.com/uptrace/bun"
)

// CollectionRepository handles database operations for collections
type CollectionRepository struct {
	db *bun.DB
}

func NewCollectionRepository(db *bun.DB) interfaces.CollectionRepository {
	return &CollectionRepository{db: db}
}

// Create adds a new collection to the database
func (r *CollectionRepository) Create(ctx context.Context, collection *models.Collection) error {
	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()

	_, err := r.db.NewInsert().
		Model(collection).
		Returning("id").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

// GetByID retrieves a collection by its ID
func (r *CollectionRepository) GetByID(ctx context.Context, id int64) (*models.Collection, error) {
	collection := &models.Collection{}
	err := r.db.NewSelect().
		Model(collection).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get collection by ID: %w", err)
	}

	return collection, nil
}

// List returns all collections with pagination
func (r *CollectionRepository) List(ctx context.Context, offset, limit int) ([]*models.Collection, error) {
	var collections []*models.Collection
	err := r.db.NewSelect().
		Model(&collections).
		OrderExpr("created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	return collections, nil
}

// Update modifies an existing collection
func (r *CollectionRepository) Update(ctx context.Context, collection *models.Collection) error {
	collection.UpdatedAt = time.Now()

	_, err := r.db.NewUpdate().
		Model(collection).
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update collection: %w", err)
	}

	return nil
}

// Delete removes a collection from the database
func (r *CollectionRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().
		Model((*models.Collection)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// GetWithRequests retrieves a collection with all its requests
func (r *CollectionRepository) GetWithRequests(ctx context.Context, id int64) (*models.Collection, error) {
	collection := &models.Collection{}
	err := r.db.NewSelect().
		Model(collection).
		Where("id = ?", id).
		Relation("Requests").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get collection with requests: %w", err)
	}

	return collection, nil
}

// Count returns the total number of collections
func (r *CollectionRepository) Count(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Collection)(nil)).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count collections: %w", err)
	}

	return count, nil
}
