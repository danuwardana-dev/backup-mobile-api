package types

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return []byte(`null`), nil
}

func (n *NullString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	err := json.Unmarshal(b, &n.String)
	n.Valid = err == nil
	return err
}

func NewNullString(s string) NullString { return NullString{sql.NullString{String: s, Valid: true}} }

func NewValidNullString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: s != ""}}
}
