-- +goose Up
CREATE TABLE repository (
    id TEXT PRIMARY KEY,
    slug VARCHAR(16)
);

-- +goose Down
DROP TABLE repository;
