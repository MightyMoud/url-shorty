package main

import (
	"fmt"
	"github.com/ms-mousa/url-shorty/middleware"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})
	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.LoggerMiddleware(router),
	}

	fmt.Println("Server listening on 3000")
	server.ListenAndServe()
}
