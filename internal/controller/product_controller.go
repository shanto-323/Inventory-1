package controller

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"inventory/internal/storage"
	"inventory/pkg/pb"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
)

type ProductController struct {
	router  *mux.Router
	service storage.Service
}

func NewProductController(router *mux.Router, service storage.Service) *ProductController {
	newRouter := router.PathPrefix("/product").Subrouter()
	return &ProductController{
		router:  newRouter,
		service: service,
	}
}

func (c *ProductController) StartProductControoler() {
	c.router.HandleFunc("/create", c.createProductHandler).Methods("GET")
}

func (c *ProductController) createProductHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, 400, err)
		return
	}
	defer r.Body.Close()

	var product pb.Product
	if err := protojson.Unmarshal(body, &product); err != nil {
		WriteJson(w, 400, err)
		return
	}

	resp, err := c.service.CreateProduct(context.Background(), &product)
	if err != nil {
		log.Println(err)
		WriteJson(w, 400, err)
		return
	}

	WriteJson(w, 200, &resp)
	return
}

func WriteJson(w http.ResponseWriter, status int, msg any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}
