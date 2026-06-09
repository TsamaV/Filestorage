package file

import (
	"bytes"
	"context"
	"filestorage/pkg/middleware"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

type MockFileService struct{}

func (m *MockFileService) Upload(ctx context.Context, userID int, fileName string, fileData io.Reader, size int64) error {
	return nil
}

func (m *MockFileService) GetFiles(ctx context.Context, page, limit int) ([]File, error) {
	return []File{
		{
			ID: 1, UserID: 1, FileName: "test.jpeg", Size: 100,
		},
	}, nil
}

func (m *MockFileService) Delete(ctx context.Context, id, userID int) error {
	return nil
}

func TestGetFiles(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/files?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	mock := &MockFileService{}
	logger, _ := zap.NewDevelopment()

	handler := NewHandler(mock, logger)
	handler.GetFiles(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpload(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, _ := writer.CreateFormFile("file", "test.jpeg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodDelete, "/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextUserID, float64(1)))
	w := httptest.NewRecorder()

	mock := &MockFileService{}
	logger, _ := zap.NewDevelopment()

	handler := NewHandler(mock, logger)
	handler.Upload(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDelete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/files/{id}", nil)
	w := httptest.NewRecorder()
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextUserID, float64(1)))

	mock := &MockFileService{}
	logger, _ := zap.NewDevelopment()

	handler := NewHandler(mock, logger)
	handler.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected %d, got %d", http.StatusNoContent, w.Code)
	}
}