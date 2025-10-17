package minio

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/enum"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"
)

type minioRepository struct {
	minioClinet      *minio.Client
	rootConfig       *config.Root
	minioPathNameMap *map[enum.MinioPathName]string
	BucketName       string
	Clogger          *helpers.CustomLogger
}

func NewMinioRepository(
	minioClinet *minio.Client,
	rootConfig *config.Root,
	BucketName string,
	Clogger *helpers.CustomLogger,
) MinioRepository {
	minioPathNameMap := enum.MinioPathNameMap
	return &minioRepository{
		minioClinet:      minioClinet,
		rootConfig:       rootConfig,
		minioPathNameMap: &minioPathNameMap,
		BucketName:       BucketName,
		Clogger:          Clogger,
	}
}

type MinioRepository interface {
	PutObject(ctx context.Context, file *multipart.FileHeader, path string, fileName *string) (*minio.UploadInfo, *string, error)
	GenerateMinioPresignedURL(ctx context.Context, fileName *string, expires time.Duration) (string, error)
}

func (m *minioRepository) PutObject(ctx context.Context, file *multipart.FileHeader, path string, fileName *string) (*minio.UploadInfo, *string, error) {
	exists, errBucketExists := m.minioClinet.BucketExists(ctx, m.BucketName)
	if errBucketExists != nil {
		m.Clogger.ErrorLogger(ctx, "PutObject.minioClinet.BucketExists", errBucketExists)
		return nil, nil, errBucketExists
	}
	if !exists {
		err := m.minioClinet.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			m.Clogger.ErrorLogger(ctx, "PutObject.minioClinet.MakeBucket", errBucketExists)
			return nil, nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	src, err := file.Open()
	if err != nil {
		m.Clogger.ErrorLogger(ctx, "PutObject.file.Open", err)
		return nil, nil, err
	}
	defer src.Close()

	var objectName = fmt.Sprintf("%s/%s%s", path, *fileName, filepath.Ext(file.Filename))
	info, err := m.minioClinet.PutObject(ctx, m.BucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		m.Clogger.ErrorLogger(ctx, "PutObject.minioClinet.PutObject", err)
		return nil, nil, err
	}
	return &info, &objectName, nil

}
func (m *minioRepository) GenerateMinioPresignedURL(ctx context.Context, fileName *string, expires time.Duration) (string, error) {
	reqParams := url.Values{}
	url, err := m.minioClinet.PresignedGetObject(ctx, m.BucketName, *fileName, expires, reqParams)
	if err != nil {
		m.Clogger.ErrorLogger(ctx, "GenerateMinioPresignedURL.minioClinet.PresignedGetObject", err)
		return "", err
	}
	return url.String(), nil
}
