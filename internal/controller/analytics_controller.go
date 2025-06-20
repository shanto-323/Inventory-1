package controller

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	resp, err := c.service.FindMinStock(ctx, stock)
	if err != nil {
		log.Println(err)
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}

func (c *AnalyticsController) searchFilterHandler(w http.ResponseWriter, r *http.Request) error {
	productFilter := &pkg.FilterModel{}
	if err := json.NewDecoder(r.Body).Decode(productFilter); err != nil && err != io.EOF {
		return err
	}
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	resp, err := c.service.GetProductBySearchFilter(ctx, productFilter)
	if err != nil {
		log.Println(err)
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}
