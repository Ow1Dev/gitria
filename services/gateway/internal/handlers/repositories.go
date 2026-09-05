package handlers

import (
	"net/http"

	"github.com/rs/zerolog"

	pb "github.com/Ow1Dev/gitria.git/gen/git/v0"
	"github.com/Ow1Dev/gitria.git/libs/httputil"
)

func HandleCreateRepostory(conn pb.GitServiceClient, logger *zerolog.Logger) http.Handler {
	type response struct {
			Path string `json:"path"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rsp, err := conn.CreateRepository(r.Context(), &pb.CreateRepositoryRequest{
			Name: "gitria",
			Owner: "ow1",
		})

		if err != nil {
			logger.Error().Err(err).Msg("Error when CreateRepository on git service")
			httputil.WriteError(w, http.StatusInternalServerError, "I got bad request back from my grpc server")
			return
		}

		httputil.Encode(w, r, http.StatusCreated, response{
			Path:  rsp.Path,
		})
	})
}
