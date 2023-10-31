package rest

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func DeleteUser(logger *slog.Logger, changer UserChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, err := uuid.Parse(chi.URLParam(r, "userId"))
		if err != nil {
			logger.Error("parsing URL parameter userId failed", err)
			render.Status(r, http.StatusBadRequest)
			return
		}
		logger.Info("Url parameter parsed", slog.Any("user_id", userId))

		err = changer.DeleteUser(userId)
		if err != nil {
			logger.Error("deleting Users failed", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		logger.Info("User deleted")
		render.Status(r, http.StatusNoContent)
	}
}
