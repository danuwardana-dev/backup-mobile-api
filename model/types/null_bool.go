package types

import (
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return []byte(`null`), nil
}

func (n *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	err := json.Unmarshal(b, &n.Bool)
	n.Valid = err == nil
	return err
}

func NewNullBool(s bool) NullBool { return NullBool{sql.NullBool{Bool: s, Valid: true}} }
