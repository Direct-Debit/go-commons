package stdext

import (
	"time"
)

// Date returns a copy of the given time, rounded to midnight UTC
func Date(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, time.UTC)
}

func FixRFC3339Nano(in string) string {
	t, err := time.Parse(time.RFC3339Nano, in)
	if err != nil {
		return in
	}
	RFC3339NanoFixed := "2006-01-02T15:04:05.000000000Z07:00"
	return t.Format(RFC3339NanoFixed)
}
