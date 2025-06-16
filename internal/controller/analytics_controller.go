package controller

import (
	"context"
	"log"
	"net/http"

	"inventory/internal/storage"
	"inventory/pkg"

	"github.com/gorilla/mux"
)

type AnalyticsController struct {
	router  *mux.Router
	service storage.Service
}

func NewAnalyticsController(router *mux.Router, service storage.Service) *AnalyticsController {
	newRouter := router.PathPrefix("/analytics").Subrouter()
	return &AnalyticsController{
		router:  newRouter,
		service: service,
	}
}

func (c *AnalyticsController) StartAnalyticsControoler() {
	c.router.HandleFunc("/stock", pkg.HandleAdapter(c.getStock)).Methods("GET")
}

func (c *AnalyticsController) getStock(w http.ResponseWriter, r *http.Request) error {
	stock := r.URL.Query().Get("level")
	if stock == "" {
		stock = "3"
	}
	resp, err := c.service.FindMinStock(context.Background(), stock)
	if err != nil {
		log.Println(err)
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}
