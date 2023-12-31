package rest

import (
	"name-details-adder/internal/db"

	"github.com/google/uuid"
)

type UserChanger interface {
	CreateUser() (uuid.UUID, error)
	DeleteUser(userId uuid.UUID) error
	UpdateUser(uuid.UUID, *db.UserPointers) error
	GetUsers(map[string]interface{}, int, int) ([]*db.User, error)
}
