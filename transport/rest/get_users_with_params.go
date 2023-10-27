package rest

import (
	"fmt"
	"log/slog"
	"name-details-adder/internal/db"
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

		var fields = map[string]bool{
			"id":          true,
			"name":        true,
			"surname":     true,
			"patronymic":  true,
			"age":         true,
			"gender":      true,
			"nationality": true,
		}

		dbParams := make(map[string]string)
		queryParams := r.URL.Query()
		for k, v := range queryParams {
			if fields[k] {
				if len(v) < 1 {
					logger.Error(fmt.Sprintf("empty value of query parameter: %s", k))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, GetResponse{})
					return
				} else if len(v) > 1 {
					logger.Error(fmt.Sprintf("too many values for query parameter: %s", k))
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, GetResponse{})
				}
				dbParams[k] = v[0]
			}
		}
		logger.Info("Url parameters parsed", slog.Any("params", dbParams))
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
