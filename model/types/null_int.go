package types

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return []byte(`null`), nil
}

func (n *NullInt64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	err := json.Unmarshal(b, &n.Int64)
	n.Valid = err == nil
	return err
}

func NewNullInt64[T int | int64 | uint](s T) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: int64(s), Valid: true}}
}
