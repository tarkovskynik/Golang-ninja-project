package s3

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ConfigFileStorage struct {
	EndPoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
}

type UploadInput struct {
	Name        string
	FilePath    string
	ContentType string
}

type FileStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
}

func NewFileStorage(ctx context.Context, cfg ConfigFileStorage) (*FileStorage, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(cfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Make a new bucket
	err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{Region: cfg.Location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.BucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("We already own %s\n", cfg.BucketName)
		} else {
			return nil, err
		}
	} else {
		fmt.Printf("Successfully created %s\n", cfg.BucketName)
	}

	return &FileStorage{
		client:   minioClient,
		bucket:   cfg.BucketName,
		endpoint: cfg.EndPoint,
	}, nil
}

func (fs *FileStorage) Upload(ctx context.Context, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	// Upload the file with FPutObject
	_, err := fs.client.FPutObject(ctx, fs.bucket, input.Name, input.FilePath, opts)
	if err != nil {
		return "", err
	}

	return fs.GenerateFileURL(ctx, input.Name)
}

func (fs *FileStorage) GenerateFileURL(ctx context.Context, filename string) (string, error) {
	// DigitalOcean Spaces URL format.
	// return fmt.Sprintf("https://%s.%s/%s", fs.bucket, fs.endpoint, filename)

	// After set the policy of specific bucket to download in minio server, we can get the resource public url as follow:
	// var publicUrl = minioClient.protocol + '//' + minioClient.host + ':' + minioClient.port + '/' + minioBucket + '/' + obj.name

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+filename+"\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := fs.client.PresignedGetObject(ctx, fs.bucket, filename, time.Second*24*60*60, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
