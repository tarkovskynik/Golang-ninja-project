package psql

import (
	"context"
	"database/sql"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"
)

type Files struct {
	db *sql.DB
}

func NewFiles(db *sql.DB) *Files {
	return &Files{db}
}

func (r *Files) StoreFileInfo(ctx context.Context, input domain.File) error {
	query := "INSERT INTO files(user_id, name, size, type, content_type, url, upload_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.ExecContext(ctx, query, input.UserID, input.Name, input.Size, input.Type, input.ContentType, input.URL, input.UploadAt)
	return err
}

func (r *Files) GetFiles(ctx context.Context, id int) ([]domain.File, error) {
	query := "SELECT id, user_id, name, size, type, content_type, url, upload_at FROM files WHERE user_id=$1"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	var files []domain.File
	for rows.Next() {
		var file domain.File
		err = rows.Scan(&file.ID, &file.UserID, &file.Name, &file.Size, &file.Type, &file.ContentType, &file.URL, &file.UploadAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, err
}
