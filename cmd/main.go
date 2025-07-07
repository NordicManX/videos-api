package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nordicmanx/videos-api/internal/api"
	"github.com/nordicmanx/videos-api/internal/repository"
)

func main() {

	databaseUrl := "user=postgres password=postgres dbname=postgres host=localhost port=5433 sslmode=disable"

	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível conectar ao banco de dados: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	categoryRepo := repository.NewCategoryRepository(dbpool)
	categoryHandler := api.NewCategoryHandler(categoryRepo)

	r.Route("/categories", func(r chi.Router) {
		r.Post("/", categoryHandler.CreateCategory)
		r.Get("/", categoryHandler.GetAllCategories)
		r.Get("/{id}", categoryHandler.GetCategoryByID)
		r.Put("/{id}", categoryHandler.UpdateCategory)
		r.Delete("/{id}", categoryHandler.DeleteCategory)
	})

	log.Println("Servidor iniciando na porta :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Não foi possível iniciar o servidor: %v", err)
	}
}
