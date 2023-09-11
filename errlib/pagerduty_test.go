package errlib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPagerDuty_ClampSeverity(t *testing.T) {
	p1 := PagerDuty{
		MinSeverity: PagerDutyWarn,
		MaxSeverity: PagerDutyError,
	}
	_, f := p1.clampSeverity(PagerDutyInfo)
	assert.False(t, f)
	e, _ := p1.clampSeverity(PagerDutyFatal)
	assert.Equal(t, PagerDutyError, e)
}
