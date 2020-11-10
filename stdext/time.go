package stdext

import "time"

// return a copy of the given time, but with hours, minutes, seconds and nanoseconds set to 0
func Date(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, in.Location())
}
