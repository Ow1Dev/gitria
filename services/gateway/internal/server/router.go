package server

import (
	"net/http"

	"github.com/Ow1Dev/gitria.git/libs/httputil"
	"github.com/Ow1Dev/gitria.git/services/gateway/internal/handlers"
	"github.com/rs/zerolog"

	pb "github.com/Ow1Dev/gitria.git/gen/git/v0"
)

func addRoutes(mux *http.ServeMux, conn pb.GitServiceClient, logger *zerolog.Logger) {
	mux.Handle("POST /repository", handlers.HandleCreateRepostory(conn, logger))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		httputil.WriteError(w, http.StatusNotFound, "not found")
	})
}

