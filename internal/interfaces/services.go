package interfaces

import (
	"context"
	"postman-api/internal/models"
)

// CollectionService defines operations for managing collections
type CollectionService interface {
	CreateCollection(ctx context.Context, collection *models.Collection) error
	GetCollection(ctx context.Context, id int64) (*models.Collection, error)
	GetCollectionWithRequests(ctx context.Context, id int64) (*models.Collection, error)
	ListCollections(ctx context.Context, page, pageSize int) ([]*models.Collection, int, error)
	UpdateCollection(ctx context.Context, collection *models.Collection) error
	DeleteCollection(ctx context.Context, id int64) error
	ImportPostmanCollection(ctx context.Context, data []byte) (int64, error)
	ExportPostmanCollection(ctx context.Context, id int64) ([]byte, error)
}

// RequestService defines operations for managing API requests
type RequestService interface {
	CreateRequest(ctx context.Context, request *models.Request) error
	GetRequest(ctx context.Context, id int64) (*models.Request, error)
	ListRequests(ctx context.Context, page, pageSize int) ([]*models.Request, int, error)
	ListRequestsByCollection(ctx context.Context, collectionID int64, page, pageSize int) ([]*models.Request, int, error)
	DeleteRequest(ctx context.Context, id int64) error
	UpdateRequestPayload(ctx context.Context, id int64, body models.JSONMap) error
	UpdateRequestHeaders(ctx context.Context, id int64, headers map[string]string) error
	UpdateRequestParams(ctx context.Context, id int64, params models.JSONMap) error
	CloneRequest(ctx context.Context, id int64, newName string) (int64, error)
}

// OpenAPIService defines operations for managing OpenAPI specifications
type OpenAPIService interface {
	CreateOpenAPISpec(ctx context.Context, spec *models.OpenAPISpec) error
	GetOpenAPISpec(ctx context.Context, id int64) (*models.OpenAPISpec, error)
	GetOpenAPISpecByTitle(ctx context.Context, title string) (*models.OpenAPISpec, error)
	ListOpenAPISpecs(ctx context.Context, page, pageSize int) ([]*models.OpenAPISpec, int, error)
	UpdateOpenAPISpec(ctx context.Context, spec *models.OpenAPISpec) error
	DeleteOpenAPISpec(ctx context.Context, id int64) error
	ImportOpenAPISpec(ctx context.Context, data []byte) (int64, error)
	ExportOpenAPISpec(ctx context.Context, id int64) ([]byte, error)
}
