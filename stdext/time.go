package stdext

import (
	"time"
)

// Date returns a copy of the given time, rounded to midnight UTC
func Date(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, time.UTC)
}

func DateEqual(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func FixRFC3339Nano(in string) string {
	t, err := time.Parse(time.RFC3339Nano, in)
	if err != nil {
		return in
	}
	RFC3339NanoFixed := "2006-01-02T15:04:05.000000000Z07:00"
	return t.Format(RFC3339NanoFixed)
}

type TimeRange struct {
	Start *time.Time
	End   *time.Time
}

func NewTimeRange(start, end time.Time) TimeRange {
	return TimeRange{&start, &end}
}

func (t TimeRange) StartAt() (time.Time, bool) {
	if t.Start == nil || t.Start == (*time.Time)(nil) {
		return time.Time{}, false
	}
	return *t.Start, true
}

func (t TimeRange) EndAt() (time.Time, bool) {
	if t.End == nil || t.End == (*time.Time)(nil) {
		return time.Time{}, false
	}
	return *t.End, true
}
