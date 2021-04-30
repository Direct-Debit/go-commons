//	RoutingKey: "82b77581eb58450dc02ee107978f8403"
package errlib

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
)

const (
	sevFatal = "fatal"
	sevError = "error"
	sevWarn  = "warning"
	sevInfo  = "info"
)

type event_traceback struct {
	Traceback string
}

type PagerDuty struct {
	Product     string
	Component   string
	Environment string
	RoutingKey  string
}

func (p PagerDuty) Fatal(s string) {
	p.createPagerdutyAlert(s, sevFatal)
}

func (p PagerDuty) Panic(s string) {
	p.createPagerdutyAlert(s, sevFatal)
}

func (p PagerDuty) Error(s string) {
	p.createPagerdutyAlert(s, sevError)
}

func (p PagerDuty) Warn(s string) {
	p.createPagerdutyAlert(s, sevWarn)
}

func (p PagerDuty) Info(s string) {
	p.createPagerdutyAlert(s, sevInfo)
}

func (p PagerDuty) Debug(s string) {
	return
}

func (p PagerDuty) Trace(s string) {
	return
}

func (p PagerDuty) createPagerdutyAlert(msg string, severity string) {
	details := event_traceback{msg}
	event := pagerduty.V2Event{
		RoutingKey: p.RoutingKey,
		Action:     "trigger",
		Payload: &pagerduty.V2Payload{
			Summary:  "[dps/error/DPSM]", //Add product/component/env?
			Source:   "DPS Monitor",
			Severity: severity,
			Details:  details,
		},
	}
	resp, err := pagerduty.ManageEvent(event)
	if err != nil {
		warnFunc(fmt.Sprintf("Unable to create pagerduty event! %v\n%v", resp, err))
	}
}
