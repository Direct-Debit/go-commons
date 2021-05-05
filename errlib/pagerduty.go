package errlib

import (
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
)
// TODO export vars for external use
const (
	sevFatal = "critical"
	sevError = "error"
	sevWarn  = "warning"
	sevInfo  = "info"
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
	maximumSeverity, _ := Find([]string{"critical", "error", "warning", "info"}, p.SeverityLevel)
	currentSeverity, isValid := Find([]string{"critical", "error", "warning", "info"}, severity)
	if !isValid {
		debugFunc(fmt.Sprintf("Invalid severity level %q passed to pagerduty logger! No pagerduty alert was created.", severity))
		return
	}
	if currentSeverity > maximumSeverity {
		return
	}

	details := event_traceback{msg}
	summary := fmt.Sprintf("[dps/error/DPSM] - %s : {", strings.ToUpper(severity))
	summary += fmt.Sprintf(" Environment : %q ,", p.Environment)
	summary += fmt.Sprintf(" Component : %q ,", p.Component)
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
		debugFunc(fmt.Sprintf("Unable to create pagerduty event! %v\n%v", resp, err))
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
