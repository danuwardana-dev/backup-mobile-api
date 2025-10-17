package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return []byte(`null`), nil
}

func (n *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}

	err := json.Unmarshal(b, &n.Time)
	n.Valid = err == nil
	return err
}

func NewNullTime(s time.Time) NullTime { return NullTime{sql.NullTime{Time: s, Valid: true}} }
