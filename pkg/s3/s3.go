package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type ItfS3 interface {
	UploadFile(session *session.Session, file *multipart.FileHeader, fileName string) (string, error)
	PresignUrl(client *s3.S3, fileName string) (string, error)
}

type s3Client struct {
	client     *s3.S3
	bucketName string
}

func New() (ItfS3, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}

	return &s3Client{client: s3.New(sess), bucketName: os.Getenv("AWS_BUCKET_NAME")}, nil
}

func (s *s3Client) UploadFile(session *session.Session, file *multipart.FileHeader, fileName string) (string, error) {
	uploader := s3manager.NewUploader(session)

	uniqueFileName, err := generateUniqueFileName(fileName)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			fmt.Println("Failed to close file")
		}
	}(src)

	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(uniqueFileName),
		Body:   src,
	})

	if err != nil {
		return "", err
	}

	return uploadOutput.Location, nil
}

func (s *s3Client) PresignUrl(client *s3.S3, fileName string) (string, error) {
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}

func (s *s3Client) DeleteFile(client *s3.S3, fileName string) error {
	_, err := client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})

	return err
}

func newSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}

func generateUniqueFileName(fileName string) (string, error) {
	uniqueFileName := fmt.Sprintf("%s-%s", strings.ReplaceAll(time.Now().String(), " ", ""), fileName)
	return uniqueFileName, nil
}
