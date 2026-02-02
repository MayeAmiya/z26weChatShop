package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var bucketName string
var publicURL string

// MinIOConfig MinIO configuration
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
	PublicURL       string
}

// InitMinIO initializes MinIO client
func InitMinIO(cfg MinIOConfig) error {
	var err error

	minioClient, err = minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create minio client: %w", err)
	}

	bucketName = cfg.BucketName
	publicURL = cfg.PublicURL

	// Create bucket if not exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		// Set bucket policy to public read
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}]
		}`, bucketName)

		err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			fmt.Printf("Warning: failed to set bucket policy: %v\n", err)
		}
	}

	fmt.Printf("MinIO initialized: endpoint=%s, bucket=%s\n", cfg.Endpoint, bucketName)
	return nil
}

// UploadFile uploads a file to MinIO
func UploadFile(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, error) {
	if minioClient == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	// Generate unique filename
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("%s/%d%s", time.Now().Format("2006/01/02"), time.Now().UnixNano(), ext)

	// Upload file
	_, err := minioClient.PutObject(ctx, bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return public URL
	url := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(publicURL, "/"), bucketName, objectName)
	return url, nil
}

// DeleteFile deletes a file from MinIO
func DeleteFile(ctx context.Context, objectURL string) error {
	if minioClient == nil {
		return fmt.Errorf("minio client not initialized")
	}

	// Extract object name from URL
	prefix := fmt.Sprintf("%s/%s/", strings.TrimSuffix(publicURL, "/"), bucketName)
	objectName := strings.TrimPrefix(objectURL, prefix)

	if objectName == objectURL {
		return fmt.Errorf("invalid object URL")
	}

	err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetDefaultMinIOConfig returns default MinIO config from environment
func GetDefaultMinIOConfig() MinIOConfig {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "z26b"
	}

	minioPublicURL := os.Getenv("MINIO_PUBLIC_URL")
	if minioPublicURL == "" {
		minioPublicURL = "http://localhost:9000"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	return MinIOConfig{
		Endpoint:        endpoint,
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		BucketName:      bucket,
		UseSSL:          useSSL,
		PublicURL:       minioPublicURL,
	}
}

// IsMinIOInitialized returns whether MinIO is initialized
func IsMinIOInitialized() bool {
	return minioClient != nil
}
