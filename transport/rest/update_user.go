package rest

import (
	"log/slog"
	"name-details-adder/internal/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UpdateResponse struct{}

func UpdateUser(logger *slog.Logger, changer UserChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, err := uuid.Parse(chi.URLParam(r, "userId"))
		if err != nil {
			logger.Error("parsing URL parameter userId failed", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, UpdateResponse{})
			return
		}
		logger.Info("Url parameter parsed", slog.Any("user_id", userId))

		req := db.UserPointers{}
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("decoding request body failed", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, UpdateResponse{})
			return
		}

		logger.Info("Request body decoded", slog.Any("request", req))
		err = changer.UpdateUser(userId, &req)
		if err != nil {
			logger.Error("changing user failed", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, UpdateResponse{})
			return
		}

		logger.Info("User changed")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, UpdateResponse{})
	}
}
