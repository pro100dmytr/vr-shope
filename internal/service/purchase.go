package service

import (
	"context"
	"errors"
	"time"
	"vr-shope/internal/models/repositories"
	"vr-shope/internal/models/services"
	"vr-shope/internal/repository"
	"vr-shope/internal/utils/uuids"
)

type PurchaseService struct {
	repo *repository.PurchaseRepository
}

func NewPurchaseService(repo *repository.PurchaseRepository) *PurchaseService {
	return &PurchaseService{repo}
}

func (s *PurchaseService) Create(ctx context.Context, purchase *services.Purchase) error {
	purchaseServ := repositories.Purchase{
		ID:     uuids.IntToUUID(int64(purchase.ID)),
		UserID: uuids.IntToUUID(int64(purchase.UserID)),
		Cost:   purchase.Cost,
		Date:   purchase.Date,
	}

	err := s.repo.Create(ctx, &purchaseServ)
	if err != nil {
		return err
	}

	return nil
}

func (s *PurchaseService) Get(ctx context.Context, id int64) (*services.Purchase, error) {
	purchaseRepo, err := s.repo.Get(ctx, uuids.IntToUUID(id))
	if err != nil {
		return nil, err
	}

	return &services.Purchase{
		ID:     uuids.UUIDToInt(purchaseRepo.ID),
		UserID: uuids.UUIDToInt(purchaseRepo.UserID),
		Cost:   purchaseRepo.Cost,
		Date:   purchaseRepo.Date,
	}, nil
}

func (s *PurchaseService) GetAll(ctx context.Context) ([]*services.Purchase, error) {
	purchasesRepo, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var purchases []*services.Purchase
	for _, p := range purchasesRepo {
		purchases = append(purchases, &services.Purchase{
			ID:     uuids.UUIDToInt(p.ID),
			UserID: uuids.UUIDToInt(p.UserID),
			Cost:   p.Cost,
			Date:   time.Now(),
		})
	}

	return purchases, nil
}

func (s *PurchaseService) Update(ctx context.Context, purchase *services.Purchase) error {
	exists, err := s.repo.ExistsByID(ctx, uuids.IntToUUID(int64(purchase.ID)))
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("purchase not found")
	}

	purchaseRepo := &repositories.Purchase{
		ID:     uuids.IntToUUID(int64(purchase.ID)),
		UserID: uuids.IntToUUID(int64(purchase.UserID)),
		Cost:   purchase.Cost,
		Date:   time.Now(),
	}

	err = s.repo.Update(ctx, purchaseRepo)
	if err != nil {
		return err
	}

	return nil
}

func (s *PurchaseService) Delete(ctx context.Context, id int64) error {
	exists, err := s.repo.ExistsByID(ctx, uuids.IntToUUID(id))
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("purchase not found")
	}

	err = s.repo.Delete(ctx, uuids.IntToUUID(id))
	if err != nil {
		return err
	}

	return nil
}
