package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
)

type Server struct {
	router *chi.Mux
}

func New() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

func (s *Server) Init(env string) {
	log.Printf("env: %s", env)
}

func (s *Server) Middleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(time.Second * 600))
}

func (s *Server) Router() {
	h := NewHandler()
	s.router.Route("/api", func(api chi.Router) {
		api.Use(Auth("db connection"))
		api.Route("/members", func(members chi.Router) {
			members.Get("/{id}", h.Show)
			members.Get("/", h.List)
		})
	})
}

func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "admin" {
				respondError(w, http.StatusUnauthorized, fmt.Errorf("Authorization Forbidden"))
			}
		})
	}
	return
}

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	type json struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	res := json{ID: id, Name: fmt.Sprint("name_", id)}
	respondJSON(w, http.StatusOK, res)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users := []struct {
		ID   int    `json:"id"`
		User string `json:"user"`
	}{
		{1, "hoge"},
		{2, "fuga"},
		{3, "piyo"},
	}
	respondJSON(w, http.StatusOK, users)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != "token" {
		respondError(w, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
		return
	}
	type json struct {
		Message string `json:message`
	}
	res := json{Message: "auth ok"}
	respondJSON(w, http.StatusOK, res)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "   ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, err error) {
	log.Printf("err: %v", err)
	if e, ok := err.(*HTTPError); ok {
		respondJSON(w, e.Code, e)
	} else if err != nil {
		he := HTTPError{
			Code:    code,
			Message: err.Error(),
		}
		respondJSON(w, code, he)
	}
}

type HTTPError struct {
	Code    int
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}

func main() {
	var (
		port   = flag.String("port", "8080", "Addr to bind")
		env    = flag.String("env", "develop", "Exec env (local, beta, production")
		gendoc = flag.Bool("gendoc", true, "Generate document")
	)
	flag.Parse()

	s := New()
	s.Init(*env)
	s.Middleware()
	s.Router()

	if *gendoc {
		doc := docgen.MarkdownRoutesDoc(s.router, docgen.MarkdownOpts{
			ProjectPath: "github.com/Asuforce/go-chi-api",
			Intro:       "generated docs.",
		})
		file, err := os.Create(`doc.md`)
		if err != nil {
			log.Printf("err: %v", err)
		}
		defer file.Close()
		file.Write(([]byte)(doc))
	}
	log.Println("starting app")
	http.ListenAndServe(fmt.Sprint(":", *port), s.router)
}
