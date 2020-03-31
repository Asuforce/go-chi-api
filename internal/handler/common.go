package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Controller struct{}

func NewController() *Controller { return &Controller{} }

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Controller) Show(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	res := User{ID: id, Name: fmt.Sprint("name_", id)}
	return http.StatusOK, res, nil
}

func (h *Controller) List(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	users := []User{
		{1, "hoge"},
		{2, "fuga"},
		{3, "piyo"},
	}
	return http.StatusOK, users, nil
}

type AuthInfo struct {
	Authorization string `json:"authorization"`
}

func (h *Controller) Login(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	token := r.URL.Query().Get("token")
	if token != "token" {
		return http.StatusUnauthorized, nil, fmt.Errorf("Invalid token")
	}
	res := AuthInfo{Authorization: "admin"}
	return http.StatusOK, res, nil
}
