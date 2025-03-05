package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found, using system environment variables.")
	}
}

// Get environment variable safely
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Global Shopify API credentials
var (
	shopifyURL  string
	accessToken string
	userEmail   string
)

// Initialize configuration
func init() {
	loadEnv()
	shopifyURL = getEnv("SHOPIFY_URL", "")
	accessToken = getEnv("SHOPIFY_ACCESS_TOKEN", "")
	userEmail = getEnv("USER_EMAIL", "")
}
