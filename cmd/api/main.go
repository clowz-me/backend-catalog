package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/clowz-me/catalog/internal/handler"
	"github.com/clowz-me/catalog/internal/repository"
	"github.com/clowz-me/catalog/internal/service"
	"github.com/clowz-me/pkg/database"
)

func main() {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "clowz"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "clowz_catalog"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database tables
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS categories (
			id UUID PRIMARY KEY,
			establishment_id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			"order" INT DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		ALTER TABLE categories ADD COLUMN IF NOT EXISTS image_url VARCHAR(500);

		CREATE INDEX IF NOT EXISTS idx_categories_establishment ON categories(establishment_id);

		CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY,
			establishment_id UUID NOT NULL,
			category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price BIGINT NOT NULL,
			image_url TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_products_establishment ON products(establishment_id);
		CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
	`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}

	repo := repository.NewPostgresCatalogRepository(db)
	svc := service.NewCatalogService(repo)
	h := handler.NewCatalogHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://admin-local.qperto.com", "https://api-local.qperto.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1/catalog", func(r chi.Router) {
		h.RegisterRoutes(r)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	port := getEnv("PORT", "8086")
	log.Printf("Catalog service starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
