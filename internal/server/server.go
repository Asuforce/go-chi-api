package server

import (
	"log"
	"time"

	h "github.com/Asuforce/go-chi-api/internal/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	c := h.NewController()
	s.router.Route("/api", func(api chi.Router) {
		api.Use(Auth("db connection"))
		api.Route("/members", func(members chi.Router) {
			members.Get("/{id}", h.Handler(c.Show).ServeHTTP)
			members.Get("/", h.Handler(c.List).ServeHTTP)
		})
	})
}

func (s *Server) GetRouter() *chi.Mux { return s.router }
