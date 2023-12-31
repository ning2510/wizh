package minio

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

func CreateBucket(bucketName string) error {
	if len(bucketName) <= 0 {
		return errors.New("bucketName invalid")
	}

	if err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
		exist, exErr := minioClient.BucketExists(context.Background(), bucketName)
		if exist && exErr != nil {
			return nil
		} else {
			return exErr
		}
	}
	return nil
}

func RemoveBucket(ctx context.Context, bucketName string) error {
	return minioClient.RemoveBucket(ctx, bucketName)
}

func UploadFileByPath(bucketName, objectName, path, contentType string) (int64, error) {
	if len(bucketName) <= 0 || len(objectName) <= 0 {
		return -1, errors.New("invalid argument")
	}

	info, err := minioClient.FPutObject(context.Background(), bucketName, objectName, path, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return -1, err
	}

	return info.Size, nil
}

func UploadFileByIO(bucketName, objectName string, reader io.Reader, size int64, contentType string) (int64, error) {
	if len(bucketName) <= 0 || len(objectName) <= 0 {
		return -1, errors.New("invalid argument")
	}

	info, err := minioClient.PutObject(context.Background(), bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return -1, err
	}

	return info.Size, nil
}

func GetFileTemporaryURL(bucketName, objectName string) (string, error) {
	if len(bucketName) <= 0 || len(objectName) <= 0 {
		return "", errors.New("invalid argument")
	}

	expires := time.Second * time.Duration(ExpireTime)
	url, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expires, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
