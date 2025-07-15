package interfaces

import (
	"context"
	"postman-api/internal/models"
)

// CollectionRepository defines operations for collection persistence
type CollectionRepository interface {
	Create(ctx context.Context, collection *models.Collection) error
	GetByID(ctx context.Context, id int64) (*models.Collection, error)
	GetWithRequests(ctx context.Context, id int64) (*models.Collection, error)
	List(ctx context.Context, offset, limit int) ([]*models.Collection, error)
	Update(ctx context.Context, collection *models.Collection) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}

// RequestRepository defines operations for request persistence
type RequestRepository interface {
	Create(ctx context.Context, request *models.Request) error
	GetByID(ctx context.Context, id int64) (*models.Request, error)
	List(ctx context.Context, offset, limit int) ([]*models.Request, error)
	ListByCollectionID(ctx context.Context, collectionID int64, offset, limit int) ([]*models.Request, error)
	Update(ctx context.Context, request *models.Request) error
	Delete(ctx context.Context, id int64) error
	DeleteByCollectionID(ctx context.Context, collectionID int64) error
	Count(ctx context.Context) (int, error)
	CountByCollectionID(ctx context.Context, collectionID int64) (int, error)
}

// OpenAPIRepository defines operations for OpenAPI spec persistence
type OpenAPIRepository interface {
	Create(ctx context.Context, spec *models.OpenAPISpec) error
	GetByID(ctx context.Context, id int64) (*models.OpenAPISpec, error)
	GetByTitle(ctx context.Context, title string) (*models.OpenAPISpec, error)
	List(ctx context.Context, offset, limit int) ([]*models.OpenAPISpec, error)
	Update(ctx context.Context, spec *models.OpenAPISpec) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}
