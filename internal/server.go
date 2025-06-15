package internal

import (
	"fmt"
	"net/http"

	"inventory/internal/controller"
	"inventory/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	ipAddr  string
	service storage.Service
}

func NewServer(ipAddr string, service storage.Service) *Server {
	ip := fmt.Sprintf(":%s", ipAddr)
	return &Server{
		ipAddr:  ip,
		service: service,
	}
}

func (s *Server) Start() error {
	mux := mux.NewRouter()
	router := mux.PathPrefix("/api/v1").Subrouter()

	productController := controller.NewProductController(router, s.service)
	productController.StartProductControoler()

	fmt.Printf("server running on port %s....\n", s.ipAddr)
	return http.ListenAndServe(s.ipAddr, mux)
}
