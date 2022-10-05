package service

import (
	"context"

	"github.com/tarkovskynik/Golang-ninja-project/pkg/logger"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/s3"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"
)

type FilesRepository interface {
	GetFiles(ctx context.Context, id int) ([]domain.File, error)
	StoreFileInfo(ctx context.Context, input domain.File) error
}

type FilesServece struct {
	filesRepo FilesRepository
	storage   *s3.FileStorage
}

func NewServeceFiles(filesRepo FilesRepository, storage *s3.FileStorage) *FilesServece {
	return &FilesServece{
		filesRepo: filesRepo,
		storage:   storage,
	}
}

func (s *FilesServece) Upload(ctx context.Context, input domain.File) (string, error) {
	return s.storage.Upload(ctx,
		s3.UploadInput{Name: input.Name,
			FilePath:    "/" + input.Name,
			ContentType: input.ContentType,
		})
}

func (s *FilesServece) GetFiles(ctx context.Context, id int) ([]domain.File, error) {
	files, err := s.filesRepo.GetFiles(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(files); i++ {
		files[i].URL, err = s.storage.GenerateFileURL(ctx, files[i].Name)
		if err != nil {
			logger.LogError("ServeceFiles.GetFiles", err)
		}
	}

	return files, nil
}

func (s *FilesServece) StoreFileInfo(ctx context.Context, input domain.File) error {
	return s.filesRepo.StoreFileInfo(ctx, input)
}
