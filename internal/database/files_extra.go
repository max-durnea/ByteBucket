package database

import (
	"context"
	"github.com/google/uuid"
)

// GetFilesByUser returns files belonging to a specific user ordered by created_at desc.
func (q *Queries) GetFilesByUser(ctx context.Context, userID uuid.UUID) ([]File, error) {
	rows, err := q.db.QueryContext(ctx, `SELECT id, user_id, object_key, file_name, mime_type, created_at FROM files WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.UserID, &f.ObjectKey, &f.FileName, &f.MimeType, &f.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}
