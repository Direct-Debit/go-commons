package errlib

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
)

const PAGERDUTY_ROUTING_KEY = "d7827f7d66574b06d0696b696e865442"

const (
	sevFatal = "fatal"
)

type event_traceback struct {
	Traceback string
}

type PagerDuty struct {
	Product     string
	Component   string
	Environment string
}

func (p PagerDuty) Fatal(s string) {
	p.createPagerdutyAlert(s, sevFatal)
}

func (p PagerDuty) Panic(s string) {
	panic("implement me")
}

func (p PagerDuty) Error(s string) {
	panic("implement me")
}

func (p PagerDuty) Warn(s string) {
	panic("implement me")
}

func (p PagerDuty) Info(s string) {
	panic("implement me")
}

func (p PagerDuty) Debug(s string) {
	panic("implement me")
}

func (p PagerDuty) Trace(s string) {
	panic("implement me")
}

func (p PagerDuty) createPagerdutyAlert(msg string, severity string) {
	details := event_traceback{msg}
	event := pagerduty.V2Event{
		RoutingKey: PAGERDUTY_ROUTING_KEY,
		Action:     "trigger",
		Payload: &pagerduty.V2Payload{
			Summary:  "[dps/error/DPSM]",
			Source:   "DPS Monitor",
			Severity: severity,
			Details:  details,
		},
	}
	message, res := pagerduty.ManageEvent(event)
	if res != nil {
		warnFunc(fmt.Sprintf("Unable to create pagerduty event! %v\n%v", res, message))
	}
}
