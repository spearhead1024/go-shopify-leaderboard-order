package main

import (
	"context"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
)

// Fetch the "flag" product variant ID
func fetchFlagProduct(client *graphql.Client) (string, error) {
	query := `
	query {
		products(first: 5, query: "title:*flag*") {
			edges {
				node {
					id
					variants(first: 1) {
						edges {
							node {
								id
							}
						}
					}
				}
			}
		}
	}`

	req := graphql.NewRequest(query)
	req.Header.Set("X-Shopify-Access-Token", accessToken)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var response struct {
		Products struct {
			Edges []struct {
				Node struct {
					ID       string `json:"id"`
					Variants struct {
						Edges []struct {
							Node struct {
								ID string `json:"id"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"variants"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"products"`
	}

	if err := client.Run(ctx, req, &response); err != nil {
		return "", fmt.Errorf("GraphQL request error: %v", err)
	}

	// Check if products were found
	if len(response.Products.Edges) == 0 {
		return "", fmt.Errorf("flag product not found - no matching products in Shopify")
	}
	product := response.Products.Edges[0].Node
	fmt.Println("Found Product ID:", product.ID)

	// If product has a variant, return its ID
	if len(product.Variants.Edges) > 0 {
		fmt.Println("Found Variant ID:", product.Variants.Edges[0].Node.ID)
		return product.Variants.Edges[0].Node.ID, nil
	}

	// If no variant exists, return the product ID itself
	fmt.Println("No variant found, using Product ID instead.")
	return product.ID, nil
}

// Create an order for the flag product
func createOrder(client *graphql.Client, productID string) error {
	mutation := `
	mutation {
		orderCreate(input: {
			email: "` + userEmail + `",
			lineItems: [{
				variantId: "` + productID + `",
				quantity: 1
			}]
		}) {
			order {
				id
			}
			userErrors {
				field
				message
			}
		}
	}`

	req := graphql.NewRequest(mutation)
	req.Header.Set("X-Shopify-Access-Token", accessToken)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var response struct {
		OrderCreate struct {
			Order struct {
				ID string `json:"id"`
			} `json:"order"`
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"orderCreate"`
	}

	if err := client.Run(ctx, req, &response); err != nil {
		return fmt.Errorf("error creating order: %v", err)
	}

	// Check for GraphQL user errors
	if len(response.OrderCreate.UserErrors) > 0 {
		return fmt.Errorf("GraphQL error: %v", response.OrderCreate.UserErrors[0].Message)
	}

	fmt.Println("Order placed successfully! Order ID:", response.OrderCreate.Order.ID)
	return nil
}

// Run Task 2
func RunOrder() {
	client := graphql.NewClient(shopifyURL)

	productID, err := fetchFlagProduct(client)
	if err != nil {
		fmt.Println("Error fetching flag product:", err)
		return
	}

	fmt.Println("Using Variant/Product ID:", productID)

	if err := createOrder(client, productID); err != nil {
		fmt.Println("Error creating order:", err)
	}
}
