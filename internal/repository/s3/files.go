package s3

import (
	"context"
	"fmt"

	"github.com/tarkovskynik/Golang-ninja-project/pkg/s3"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"
)

type S3FilesStorage struct {
	storage *s3.FileStorage
}

func NewS3FilesStorage(storage *s3.FileStorage) *S3FilesStorage {
	return &S3FilesStorage{storage}
}

func (r *S3FilesStorage) Upload(ctx context.Context, input domain.File) (string, error) {
	return r.storage.Upload(ctx,
		s3.UploadInput{Name: input.Name,
			FilePath:    fmt.Sprintf("temp.file.%d-%s", input.UserID, input.Name),
			ContentType: input.ContentType,
		})
}

func (r *S3FilesStorage) GenerateFileURL(ctx context.Context, filename string) (string, error) {
	return r.storage.GenerateFileURL(ctx, filename)
}
