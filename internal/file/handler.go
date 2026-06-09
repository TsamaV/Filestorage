package file

import (
	"context"
	"filestorage/pkg/middleware"
	pkg "filestorage/pkg/response"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type FileServiceInterface interface {
	Upload(ctx context.Context, userID int, fileName string, fileData io.Reader, size int64) error
	GetFiles(ctx context.Context, page, limit int) ([]File, error)
	Delete(ctx context.Context, id, userID int) error
}

type Handler struct {
	FileService FileServiceInterface
	Logger      *zap.Logger
}

func NewHandler(fileService FileServiceInterface, logger *zap.Logger) *Handler {
	return &Handler{
		FileService: fileService,
		Logger:      logger,
	}
}

func (h *Handler) Routes(mux *http.ServeMux, m *middleware.AuthMiddleware) {
	mux.Handle("POST /upload", m.IsAuthed(http.HandlerFunc(h.Upload)))
	mux.Handle("GET /files", m.IsAuthed(http.HandlerFunc(h.GetFiles)))
	mux.Handle("DELETE /files/{id}", m.IsAuthed(http.HandlerFunc(h.Delete)))
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.Logger.Error("parse form error", zap.Error(err))
		pkg.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.Logger.Error("get file error", zap.Error(err))
		pkg.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	h.Logger.Info("upload started", zap.String("filename", header.Filename), zap.Int64("size", header.Size))

	contentType := header.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" {
		h.Logger.Warn("invalid file format", zap.String("content_type", contentType))
		pkg.WriteError(w, http.StatusBadRequest, "only png and jpeg allowed")
		return
	}

	userID := int(r.Context().Value(middleware.ContextUserID).(float64))

	err = h.FileService.Upload(r.Context(), userID, header.Filename, file, header.Size)
	if err != nil {
		h.Logger.Error("upload failed", zap.Error(err))
		pkg.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "file uploaded"})

	h.Logger.Info("file uploaded", zap.String("filename", header.Filename), zap.Int("user_id", userID))
}

func (h *Handler) GetFiles(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("getting all files")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 1
	}

	files, err := h.FileService.GetFiles(r.Context(), page, limit)
	if err != nil {
		h.Logger.Error("get files failed", zap.Error(err))
		pkg.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	pkg.WriteJSON(w, http.StatusOK, files)

	h.Logger.Info("files returned", zap.Int("count", len(files)))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fileID, _ := strconv.Atoi(id)
	userID := int(r.Context().Value(middleware.ContextUserID).(float64))

	h.Logger.Info("deleting file", zap.Int("file_id", fileID), zap.Int("user_id", userID), zap.Int("User_id", userID))

	if err := h.FileService.Delete(r.Context(), fileID, userID); err != nil {
		h.Logger.Error("delete failed", zap.Error(err))
		if err.Error() == "forbidden" {
			pkg.WriteError(w, http.StatusForbidden, "not your file")
			return
		}
		pkg.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.Logger.Info("file deleted", zap.Int("file_id", fileID))
}
