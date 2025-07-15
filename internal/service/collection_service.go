package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
)

// CollectionService handles business logic for collections
type CollectionService struct {
	collectionRepo interfaces.CollectionRepository
	requestRepo    interfaces.RequestRepository
}

// NewCollectionService creates a new collection service
func NewCollectionService(
	collectionRepo interfaces.CollectionRepository,
	requestRepo interfaces.RequestRepository,
) interfaces.CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
		requestRepo:    requestRepo,
	}
}

// CreateCollection creates a new collection
func (s *CollectionService) CreateCollection(ctx context.Context, collection *models.Collection) error {
	return s.collectionRepo.Create(ctx, collection)
}

// GetCollection retrieves a collection by ID
func (s *CollectionService) GetCollection(ctx context.Context, id int64) (*models.Collection, error) {
	return s.collectionRepo.GetByID(ctx, id)
}

// GetCollectionWithRequests retrieves a collection with all its requests
func (s *CollectionService) GetCollectionWithRequests(ctx context.Context, id int64) (*models.Collection, error) {
	return s.collectionRepo.GetWithRequests(ctx, id)
}

// ListCollections returns all collections with pagination
func (s *CollectionService) ListCollections(ctx context.Context, page, pageSize int) ([]*models.Collection, int, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	collections, err := s.collectionRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.collectionRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return collections, total, nil
}

// UpdateCollection updates an existing collection
func (s *CollectionService) UpdateCollection(ctx context.Context, collection *models.Collection) error {
	existingCollection, err := s.collectionRepo.GetByID(ctx, collection.ID)
	if err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	collection.Items = existingCollection.Items

	return s.collectionRepo.Update(ctx, collection)
}

// DeleteCollection removes a collection and all its requests
func (s *CollectionService) DeleteCollection(ctx context.Context, id int64) error {
	err := s.requestRepo.DeleteByCollectionID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete requests in collection: %w", err)
	}

	return s.collectionRepo.Delete(ctx, id)
}

// ImportPostmanCollection imports a Postman collection from JSON
func (s *CollectionService) ImportPostmanCollection(ctx context.Context, data []byte) (int64, error) {
	var postmanCollection models.PostmanCollection
	if err := json.Unmarshal(data, &postmanCollection); err != nil {
		return 0, fmt.Errorf("invalid Postman collection format: %w", err)
	}

	if postmanCollection.Info.Name == "" {
		return 0, errors.New("collection name is required")
	}

	variables := make(models.JSONMap)
	for _, v := range postmanCollection.Variable {
		variables[v.Key] = v.Value
	}

	var auth models.JSONMap
	if postmanCollection.Auth != nil {
		if err := json.Unmarshal(postmanCollection.Auth, &auth); err == nil {
			// Successful unmarshaling
		}
	}

	var events models.JSONMap
	if len(postmanCollection.Event) > 0 {
		eventsBytes, err := json.Marshal(postmanCollection.Event)
		if err == nil {
			if err := json.Unmarshal(eventsBytes, &events); err == nil {
				// Successful unmarshaling
			}
		}
	}

	var items models.JSONMap
	itemsBytes, err := json.Marshal(postmanCollection.Item)
	if err == nil {
		if err := json.Unmarshal(itemsBytes, &items); err == nil {
			// Successful unmarshaling
		}
	}

	collection := &models.Collection{
		Name:        postmanCollection.Info.Name,
		Description: postmanCollection.Info.Description,
		Schema:      postmanCollection.Schema,
		Variables:   variables,
		Auth:        auth,
		Events:      events,
		Items:       items,
		PostmanID:   postmanCollection.Info.PostmanID,
		ExporterID:  postmanCollection.Info.ExporterID,
	}

	if err := s.collectionRepo.Create(ctx, collection); err != nil {
		return 0, fmt.Errorf("failed to create collection: %w", err)
	}

	if err := s.processPostmanItems(ctx, postmanCollection.Item, collection.ID, ""); err != nil {
		return 0, err
	}

	return collection.ID, nil
}

// processPostmanItems processes items in a Postman collection, handling nested folders
func (s *CollectionService) processPostmanItems(ctx context.Context, items []models.PostmanItem, collectionID int64, parentPath string) error {
	for _, item := range items {
		currentPath := parentPath
		if currentPath != "" {
			currentPath += "/"
		}
		currentPath += item.Name

		if len(item.Item) > 0 {
			if err := s.processPostmanItems(ctx, item.Item, collectionID, currentPath); err != nil {
				return err
			}
			continue
		}

		if item.Request == nil {
			continue
		}

		request := &models.Request{
			CollectionID: collectionID,
			Name:         item.Name,
			Description:  item.Description,
			FolderPath:   parentPath,
			Method:       item.Request.Method,
			PostmanID:    item.PostmanID,
		}

		var urlMap models.JSONMap

		switch v := item.Request.URL.(type) {
		case string:
			if v != "" {
				urlMap = models.JSONMap{
					"raw": v,
				}
			} else {
				urlMap = models.JSONMap{}
			}
		default:
			urlBytes, err := json.Marshal(item.Request.URL)
			if err == nil {
				if err := json.Unmarshal(urlBytes, &urlMap); err != nil {
					urlMap = models.JSONMap{}
				}
			} else {
				urlMap = models.JSONMap{}
			}
		}

		if urlMap == nil {
			urlMap = models.JSONMap{}
		}

		request.URL = urlMap

		if len(item.Request.Header) > 0 {
			headers := make(map[string]string)
			for _, kv := range item.Request.Header {
				headers[kv.Key] = kv.Value
			}
			request.Headers = headers
		}

		bodyBytes, err := json.Marshal(item.Request.Body)
		if err == nil {
			var bodyMap models.JSONMap
			if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
				request.Body = bodyMap
			}
		}

		if item.Request.Auth != nil {
			var authMap models.JSONMap
			authBytes, err := json.Marshal(item.Request.Auth)
			if err == nil {
				if err := json.Unmarshal(authBytes, &authMap); err == nil {
					request.Auth = authMap
				}
			}
		}

		if len(item.Event) > 0 {
			eventsBytes, err := json.Marshal(item.Event)
			if err == nil {
				var eventsMap models.JSONMap
				if err := json.Unmarshal(eventsBytes, &eventsMap); err == nil {
					request.Events = eventsMap
				}
			}
		}

		if len(item.Response) > 0 {
			responsesBytes, err := json.Marshal(item.Response)
			if err == nil {
				var responsesMap models.JSONMap
				if err := json.Unmarshal(responsesBytes, &responsesMap); err == nil {
					request.Responses = responsesMap
				}
			}
		}

		if err := s.requestRepo.Create(ctx, request); err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
	}

	return nil
}

// ExportPostmanCollection exports a collection to Postman format
func (s *CollectionService) ExportPostmanCollection(ctx context.Context, id int64) ([]byte, error) {
	collection, err := s.GetCollection(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	postmanCollection := models.PostmanCollection{
		Info: models.CollectionInfo{
			Name:        collection.Name,
			Description: collection.Description,
			Schema:      collection.Schema,
			PostmanID:   collection.PostmanID,
			ExporterID:  collection.ExporterID,
		},
		Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
	}

	if len(collection.Items) > 0 {
		itemsBytes, err := json.Marshal(collection.Items)
		if err == nil {
			var items []models.PostmanItem
			if err := json.Unmarshal(itemsBytes, &items); err == nil {
				postmanCollection.Item = items

				if collection.Variables != nil {
					for k, v := range collection.Variables {
						postmanCollection.Variable = append(postmanCollection.Variable, models.KeyValuePair{
							Key:   k,
							Value: fmt.Sprintf("%v", v),
						})
					}
				}

				if collection.Auth != nil {
					authBytes, _ := json.Marshal(collection.Auth)
					postmanCollection.Auth = authBytes
				}

				if collection.Events != nil {
					eventsBytes, _ := json.Marshal(collection.Events)
					json.Unmarshal(eventsBytes, &postmanCollection.Event)
				}

				return json.MarshalIndent(postmanCollection, "", "  ")
			}
		}
	}

	requests, err := s.requestRepo.ListByCollectionID(ctx, id, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}

	folderMap := make(map[string][]models.PostmanItem)
	for _, req := range requests {
		postmanReq := &models.PostmanRequest{
			Method:      req.Method,
			Description: req.Description,
		}

		if req.URL != nil {
			if urlBytes, err := json.Marshal(req.URL); err == nil {
				json.Unmarshal(urlBytes, &postmanReq.URL)
			} else {
				postmanReq.URL = ""
			}
		}

		if req.Headers != nil {
			var headerArr []models.KeyValuePair
			for k, v := range req.Headers {
				headerArr = append(headerArr, models.KeyValuePair{Key: k, Value: v})
			}
			postmanReq.Header = headerArr
		}

		if req.Body != nil {
			bodyBytes, _ := json.Marshal(req.Body)
			json.Unmarshal(bodyBytes, &postmanReq.Body)
		}

		if req.Auth != nil {
			authBytes, _ := json.Marshal(req.Auth)
			postmanReq.Auth = authBytes
		}

		item := models.PostmanItem{
			Name:        req.Name,
			Description: req.Description,
			PostmanID:   req.PostmanID,
			Request:     postmanReq,
		}

		if req.Events != nil {
			eventsBytes, _ := json.Marshal(req.Events)
			json.Unmarshal(eventsBytes, &item.Event)
		}

		if req.Responses != nil {
			responsesBytes, _ := json.Marshal(req.Responses)
			json.Unmarshal(responsesBytes, &item.Response)
		}

		folderPath := req.FolderPath
		folderMap[folderPath] = append(folderMap[folderPath], item)
	}

	postmanCollection.Item = folderMap[""]

	for path, items := range folderMap {
		if path == "" {
			continue
		}

		folder := models.PostmanItem{
			Name: path,
			Item: items,
		}

		postmanCollection.Item = append(postmanCollection.Item, folder)
	}

	if collection.Variables != nil {
		for k, v := range collection.Variables {
			postmanCollection.Variable = append(postmanCollection.Variable, models.KeyValuePair{
				Key:   k,
				Value: fmt.Sprintf("%v", v),
			})
		}
	}

	if collection.Auth != nil {
		authBytes, _ := json.Marshal(collection.Auth)
		postmanCollection.Auth = authBytes
	}

	if collection.Events != nil {
		eventsBytes, _ := json.Marshal(collection.Events)
		json.Unmarshal(eventsBytes, &postmanCollection.Event)
	}

	return json.MarshalIndent(postmanCollection, "", "  ")
}
