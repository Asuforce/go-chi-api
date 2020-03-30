package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	res := User{ID: id, Name: fmt.Sprint("name_", id)}
	respondJSON(w, http.StatusOK, res)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{1, "hoge"},
		{2, "fuga"},
		{3, "piyo"},
	}
	respondJSON(w, http.StatusOK, users)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != "token" {
		RespondError(w, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
		return
	}
	type json struct {
		Message string `json:message`
	}
	res := json{Message: "auth ok"}
	respondJSON(w, http.StatusOK, res)
}
