package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vr-shope/internal/models"
	"vr-shope/internal/repository"
	"vr-shope/internal/utils/uuids"
)

type PurchaseService struct {
	repo *repository.PurchaseRepository
}

func NewPurchaseService(repo *repository.PurchaseRepository) *PurchaseService {
	return &PurchaseService{repo}
}

func (s *PurchaseService) Create(ctx context.Context, purchase *models.Purchase) error {
	purchaseServ := repository.Purchase{
		UserID:    uuids.IntToUUID(int64(purchase.UserID)),
		ProductID: uuids.IntToUUID(int64(purchase.ProductID)),
		Date:      purchase.Date,
	}

	purchases, err := s.repo.CheckMeans(ctx, purchaseServ.UserID, purchaseServ.ProductID)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, purchases)
	if err != nil {
		return err
	}

	return nil
}

func (s *PurchaseService) Get(ctx context.Context, id int64) (*models.Purchase, error) {
	purchaseRepo, err := s.repo.Get(ctx, uuids.IntToUUID(id))
	if err != nil {
		return nil, err
	}

	return &models.Purchase{
		ID:         uuids.UUIDToInt(purchaseRepo.ID),
		UserID:     uuids.UUIDToInt(purchaseRepo.UserID),
		ProductID:  uuids.UUIDToInt(purchaseRepo.ProductID),
		Date:       purchaseRepo.Date,
		WalletUSDT: purchaseRepo.WalletUSDT,
		Cost:       purchaseRepo.Cost,
	}, nil
}

func (s *PurchaseService) GetAll(ctx context.Context) ([]*models.Purchase, error) {
	purchasesRepo, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var purchases []*models.Purchase
	for _, purchase := range purchasesRepo {
		purchases = append(purchases, &models.Purchase{
			ID:         uuids.UUIDToInt(purchase.ID),
			UserID:     uuids.UUIDToInt(purchase.UserID),
			ProductID:  uuids.UUIDToInt(purchase.ProductID),
			Date:       time.Now(),
			WalletUSDT: purchase.WalletUSDT,
			Cost:       purchase.Cost,
		})
	}

	return purchases, nil
}

func (s *PurchaseService) Update(ctx context.Context, purchase *models.Purchase) error {
	exists, err := s.repo.ExistsByID(ctx, uuids.IntToUUID(int64(purchase.ID)))
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("purchase not exist")
	}

	purchaseRepo := &repository.Purchase{
		ID:         uuids.IntToUUID(int64(purchase.ID)),
		UserID:     uuids.IntToUUID(int64(purchase.UserID)),
		ProductID:  uuids.IntToUUID(int64(purchase.ProductID)),
		Date:       time.Now(),
		WalletUSDT: purchase.WalletUSDT,
		Cost:       purchase.Cost,
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
