package errlib

import (
	"fmt"
	"strings"

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
	Traceback string
}

type PagerDuty struct {
	Environment   string
	SeverityLevel string
	RoutingKey    string
	Component     string
	Function      string
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
	maximumSeverity, _ := Find([]string{"critical", "error", "warning", "info"}, p.SeverityLevel)
	currentSeverity, isValid := Find([]string{"critical", "error", "warning", "info"}, severity)
	if !isValid {
		log.Errorf("Invalid severity level %q passed to pagerduty logger! No pagerduty alert was created.", severity)
		return
	}
	if currentSeverity > maximumSeverity {
		return
	}

	details := event_traceback{msg}
	summary := fmt.Sprintf("[dps/error/DPSM] - %s : {", strings.ToUpper(severity))
	summary += fmt.Sprintf(" Environment : %q,", p.Environment)
	summary += fmt.Sprintf(" Component : %q,", p.Component)
	summary += fmt.Sprintf(" Function : %q }", p.Function)

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

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
