package controller

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"inventory/internal/storage"
	"inventory/pkg"
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
	c.router.HandleFunc("/{id}", pkg.HandleAdapter(c.getProductById)).Methods("GET")
	c.router.HandleFunc("/{id}", pkg.HandleAdapter(c.updateProductHandler)).Methods("PUT")
	c.router.HandleFunc("/{id}", pkg.HandleAdapter(c.deleteProduct)).Methods("DELETE")

	c.router.HandleFunc("", pkg.HandleAdapter(c.createProductHandler)).Methods("POST")
	c.router.HandleFunc("", pkg.HandleAdapter(c.getAllProductHandler)).Methods("GET")
}

func (c *ProductController) createProductHandler(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var product pb.Product
	if err := protojson.Unmarshal(body, &product); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	resp, err := c.service.CreateProduct(ctx, &product)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}

func (c *ProductController) getProductById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	resp, err := c.service.GetProductById(ctx, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}

func (c *ProductController) getAllProductHandler(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	resp, err := c.service.GetAllProducts(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}

func (c *ProductController) updateProductHandler(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var product pb.Product
	if err := protojson.Unmarshal(body, &product); err != nil {
		return err
	}
	product.Id = mux.Vars(r)["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	resp, err := c.service.UpdateProduct(ctx, &product)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, 200, &resp)
}

func (c *ProductController) deleteProduct(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	err := c.service.DeleteProduct(context.Background(), id)
	if err != nil {
		return err
	}

	return pkg.WriteJson(w, 200, "Deleted")
}
