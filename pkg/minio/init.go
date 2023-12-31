package minio

import (
	"wizh/pkg/viper"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient               *minio.Client
	minioConfig               = viper.InitConf("minio")
	MinioEndPoint             = minioConfig.Viper.GetString("minio.Endpoint")
	MinioAccessKeyId          = minioConfig.Viper.GetString("minio.AccessKeyId")
	MinioSecretAccessKey      = minioConfig.Viper.GetString("minio.SecretAccessKey")
	UseSSL                    = minioConfig.Viper.GetBool("minio.UseSSL")
	VideoBucketName           = minioConfig.Viper.GetString("minio.VideoBucketName")
	CoverBucketName           = minioConfig.Viper.GetString("minio.CoverBucketName")
	AvatarBucketName          = minioConfig.Viper.GetString("minio.AvatarBucketName")
	BackgroungImageBucketName = minioConfig.Viper.GetString("minio.BackgroungImageBucketName")
	ExpireTime                = minioConfig.Viper.GetUint32("minio.ExpireTime")
)

func init() {
	var err error
	minioClient, err = minio.New(MinioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioAccessKeyId, MinioSecretAccessKey, ""),
		Secure: UseSSL,
	})

	if err != nil {
		panic(err)
	}

	if err = CreateBucket(VideoBucketName); err != nil {
		panic(err)
	}

	if err = CreateBucket(CoverBucketName); err != nil {
		panic(err)
	}

	if err = CreateBucket(AvatarBucketName); err != nil {
		panic(err)
	}

	if err = CreateBucket(BackgroungImageBucketName); err != nil {
		panic(err)
	}
}
