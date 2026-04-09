package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/vadim/filestorage/internal/model"
	"github.com/vadim/filestorage/internal/storage"
)

var ErrQuotaExceeded = errors.New("storage quota exceeded")

type FileService struct {
	store      *storage.FileStorage
	uploadDir  string
	maxQuotaMB int
	log        *slog.Logger
}

func NewFileService(store *storage.FileStorage, uploadDir string, maxQuotaMB int, log *slog.Logger) *FileService {
	return &FileService{
		store:      store,
		uploadDir:  uploadDir,
		maxQuotaMB: maxQuotaMB,
		log:        log,
	}
}

func (s *FileService) Upload(ctx context.Context, userID int, filename string, fileData io.Reader, size int64) (model.File, error) {
	used, err := s.store.GetUsedSpace(ctx, userID)
	if err != nil {
		return model.File{}, err
	}

	maxBytes := int64(s.maxQuotaMB) * 1024 * 1024
	if used+size > maxBytes {
		return model.File{}, ErrQuotaExceeded
	}

	// TODO: добавить проверку расширения файла (запретить .exe и т.д.)

	userDir := filepath.Join(s.uploadDir, strconv.Itoa(userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return model.File{}, err
	}

	// если файл с таким именем уже есть, перезапишется. Пока ок.
	dst, err := os.Create(filepath.Join(userDir, filename))
	if err != nil {
		return model.File{}, err
	}
	defer dst.Close()

	written, err := io.Copy(dst, fileData)
	if err != nil {
		return model.File{}, err
	}

	f := model.File{
		UserID: userID,
		Name:   filename,
		Size:   written,
		Path:   filepath.Join(userDir, filename),
	}

	return s.store.Save(ctx, f)
}

func (s *FileService) List(ctx context.Context, userID int, search string) ([]model.File, error) {
	return s.store.ListByUser(ctx, userID, search)
}

func (s *FileService) Get(ctx context.Context, fileID int, userID int) (model.File, error) {
	f, err := s.store.GetByID(ctx, fileID)
	if err != nil {
		return model.File{}, err
	}
	if f.UserID != userID {
		return model.File{}, storage.ErrFileNotFound
	}
	return f, nil
}

func (s *FileService) Delete(ctx context.Context, fileID int, userID int) error {
	f, err := s.Get(ctx, fileID, userID)
	if err != nil {
		return err
	}

	if err := os.Remove(f.Path); err != nil && !os.IsNotExist(err) {
		s.log.Warn("failed to remove file from disk", "path", f.Path, "err", err)
	}

	return s.store.Delete(ctx, fileID)
}
