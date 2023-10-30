package sns

import (
	"encoding/json"
	"fmt"
	"github.com/Direct-Debit/go-commons/stdext"
	"runtime/debug"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
)

const snsMaxSubjectLength = 100

const (
	AlertFatal = "critical"
	AlertError = "error"
	AlertWarn  = "warning"
	AlertInfo  = "info"
)

var severityMap = map[string]int{
	AlertFatal: 1000,
	AlertError: 100,
	AlertWarn:  10,
	AlertInfo:  1,
	"":         0,
}

func validSeverity(severity string) bool {
	_, valid := stdext.FindInSlice([]string{AlertFatal, AlertError, AlertWarn, AlertInfo}, severity)
	return valid
}

func severeEnough(severity string, minSeverity string) bool {
	return severityMap[severity] >= severityMap[minSeverity]
}

type EventTraceback struct {
	Time      string
	LogInfo   string
	Traceback string
}

// Alerts implements methods to send notifications to SNS with subject:
// [${Product}] ${Severity} event @ ${LogReference} - ${Timestamp}.
// The body format will be as per Pagerduty payload spec for easy integration.
type Alerts struct {
	Client Client
	// Some kind of reference to help us know which logs to check
	LogReference string
	// The affected product
	Product string
	// The unique location of the affected system, preferably a hostname or FQDN.
	Source string
	// The minimum severity to raise alerts for
	MinSeverity string
	// The maximum severity to set when raising alerts.
	// If MaxSeverity is less than MinSeverity, alerts will still be raised at MaxSeverity
	MaxSeverity string
}

func (a Alerts) clampSeverity(severity string) (string, bool) {
	if len(a.MaxSeverity) == 0 || !validSeverity(a.MaxSeverity) {
		a.MaxSeverity = AlertFatal
	}
	if len(a.MinSeverity) == 0 || !validSeverity(a.MinSeverity) {
		a.MinSeverity = AlertInfo
	}

	if !severeEnough(severity, a.MinSeverity) {
		log.Infof("%s is not severe enough, skipping alert", severity)
		return "", false
	}
	if severityMap[severity] > severityMap[a.MaxSeverity] {
		severity = a.MaxSeverity
	}
	return severity, true
}

func (a Alerts) Fatal(s string) {
	a.createAlert(s, AlertFatal)
}

func (a Alerts) Panic(s string) {
	a.createAlert(s, AlertFatal)
}

func (a Alerts) Error(s string) {
	a.createAlert(s, AlertError)
}

func (a Alerts) Warn(s string) {
	a.createAlert(s, AlertWarn)
}

func (a Alerts) Info(s string) {
	a.createAlert(s, AlertInfo)
}

func (a Alerts) Debug(_ string) {}

func (a Alerts) Trace(_ string) {}

func (a Alerts) createAlert(msg string, severity string) {
	if !validSeverity(severity) {
		log.Errorf("Invalid pagerduty severity %q used when trying to create a new alert.", severity)
		return
	}
	var ok bool
	severity, ok = a.clampSeverity(severity)
	if !ok {
		return
	}

	summary := fmt.Sprintf("[%s] %s", a.Product, msg)
	details := EventTraceback{
		Time:      time.Now().Format(time.RFC3339),
		LogInfo:   fmt.Sprintf("%s event @ %s", strings.ToUpper(severity), a.LogReference),
		Traceback: string(debug.Stack()),
	}

	event := pagerduty.V2Payload{
		Summary:  summary,
		Source:   a.Source,
		Severity: severity,
		Details:  details,
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Errorf("Unable to create json for sns alert event: %v", err)
	}

	subj := fmt.Sprintf("%s had a %s event", a.Product, strings.ToUpper(severity))
	if len(subj) > snsMaxSubjectLength {
		subj = subj[:snsMaxSubjectLength]
	}
	if err = a.Client.Publish(subj, string(body)); err != nil {
		log.Errorf("Failed to publish alert to SNS: %v", err)
	}
}
