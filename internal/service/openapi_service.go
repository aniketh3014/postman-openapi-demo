package service

import (
	"context"
	"encoding/json"
	"fmt"
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
	"time"
)

// OpenAPIService handles business logic for OpenAPI specifications
type OpenAPIService struct {
	openAPIRepo interfaces.OpenAPIRepository
}

// NewOpenAPIService creates a new OpenAPI service
func NewOpenAPIService(
	openAPIRepo interfaces.OpenAPIRepository,
) interfaces.OpenAPIService {
	return &OpenAPIService{
		openAPIRepo: openAPIRepo,
	}
}

// CreateOpenAPISpec creates a new OpenAPI specification
func (s *OpenAPIService) CreateOpenAPISpec(ctx context.Context, spec *models.OpenAPISpec) error {
	return s.openAPIRepo.Create(ctx, spec)
}

// GetOpenAPISpec retrieves an OpenAPI specification by ID
func (s *OpenAPIService) GetOpenAPISpec(ctx context.Context, id int64) (*models.OpenAPISpec, error) {
	return s.openAPIRepo.GetByID(ctx, id)
}

// GetOpenAPISpecByTitle retrieves an OpenAPI specification by title
func (s *OpenAPIService) GetOpenAPISpecByTitle(ctx context.Context, title string) (*models.OpenAPISpec, error) {
	return s.openAPIRepo.GetByTitle(ctx, title)
}

// ListOpenAPISpecs returns all OpenAPI specifications with pagination
func (s *OpenAPIService) ListOpenAPISpecs(ctx context.Context, page, pageSize int) ([]*models.OpenAPISpec, int, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	specs, err := s.openAPIRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.openAPIRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return specs, total, nil
}

// UpdateOpenAPISpec updates an existing OpenAPI specification
func (s *OpenAPIService) UpdateOpenAPISpec(ctx context.Context, spec *models.OpenAPISpec) error {
	existingSpec, err := s.openAPIRepo.GetByID(ctx, spec.ID)
	if err != nil {
		return fmt.Errorf("OpenAPI specification not found: %w", err)
	}

	spec.CreatedAt = existingSpec.CreatedAt
	spec.UpdatedAt = time.Now()

	return s.openAPIRepo.Update(ctx, spec)
}

// DeleteOpenAPISpec removes an OpenAPI specification
func (s *OpenAPIService) DeleteOpenAPISpec(ctx context.Context, id int64) error {
	return s.openAPIRepo.Delete(ctx, id)
}

// ImportOpenAPISpec imports an OpenAPI specification from JSON
func (s *OpenAPIService) ImportOpenAPISpec(ctx context.Context, data []byte) (int64, error) {
	var content models.JSONMap
	if err := json.Unmarshal(data, &content); err != nil {
		return 0, fmt.Errorf("invalid OpenAPI format: %w", err)
	}

	info, ok := content["info"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("invalid OpenAPI format: missing or invalid 'info' object")
	}

	title, ok := info["title"].(string)
	if !ok || title == "" {
		return 0, fmt.Errorf("invalid OpenAPI format: missing or invalid 'title'")
	}

	version, ok := info["version"].(string)
	if !ok || version == "" {
		return 0, fmt.Errorf("invalid OpenAPI format: missing or invalid 'version'")
	}

	description := ""
	if desc, ok := info["description"].(string); ok {
		description = desc
	}

	spec := &models.OpenAPISpec{
		Title:       title,
		Description: description,
		Version:     version,
		Content:     content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.openAPIRepo.Create(ctx, spec); err != nil {
		return 0, fmt.Errorf("failed to create OpenAPI spec: %w", err)
	}

	return spec.ID, nil
}

// ExportOpenAPISpec exports an OpenAPI specification to JSON
func (s *OpenAPIService) ExportOpenAPISpec(ctx context.Context, id int64) ([]byte, error) {
	spec, err := s.GetOpenAPISpec(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI spec: %w", err)
	}

	if spec.Content == nil {
		return nil, fmt.Errorf("OpenAPI spec has no content")
	}

	return json.MarshalIndent(spec.Content, "", "  ")
}
