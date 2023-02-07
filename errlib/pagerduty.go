package errlib

import (
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

type EventTraceback struct {
	Message   string
	Time      string
	LogInfo   string
	Traceback string
}

// PagerDuty implements methods to send notifications to pagerduty in the format:
// [${Product}] ${Severity} event @ ${LogReference} - ${Timestamp}
type PagerDuty struct {
	// The Routing Key used to connect to Pagerduty
	RoutingKey string
	// Some kind of reference to help us know which logs to check
	LogReference string
	// The affected product
	Product string
	// The unique location of the affected system, preferably a hostname or FQDN.
	Source string
}

func validPDSeverities() []string {
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

func (p PagerDuty) Debug(_ string) {}

func (p PagerDuty) Trace(_ string) {}

func (p PagerDuty) createPagerdutyAlert(msg string, severity string) {
	_, isValid := stdext.FindInStrSlice(validPDSeverities(), severity)
	if !isValid {
		log.Errorf("Invalid pagerduty severity %q used when trying to create a new alert.", severity)
		return
	}

	summary := fmt.Sprintf("[%s] %s", p.Product, msg)
	details := EventTraceback{
		Message:   msg,
		Time:      time.Now().Format(time.RFC3339),
		LogInfo:   fmt.Sprintf("%s event @ %s", strings.ToUpper(severity), p.LogReference),
		Traceback: string(debug.Stack()),
	}

	event := pagerduty.V2Event{
		RoutingKey: p.RoutingKey,
		Action:     "trigger",
		Payload: &pagerduty.V2Payload{
			Summary:  summary,
			Source:   p.Source,
			Severity: severity,
			Details:  details,
		},
	}
	resp, err := pagerduty.ManageEvent(event)
	if err != nil {
		log.Errorf("Unable to create pagerduty event! %v\n%v", resp, err)
	}
}
