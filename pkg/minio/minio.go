package minio

import (
	"context"
	"mime/multipart"
	log "qiu/blog/pkg/logging"
	"qiu/blog/pkg/setting"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Ctx context.Context
var MinioClient *minio.Client

func Setup() {
	Ctx = context.Background()
	endpoint := setting.MinioSetting.EndPoint
	accessKeyID := setting.MinioSetting.AccessKeyID
	secretAccessKey := setting.MinioSetting.SecretAccessKey
	useSSL := false
	var err error
	// Initialize minio client object.
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic("Minio Setup Fail!")
		// log.Logger.Fatal(err)
	}
}

func PutImage(bucketName string, fileName string, mf *multipart.FileHeader) error {
	src, err := mf.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = MinioClient.PutObject(context.Background(), bucketName, fileName, src, -1, minio.PutObjectOptions{ContentType: "application/image"})
	if err != nil {
		log.Logger.Error(err)
	}
	return nil
}
