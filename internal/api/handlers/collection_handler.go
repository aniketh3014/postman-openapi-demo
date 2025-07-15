package handlers

import (
	"fmt"
	"io"
	"net/http"
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CollectionHandler handles HTTP requests for collections
type CollectionHandler struct {
	collectionService interfaces.CollectionService
	openAPIService    interfaces.OpenAPIService
}

// NewCollectionHandler creates a new collection handler
func NewCollectionHandler(collectionService interfaces.CollectionService, openAPIService interfaces.OpenAPIService) *CollectionHandler {
	return &CollectionHandler{
		collectionService: collectionService,
		openAPIService:    openAPIService,
	}
}

// Create handles creating a new collection
func (h *CollectionHandler) Create(c *gin.Context) {
	var collection models.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		SendBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.collectionService.CreateCollection(c.Request.Context(), &collection); err != nil {
		SendInternalError(c, "Failed to create collection: "+err.Error())
		return
	}

	SendCreated(c, collection)
}

// Get retrieves a collection by ID
func (h *CollectionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	collection, err := h.collectionService.GetCollection(c.Request.Context(), id)
	if err != nil {
		SendNotFound(c, "Collection not found")
		return
	}

	SendSuccess(c, collection)
}

// GetWithRequests retrieves a collection with all its requests
func (h *CollectionHandler) GetWithRequests(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	collection, err := h.collectionService.GetCollectionWithRequests(c.Request.Context(), id)
	if err != nil {
		SendNotFound(c, "Collection not found")
		return
	}

	SendSuccess(c, collection)
}

// List returns all collections with pagination
func (h *CollectionHandler) List(c *gin.Context) {
	page, pageSize := GetPaginationParams(c)

	collections, total, err := h.collectionService.ListCollections(c.Request.Context(), page, pageSize)
	if err != nil {
		SendInternalError(c, "Failed to list collections: "+err.Error())
		return
	}

	SendPaginated(c, collections, page, pageSize, total)
}

// Update updates an existing collection
func (h *CollectionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	var collection models.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		SendBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	collection.ID = id

	if err := h.collectionService.UpdateCollection(c.Request.Context(), &collection); err != nil {
		SendInternalError(c, "Failed to update collection: "+err.Error())
		return
	}

	SendSuccess(c, collection)
}

// Delete removes a collection and all its requests
func (h *CollectionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	if err := h.collectionService.DeleteCollection(c.Request.Context(), id); err != nil {
		SendInternalError(c, "Failed to delete collection: "+err.Error())
		return
	}

	SendSuccess(c, map[string]string{"message": "Collection deleted successfully"})
}

// Import imports a Postman collection from JSON
func (h *CollectionHandler) Import(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		SendBadRequest(c, "Invalid file: "+err.Error())
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		SendInternalError(c, "Failed to read file: "+err.Error())
		return
	}

	collectionID, err := h.collectionService.ImportPostmanCollection(c.Request.Context(), data)
	if err != nil {
		SendBadRequest(c, "Failed to import collection: "+err.Error())
		return
	}

	SendCreated(c, map[string]int64{"id": collectionID})
}

// Export exports a collection to Postman format
func (h *CollectionHandler) Export(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	collection, err := h.collectionService.GetCollection(c.Request.Context(), id)
	if err != nil {
		SendNotFound(c, "Collection not found")
		return
	}

	data, err := h.collectionService.ExportPostmanCollection(c.Request.Context(), id)
	if err != nil {
		SendInternalError(c, "Failed to export collection: "+err.Error())
		return
	}

	filename := fmt.Sprintf("%s.postman_collection.json", collection.Name)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/json", data)
}
