package stdext

import "time"

// Date returns a copy of the given time, rounded to midnight UTC
func Date(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, time.UTC)
}
