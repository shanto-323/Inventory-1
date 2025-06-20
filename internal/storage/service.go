package storage

import (
	"context"
	"fmt"
	"strconv"

	"inventory/pkg"
	"inventory/pkg/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	CreateProduct(ctx context.Context, product *pb.Product) (*pb.Product, error)
	GetProductById(ctx context.Context, productId string) (*pb.Product, error)
	UpdateProduct(ctx context.Context, product *pb.Product) (*pb.Product, error)
	DeleteProduct(ctx context.Context, productId string) error

	// Analytics
	FindMinStock(ctx context.Context, levelString string) ([]*pb.Product, error)
	GetProductBySearchFilter(ctx context.Context, filterModel *pkg.FilterModel) ([]*pb.Product, error)
}

type productService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, product *pb.Product) (*pb.Product, error) {
	Id := uuid.New().String()
	product.Id = Id
	product.DateAdded = timestamppb.Now()

	if err := s.repo.Upsert(ctx, product, Id); err != nil {
		return nil, returnServiceString(err)
	}

	return product, nil
}

func (s *productService) GetProductById(ctx context.Context, productId string) (*pb.Product, error) {
	resp, err := s.repo.Product(ctx, productId)
	if err != nil {
		return nil, returnServiceString(err)
	}

	return resp, nil
}

func (s *productService) UpdateProduct(ctx context.Context, product *pb.Product) (*pb.Product, error) {
	resp, err := s.GetProductById(ctx, product.Id)
	if err != nil {
		return nil, returnServiceString(err)
	}
	mutationHelper(resp, product)

	err = s.repo.Upsert(ctx, resp, resp.Id)
	if err != nil {
		return nil, returnServiceString(err)
	}
	return resp, nil
}

func (s *productService) DeleteProduct(ctx context.Context, productId string) error {
	return s.repo.Delete(ctx, productId)
}

// Analytics
func (s *productService) FindMinStock(ctx context.Context, levelString string) ([]*pb.Product, error) {
	level, err := strconv.ParseInt(levelString, 10, 64)
	if err != nil {
		return nil, returnServiceString(err)
	}
	return s.repo.MinStock(ctx, int(level))
}

func (s *productService) GetProductBySearchFilter(ctx context.Context, filterModel *pkg.FilterModel) ([]*pb.Product, error) {
	resp, err := s.repo.SearchWithFilter(ctx, filterModel)
	if err != nil {
		return nil, returnServiceString(err)
	}

	return resp, nil
}

func mutationHelper(dbData *pb.Product, product *pb.Product) {
	if product.Type != "" {
		dbData.Type = product.Type
	}
	if product.Brand != "" {
		dbData.Brand = product.Brand
	}
	if product.Name != "" {
		dbData.Name = product.Name
	}
	if product.Model != "" {
		dbData.Model = product.Model
	}
	if product.Stock != 0 {
		dbData.Stock = product.Stock
	}
	if product.Warranty != "" {
		dbData.Warranty = product.Warranty
	}
	if product.Supplier != "" {
		dbData.Supplier = product.Supplier
	}
	if product.DateAdded != nil {
		dbData.DateAdded = product.DateAdded
	}
	if product.Note != "" {
		dbData.Note = product.Note
	}

	if product.Specs != nil {
		if dbData.Specs == nil {
			dbData.Specs = make(map[string]string)
		}
		for k, v := range product.Specs {
			if v != "" {
				dbData.Specs[k] = v
			}
		}
	}
}

func returnServiceString(m any) error {
	return fmt.Errorf("service: %s\n", m)
}
