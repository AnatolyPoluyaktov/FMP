package main

import (
	"context"
	"fmt"
	"log"

	"minapp-backend/internal/generated"
)

func main() {
	// Create client
	client, err := generated.NewClientWithResponses("http://localhost:8080/api/v1")
	if err != nil {
		log.Fatal(err)
	}

	// Example: Get categories
	categoriesResp, err := client.GetCategoriesWithResponse(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if categoriesResp.StatusCode() == 200 {
		categories := *categoriesResp.JSON200
		fmt.Printf("Found %d categories\n", len(categories))

		for _, category := range categories {
			fmt.Printf("- %s: %s\n", category.Name, category.Description)
		}
	}

	// Example: Create a new category
	createReq := generated.CreateCategoryRequest{
		Name:        "Food",
		Description: "Food and dining expenses",
	}

	createResp, err := client.CreateCategoryWithResponse(context.Background(), createReq)
	if err != nil {
		log.Fatal(err)
	}

	if createResp.StatusCode() == 201 {
		category := *createResp.JSON201
		fmt.Printf("Created category: %s\n", category.Name)
	}
}
