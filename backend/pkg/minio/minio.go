package minio

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/setting"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

func ReadFrameAndPutJpeg(inFileName string, outFileName string, frameNum int) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	_, err = MinioClient.PutObject(context.Background(), "video", "preview/"+outFileName, buf, -1, minio.PutObjectOptions{ContentType: "application/image"})
	if err != nil {
		panic(err)
	}
}

func PutVideoAndPreview(inputFilePath string, outFileName string) error {

	buf := bytes.NewBuffer(nil)
	img_buf := bytes.NewBuffer(nil)
	s1 := ffmpeg.Input(inputFilePath, nil)

	err := s1.Output("pipe:", ffmpeg.KwArgs{"format": "mpegts"}).WithOutput(buf, os.Stdout).Run()
	if err != nil {
		panic(err)
	}
	err = s1.Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 5)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithInput(buf).
		WithOutput(img_buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}

	_, err = MinioClient.PutObject(context.Background(), "video", outFileName+".ts", buf, -1, minio.PutObjectOptions{ContentType: "application/video"})
	if err != nil {
		panic(err)
	}
	_, err = MinioClient.PutObject(context.Background(), "preview", outFileName+".jpg", img_buf, -1, minio.PutObjectOptions{ContentType: "application/image"})
	if err != nil {
		panic(err)
	}
	return nil
	// ReadFrameAsJpeg(buf, "test01.jpg", 1)
}
