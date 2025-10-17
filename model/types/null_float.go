package types

import (
	"database/sql"
	"encoding/json"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}
	return []byte(`null`), nil
}

func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	err := json.Unmarshal(b, &n.Float64)
	n.Valid = err == nil
	return err
}

func NewNullFloat64(s float64) NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: s, Valid: true}}
}

func NewValidNullFloat64(s float64) NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: s, Valid: s != 0}}
}
