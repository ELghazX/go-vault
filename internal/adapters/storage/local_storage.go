package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	basePath string
}

func NewLocalFileStorage(basePath string) *LocalFileStorage {
	os.MkdirAll(basePath, 0o755)
	return &LocalFileStorage{basePath: basePath}
}

func (s *LocalFileStorage) SaveFile(ctx context.Context, reader io.Reader, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (s *LocalFileStorage) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.Open(fullPath)
}

func (s *LocalFileStorage) DeleteFile(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}
