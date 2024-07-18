package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/snwfdhmp/errlog"
)

type Blog struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func init() {
	if err := godotenv.Load(); errlog.Debug(err) {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
func gracefulShutdown(server *http.Server) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Service interrupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); errlog.Debug(err) {
			log.Printf("HTTP server Shutdown error: %v", err)
		}
		log.Println("Shutdown complete")
		close(idleConnsClosed)
	}()

	<-idleConnsClosed
	log.Println("Service stopped")
}

func startServer() *http.Server {
	router := httprouter.New()
	//routes.Routes(router)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return server
}

func main() {

	main_server := startServer()

	log.Printf("Server is starting at port %s", main_server.Addr)
	go func() {
		if err := main_server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
	gracefulShutdown(main_server)
}
