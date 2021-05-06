package errlib

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Direct-Debit/go-commons/stdext"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
)

const (
	PagerDutyFatal = "critical"
	PagerDutyError = "error"
	PagerDutyWarn  = "warning"
	PagerDutyInfo  = "info"
)

type event_traceback struct {
	Message   string
	Traceback string
}

type PagerDuty struct {
	LogReference string
	RoutingKey   string
	Product      string
}

func getValidSeverities() []string {
	return []string{PagerDutyFatal, PagerDutyError, PagerDutyWarn, PagerDutyInfo}
}

func (p PagerDuty) Fatal(s string) {
	p.createPagerdutyAlert(s, PagerDutyFatal)
}

func (p PagerDuty) Panic(s string) {
	p.createPagerdutyAlert(s, PagerDutyFatal)
}

func (p PagerDuty) Error(s string) {
	p.createPagerdutyAlert(s, PagerDutyError)
}

func (p PagerDuty) Warn(s string) {
	p.createPagerdutyAlert(s, PagerDutyWarn)
}

func (p PagerDuty) Info(s string) {
	p.createPagerdutyAlert(s, PagerDutyInfo)
}

func (p PagerDuty) Debug(s string) {
	return
}

func (p PagerDuty) Trace(s string) {
	return
}

func (p PagerDuty) createPagerdutyAlert(msg string, severity string) {
	_, isValid := stdext.Contains(getValidSeverities(), severity)
	if !isValid {
		ErrorError(errors.New("Value Error"), "%s", fmt.Sprintf("Invalid pagerduty severity %q used when trying to create a new alert.", severity))
		return
	}

	details := event_traceback{msg, string(debug.Stack())}
	summary := fmt.Sprintf("[%s] %s event @ %s - %s", p.Product, strings.ToUpper(severity), p.LogReference, time.Now().Format(time.RFC3339))

	event := pagerduty.V2Event{
		RoutingKey: p.RoutingKey,
		Action:     "trigger",
		Payload: &pagerduty.V2Payload{
			Summary:  summary,
			Source:   "DPS Monitor",
			Severity: severity,
			Details:  details,
		},
	}
	resp, err := pagerduty.ManageEvent(event)
	if err != nil {
		log.Errorf("Unable to create pagerduty event! %v\n%v", resp, err)
	}
}
