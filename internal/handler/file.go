package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vadim/filestorage/internal/auth"
	"github.com/vadim/filestorage/internal/service"
	"github.com/vadim/filestorage/internal/storage"
)

type FileHandler struct {
	svc *service.FileService
	log *slog.Logger
}

func NewFileHandler(svc *service.FileService, log *slog.Logger) *FileHandler {
	return &FileHandler{svc: svc, log: log}
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		sendError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	saved, err := h.svc.Upload(r.Context(), userID, header.Filename, file, header.Size)
	if err != nil {
		if err == service.ErrQuotaExceeded {
			sendError(w, http.StatusRequestEntityTooLarge, "quota exceeded")
			return
		}
		h.log.Error("upload failed", "err", err)
		sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	sendJSON(w, http.StatusCreated, saved)
}

func (h *FileHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	files, err := h.svc.List(r.Context(), userID, r.URL.Query().Get("search"))
	if err != nil {
		h.log.Error("list files", "err", err)
		sendError(w, http.StatusInternalServerError, "internal error")
		return
	}
	sendJSON(w, http.StatusOK, files)
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	f, err := h.svc.Get(r.Context(), id, userID)
	if err != nil {
		if err == storage.ErrFileNotFound {
			sendError(w, http.StatusNotFound, "not found")
			return
		}
		h.log.Error("download failed", "err", err)
		sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Header().Set("Content-Disposition", `attachment; filename="`+f.Name+`"`)
	http.ServeFile(w, r, f.Path)
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.svc.Delete(r.Context(), id, userID)
	if err == storage.ErrFileNotFound {
		sendError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		h.log.Error("delete failed", "err", err)
		sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
