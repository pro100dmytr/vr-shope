package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"vr-shope/internal/config"
	"vr-shope/internal/models/repositories"
	"vr-shope/internal/models/services"
	"vr-shope/internal/storage/postgresql"
)

type UserRepository struct {
	db *sql.DB
}

func (s *UserRepository) Close() error {
	return postgresql.CloseConn(s.db)
}

func NewUserStorage(cfg *config.Config) (*UserRepository, error) {
	db, err := postgresql.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Create(ctx context.Context, userRepo *repositories.User) error {
	query := `INSERT INTO users (id, login, name, last_name, phone_number, password, email, salt) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		userRepo.ID,
		userRepo.Login,
		userRepo.Name,
		userRepo.LastName,
		userRepo.PhoneNumber,
		userRepo.Password,
		userRepo.Email,
		userRepo.Salt,
	).Scan(&userRepo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*services.User, error) {
	query := `SELECT id, login, name, last_name, phone_number, password, email, wallet_usdt,  
			  FROM users WHERE id = $1`

	user := &services.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.LastName,
		&user.PhoneNumber,
		&user.Password,
		&user.Email,
		&user.WalletUSDT,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*services.User, error) {
	query := `SELECT id, login, name, last_name, phone_number, password, email, wallet_usdt FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*services.User
	for rows.Next() {
		user := &services.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Name,
			&user.LastName,
			&user.PhoneNumber,
			&user.Password,
			&user.Email,
			&user.WalletUSDT,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, userServ *services.User) error {
	query := `UPDATE users SET login = $2, name = $3, last_name = $4, phone_number = $5, password = $6, email = $7, wallet_usdt = $8 WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query,
		userServ.ID,
		userServ.Login,
		userServ.Name,
		userServ.LastName,
		userServ.PhoneNumber,
		userServ.Password,
		userServ.Email,
		userServ.WalletUSDT,
	).Scan(&userServ.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = $1`

	var exists int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = $1`

	var exists int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*services.User, error) {
	query := `SELECT id, login, name, last_name, phone_number, password, email, wallet_usdt, 
			  FROM users WHERE id = $1`

	user := &services.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.LastName,
		&user.PhoneNumber,
		&user.Password,
		&user.Email,
		&user.WalletUSDT,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserRepository) GetUserByLogin(ctx context.Context, login string) (*repositories.User, error) {
	user := &repositories.User{}

	const query = `SELECT id, login, password_hash, salt FROM users WHERE login = $1`
	err := s.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (s *UserRepository) GetUsers(ctx context.Context, offset, limit int) ([]*repositories.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `
        SELECT id, login, email FROM users OFFSET $1 LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*repositories.User
	for rows.Next() {
		user := &repositories.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return users, rows.Err()
}
