package errlib

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
)

const PAGERDUTY_ROUTING_KEY = "d7827f7d66574b06d0696b696e865442"

type event_traceback struct {
	Traceback string
}

func CreatePagerdutyAlert(err error, severity string, format string, a ...interface{}) {
	details := event_traceback{fmt.Sprintf(format, a...)}
	event := pagerduty.V2Event{
		RoutingKey: PAGERDUTY_ROUTING_KEY,
		Action:     "trigger",
		Payload: &pagerduty.V2Payload{
			Summary:  "[dps/error/DPSM] " + fmt.Sprintf("%v", err),
			Source:   "DPS Monitor",
			Severity: severity,
			Details:  details,
		},
	}
	msg, res := pagerduty.ManageEvent(event)
	if res != nil {
		warnFunc(fmt.Sprintf("Unable to create pagerduty event! %v\n%v", res, msg))
	}
}
