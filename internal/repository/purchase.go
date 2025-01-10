package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"vr-shope/internal/config"
	"vr-shope/internal/models/repositories"
	"vr-shope/internal/storage/postgresql"
)

type PurchaseRepository struct {
	db *sql.DB
}

func (s *PurchaseRepository) Close() error {
	return postgresql.CloseConn(s.db)
}

func NewPurchaseStorage(cfg *config.Config) (*PurchaseRepository, error) {
	db, err := postgresql.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &PurchaseRepository{db: db}, nil
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *repositories.Purchase) error {
	query := `
		INSERT INTO purchases (id, user_id, cost, date)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, purchase.ID, purchase.UserID, purchase.Cost, purchase.Date)
	return err
}

func (r *PurchaseRepository) Get(ctx context.Context, id uuid.UUID) (*repositories.Purchase, error) {
	query := `
		SELECT id, user_id, cost, date
		FROM purchases
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var purchase repositories.Purchase
	err := row.Scan(&purchase.ID, &purchase.UserID, &purchase.Cost, &purchase.Date)
	if err == sql.ErrNoRows {
		return nil, errors.New("purchase not found")
	} else if err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (r *PurchaseRepository) GetAll(ctx context.Context) ([]*repositories.Purchase, error) {
	query := `
		SELECT id, user_id, cost, date
		FROM purchases`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*repositories.Purchase
	for rows.Next() {
		var purchase repositories.Purchase
		err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.Cost, &purchase.Date)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, &purchase)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return purchases, nil
}

func (r *PurchaseRepository) Update(ctx context.Context, purchase *repositories.Purchase) error {
	query := `
		UPDATE purchases
		SET user_id = $1, cost = $2, date = $3
		WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, purchase.UserID, purchase.Cost, purchase.Date, purchase.ID)
	return err
}

func (r *PurchaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM purchases
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PurchaseRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM purchases
			WHERE id = $1
		)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}
