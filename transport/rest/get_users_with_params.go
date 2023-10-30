package rest

import (
	"log/slog"
	"name-details-adder/internal/db"
	"name-details-adder/utils/validators"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type GetResponse struct {
	Users []*db.User `json:"users"`
}

func GetUsers(logger *slog.Logger, changer UserChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		queryParams := r.URL.Query()
		logger.Info("Url parameters parsed", slog.Any("params", queryParams))

		dbParams, err := validators.ValidateQuery(queryParams)
		if err != nil {
			logger.Error("validating query failed", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, GetResponse{})
			return
		}
		logger.Info("Valid url parameters", slog.Any("params", dbParams))

		users, err := changer.GetUsers(dbParams)
		if err != nil {
			logger.Error("getting Users failed", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, GetResponse{})
			return
		}

		logger.Info("Users got")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, GetResponse{Users: users})
	}
}
