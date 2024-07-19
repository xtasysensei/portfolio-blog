package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"portfolio-blog/internal"
	"portfolio-blog/pkg/handlers"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func Routes(route *chi.Mux) {
	route.Get("/", handlers.Index().ServeHTTP)
	route.Get("/health", Health)
	route.Get("/show-blogs", handlers.ShowBlogPost().ServeHTTP)
	route.Get("/admin-login", handlers.AdminLogin().ServeHTTP)
	route.Post("/create-admin", internal.CreateAdmin)
	route.Post("/create-post", handlers.CreateBlogPost)

	route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func GracefulShutdown(server *http.Server) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Service interrupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown error: %v", err)
		}
		log.Println("Shutdown complete")
		close(idleConnsClosed)
	}()

	<-idleConnsClosed
	log.Println("Service stopped")
}

func StartServer() (*http.Server, *chi.Mux) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	Routes(router)
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return server, router
}
