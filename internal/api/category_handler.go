package api

import (
	"encoding/json"
	"log" // Adicionado para logging de erros
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nordicmanx/videos-api/internal/repository" // Corrigido para o seu módulo
)

// CategoryHandler lida com as requisições HTTP para categorias.
type CategoryHandler struct {
	Repo *repository.CategoryRepository
}

// NewCategoryHandler cria uma nova instância do handler de categorias.
func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{Repo: repo}
}

// CreateCategory manipula a criação de uma nova categoria.
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	if requestBody.Name == "" {
		http.Error(w, "O nome da categoria é obrigatório", http.StatusBadRequest)
		return
	}

	category, err := h.Repo.CreateCategory(r.Context(), requestBody.Name)
	if err != nil {
		// ATUALIZAÇÃO: Loga o erro detalhado do banco de dados na consola
		log.Printf("Erro ao criar categoria no repositório: %v", err)
		http.Error(w, "Erro ao criar categoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GetAllCategories manipula a listagem de todas as categorias.
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Repo.GetAllCategories(r.Context())
	if err != nil {
		// ATUALIZAÇÃO: Loga o erro detalhado do banco de dados na consola
		log.Printf("Erro ao buscar categorias no repositório: %v", err)
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetCategoryByID manipula a busca de uma categoria por ID.
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	category, err := h.Repo.GetCategoryByID(r.Context(), id)
	if err != nil {
		// ATUALIZAÇÃO: Loga o erro detalhado do banco de dados na consola
		log.Printf("Erro ao buscar categoria por ID no repositório: %v", err)
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory manipula a atualização de uma categoria.
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	if err := h.Repo.UpdateCategory(r.Context(), id, requestBody.Name); err != nil {
		// ATUALIZAÇÃO: Loga o erro detalhado do banco de dados na consola
		log.Printf("Erro ao atualizar categoria no repositório: %v", err)
		http.Error(w, "Erro ao atualizar categoria", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteCategory manipula a exclusão de uma categoria.
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DeleteCategory(r.Context(), id); err != nil {
		// ATUALIZAÇÃO: Loga o erro detalhado do banco de dados na consola
		log.Printf("Erro ao deletar categoria no repositório: %v", err)
		http.Error(w, "Erro ao deletar categoria", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
