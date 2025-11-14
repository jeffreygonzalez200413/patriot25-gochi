package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/juhun32/patriot25-gochi/go/middleware"
	"github.com/juhun32/patriot25-gochi/go/repo"
)

type TodosHandler struct {
	TodoRepo *repo.TodoRepo
}

func NewTodosHandler(repo *repo.TodoRepo) *TodosHandler {
	return &TodosHandler{TodoRepo: repo}
}

func (h *TodosHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	todos, err := h.TodoRepo.ListTodos(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to list todos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"todos": todos})
}

func (h *TodosHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Text  string `json:"text"`
		DueAt *int64 `json:"dueAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	todoID := uuid.NewString()
	todo, err := h.TodoRepo.CreateTodo(r.Context(), userID, todoID, body.Text, body.DueAt)
	if err != nil {
		http.Error(w, "failed to create todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
