package httpserver

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/zerolog"
)

func NewHttpServer(logger zerolog.Logger) http.Handler {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("gitrea", "1.00"))

	addRoutes(api, mux, logger)

	var h http.Handler = mux

	return h
}
