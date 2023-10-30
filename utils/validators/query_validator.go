package validators

import (
	"fmt"
	"name-details-adder/internal/db"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

func ValidateQuery(queryParams url.Values) (map[string]interface{}, error) {
	scope := "utils.query_validator.ValidateQuery"

	dbParams := make(map[string]interface{})
	for k, v := range queryParams {
		if db.Fields[k] {
			if len(v) < 1 {
				return nil, fmt.Errorf("%s:: empty value of query parameter: %s", scope, k)
			} else if len(v) > 1 {
				return nil, fmt.Errorf("%s:: too many values for query parameter: %s", scope, k)
			}
			switch k {
			case "age":
				age, err := strconv.Atoi(v[0])
				if err != nil {
					return nil, fmt.Errorf("%s:: invalid value of query parameter: age", scope)
				}
				dbParams[k] = age
			case "id":
				userId, err := uuid.Parse(v[0])
				if err != nil {
					return nil, fmt.Errorf("%s:: invalid value of query parameter: id", scope)
				}
				dbParams[k] = userId
			case "patronymic":
				if v[0] == "null" {
					dbParams[k] = nil
				} else {
					dbParams[k] = v[0]
				}
			default:
				dbParams[k] = v[0]
			}
		}
	}
	return dbParams, nil
}
