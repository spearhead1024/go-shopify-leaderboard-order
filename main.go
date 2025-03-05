package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <leaderboard|order>")
		return
	}

	switch os.Args[1] {
	case "leaderboard":
		RunLeaderboard()
	case "order":
		RunOrder()
	default:
		fmt.Println("Invalid task. Use 'leaderboard' or 'order'.")
	}
}
