package api

import (
	"database/sql"
	"log"
	"net/http"
	"serverAPI/service/user"

	"github.com/gorilla/mux"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	log.Println("Starting server...")

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(subrouter)

	log.Printf("Listening on %v", s.addr)

	return http.ListenAndServe(s.addr, router)
}
