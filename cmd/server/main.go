package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/MikaelLennart/metrics.git/internal/router"
	"github.com/MikaelLennart/metrics.git/internal/store"
)

// Server Main ...
func main() {
	address := flag.String("a", "8080", "server port adress")
	flag.Parse()
	port := ":" + *address
	s := store.NewMemStorage()
	r := router.NewRouter(s)

	fmt.Printf("Server started at %s\r\n", port)
	http.ListenAndServe(port, r)

}
