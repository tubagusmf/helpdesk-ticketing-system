package helper

import (
	"time"
)

type JSONTime time.Time

const jsonTimeFormat = "2006-01-02 15:04:05"

func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := "\"" + time.Time(t).Format(jsonTimeFormat) + "\""
	return []byte(formatted), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	parsed, err := time.Parse(`"`+jsonTimeFormat+`"`, string(b))
	if err != nil {
		return err
	}
	*t = JSONTime(parsed)
	return nil
}

func (t JSONTime) String() string {
	return time.Time(t).Format(jsonTimeFormat)
}
