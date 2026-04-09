package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vadim/filestorage/internal/model"
)

var ErrFileNotFound = errors.New("file not found")

type FileStorage struct {
	pool *pgxpool.Pool
}

func NewFileStorage(pool *pgxpool.Pool) *FileStorage {
	return &FileStorage{pool: pool}
}

func (s *FileStorage) Save(ctx context.Context, f model.File) (model.File, error) {
	err := s.pool.QueryRow(ctx,
		`INSERT INTO files (user_id, name, size, path)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		f.UserID, f.Name, f.Size, f.Path,
	).Scan(&f.ID, &f.CreatedAt)
	return f, err
}

func (s *FileStorage) GetByID(ctx context.Context, id int) (model.File, error) {
	var f model.File
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, size, path, created_at FROM files WHERE id = $1`, id,
	).Scan(&f.ID, &f.UserID, &f.Name, &f.Size, &f.Path, &f.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.File{}, ErrFileNotFound
	}
	return f, err
}

func (s *FileStorage) ListByUser(ctx context.Context, userID int, search string) ([]model.File, error) {
	query := `SELECT id, user_id, name, size, path, created_at FROM files WHERE user_id = $1`
	args := []any{userID}

	if search != "" {
		query += ` AND name ILIKE $2`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []model.File
	for rows.Next() {
		var f model.File
		if err := rows.Scan(&f.ID, &f.UserID, &f.Name, &f.Size, &f.Path, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, rows.Err()
}

func (s *FileStorage) Delete(ctx context.Context, id int) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM files WHERE id = $1`, id)
	return err
}

func (s *FileStorage) GetUsedSpace(ctx context.Context, userID int) (int64, error) {
	var total int64
	err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(size), 0) FROM files WHERE user_id = $1`, userID,
	).Scan(&total)
	return total, err
}
