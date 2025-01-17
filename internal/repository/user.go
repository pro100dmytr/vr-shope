package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) (*UserStorage, error) {
	return &UserStorage{db: db}, nil
}

func (r *UserStorage) Create(ctx context.Context, userRepo *User) error {
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

func (r *UserStorage) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
	SELECT 
	    id,
	    login,
	    name, 
	    last_name, 
	    phone_number, 
	    password, 
	    email, 
	    wallet_usdt,  
	FROM 
		users 
	WHERE 
	    id = $1
	    `

	user := &User{}
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
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserStorage) GetAll(ctx context.Context) ([]*User, error) {
	query := `SELECT id, login, name, last_name, phone_number, password, email, wallet_usdt FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
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

func (r *UserStorage) Update(ctx context.Context, userServ *User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
	UPDATE
    	users 
	SET
	    login = $1, 
	    name = $2,
	    last_name = $3,
	    phone_number = $4,
	    password = $5, 
	    email = $6, 
	    wallet_usdt = $7
	WHERE 
	    id = $8
	    `

	err = r.db.QueryRowContext(ctx, query,
		userServ.Login,
		userServ.Name,
		userServ.LastName,
		userServ.PhoneNumber,
		userServ.Password,
		userServ.Email,
		userServ.WalletUSDT,
		userServ.ID,
	).Scan(&userServ.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *UserStorage) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *UserStorage) ExistsByEmail(ctx context.Context, email string) (bool, error) {
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

func (r *UserStorage) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = $1`

	var exists uuid.UUID
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserStorage) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, login, name, last_name, phone_number, password, email, wallet_usdt, 
			  FROM users WHERE id = $1`

	user := &User{}
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

func (s *UserStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	user := &User{}

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

func (s *UserStorage) GetUsers(ctx context.Context, offset, limit int) ([]*User, error) {
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

	var users []*User
	for rows.Next() {
		user := &User{}
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
