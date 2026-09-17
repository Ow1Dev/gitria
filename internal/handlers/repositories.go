package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ow1Dev/gitria-git/internal/sqliteStore"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"

	"github.com/oklog/ulid/v2"
)

type CreateRepositoryInput struct {
	Body struct {
		Name  string `json:"name" minLength:"1" maxLength:"16" example:"gitria"`
	}
}

type CreateRepositoryOutput struct {
	Body struct {
		Id string `json:"id" example:"gitria"`
	}
}

func CreateRepository(db *sql.DB, queries sqliteStore.Queries, logger zerolog.Logger) func(context.Context, *CreateRepositoryInput) (*CreateRepositoryOutput, error) {
	return func(
		ctx context.Context,
		input *CreateRepositoryInput,
	) (*CreateRepositoryOutput, error) {

		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()
		qtx := queries.WithTx(tx)

		id := ulid.Make()
		val, err := qtx.CreateRepository(ctx, 
			sqliteStore.CreateRepositoryParams{
				ID: id.String(),
				Slug: input.Body.Name,
			},
		)
		if err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
				return nil, fmt.Errorf("repository with slug %q already exists", input.Body.Name)
			}

			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}

		rsp := &CreateRepositoryOutput{}
		rsp.Body.Id = val.ID

		return rsp, nil
	}
}
