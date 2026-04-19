package config

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinIOClient *minio.Client

func InitMinIO() error {
	cfg := GlobalConfig.MinIO
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("init minio err: %v", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return fmt.Errorf("check bucket failed: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("create bucket failed: %v", err)
		}
	}
	MinIOClient = client
	return nil
}

func UploadFile(ctx context.Context, objectName string, filePath string) (string, error) {
	cfg := GlobalConfig.MinIO
	info, err := MinIOClient.FPutObject(ctx, cfg.Bucket, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("upload file failed: %v", err)
	}
	url := fmt.Sprintf("%s/%s/%s", cfg.Endpoint, info.Bucket, info.Key)
	return url, nil
}

func UploadFileReader(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	cfg := GlobalConfig.MinIO
	_, err := MinIOClient.PutObject(ctx, cfg.Bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload file failed: %v", err)
	}
	url := fmt.Sprintf("%s/%s/%s", cfg.Endpoint, cfg.Bucket, objectName)
	return url, nil
}
