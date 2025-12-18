package repositories

import (
	"context"
	"database/sql"

	"github.com/elghazx/go-vault/internal/core/domain"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Username, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *SQLiteUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = ?`
	row := r.db.QueryRowContext(ctx, query, username)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *SQLiteUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
