package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// Date is storing only date with ISO8601
type Date struct {
	dateString string
	// TODO: change to private
	Time time.Time
}

func (d Date) String() string { return d.dateString }

func (d Date) ToTime() time.Time { return d.Time }

func (d *Date) parseTime() error {
	var err error
	d.Time, err = time.Parse("2006-01-02", d.dateString)
	if err != nil {
		return err
	}
	return nil
}

// Scan implements the Scanner interface for  Date
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var err error

	switch value.(type) {
	case []byte:
		err = json.Unmarshal(value.([]byte), &d.dateString)
	case time.Time:
		d.Time = value.(time.Time)
		d.dateString = d.Time.Format("2006-01-02")
	case *time.Time:
		ts := value.(*time.Time)
		if ts == nil {
			return nil
		}
		d.Time = *ts
		d.dateString = d.Time.Format("2006-01-02")
	case NullTime:
		ts := value.(NullTime)
		if !ts.Valid {
			return nil
		}
		d.Time = ts.Time
		d.dateString = d.Time.Format("2006-01-02")
	}

	if err != nil {
		return err
	}

	return d.parseTime()
}

func (d Date) Value() (driver.Value, error) {
	return json.Marshal(d.dateString)
}

// MarshalJSON for  Date
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.dateString)
}

// UnmarshalJSON for  Date
func (d *Date) UnmarshalJSON(b []byte) error {
	if string(b) == "" {
		return nil
	}
	if err := json.Unmarshal(b, &d.dateString); err != nil {
		return err
	}

	return d.parseTime()
}

// MarshalBinary for  Date
func (d Date) MarshalBinary() ([]byte, error) {
	return d.MarshalJSON()
}

// UnmarshalBinary for  Date
func (d *Date) UnmarshalBinary(b []byte) error {
	return d.UnmarshalJSON(b)
}

// UnmarshalParam for  Date
func (d *Date) UnmarshalParam(data string) error {
	return d.UnmarshalJSON([]byte(`"` + data + `"`))
}

func NewDate(date string) (Date, error) {
	nd := Date{dateString: date}
	err := nd.parseTime()
	return nd, err
}

func (d Date) IsZero() bool { return d.dateString == "" || d.Time.IsZero() }

func (d Date) After(u Date) bool  { return d.Time.After(u.Time) }
func (d Date) Before(u Date) bool { return d.Time.Before(u.Time) }
func (d Date) Equal(u Date) bool  { return d.Time.Equal(u.Time) }
func (d Date) BehindNow() bool {
	if d.Time.IsZero() {
		return false
	}
	return d.Before(NewDateNow())
}

func NewDateFromTime(date time.Time) Date {
	d := Date{
		dateString: date.Format("2006-01-02"),
		Time:       date,
	}
	return d
}

func NewDateNow() Date { return NewDateFromTime(time.Now()) }

type NullDate struct {
	Valid bool
	Date
}

// Scan implements the Scanner interface for  Date
func (d *NullDate) Scan(value interface{}) error {
	if d == nil {
		*d = NullDate{}
	}

	err := d.Date.Scan(value)
	d.Valid = err == nil
	return err
}

func (d NullDate) Value() (driver.Value, error) {
	if d.Valid {
		return d.Date.Value()
	}
	return nil, nil
}

// MarshalJSON for  Date
func (d NullDate) MarshalJSON() ([]byte, error) {
	if d.Valid {
		return d.Date.MarshalJSON()
	}
	return []byte(`null`), nil
}

// UnmarshalJSON for  Date
func (d *NullDate) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || string(b) == `""` {
		return nil
	}

	if d == nil {
		*d = NullDate{}
	}

	err := d.Date.UnmarshalJSON(b)
	d.Valid = err == nil

	return err
}

// MarshalBinary for  Date
func (d NullDate) MarshalBinary() ([]byte, error) {
	return d.MarshalJSON()
}

// UnmarshalBinary for  Date
func (d *NullDate) UnmarshalBinary(b []byte) error {
	return d.UnmarshalJSON(b)
}

// UnmarshalParam for  Date
func (d *NullDate) UnmarshalParam(data string) error {
	return d.UnmarshalJSON([]byte(`"` + data + `"`))
}

func (d *NullDate) TimeWithDefault(defaultValue time.Time) time.Time {
	if d.Valid {
		return d.Time
	}

	return defaultValue
}

func NewNullDate(date string) (NullDate, error) {
	if date == "" {
		return NullDate{}, nil
	}

	nulldate, err := NewDate(date)
	if err != nil {
		return NullDate{}, err
	}

	return NullDate{
		Valid: true,
		Date:  nulldate,
	}, nil
}

func NewNullDateFromTime(date time.Time) NullDate {
	return NullDate{
		Valid: true,
		Date:  NewDateFromTime(date),
	}
}

// NullEmptyDate to handle user's unintentional sending empty string date which are not intended :)
type NullEmptyDate struct{ NullDate }

func (date *NullEmptyDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	s = strings.Trim(s, `'`)
	if strings.ToUpper(s) == "NULL" || s == "" {
		*date = NullEmptyDate{NullDate{Valid: false}}
		return nil
	}
	return date.NullDate.UnmarshalJSON(b)
}

func (date NullEmptyDate) MarshalJSON() ([]byte, error) { return date.NullDate.MarshalJSON() }
