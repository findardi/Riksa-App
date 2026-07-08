package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/findardi/Wadi/server/internal/platform/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinio(cfg config.MinioConfig) (*MinioStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.SslMode,
	})

	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	s := &MinioStorage{
		client: client,
		bucket: cfg.BucketName,
	}

	if err := s.ensureBucket(context.Background()); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *MinioStorage) ensureBucket(ctx context.Context) error {
	exist, err := m.client.BucketExists(ctx, m.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}

	if exist {
		return nil
	}

	if err := m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("make bucket: %w", err)
	}

	return nil
}

func (m *MinioStorage) PresignedPut(ctx context.Context, key string, expiry time.Duration) (string, error) {
	u, err := m.client.PresignedPutObject(ctx, m.bucket, key, expiry)
	if err != nil {
		return "", fmt.Errorf("presign put: %w", err)
	}

	return u.String(), nil
}

func (m *MinioStorage) PresignedGet(ctx context.Context, key, filename string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	u, err := m.client.PresignedGetObject(ctx, m.bucket, key, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("presign get: %w", err)
	}

	return u.String(), nil
}

func (m *MinioStorage) Stat(ctx context.Context, key string) (size int64, contentType string, err error) {
	info, err := m.client.StatObject(ctx, m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return 0, "", fmt.Errorf("stat object: %w", err)
	}

	return info.Size, info.ContentType, nil
}

func (m *MinioStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})

	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}

	return obj, nil
}

func (m *MinioStorage) Delete(ctx context.Context, key string) error {
	if err := m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}
