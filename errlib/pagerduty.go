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

var severityMap = map[string]int{
	PagerDutyFatal: 1000,
	PagerDutyError: 100,
	PagerDutyWarn:  10,
	PagerDutyInfo:  1,
	"":             0,
}

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
	// The minimum severity to raise alerts for
	MinSeverity string
	// Tha maximum severity to set when raising alerts.
	// If MaxSeverity is less than MinSeverity, alerts will still be raised at MaxSeverity
	MaxSeverity string
}

// The list of all severities in order of most to least severe
func validSeverity(severity string) bool {
	_, valid := stdext.FindInSlice([]string{PagerDutyFatal, PagerDutyError, PagerDutyWarn, PagerDutyInfo}, severity)
	return valid
}

func severeEnough(severity string, minSeverity string) bool {
	return severityMap[severity] >= severityMap[minSeverity]
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
	if !validSeverity(severity) {
		log.Errorf("Invalid pagerduty severity %q used when trying to create a new alert.", severity)
		return
	}
	if !severeEnough(severity, p.MinSeverity) {
		log.Infof("%s is not severe enough, skipping alert", severity)
	}
	if len(p.MaxSeverity) == 0 {
		p.MaxSeverity = PagerDutyFatal
	}
	if validSeverity(p.MaxSeverity) && severityMap[severity] < severityMap[p.MaxSeverity] {
		severity = p.MaxSeverity
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
