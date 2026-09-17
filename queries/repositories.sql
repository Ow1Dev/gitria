-- name: CreateRepository :one
INSERT INTO repository (
  id, slug
) VALUES (
  ?, ?
)
RETURNING *;
