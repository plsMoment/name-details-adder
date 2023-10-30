package utils

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type NullString sql.NullString

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = string(b) != "null"
	return err
}

func (ns *NullString) TextValue() (pgtype.Text, error) {
	return pgtype.Text{
		String: ns.String,
		Valid:  ns.Valid,
	}, nil
}

func (ns *NullString) ScanText(v pgtype.Text) error {
	*ns = NullString(v)
	return nil
}
