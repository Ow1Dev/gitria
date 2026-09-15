
package handlers

import (
	"context"

	"github.com/rs/zerolog"
)

type CreateRepositoryInput struct {
	Body struct {
		Name  string `json:"name" minLength:"1" maxLength:"100" example:"gitria"`
		Owner string `json:"owner" minLength:"1" maxLength:"100" example:"ow1"`
	}
}

type CreateRepositoryOutput struct {
	Body struct {
		Path string `json:"path" example:"gitria"`
	}
}

func CreateRepository(logger zerolog.Logger) func(context.Context, *CreateRepositoryInput) (*CreateRepositoryOutput, error) {
	return func(
		ctx context.Context,
		input *CreateRepositoryInput,
	) (*CreateRepositoryOutput, error) {

		resp := &CreateRepositoryOutput{}
		resp.Body.Path = "gitria" 

		return resp, nil
	}
}
