-- +goose Up
CREATE TABLE repository (
    id TEXT PRIMARY KEY,
    slug VARCHAR(16) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE repository;
