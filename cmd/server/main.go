package main

import (
	"fmt"
	"net/http"
	"os"

	// "github.com/MikaelLennart/metrics.git/config"
	"github.com/MikaelLennart/metrics.git/config"
	"github.com/MikaelLennart/metrics.git/internal/router"
	"github.com/MikaelLennart/metrics.git/internal/store"
)

// Server Main ...
func main() {
	LetServerAddress := os.Getenv("SERVER_ADDRESS")
	fmt.Printf("Считанное значение SERVER_ADDRESS: [%s]\n", LetServerAddress)

	cfg := config.ServerConfig()
	s := store.NewMemStorage()
	r := router.NewRouter(s)
	fmt.Printf("Server started at %s\r\n", cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)

}
