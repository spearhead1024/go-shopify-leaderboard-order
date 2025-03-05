package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/machinebox/graphql"
)

// Structs for API response
type AmountSpent struct {
	Amount string `json:"amount"`
}

type Customer struct {
	ID          string      `json:"id"`
	DisplayName string      `json:"displayName"`
	Email       string      `json:"email"`
	AmountSpent AmountSpent `json:"amountSpent"`
}

type CustomersResponse struct {
	Customers struct {
		Edges []struct {
			Node Customer `json:"node"`
		} `json:"edges"`
	} `json:"customers"`
}

// Fetch top customers
func fetchTopCustomers(client *graphql.Client) ([]Customer, error) {
	query := `
	query {
		customers(first: 100, query: "tag:task1 AND tag:level:3") {
			edges {
				node {
					id
					displayName
					email
					amountSpent {
						amount
					}
				}
			}
		}
	}`
	req := graphql.NewRequest(query)
	req.Header.Set("X-Shopify-Access-Token", accessToken)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var response CustomersResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	customers := make([]Customer, 0)
	for _, edge := range response.Customers.Edges {
		customers = append(customers, edge.Node)
	}

	sort.Slice(customers, func(i, j int) bool {
		amountI, _ := strconv.ParseFloat(customers[i].AmountSpent.Amount, 64)
		amountJ, _ := strconv.ParseFloat(customers[j].AmountSpent.Amount, 64)
		return amountI > amountJ
	})

	if len(customers) > 50 {
		customers = customers[:50]
	}

	return customers, nil
}

// Export customers to CSV
func exportToCSV(customers []Customer) error {
	file, err := os.Create("leaderboard.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"ID", "Name", "Email", "Amount Spent"})

	for _, customer := range customers {
		writer.Write([]string{customer.ID, customer.DisplayName, customer.Email, customer.AmountSpent.Amount})
	}

	fmt.Println("CSV Exported Successfully: leaderboard.csv")
	return nil
}

// Run Task 1
func RunLeaderboard() {
	client := graphql.NewClient(shopifyURL)

	customers, err := fetchTopCustomers(client)
	if err != nil {
		fmt.Println("Error fetching customers:", err)
		return
	}

	if err := exportToCSV(customers); err != nil {
		fmt.Println("Error writing CSV:", err)
	}
}
