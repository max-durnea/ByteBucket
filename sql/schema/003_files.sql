-- +goose Up

CREATE TABLE files(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE files;