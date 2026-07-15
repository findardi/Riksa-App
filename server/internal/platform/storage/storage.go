package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	PresignedPut(ctx context.Context, key string, expiry time.Duration) (string, error)
	PresignedGet(ctx context.Context, key, filename string, expiry time.Duration) (string, error)
	Stat(ctx context.Context, key string) (size int64, contentType string, err error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
}
