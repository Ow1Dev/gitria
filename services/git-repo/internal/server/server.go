package server

import (
	"context"
	v0 "github.com/Ow1Dev/gitria.git/gen/git/v0"
)

type Server struct {
	v0.UnimplementedGitServiceServer
}

// CreateRepository implements [v0.GitServiceServer].
func (s *Server) CreateRepository(context.Context, *v0.CreateRepositoryRequest) (*v0.Repository, error) {
	return &v0.Repository{
		Path: "ow1/gitria",
	}, nil
}

func NewServer() *Server {
	return &Server{}
}
