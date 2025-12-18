package ports

import (
	"context"
	"io"

	"github.com/elghazx/go-vault/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type FileRepository interface {
	SaveMetadata(ctx context.Context, file *domain.File) error
	GetMetadata(ctx context.Context, uuid string) (*domain.File, error)
	DeleteMetadata(ctx context.Context, uuid string) error
	GetByOwnerID(ctx context.Context, ownerID int64) ([]*domain.File, error)
	IncrementDownloadCount(ctx context.Context, uuid string) error
}

type FileStorage interface {
	SaveFile(ctx context.Context, reader io.Reader, path string) error
	GetFile(ctx context.Context, path string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, path string) error
}
