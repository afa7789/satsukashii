package main

import (
	"github.com/afa7789/satsukashii/internal/server"
)

// Main is the command function that setup initial structs and value.
// After it, it will start the server.
func main() {
	// Setup Repositories
	// r := database.NewRepositories()
	// if r == nil {
	// 	print("db as nil")
	// }

	// si := &domain.ServerInput{
	// 	// Reps: r,
	// }

	// Setup and start server
	s := server.New()
	s.Start(8080)
}
