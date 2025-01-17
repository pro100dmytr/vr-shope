package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type PurchaseRepository struct {
	db *sql.DB
}

func NewPurchaseStorage(db *sql.DB) (*PurchaseRepository, error) {
	return &PurchaseRepository{db: db}, nil
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *Purchase) error {
	query := `
		INSERT INTO purchases (id, user_id, product_id, created_at, wallet_usdt, cost)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(
		ctx,
		query,
		purchase.ID,
		purchase.UserID,
		purchase.ProductID,
		purchase.Date,
		purchase.WalletUSDT,
		purchase.Cost,
	)
	return err
}

func (r *PurchaseRepository) Get(ctx context.Context, id uuid.UUID) (*Purchase, error) {
	query := `
		SELECT id, user_id, product_id, created_at, wallet_usdt, cost date
		FROM purchases
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var purchase Purchase
	err := row.Scan(
		&purchase.ID,
		&purchase.UserID,
		&purchase.ProductID,
		&purchase.Date,
		&purchase.WalletUSDT,
		&purchase.Cost,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("purchase not found")
	} else if err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (r *PurchaseRepository) GetAll(ctx context.Context) ([]*Purchase, error) {
	query := `
		SELECT id, user_id, product_id, created_at, wallet_usdt, cost,  date
		FROM purchases`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*Purchase
	for rows.Next() {
		var purchase Purchase
		err := rows.Scan(
			&purchase.ID,
			&purchase.UserID,
			&purchase.ProductID,
			&purchase.Date,
			&purchase.WalletUSDT,
			&purchase.Cost,
		)
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

func (r *PurchaseRepository) Update(ctx context.Context, purchase *Purchase) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
		UPDATE purchases
		SET user_id = $1, product_id = $2, created_at = $3, wallet_usdt = $4, cost = $5
		WHERE id = $6`
	_, err = r.db.ExecContext(
		ctx,
		query,
		purchase.UserID,
		purchase.ProductID,
		purchase.Date,
		purchase.WalletUSDT,
		purchase.Cost,
		purchase.ID,
	)
	if err != nil {
		fmt.Errorf("failed to update purchase")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return err
}

func (r *PurchaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
		DELETE FROM purchases
		WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Errorf("failed to delete purchase")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

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

func (r *PurchaseRepository) CheckMeans(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*Purchase, error) {
	const queryUserMoney = `SELECT wallet_usdt FROM users WHERE user_id = $1`

	row := r.db.QueryRowContext(ctx, queryUserMoney, userID)

	var purchase *Purchase
	err := row.Scan(&purchase.WalletUSDT)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	const queryProductCost = `SELECT cost FROM products WHERE product_id = $1`

	row = r.db.QueryRowContext(ctx, queryProductCost, productID)

	err = row.Scan(&purchase.Cost)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return purchase, nil
}
