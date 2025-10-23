package main

import (
	"log"
	"net/http"

	"fmp-core/internal/generated"

	"github.com/gin-gonic/gin"
)

// ServerInterface implementation
type Server struct {
	// Add your dependencies here
}

// Implement all the generated interface methods
func (s *Server) GetCategories(ctx *gin.Context) {
	// Implementation here
	ctx.JSON(http.StatusOK, []generated.Category{})
}

func (s *Server) CreateCategory(ctx *gin.Context) {
	var req generated.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, generated.Error{Error: err.Error()})
		return
	}

	// Implementation here
	category := generated.Category{
		Id:        "generated-uuid",
		Name:      req.Name,
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-01T00:00:00Z",
	}

	ctx.JSON(http.StatusCreated, category)
}

// ... implement all other methods from the generated interface

func main() {
	// Create server instance
	server := &Server{}

	// Create Gin router
	router := gin.Default()

	// Register generated routes
	generated.RegisterHandlers(router, server)

	// Start server
	log.Fatal(router.Run(":8080"))
}
