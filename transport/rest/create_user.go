package rest

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type CreateResponse struct {
	Id string `json:"id,omitempty"`
}

func CreateUser(logger *slog.Logger, changer UserChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, err := changer.CreateUser()
		if err != nil {
			logger.Error("creating User failed", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, CreateResponse{})
			return
		}

		logger.Info("User created")
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, CreateResponse{Id: userId.String()})
	}
}
