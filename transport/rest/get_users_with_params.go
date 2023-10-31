package rest

import (
	"log/slog"
	"name-details-adder/internal/db"
	"name-details-adder/utils/validators"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type GetResponse struct {
	Users      []*db.User `json:"users"`
	NextPageId int        `json:"next_page_id,omitempty"`
}

func GetUsers(logger *slog.Logger, changer UserChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		queryParams := r.URL.Query()
		logger.Info("Url parameters parsed", slog.Any("params", queryParams))

		pageIdStr := queryParams.Get("page_id")
		pageId := 1
		var err error
		if pageIdStr != "" {
			pageId, err = strconv.Atoi(pageIdStr)
			if err != nil {
				logger.Error("parsing query parameter age_id failed", err)
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, GetResponse{})
				return
			}
		}
		logger.Info("Url parameter page_id parsed", slog.Any("page_id", pageId))

		dbParams, err := validators.ValidateQuery(queryParams)
		if err != nil {
			logger.Error("validating query failed", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, GetResponse{})
			return
		}
		logger.Info("Valid url parameters", slog.Any("params", dbParams))

		pageSize := 5
		users, err := changer.GetUsers(dbParams, pageId, pageSize)
		if err != nil {
			logger.Error("getting Users failed", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, GetResponse{})
			return
		}

		logger.Info("Users got")
		render.Status(r, http.StatusOK)
		if len(users) == pageSize {
			render.JSON(w, r, GetResponse{Users: users, NextPageId: pageId + 1})
		} else {
			render.JSON(w, r, GetResponse{Users: users})
		}
	}
}
