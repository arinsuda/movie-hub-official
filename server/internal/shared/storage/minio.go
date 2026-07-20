package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/arinsuda/movie-hub/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	presignDuration = 24 * time.Hour
)

type MinIOClient struct {
	client     *minio.Client
	bucketName string
	useSSL     bool
}

func NewMinIOClient(cfg config.MinIOConfig) (*MinIOClient, error) {

	endpoint := strings.TrimPrefix(cfg.Endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")

	if strings.HasPrefix(cfg.Endpoint, "https://") {
		cfg.UseSSL = true
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: create client: %w", err)
	}

	mc := &MinIOClient{
		client:     client,
		bucketName: cfg.BucketName,
		useSSL:     cfg.UseSSL,
	}

	if err := mc.ensureBucket(context.Background()); err != nil {
		return nil, err
	}

	return mc, nil
}

func (m *MinIOClient) ensureBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return fmt.Errorf("minio: check bucket: %w", err)
	}

	if !exists {
		if err := m.client.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("minio: make bucket: %w", err)
		}
	}

	return nil
}

func (m *MinIOClient) UploadAvatar(ctx context.Context, userID uint, file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	objectKey := fmt.Sprintf("avatars/%d/%s%s", userID, uuid.NewString(), ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := m.client.PutObject(ctx, m.bucketName, objectKey, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio: upload avatar: %w", err)
	}

	return objectKey, nil
}

func (m *MinIOClient) UploadBuffer(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio: upload buffer: %w", err)
	}
	return objectKey, nil
}

func (m *MinIOClient) PresignURL(ctx context.Context, objectKey string) (string, error) {
	if objectKey == "" {
		return "", nil
	}

	url, err := m.client.PresignedGetObject(ctx, m.bucketName, objectKey, presignDuration, nil)
	if err != nil {
		return "", fmt.Errorf("minio: presign: %w", err)
	}

	return url.String(), nil
}

func (m *MinIOClient) DeleteObject(ctx context.Context, objectKey string) error {
	if objectKey == "" {
		return nil
	}

	err := m.client.RemoveObject(ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio: delete object: %w", err)
	}

	return nil
}
