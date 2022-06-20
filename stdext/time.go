package stdext

import (
	"strings"
	"time"
)

// Date returns a copy of the given time, rounded to midnight UTC
func Date(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, time.UTC)
}

func FixRFC3339Nano(in string) string {
	_, err := time.Parse(time.RFC3339Nano, in)
	if err != nil {
		return in
	}

	var splits []string
	var tz string
	switch {
	case strings.Contains(in, "Z"):
		splits = strings.SplitN(in, "Z", 2)
		tz = "Z"
		break
	case strings.Contains(in, "+"):
		splits = strings.SplitN(in, "+", 2)
		tz = "+"
		break
	case strings.Contains(in, "-"):
		splits = strings.SplitN(in, "-", 2)
		tz = "-"
		break
	}

	tp := splits[0]
	target := len(time.RFC3339Nano) - 6

	if len(tp) >= target {
		return in
	}

	builder := strings.Builder{}
	builder.WriteString(tp)
	for builder.Len() < target {
		builder.WriteString("0")
	}
	builder.WriteString(tz)
	builder.WriteString(splits[1])
	return builder.String()
}
