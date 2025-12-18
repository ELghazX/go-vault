package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"path/filepath"
	"time"

	"github.com/elghazx/go-vault/internal/core/domain"
	"github.com/elghazx/go-vault/internal/core/ports"
)

type FileService struct {
	repo    ports.FileRepository
	storage ports.FileStorage
}

func NewFileService(repo ports.FileRepository, storage ports.FileStorage) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

func (s *FileService) GetFileMetadata(ctx context.Context, uuid string) (*domain.File, error) {
	return s.repo.GetMetadata(ctx, uuid)
}

func (s *FileService) UploadFile(ctx context.Context, reader io.Reader, filename, contentType string, ownerID int64, isOneTime bool) (*domain.File, error) {
	uuid := generateUUID()
	ext := filepath.Ext(filename)
	filePath := uuid + ext

	if err := s.storage.SaveFile(ctx, reader, filePath); err != nil {
		return nil, err
	}

	file := &domain.File{
		UUID:        uuid,
		FileName:    filename,
		FilePath:    filePath,
		ContentType: contentType,
		OwnerID:     ownerID,
		IsOneTime:   isOneTime,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		CreatedAt:   time.Now(),
	}

	if err := s.repo.SaveMetadata(ctx, file); err != nil {
		s.storage.DeleteFile(ctx, filePath)
		return nil, err
	}

	return file, nil
}

func (s *FileService) DownloadFile(ctx context.Context, uuid string) (io.ReadCloser, *domain.File, error) {
	file, err := s.repo.GetMetadata(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}

	if file.IsExpired() {
		s.BurnFile(ctx, uuid)
		return nil, nil, errors.New("file expired")
	}

	reader, err := s.storage.GetFile(ctx, file.FilePath)
	if err != nil {
		return nil, nil, err
	}

	// Increment download count first
	s.repo.IncrementDownloadCount(ctx, uuid)

	// Check if should burn AFTER incrementing
	if file.IsOneTime {
		go s.BurnFile(context.Background(), uuid)
	}

	return reader, file, nil
}

func (s *FileService) GetUserFiles(ctx context.Context, ownerID int64) ([]*domain.File, error) {
	return s.repo.GetByOwnerID(ctx, ownerID)
}

func (s *FileService) BurnFile(ctx context.Context, uuid string) error {
	file, err := s.repo.GetMetadata(ctx, uuid)
	if err != nil {
		return err
	}

	s.storage.DeleteFile(ctx, file.FilePath)
	return s.repo.DeleteMetadata(ctx, uuid)
}

func generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
