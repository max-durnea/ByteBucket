-- name: CreateFile :one

INSERT INTO files(id,user_id,object_key,file_name,mime_type,created_at)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;