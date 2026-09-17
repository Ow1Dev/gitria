package httpserver

import (
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Ow1Dev/gitria-git/internal/handlers"
	"github.com/Ow1Dev/gitria-git/internal/sqliteStore"
	"github.com/Ow1Dev/gitria-git/pkgs/httputil"
	"github.com/rs/zerolog"
)

func addRoutes(
	api huma.API,
	db *sql.DB,
	queries sqliteStore.Queries,
	mux *http.ServeMux,
	logger zerolog.Logger) {
	huma.Post(api, "/repository", handlers.CreateRepository(db, queries, logger)) 

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		httputil.WriteError(w, http.StatusNotFound, "about:blank", "The requested resource was not found.")
	})
}
