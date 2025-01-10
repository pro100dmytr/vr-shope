package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"vr-shope/internal/config"
	"vr-shope/internal/models/repositories"
	"vr-shope/internal/storage/postgresql"
)

type ProductRepository struct {
	db *sql.DB
}

func (s *ProductRepository) Close() error {
	return postgresql.CloseConn(s.db)
}

func NewProductStorage(cfg *config.Config) (*ProductRepository, error) {
	db, err := postgresql.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &ProductRepository{db: db}, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *repositories.Product) error {
	query := `
		INSERT INTO products (id, name, cost, quantity_stock, guarantees, country)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Cost,
		product.QuantityStock,
		product.Guarantees,
		product.Country,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Get(ctx context.Context, id uuid.UUID) (*repositories.Product, error) {
	query := `
		SELECT id, name, cost, quantity_stock, guarantees, country, likes
		FROM products
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var product repositories.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Cost,
		&product.QuantityStock,
		&product.Guarantees,
		&product.Country,
		&product.Like,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]*repositories.Product, error) {
	query := `
		SELECT id, name, cost, quantity_stock, guarantees, country, likes
		FROM products
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*repositories.Product
	for rows.Next() {
		var product repositories.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Cost,
			&product.QuantityStock,
			&product.Guarantees,
			&product.Country,
			&product.Like,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *repositories.Product) (*repositories.Product, error) {
	query := `
		UPDATE products
		SET name = $2, cost = $3, quantity_stock = $4, guarantees = $5, country = $6, likes = $7
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Cost,
		product.QuantityStock,
		product.Guarantees,
		product.Country,
		product.Like,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no rows were updated")
	}

	return product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted")
	}

	return nil
}

func (s *ProductRepository) AddLike(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const checkQuery = `SELECT 1 FROM products WHERE id = $1`
	var exists bool
	err = s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to check if products exists: %w", err)
	}

	const query = `UPDATE products SET likes = likes + 1 WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ProductRepository) RemoveLike(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const checkQuery = `SELECT 1 FROM products WHERE id = $1`
	var exists bool
	err = s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to check if products exists: %w", err)
	}

	const query = `UPDATE products SET likes = likes - 1 WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ProductRepository) GetForName(ctx context.Context, name string) ([]*repositories.Product, error) {
	const query = `SELECT * FROM products WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*repositories.Product
	for rows.Next() {
		product := &repositories.Product{}
		if err := rows.Scan(
			product.ID,
			product.Name,
			product.Cost,
			product.QuantityStock,
			product.Guarantees,
			product.Country,
			product.Like,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (s *ProductRepository) GetProducts(ctx context.Context, offset, limit int) ([]*repositories.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `
        SELECT id, name, cost, quantityStock, guarantees, country, likes FROM users OFFSET $1 LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*repositories.Product
	for rows.Next() {
		product := &repositories.Product{}
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Cost,
			&product.QuantityStock,
			&product.Guarantees,
			&product.Country,
			&product.Like,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return products, rows.Err()
}
