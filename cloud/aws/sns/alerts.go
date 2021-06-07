package sns

import (
	"encoding/json"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"strings"
	"time"
)

const snsMaxSubjectLength = 100

const (
	AlertFatal = "critical"
	AlertError = "error"
	AlertWarn  = "warning"
	AlertInfo  = "info"
)

type EventTraceback struct {
	Message   string
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
	details := EventTraceback{msg, string(debug.Stack())}
	summary := fmt.Sprintf("[%s] %s event @ %s - %s", a.Product, strings.ToUpper(severity), a.LogReference, time.Now().Format(time.RFC3339))

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
