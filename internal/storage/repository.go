package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"inventory/pkg/pb"

	"github.com/elastic/go-elasticsearch/v9"
)

type Repository interface {
	Upsert(ctx context.Context, product *pb.Product, productId string) error
	Product(ctx context.Context, productId string) (*pb.Product, error)
	Products(ctx context.Context) ([]*pb.Product, error)
	Delete(ctx context.Context, productId string) error

	// Analytics
	MinStock(ctx context.Context, level int) ([]*pb.Product, error)
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

type allDocument struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []document `json:"hits"`
	} `json:"hits"`
}

const (
	INVENTORY_INDEX = "inventroy"
)

func (r *inventoryRepository) Upsert(ctx context.Context, product *pb.Product, productId string) error {
	doc := map[string]interface{}{
		"doc": map[string]interface{}{
			"type":       product.GetType(),
			"brand":      product.GetBrand(),
			"name":       product.GetName(),
			"model":      product.GetModel(),
			"stock":      product.GetStock(),
			"specs":      product.GetSpecs(),
			"warranty":   product.GetWarranty(),
			"supplier":   product.GetSupplier(),
			"date_added": product.GetDateAdded(),
			"note":       product.GetNote(),
		},
		"doc_as_upsert": true,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("json encode error: %w", err)
	}

	_, err := r.client.Update(
		INVENTORY_INDEX,
		productId,
		&buf,
		r.client.Update.WithContext(ctx),
		r.client.Update.WithRefresh("true"),
	)
	if err != nil {
		return returnString(err)
	}
	return nil
}

func (r *inventoryRepository) Delete(ctx context.Context, productId string) error {
	_, err := r.client.Delete(
		INVENTORY_INDEX,
		productId,
		r.client.Delete.WithContext(ctx),
		r.client.Delete.WithRefresh("true"),
	)
	if err != nil {
		return returnString(err)
	}
	return nil
}

func (r *inventoryRepository) Product(ctx context.Context, productId string) (*pb.Product, error) {
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

	document.Source.Id = document.Id
	return &document.Source, nil
}

func (r *inventoryRepository) Products(ctx context.Context) ([]*pb.Product, error) {
	stringQuery := `{
		"query": {
			"match_all": {}
		}
	}`

	return r.searchResult(ctx, stringQuery)
}

// Analytics
func (r *inventoryRepository) MinStock(ctx context.Context, level int) ([]*pb.Product, error) {
	stringQuery := fmt.Sprintf(`{
		"query": {
			"range": {
				"stock": {
					"lte" : %d
				}
			}
		}
	}`, level)

	return r.searchResult(ctx, stringQuery)
}

func (r *inventoryRepository) searchResult(ctx context.Context, query string) ([]*pb.Product, error) {
	resp, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(INVENTORY_INDEX),
		r.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, returnString(err)
	}
	if resp.IsError() {
		return nil, returnString("field is empty")
	}

	var allDocument allDocument
	if err := json.NewDecoder(resp.Body).Decode(&allDocument); err != nil {
		return nil, returnString(err)
	}

	var products []*pb.Product
	for i := range allDocument.Hits.Hits {
		product := &allDocument.Hits.Hits[i].Source
		product.Id = allDocument.Hits.Hits[i].Id
		products = append(products, product)
	}

	return products, nil
}

func returnString(m any) error {
	return fmt.Errorf("repository: %s\n", m)
}
