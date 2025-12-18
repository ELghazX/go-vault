package domain

import "time"

type File struct {
	UUID          string
	FileName      string
	FilePath      string
	FileSize      int64
	ContentType   string
	OwnerID       int64
	IsOneTime     bool
	ExpiresAt     time.Time
	DownloadCount int
	CreatedAt     time.Time
}

func (f *File) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

func (f *File) ShouldBeBurned() bool {
	return f.IsOneTime && f.DownloadCount >= 1
}
