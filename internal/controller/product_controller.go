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
	newRouter := router.PathPrefix("/products").Subrouter()
	return &ProductController{
		router:  newRouter,
		service: service,
	}
}

func (c *ProductController) StartProductControoler() {
	c.router.HandleFunc("/{id}", c.getProductById).Methods("GET")
	c.router.HandleFunc("/{id}", c.updateProductHandler).Methods("PUT")
	c.router.HandleFunc("/{id}", c.deleteProduct).Methods("DELETE")

	c.router.HandleFunc("", c.createProductHandler).Methods("POST")
	c.router.HandleFunc("", c.getAllProductHandler).Methods("GET")
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

func (c *ProductController) getProductById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	resp, err := c.service.GetProductById(context.Background(), id)
	if err != nil {
		log.Println(err)
		WriteJson(w, 400, err)
		return
	}

	WriteJson(w, 200, &resp)
	return
}

func (c *ProductController) getAllProductHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := c.service.GetAllProducts(context.Background())
	if err != nil {
		log.Println(err)
		WriteJson(w, 400, err)
		return
	}

	WriteJson(w, 200, &resp)
	return
}

func (c *ProductController) updateProductHandler(w http.ResponseWriter, r *http.Request) {
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
	product.Id = mux.Vars(r)["id"]
	resp, err := c.service.UpdateProduct(context.Background(), &product)
	if err != nil {
		log.Println(err)
		WriteJson(w, 400, err)
		return
	}

	WriteJson(w, 200, &resp)
	return
}

func (c *ProductController) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := c.service.DeleteProduct(context.Background(), id)
	if err != nil {
		log.Println(err)
		WriteJson(w, 400, err)
		return
	}

	WriteJson(w, 200, "deleted")
	return
}

func WriteJson(w http.ResponseWriter, status int, msg any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}
