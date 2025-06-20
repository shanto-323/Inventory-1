package controller

import (
	"context"
	"encoding/json"
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
	c.router.HandleFunc("/stock", pkg.HandleAdapter(c.getStockHandler)).Methods("GET")
	c.router.HandleFunc("/search", pkg.HandleAdapter(c.searchFilterHandler)).Methods("GET")
}

func (c *AnalyticsController) getStockHandler(w http.ResponseWriter, r *http.Request) error {
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

func (c *AnalyticsController) searchFilterHandler(w http.ResponseWriter, r *http.Request) error {
	var productFilter *pkg.FilterModel
	if err := json.NewDecoder(r.Body).Decode(productFilter); err != nil {
		return err
	}
	resp, err := c.service.GetProductBySearchFilter(context.Background(), productFilter)
	if err != nil {
		log.Println(err)
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}
