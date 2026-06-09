package file

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type FileRepo interface {
	Create(ctx context.Context, userID int, url, fileName string, size int) (int, error)
	GetFiles(ctx context.Context, page, limit int) ([]File, error)
	GetByID(ctx context.Context, id int) (File, error)
	Delete(ctx context.Context, id int) error
}

type FileService struct {
	fileRepo FileRepo
	s3       *minio.Client
	bucket   string
}

func NewFileService(file FileRepo, s3 *minio.Client, bucket string) *FileService {
	return &FileService{
		fileRepo: file,
		s3:       s3,
		bucket:   bucket,
	}
}

func (service *FileService) Upload(ctx context.Context, userID int, fileName string, fileData io.Reader, size int64) error {
	_, err := service.s3.PutObject(ctx, service.bucket, fileName, fileData, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:9000/%s/%s", service.bucket, fileName)

	_, err = service.fileRepo.Create(ctx, userID, url, fileName, int(size))

	return err
}

func (service *FileService) GetFiles(ctx context.Context, page, limit int) ([]File, error) {
	return service.fileRepo.GetFiles(ctx, page, limit)
}

func (service *FileService) Delete(ctx context.Context, id, userID int) error {
	f, err := service.fileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if f.UserID != userID {
		return errors.New("forbidden")
	}

	if err := service.s3.RemoveObject(ctx, service.bucket, f.FileName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	
	return service.fileRepo.Delete(ctx, id)
}