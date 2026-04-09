package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vadim/filestorage/internal/auth"
)

func NewRouter(authH *AuthHandler, fileH *FileHandler, jwtSecret string, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(requestLogger(log))

	r.Post("/api/register", authH.Register)
	r.Post("/api/login", authH.Login)

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(jwtSecret))

		r.Post("/api/files", fileH.Upload)
		r.Get("/api/files", fileH.List)
		r.Get("/api/files/{id}/download", fileH.Download)
		r.Delete("/api/files/{id}", fileH.Delete)
	})

	return r
}

func requestLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start).String(),
			)
		})
	}
}

type errorBody struct {
	Error string `json:"error"`
}

func sendJSON(w http.ResponseWriter, status int, v any) {
	buf, err := json.Marshal(v)
	if err != nil {
		http.Error(w, `{"error":"marshal failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, errorBody{Error: msg})
}
