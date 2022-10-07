package service

import (
	"context"

	"github.com/tarkovskynik/Golang-ninja-project/pkg/logger"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"
)

type FilesRepository interface {
	GetFiles(ctx context.Context, id int) ([]domain.File, error)
	StoreFileInfo(ctx context.Context, input domain.File) error
}

type S3FilesStorage interface {
	Upload(ctx context.Context, input domain.File) (string, error)
	GenerateFileURL(ctx context.Context, filename string) (string, error)
}

type FilesService struct {
	filesRepo FilesRepository
	storage   S3FilesStorage
}

func NewServiceFiles(filesRepo FilesRepository, storage S3FilesStorage) *FilesService {
	return &FilesService{
		filesRepo: filesRepo,
		storage:   storage,
	}
}

func (s *FilesService) Upload(ctx context.Context, input domain.File) (string, error) {
	return s.storage.Upload(ctx, input)
}

func (s *FilesService) GetFiles(ctx context.Context, id int) ([]domain.File, error) {
	files, err := s.filesRepo.GetFiles(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(files); i++ {
		files[i].URL, err = s.storage.GenerateFileURL(ctx, files[i].Name)
		if err != nil {
			logger.LogError("ServiceFiles.GetFiles", err)
		}
	}

	return files, nil
}

func (s *FilesService) StoreFileInfo(ctx context.Context, input domain.File) error {
	return s.filesRepo.StoreFileInfo(ctx, input)
}
