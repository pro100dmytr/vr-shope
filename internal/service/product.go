package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"vr-shope/internal/models/repositories"
	"vr-shope/internal/models/services"
	"vr-shope/internal/repository"
	"vr-shope/internal/utils/uuids"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo}
}

func (s *ProductService) Create(ctx context.Context, product *services.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}

	productID := uuid.New()
	repoProduct := &repositories.Product{
		ID:            productID,
		Name:          product.Name,
		Cost:          product.Cost,
		QuantityStock: product.QuantityStock,
		Guarantees:    product.Guarantees,
		Country:       product.Country,
		Like:          product.Like,
	}

	err := s.repo.Create(ctx, repoProduct)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) Get(ctx context.Context, id int) (*services.Product, error) {
	repoProduct, err := s.repo.Get(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		return nil, err
	}

	return &services.Product{
		ID:            uuids.UUIDToInt(repoProduct.ID),
		Name:          repoProduct.Name,
		Cost:          repoProduct.Cost,
		QuantityStock: repoProduct.QuantityStock,
		Guarantees:    repoProduct.Guarantees,
		Country:       repoProduct.Country,
		Like:          repoProduct.Like,
	}, nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]*services.Product, error) {
	repoProducts, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var products []*services.Product
	for _, repoProduct := range repoProducts {
		products = append(products, &services.Product{
			ID:            uuids.UUIDToInt(repoProduct.ID),
			Name:          repoProduct.Name,
			Cost:          repoProduct.Cost,
			QuantityStock: repoProduct.QuantityStock,
			Guarantees:    repoProduct.Guarantees,
			Country:       repoProduct.Country,
			Like:          repoProduct.Like,
		})
	}

	return products, nil
}

func (s *ProductService) Update(ctx context.Context, product *services.Product) (*services.Product, error) {
	if product.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}

	repoProduct := &repositories.Product{
		ID:            uuids.IntToUUID(int64(product.ID)),
		Name:          product.Name,
		Cost:          product.Cost,
		QuantityStock: product.QuantityStock,
		Guarantees:    product.Guarantees,
		Country:       product.Country,
		Like:          product.Like,
	}

	updatedProduct, err := s.repo.Update(ctx, repoProduct)
	if err != nil {
		return nil, err
	}

	return &services.Product{
		ID:            uuids.UUIDToInt(updatedProduct.ID),
		Name:          updatedProduct.Name,
		Cost:          updatedProduct.Cost,
		QuantityStock: updatedProduct.QuantityStock,
		Guarantees:    updatedProduct.Guarantees,
		Country:       updatedProduct.Country,
		Like:          updatedProduct.Like,
	}, nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, uuids.IntToUUID(int64(id)))
}

func (s *ProductService) AddLike(ctx context.Context, id int) error {
	err := s.repo.AddLike(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}

		return err
	}

	return nil
}

func (s *ProductService) RemoveLike(ctx context.Context, id int) error {
	err := s.repo.RemoveLike(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}

		return err
	}

	return nil
}

func (s *ProductService) GetProductByName(ctx context.Context, name string) ([]*services.Product, error) {
	repoProducts, err := s.repo.GetForName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var products []*services.Product
	for _, repoProduct := range repoProducts {
		product := &services.Product{
			ID:            uuids.UUIDToInt(repoProduct.ID),
			Name:          repoProduct.Name,
			Cost:          repoProduct.Cost,
			QuantityStock: repoProduct.QuantityStock,
			Guarantees:    repoProduct.Guarantees,
			Country:       repoProduct.Country,
			Like:          repoProduct.Like,
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *ProductService) GetProductsWithPagination(ctx context.Context, limit, offset string) ([]*services.Product, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		return nil, err
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		return nil, err
	}

	repoProducts, err := s.repo.GetProducts(ctx, offsetInt, limitInt)
	if err != nil {
		return nil, err
	}

	var products []*services.Product
	for _, repoProduct := range repoProducts {
		product := &services.Product{
			ID:            uint64(uuids.UUIDToInt(repoProduct.ID)),
			Name:          repoProduct.Name,
			Cost:          repoProduct.Cost,
			QuantityStock: repoProduct.QuantityStock,
			Guarantees:    repoProduct.Guarantees,
			Country:       repoProduct.Country,
			Like:          repoProduct.Like,
		}
		products = append(products, product)
	}

	return products, nil
}
