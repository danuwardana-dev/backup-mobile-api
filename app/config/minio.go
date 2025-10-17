package config

import (
	"github.com/labstack/gommon/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

type Minio struct {
	Endpoint        string `envconfig:"MINIO_ENDPOINT"`
	AccessKeyID     string `envconfig:"MINIO_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"MINIO_SECRET_ACCESS_KEY"`
	Token           string `envconfig:"MINIO_TOKEN"`

	Bucket string `envconfig:"MINIO_BUCKET"`

	MinioPresignedDuration time.Duration `envconfig:"MINIO_PRESIGNED_DURATION" default:"0h5m0s"`
}

func (m *Minio) MinioClientSet() (*minio.Client, error) {
	minioClinet, err := minio.New(m.Endpoint, &minio.Options{
		Creds:              credentials.NewStaticV4(m.AccessKeyID, m.SecretAccessKey, m.Token),
		Secure:             false,
		Transport:          nil,
		Trace:              nil,
		Region:             "",
		BucketLookup:       0,
		CustomRegionViaURL: nil,
		BucketLookupViaURL: nil,
		TrailingHeaders:    false,
		CustomMD5:          nil,
		CustomSHA256:       nil,
		MaxRetries:         0,
	})

	if err != nil {
		log.Errorf("MinioClientSet(): %v", err)
		return nil, err
	}
	return minioClinet, nil
}
