package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"inventory/pkg/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	CreateProduct(
		ctx context.Context,
		product *pb.Product,
	) (*pb.Product, error)
	GetProductById(ctx context.Context, productId string) (*pb.Product, error)
}

type productService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(
	ctx context.Context,
	product *pb.Product,
) (*pb.Product, error) {
	Id := uuid.New().String()
	product.Id = Id
	product.DateAdded = timestamppb.Now()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(product); err != nil {
		return nil, returnServiceString(err)
	}
	if err := s.repo.Create(ctx, buf, Id); err != nil {
		return nil, returnServiceString(err)
	}

	return product, nil
}

func (s *productService) GetProductById(ctx context.Context, productId string) (*pb.Product, error) {
	resp, err := s.repo.GetProduct(ctx, productId)
	if err != nil {
		return nil, returnServiceString(err)
	}

	return resp, nil
}

func returnServiceString(m any) error {
	return fmt.Errorf("service: %s\n", m)
}
