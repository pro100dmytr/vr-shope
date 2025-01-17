package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"vr-shope/internal/models"
	"vr-shope/internal/repository"
	"vr-shope/internal/uuids"

	"github.com/google/uuid"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo}
}

func (s *ProductService) Create(ctx context.Context, product *models.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}

	productID := uuid.New()
	repoProduct := &repository.Product{
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

func (s *ProductService) Get(ctx context.Context, id int) (*models.Product, error) {
	repoProduct, err := s.repo.Get(ctx, uuids.IntToUUID(int64(id)))
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:            uuids.UUIDToInt(repoProduct.ID),
		Name:          repoProduct.Name,
		Cost:          repoProduct.Cost,
		QuantityStock: repoProduct.QuantityStock,
		Guarantees:    repoProduct.Guarantees,
		Country:       repoProduct.Country,
		Like:          repoProduct.Like,
	}, nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]*models.Product, error) {
	repoProducts, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var products []*models.Product
	for _, repoProduct := range repoProducts {
		products = append(products, &models.Product{
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

func (s *ProductService) Update(ctx context.Context, product *models.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}

	repoProduct := &repository.Product{
		ID:            uuids.IntToUUID(int64(product.ID)),
		Name:          product.Name,
		Cost:          product.Cost,
		QuantityStock: product.QuantityStock,
		Guarantees:    product.Guarantees,
		Country:       product.Country,
		Like:          product.Like,
	}

	err := s.repo.Update(ctx, repoProduct)
	if err != nil {
		return err
	}

	return nil
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

func (s *ProductService) GetProductByName(ctx context.Context, name string) ([]*models.Product, error) {
	repoProducts, err := s.repo.GetForName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var products []*models.Product
	for _, repoProduct := range repoProducts {
		product := &models.Product{
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

func (s *ProductService) GetProductsWithPagination(ctx context.Context, limit, offset string) ([]*models.Product, error) {
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

	var products []*models.Product
	for _, repoProduct := range repoProducts {
		product := &models.Product{
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
