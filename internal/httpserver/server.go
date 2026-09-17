package httpserver

import (
	"database/sql"
	"net/http"

	"github.com/Ow1Dev/gitria-git/internal/sqliteStore"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/zerolog"
)

	
func NewHttpServer(db *sql.DB, queries sqliteStore.Queries, logger zerolog.Logger) http.Handler {
	mux := http.NewServeMux()

	config := huma.DefaultConfig("gitrea", "1.00")
	config.DocsPath = "";
	config.OpenAPIPath = "";
	
	api := humago.New(mux, config)

	addRoutes(api, db, queries, mux, logger)

	var h http.Handler = mux

	return h
}
