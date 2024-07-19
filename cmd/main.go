package main

import (
	"log"
	"net/http"
	"portfolio-blog/pkg/routes"
)

func main() {
	server, _ := routes.StartServer()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	log.Printf("Server is ready to handle requests at %s", server.Addr)
	routes.GracefulShutdown(server)
}
