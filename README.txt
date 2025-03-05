# Shopify Backend Tasks

This project contains two tasks:
1. Task 1: Leaderboard - Fetches top 50 highest-spending customers with specific tags.
2. Task 2: Capture the Flag - Places an order for a flag product.

## Setup Instructions
1. Install Go from https://go.dev/dl/
2. Clone this project: `git clone <repo-url>`
3. Install dependencies: `go mod tidy`
4. Create a `.env` file and add:
   SHOPIFY_URL=<your-shopify-url>
   SHOPIFY_ACCESS_TOKEN=<your-shopify-access-token>
   USER_EMAIL=<your-email>

## Running the Tasks
- Run Task 1: `go run main.go leaderboard`
- Run Task 2: `go run main.go order`
