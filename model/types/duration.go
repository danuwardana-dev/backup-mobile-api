package types

import (
	"encoding/json"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(bb []byte) error {
	var durationString string
	err := json.Unmarshal(bb, &durationString)
	if err != nil {
		return err
	}

	d.Duration, err = time.ParseDuration(durationString)
	return err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
