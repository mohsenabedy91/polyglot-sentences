package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type Client struct {
	log    logger.Logger
	conf   config.Minio
	client *minio.Client
}

func NewMinioClient(ctx context.Context, log logger.Logger, conf config.Minio) (*Client, error) {
	endpoint := fmt.Sprintf("%s:%s", conf.Endpoint, conf.Port)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.ID, conf.Secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	if err = minioClient.MakeBucket(ctx, conf.BucketName, minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, conf.BucketName)
		if errBucketExists == nil && exists {
			log.Warn(logger.Minio, logger.MinioCreateBucket, fmt.Sprintf("We already own %s", conf.BucketName), nil)
		} else {
			log.Error(logger.Minio, logger.Startup, err.Error(), nil)
			return nil, err
		}

	}

	log.Info(logger.Minio, logger.Startup, fmt.Sprintf("Successfully created %s", conf.BucketName), nil)

	return &Client{
		log:    log,
		conf:   conf,
		client: minioClient,
	}, nil
}

func (r *Client) UploadFile(ctx context.Context, objectName, filePath, contentType string) (string, error) {
	_, err := r.client.FPutObject(ctx, r.conf.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		r.log.Error(logger.Minio, logger.MinioUpload, err.Error(), map[logger.ExtraKey]interface{}{
			"objectName":  objectName,
			"filePath":    filePath,
			"contentType": contentType,
		})
		return "", err
	}

	return r.client.EndpointURL().String() + "/" + r.conf.BucketName + "/" + objectName, nil
}
