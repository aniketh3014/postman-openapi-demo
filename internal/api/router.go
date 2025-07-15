package api

import (
	"postman-api/internal/api/handlers"
	"postman-api/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine            *gin.Engine
	collectionHandler *handlers.CollectionHandler
	requestHandler    *handlers.RequestHandler
	openAPIHandler    *handlers.OpenAPIHandler
}

func NewRouter(
	collectionService interfaces.CollectionService,
	requestService interfaces.RequestService,
	openAPIService interfaces.OpenAPIService,
) *Router {
	return &Router{
		engine:            gin.Default(),
		collectionHandler: handlers.NewCollectionHandler(collectionService, openAPIService),
		requestHandler:    handlers.NewRequestHandler(requestService),
		openAPIHandler:    handlers.NewOpenAPIHandler(openAPIService),
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.engine.Group("/api/v1")
	{
		// Collection endpoints
		collections := api.Group("/postman")
		{
			collections.GET("", r.collectionHandler.List)
			collections.GET("/:id", r.collectionHandler.Get)
			collections.GET("/:id/with-requests", r.collectionHandler.GetWithRequests)
			collections.PUT("/:id", r.collectionHandler.Update)
			collections.DELETE("/:id", r.collectionHandler.Delete)
			collections.POST("/import", r.collectionHandler.Import)
			collections.GET("/:id/export", r.collectionHandler.Export)
		}

		// Request endpoints
		requests := api.Group("/requests")
		{
			requests.POST("", r.requestHandler.Create)
			requests.GET("", r.requestHandler.List)
			requests.GET("/:id", r.requestHandler.Get)
			requests.DELETE("/:id", r.requestHandler.Delete)
			requests.PUT("/:id/payload", r.requestHandler.UpdatePayload)
			requests.PUT("/:id/headers", r.requestHandler.UpdateHeaders)
			requests.PUT("/:id/params", r.requestHandler.UpdateParams)
			requests.POST("/:id/clone", r.requestHandler.Clone)
		}

		api.GET("/postman/:id/requests", r.requestHandler.ListByCollection)

		// OpenAPI specification endpoints
		openapi := api.Group("/openapi")
		{
			openapi.GET("", r.openAPIHandler.List)
			openapi.GET("/:id", r.openAPIHandler.Get)
			openapi.PUT("/:id", r.openAPIHandler.Update)
			openapi.DELETE("/:id", r.openAPIHandler.Delete)
			openapi.POST("/import", r.openAPIHandler.Import)
			openapi.GET("/:id/export", r.openAPIHandler.Export)
		}
	}

	return r.engine
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
