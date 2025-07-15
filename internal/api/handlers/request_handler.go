package handlers

import (
	"postman-api/internal/interfaces"
	"postman-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RequestHandler handles HTTP requests for API requests
type RequestHandler struct {
	requestService interfaces.RequestService
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(requestService interfaces.RequestService) *RequestHandler {
	return &RequestHandler{
		requestService: requestService,
	}
}

// Create handles creating a new request
func (h *RequestHandler) Create(c *gin.Context) {
	var request models.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		SendBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.requestService.CreateRequest(c.Request.Context(), &request); err != nil {
		SendInternalError(c, "Failed to create request: "+err.Error())
		return
	}

	SendCreated(c, request)
}

// Get retrieves a request by ID
func (h *RequestHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	request, err := h.requestService.GetRequest(c.Request.Context(), id)
	if err != nil {
		SendNotFound(c, "Request not found")
		return
	}

	SendSuccess(c, request)
}

// List returns all requests with pagination
func (h *RequestHandler) List(c *gin.Context) {
	page, pageSize := GetPaginationParams(c)

	requests, total, err := h.requestService.ListRequests(c.Request.Context(), page, pageSize)
	if err != nil {
		SendInternalError(c, "Failed to list requests: "+err.Error())
		return
	}

	SendPaginated(c, requests, page, pageSize, total)
}

// ListByCollection returns all requests for a collection with pagination
func (h *RequestHandler) ListByCollection(c *gin.Context) {
	collectionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid collection ID format")
		return
	}

	page, pageSize := GetPaginationParams(c)

	requests, total, err := h.requestService.ListRequestsByCollection(c.Request.Context(), collectionID, page, pageSize)
	if err != nil {
		SendInternalError(c, "Failed to list requests: "+err.Error())
		return
	}

	SendPaginated(c, requests, page, pageSize, total)
}

// // Update updates an entire request (all fields)
// func (h *RequestHandler) Update(c *gin.Context) {
// 	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if err != nil {
// 		SendBadRequest(c, "Invalid ID format")
// 		return
// 	}

// 	var request models.Request
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		SendBadRequest(c, "Invalid request body: "+err.Error())
// 		return
// 	}

// 	request.ID = id

// 	if err := h.requestService.UpdateRequest(c.Request.Context(), &request); err != nil {
// 		SendInternalError(c, "Failed to update request: "+err.Error())
// 		return
// 	}

// 	SendSuccess(c, request)
// }

// UpdatePayload updates only the payload of a request
func (h *RequestHandler) UpdatePayload(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	var body models.JSONMap
	if err := c.ShouldBindJSON(&body); err != nil {
		SendBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.requestService.UpdateRequestPayload(c.Request.Context(), id, body); err != nil {
		SendInternalError(c, "Failed to update request payload: "+err.Error())
		return
	}

	SendSuccess(c, map[string]string{"message": "Request payload updated successfully"})
}

// UpdateHeaders updates only the headers of a request
func (h *RequestHandler) UpdateHeaders(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	var headers []models.KeyValuePair
	if err := c.ShouldBindJSON(&headers); err != nil {
		SendBadRequest(c, "Invalid headers body: "+err.Error())
		return
	}

	if err := h.requestService.UpdateRequestHeaders(c.Request.Context(), id, headers); err != nil {
		SendInternalError(c, "Failed to update request headers: "+err.Error())
		return
	}

	SendSuccess(c, map[string]string{"message": "Request headers updated successfully"})
}

// UpdateParams updates only the query parameters of a request
func (h *RequestHandler) UpdateParams(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	var params models.JSONMap
	if err := c.ShouldBindJSON(&params); err != nil {
		SendBadRequest(c, "Invalid params body: "+err.Error())
		return
	}

	if err := h.requestService.UpdateRequestParams(c.Request.Context(), id, params); err != nil {
		SendInternalError(c, "Failed to update request params: "+err.Error())
		return
	}

	SendSuccess(c, map[string]string{"message": "Request parameters updated successfully"})
}

// Delete removes a request
func (h *RequestHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	if err := h.requestService.DeleteRequest(c.Request.Context(), id); err != nil {
		SendInternalError(c, "Failed to delete request: "+err.Error())
		return
	}

	SendSuccess(c, map[string]string{"message": "Request deleted successfully"})
}

// Clone creates a copy of an existing request
func (h *RequestHandler) Clone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	var body struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		SendBadRequest(c, "Invalid request body, name is required")
		return
	}

	newID, err := h.requestService.CloneRequest(c.Request.Context(), id, body.Name)
	if err != nil {
		SendInternalError(c, "Failed to clone request: "+err.Error())
		return
	}

	SendCreated(c, map[string]int64{"id": newID})
}
