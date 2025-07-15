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

// OpenAPIHandler handles HTTP requests for OpenAPI specifications
type OpenAPIHandler struct {
	openAPIService interfaces.OpenAPIService
}

// NewOpenAPIHandler creates a new OpenAPI handler
func NewOpenAPIHandler(openAPIService interfaces.OpenAPIService) *OpenAPIHandler {
	return &OpenAPIHandler{
		openAPIService: openAPIService,
	}
}

// Create handles creating a new OpenAPI specification
func (h *OpenAPIHandler) Create(c *gin.Context) {
	var spec models.OpenAPISpec
	if err := c.ShouldBindJSON(&spec); err != nil {
		SendBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.openAPIService.CreateOpenAPISpec(c.Request.Context(), &spec); err != nil {
		SendInternalError(c, "Failed to create OpenAPI specification: "+err.Error())
		return
	}

	SendCreated(c, spec)
}

// Get retrieves an OpenAPI specification by ID
func (h *OpenAPIHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	spec, err := h.openAPIService.GetOpenAPISpec(c.Request.Context(), id)
	if err != nil {
		SendNotFound(c, "OpenAPI specification not found")
		return
	}

	SendSuccess(c, spec)
}

// List returns all OpenAPI specifications with pagination
func (h *OpenAPIHandler) List(c *gin.Context) {
	page, pageSize := GetPaginationParams(c)

	specs, total, err := h.openAPIService.ListOpenAPISpecs(c.Request.Context(), page, pageSize)
	if err != nil {
		SendInternalError(c, "Failed to list OpenAPI specifications: "+err.Error())
		return
	}

	SendPaginated(c, specs, page, pageSize, total)
}

// Update updates an existing OpenAPI specification
func (h *OpenAPIHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	var spec models.OpenAPISpec
	if err := c.ShouldBindJSON(&spec); err != nil {
		SendBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	spec.ID = id

	if err := h.openAPIService.UpdateOpenAPISpec(c.Request.Context(), &spec); err != nil {
		SendInternalError(c, "Failed to update OpenAPI specification: "+err.Error())
		return
	}

	SendSuccess(c, spec)
}

// Delete removes an OpenAPI specification
func (h *OpenAPIHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	if err := h.openAPIService.DeleteOpenAPISpec(c.Request.Context(), id); err != nil {
		SendInternalError(c, "Failed to delete OpenAPI specification: "+err.Error())
		return
	}

	SendSuccess(c, map[string]string{"message": "OpenAPI specification deleted successfully"})
}

// Import imports an OpenAPI specification from JSON
func (h *OpenAPIHandler) Import(c *gin.Context) {
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

	specID, err := h.openAPIService.ImportOpenAPISpec(c.Request.Context(), data)
	if err != nil {
		SendBadRequest(c, "Failed to import OpenAPI specification: "+err.Error())
		return
	}

	SendCreated(c, map[string]int64{"id": specID})
}

// Export exports an OpenAPI specification to JSON
func (h *OpenAPIHandler) Export(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendBadRequest(c, "Invalid ID format")
		return
	}

	spec, err := h.openAPIService.GetOpenAPISpec(c.Request.Context(), id)
	if err != nil {
		SendNotFound(c, "OpenAPI specification not found")
		return
	}

	data, err := h.openAPIService.ExportOpenAPISpec(c.Request.Context(), id)
	if err != nil {
		SendInternalError(c, "Failed to export OpenAPI specification: "+err.Error())
		return
	}

	filename := fmt.Sprintf("%s.openapi.json", spec.Title)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/json", data)
}
