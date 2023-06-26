package stdext

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_FixRFC3339Nano(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{
			in:       "2022-06-20T08:32:58.26840593Z",
			expected: "2022-06-20T08:32:58.268405930Z",
		}, {
			in:       "2022-06-20T08:32:58.268405930Z",
			expected: "2022-06-20T08:32:58.268405930Z",
		}, {
			in:       "2022-06-20T08:32:58.26840593+07:00",
			expected: "2022-06-20T08:32:58.268405930+07:00",
		}, {
			in:       "2022-06-20T08:32:58.268405930+07:00",
			expected: "2022-06-20T08:32:58.268405930+07:00",
		}, {
			in:       "This is not a valid RFC3339Nano time",
			expected: "This is not a valid RFC3339Nano time",
		}, {
			in:       "invalid",
			expected: "invalid",
		},
	}

	for _, c := range cases {
		r := FixRFC3339Nano(c.in)
		assert.Equal(t, c.expected, r)
	}
}

func TestDateEqual(t *testing.T) {
	t1 := time.Now()
	t2 := time.Now().Add(time.Minute)

	assert.True(t, DateEqual(t1, t2))
	t2 = t2.AddDate(0, 0, 1)
	assert.False(t, DateEqual(t1, t2))
}
