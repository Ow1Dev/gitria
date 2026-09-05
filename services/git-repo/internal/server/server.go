package server

import (
	"context"
	"fmt"

	v0 "github.com/Ow1Dev/gitria.git/gen/git/v0"
)

type Server struct {
	v0.UnimplementedGitServiceServer
}

// CreateRepository implements [v0.GitServiceServer].
func (s *Server) CreateRepository(_ context.Context, req *v0.CreateRepositoryRequest) (*v0.Repository, error) {
	return &v0.Repository{
		Path: fmt.Sprintf("%s/%s", req.Owner, req.Name),
	}, nil
}

func NewServer() *Server {
	return &Server{}
}
