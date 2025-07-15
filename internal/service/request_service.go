package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
)

// RequestService handles business logic for API requests
type RequestService struct {
	requestRepo    interfaces.RequestRepository
	collectionRepo interfaces.CollectionRepository
}

// NewRequestService creates a new request service
func NewRequestService(
	requestRepo interfaces.RequestRepository,
	collectionRepo interfaces.CollectionRepository,
) interfaces.RequestService {
	return &RequestService{
		requestRepo:    requestRepo,
		collectionRepo: collectionRepo,
	}
}

// CreateRequest creates a new API request
func (s *RequestService) CreateRequest(ctx context.Context, request *models.Request) error {
	_, err := s.collectionRepo.GetByID(ctx, request.CollectionID)
	if err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Validate URL is valid JSON
	if request.URL != nil {
		if urlStr, ok := request.URL["raw"].(string); ok && urlStr != "" {
			request.URL = models.JSONMap{
				"raw": urlStr,
			}
		}

		if _, err := json.Marshal(request.URL); err != nil {
			request.URL = models.JSONMap{}
		}
	} else {
		request.URL = models.JSONMap{}
	}

	return s.requestRepo.Create(ctx, request)
}

// GetRequest retrieves a request by ID
func (s *RequestService) GetRequest(ctx context.Context, id int64) (*models.Request, error) {
	return s.requestRepo.GetByID(ctx, id)
}

// ListRequests returns all requests with pagination
func (s *RequestService) ListRequests(ctx context.Context, page, pageSize int) ([]*models.Request, int, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	requests, err := s.requestRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.requestRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// ListRequestsByCollection returns all requests in a collection with pagination
func (s *RequestService) ListRequestsByCollection(ctx context.Context, collectionID int64, page, pageSize int) ([]*models.Request, int, error) {
	_, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return nil, 0, fmt.Errorf("collection not found: %w", err)
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	requests, err := s.requestRepo.ListByCollectionID(ctx, collectionID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.requestRepo.CountByCollectionID(ctx, collectionID)
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// // UpdateRequest updates an existing request
// func (s *RequestService) UpdateRequest(ctx context.Context, request *models.Request) error {
// 	existingRequest, err := s.requestRepo.GetByID(ctx, request.ID)
// 	if err != nil {
// 		return fmt.Errorf("request not found: %w", err)
// 	}

// 	if existingRequest.CollectionID != request.CollectionID {
// 		_, err := s.collectionRepo.GetByID(ctx, request.CollectionID)
// 		if err != nil {
// 			return fmt.Errorf("target collection not found: %w", err)
// 		}
// 	}

// 	// Validate URL is valid JSON
// 	if request.URL != nil {
// 		if urlStr, ok := request.URL["raw"].(string); ok && urlStr != "" {
// 			request.URL = models.JSONMap{
// 				"raw": urlStr,
// 			}
// 		}

// 		if _, err := json.Marshal(request.URL); err != nil {
// 			request.URL = models.JSONMap{}
// 		}
// 	} else {
// 		request.URL = models.JSONMap{}
// 	}

// 	return s.requestRepo.Update(ctx, request)
// }

// DeleteRequest removes a request
func (s *RequestService) DeleteRequest(ctx context.Context, id int64) error {
	_, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	return s.requestRepo.Delete(ctx, id)
}

// UpdateRequestPayload updates only the payload (body) of a request
func (s *RequestService) UpdateRequestPayload(ctx context.Context, id int64, body models.JSONMap) error {
	if body == nil {
		return errors.New("body cannot be nil")
	}

	request, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	request.Body = body
	return s.requestRepo.Update(ctx, request)
}

// UpdateRequestHeaders updates only the headers of a request
func (s *RequestService) UpdateRequestHeaders(ctx context.Context, id int64, headers map[string]string) error {
	if headers == nil {
		return errors.New("headers cannot be nil")
	}

	request, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	request.Headers = headers
	return s.requestRepo.Update(ctx, request)
}

// UpdateRequestParams updates only the query parameters of a request
func (s *RequestService) UpdateRequestParams(ctx context.Context, id int64, params models.JSONMap) error {
	if params == nil {
		return errors.New("params cannot be nil")
	}

	request, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	request.Params = params
	return s.requestRepo.Update(ctx, request)
}

// CloneRequest creates a copy of an existing request
func (s *RequestService) CloneRequest(ctx context.Context, id int64, newName string) (int64, error) {
	original, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("request not found: %w", err)
	}

	urlData := models.JSONMap{}
	if original.URL != nil {
		if _, err := json.Marshal(original.URL); err == nil {
			urlData = original.URL
		}
	}

	cloned := &models.Request{
		CollectionID: original.CollectionID,
		Name:         newName,
		Description:  original.Description + " (Cloned)",
		URL:          urlData,
		Method:       original.Method,
		Headers:      original.Headers,
		Params:       original.Params,
		Body:         original.Body,
	}

	if err := s.requestRepo.Create(ctx, cloned); err != nil {
		return 0, fmt.Errorf("failed to clone request: %w", err)
	}

	return cloned.ID, nil
}
