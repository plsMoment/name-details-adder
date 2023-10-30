package db

import (
	"context"
	"fmt"
	"name-details-adder/utils"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	Id          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Surname     string           `json:"surname"`
	Patronymic  utils.NullString `json:"patronymic,omitempty"`
	Age         int              `json:"age"`
	Gender      string           `json:"gender"`
	Nationality string           `json:"nationality"`
}

type UserPointers struct {
	Name        *string           `json:"name,omitempty" sql:"name"`
	Surname     *string           `json:"surname,omitempty" sql:"surname"`
	Patronymic  *utils.NullString `json:"patronymic,omitempty" sql:"patronymic"`
	Age         *int              `json:"age,omitempty" sql:"age"`
	Gender      *string           `json:"gender,omitempty" sql:"gender"`
	Nationality *string           `json:"nationality,omitempty" sql:"nationality"`
}

// Use ONLY for READING!!!
// Safe for concurrent read, allocate in init()

var Fields = map[string]bool{
	"id":          true,
	"name":        true,
	"surname":     true,
	"patronymic":  true,
	"age":         true,
	"gender":      true,
	"nationality": true,
}

func (s *Storage) CreateUser() (uuid.UUID, error) {
	scope := "internal.db.queries.CreateUser"

	userData, err := utils.GenerateUser()
	if err != nil {
		return uuid.UUID{}, err
	}
	userDetails, err := utils.NameDetails(userData.Name)
	if err != nil {
		return uuid.UUID{}, err
	}

	id := uuid.New()
	var patronymic *string
	if userData.Patronymic.Valid {
		patronymic = &userData.Patronymic.String
	}
	_, err = s.connPool.Exec(context.TODO(),
		"INSERT INTO users (id, name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, userData.Name, userData.Surname, patronymic, userDetails.Age, userDetails.Gender, userDetails.Nationality,
	)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", scope, err)
	}
	return id, nil
}

func (s *Storage) DeleteUser(userId uuid.UUID) error {
	scope := "internal.db.queries.DeleteUser"

	res, err := s.connPool.Exec(context.TODO(), "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return fmt.Errorf("%s: %w", scope, err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: such user doesn't exist", scope)
	}

	return nil
}

func (s *Storage) UpdateUser(userId uuid.UUID, data *UserPointers) error {
	scope := "internal.db.queries.UpdateUser"

	// Boring code that analyze input struct fields using reflect
	strBuilder := strings.Builder{}
	values := []interface{}{userId}
	strBuilder.WriteString("UPDATE users SET ")
	rt := reflect.TypeOf(*data)
	rv := reflect.ValueOf(*data)
	setIndex := 1
	for i := 0; i < rt.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			var val interface{}
			if rv.Field(i).Elem().Type().String() == "utils.NullString" {
				if rv.Field(i).Elem().Interface().(utils.NullString).Valid {
					val = rv.Field(i).Elem().Interface().(utils.NullString).String
				} else {
					val = nil
				}
			} else {
				val = rv.Field(i).Elem().Interface()
			}
			setIndex++
			sqlCol := rt.Field(i).Tag.Get("sql") // example: "name"
			if sqlCol == "" {
				return fmt.Errorf("%s: sql tag is missing or has empty string value", scope)
			}
			strBuilder.WriteString(sqlCol)
			strBuilder.WriteString(fmt.Sprintf(" = $%d, ", setIndex)) // example: " = $2, "

			values = append(values, val)
		}
	}

	if setIndex == 1 {
		return fmt.Errorf("%s: update request is empty", scope)
	}

	if setIndex != len(values) {
		return fmt.Errorf("%s: something went wrong when update struct was analyzing with reflect", scope)
	}
	setStr := strBuilder.String()         // string example: "UPDATE users SET name=$2, age=$3, "
	setStr = setStr[:len(setStr)-2]       // string example: "UPDATE users SET name=$2, age=$3"
	queryStr := setStr + " WHERE id = $1" // string example: "UPDATE users SET name=$2, age=$3 WHERE id = $1"

	// Query to database
	_, err := s.connPool.Exec(context.TODO(), queryStr, values...)
	if err != nil {
		return fmt.Errorf("%s: %w", scope, err)
	}

	return nil
}

func (s *Storage) GetUsers(dbParams map[string]interface{}) ([]*User, error) {
	scope := "internal.db.queries.GetUsers"

	var queryStr string
	values := make([]interface{}, 0, len(dbParams))
	if len(dbParams) != 0 {
		strBuilder := strings.Builder{}
		strBuilder.WriteString("SELECT * FROM users WHERE ")
		j := 1
		for k, v := range dbParams {
			if k == "patronymic" && v == nil {
				strBuilder.WriteString(fmt.Sprintf("%s IS NULL and ", k))
			} else {
				strBuilder.WriteString(fmt.Sprintf("%s = $%d and ", k, j)) // example: SELECT * FROM users WHERE id = $1 and
				j++
				values = append(values, v)
			}
		}
		queryStr = strBuilder.String()        // example: SELECT * FROM users WHERE id = $1 and age = $2 and
		queryStr = queryStr[:len(queryStr)-5] // example: SELECT * FROM users WHERE id = $1 and age = $2
	} else {
		queryStr = "SELECT * FROM users"
	}

	rows, err := s.connPool.Query(context.TODO(), queryStr, values...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", scope, err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", scope, err)
	}

	return users, nil
}
