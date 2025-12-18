package repositories

import (
	"context"
	"database/sql"

	"github.com/elghazx/go-vault/internal/core/domain"
)

type SQLiteFileRepository struct {
	db *sql.DB
}

func NewSQLiteFileRepository(db *sql.DB) *SQLiteFileRepository {
	return &SQLiteFileRepository{db: db}
}

func (r *SQLiteFileRepository) SaveMetadata(ctx context.Context, file *domain.File) error {
	query := `INSERT INTO files (uuid, filename, filepath, filesize, content_type, owner_id, is_onetime, expires_at, download_count, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		file.UUID, file.FileName, file.FilePath, file.FileSize,
		file.ContentType, file.OwnerID, file.IsOneTime,
		file.ExpiresAt, file.DownloadCount, file.CreatedAt)

	return err
}

func (r *SQLiteFileRepository) GetMetadata(ctx context.Context, uuid string) (*domain.File, error) {
	query := `SELECT uuid, filename, filepath, filesize, content_type, owner_id, is_onetime, expires_at, download_count, created_at 
			  FROM files WHERE uuid = ?`

	row := r.db.QueryRowContext(ctx, query, uuid)

	var file domain.File
	err := row.Scan(&file.UUID, &file.FileName, &file.FilePath, &file.FileSize,
		&file.ContentType, &file.OwnerID, &file.IsOneTime,
		&file.ExpiresAt, &file.DownloadCount, &file.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *SQLiteFileRepository) DeleteMetadata(ctx context.Context, uuid string) error {
	query := `DELETE FROM files WHERE uuid = ?`
	_, err := r.db.ExecContext(ctx, query, uuid)
	return err
}

func (r *SQLiteFileRepository) GetByOwnerID(ctx context.Context, ownerID int64) ([]*domain.File, error) {
	query := `SELECT uuid, filename, filepath, filesize, content_type, owner_id, is_onetime, expires_at, download_count, created_at 
			  FROM files WHERE owner_id = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*domain.File
	for rows.Next() {
		var file domain.File
		err := rows.Scan(&file.UUID, &file.FileName, &file.FilePath, &file.FileSize,
			&file.ContentType, &file.OwnerID, &file.IsOneTime,
			&file.ExpiresAt, &file.DownloadCount, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

func (r *SQLiteFileRepository) IncrementDownloadCount(ctx context.Context, uuid string) error {
	query := `UPDATE files SET download_count = download_count + 1 WHERE uuid = ?`
	_, err := r.db.ExecContext(ctx, query, uuid)
	return err
}
