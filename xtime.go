package xtime

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

var (
	_defaultOptions = options{
		timeLayout: "2006-01-02 15:04:05",
		nullLayout: "null",
	}
	_options options
)

func Init(opts ...Option) {
	_options = _defaultOptions
	for _, o := range opts {
		o.apply(&_options)
	}
}

type Time struct {
	Time  time.Time
	Valid bool
}

func New(t time.Time, valid bool) Time {
	return Time{
		Time:  t,
		Valid: valid,
	}
}

func From(t time.Time) Time {
	return New(t, true)
}

func FromPtr(t *time.Time) Time {
	if t == nil {
		return New(time.Time{}, false)
	}
	return New(*t, true)
}

func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case []byte:
		t.Time, err = time.ParseInLocation(_options.timeLayout, string(x), time.Local)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("scan type %T into xtime.Time: %v", value, value)
	}
	t.Valid = err == nil
	return err
}

func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte(_options.nullLayout), nil
	}

	s := `"` + t.Time.Format(_options.timeLayout) + `"`
	return []byte(s), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		if x == _options.nullLayout {
			t.Valid = false
			return nil
		}
		t.Time, err = time.ParseInLocation(_options.timeLayout, x, time.Local)
		t.Valid = err == nil
		return err
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type xtime.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		t.Valid = valid
		t.Time, err = time.ParseInLocation(_options.timeLayout, ti, time.Local)
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		t.Valid = false
		return fmt.Errorf("json: cannot unmarshal %v into Go value of type xtime.Time", reflect.TypeOf(v).Name())
	}
}

func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte(_options.nullLayout), nil
	}
	return []byte(t.Time.Format(_options.timeLayout)), nil
}

func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	if str == _options.nullLayout {
		t.Valid = false
		return nil
	}

	var err error
	t.Time, err = time.ParseInLocation(_options.timeLayout, string(text), time.Local)
	t.Valid = err == nil
	return err
}

func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func (t Time) IsZero() bool {
	return !t.Valid
}

func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

func (t Time) String() string {
	if !t.Valid {
		return _options.nullLayout
	}
	return t.Time.Format(_options.timeLayout)
}