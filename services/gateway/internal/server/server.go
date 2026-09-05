package server

import (
	"net/http"

	"github.com/rs/zerolog"
	pb "github.com/Ow1Dev/gitria.git/gen/git/v0"
)

func New(c pb.GitServiceClient, logger *zerolog.Logger) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, c, logger)

	var h http.Handler = mux
	return h
}
