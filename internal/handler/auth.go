package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/vadim/filestorage/internal/model"
	"github.com/vadim/filestorage/internal/service"
	"github.com/vadim/filestorage/internal/storage"
)

type AuthHandler struct {
	svc *service.AuthService
	log *slog.Logger
}

func NewAuthHandler(svc *service.AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, log: log}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if req.Username == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "username and password required")
		return
	}

	if len(req.Password) < 6 {
		sendError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	user, err := h.svc.Register(r.Context(), req)
	if err != nil {
		if err == storage.ErrUserExists {
			sendError(w, http.StatusConflict, "username taken")
			return
		}
		h.log.Error("register failed", "err", err)
		sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	sendJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid body")
		return
	}

	token, err := h.svc.Login(r.Context(), req)
	if err != nil {
		if err == service.ErrWrongCredentials {
			sendError(w, http.StatusUnauthorized, "wrong credentials")
			return
		}
		h.log.Error("login failed", "err", err)
		sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	sendJSON(w, http.StatusOK, model.TokenResponse{Token: token})
}
