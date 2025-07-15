package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response is a common API response structure
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Meta contains metadata for paginated responses
type Meta struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	TotalRows int `json:"totalRows"`
	TotalPage int `json:"totalPage"`
}

// SuccessResponse creates a success response with data
func SuccessResponse(data any) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse creates an error response with message
func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(data any, page, pageSize, total int) Response {
	totalPage := total / pageSize
	if total%pageSize > 0 {
		totalPage++
	}

	return Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:      page,
			PageSize:  pageSize,
			TotalRows: total,
			TotalPage: totalPage,
		},
	}
}

// GetPaginationParams extracts pagination parameters from the request
func GetPaginationParams(c *gin.Context) (page int, pageSize int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}

// SendJSON is a helper function to send JSON responses
func SendJSON(c *gin.Context, statusCode int, response Response) {
	c.JSON(statusCode, response)
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, data any) {
	SendJSON(c, http.StatusOK, SuccessResponse(data))
}

// SendCreated sends a created response
func SendCreated(c *gin.Context, data any) {
	SendJSON(c, http.StatusCreated, SuccessResponse(data))
}

// SendError sends an error response
func SendError(c *gin.Context, statusCode int, message string) {
	SendJSON(c, statusCode, ErrorResponse(message))
}

// SendBadRequest sends a bad request error
func SendBadRequest(c *gin.Context, message string) {
	SendError(c, http.StatusBadRequest, message)
}

// SendNotFound sends a not found error
func SendNotFound(c *gin.Context, message string) {
	SendError(c, http.StatusNotFound, message)
}

// SendInternalError sends an internal server error
func SendInternalError(c *gin.Context, message string) {
	SendError(c, http.StatusInternalServerError, message)
}

// SendPaginated sends a paginated response
func SendPaginated(c *gin.Context, data any, page, pageSize, total int) {
	SendJSON(c, http.StatusOK, PaginatedResponse(data, page, pageSize, total))
}
