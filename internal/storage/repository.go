package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"inventory/pkg/pb"

	"github.com/elastic/go-elasticsearch/v9"
)

type Repository interface {
	Create(ctx context.Context, product bytes.Buffer, productId string) error
	GetProduct(ctx context.Context, productId string) (*pb.Product, error)
}

type inventoryRepository struct {
	client *elasticsearch.Client
}

func NewRepository(dsn []string) (Repository, error) {
	client, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: dsn,
		},
	)
	if err != nil {
		return nil, err
	}

	return &inventoryRepository{
		client: client,
	}, nil
}

type document struct {
	Index  string     `json:"_index"`
	Id     string     `json:"_id"`
	Found  bool       `json:"found"`
	Source pb.Product `json:"_source"`
}

const (
	INVENTORY_INDEX = "inventroy"
)

func (r *inventoryRepository) Create(ctx context.Context, product bytes.Buffer, productId string) error {
	_, err := r.client.Index(
		INVENTORY_INDEX,
		&product,
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(productId),
		r.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return returnString(err)
	}
	return nil
}

func (r *inventoryRepository) GetProduct(ctx context.Context, productId string) (*pb.Product, error) {
	resp, err := r.client.Get(
		INVENTORY_INDEX,
		productId,
		r.client.Get.WithContext(ctx),
		r.client.Get.WithRealtime(true),
	)
	if err != nil {
		return nil, returnString(err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, returnString("field is empty")
	}

	var document document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, returnString(err)
	}

	return &document.Source, nil
}

func returnString(m any) error {
	return fmt.Errorf("repository: %s\n", m)
}
