package config

import (
	"os"
	"strconv"
)

type minio struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

func loadMinioConfig() minio {
	useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))

	return minio{
		AccessKeyID:     os.Getenv("MINIO_ACCESS_KEY_ID"),
		Endpoint:        os.Getenv("MINIO_ENDPOINT"),
		SecretAccessKey: os.Getenv("MINIO_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("MINIO_BUCKET_NAME"),
		UseSSL:          useSSL,
	}
}
