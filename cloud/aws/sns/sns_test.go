package sns

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Publish(t *testing.T) {
	log.SetLevel(log.TraceLevel)

	cl := NewClient("arn:aws:sns:af-south-1:733171151776:dps-reports", "test")
	err := cl.Publish("test-message", "This is a test message from Go")
	assert.Nil(t, err)

	err = cl.Publish("test-message (with parenthesis)", "This is a test message from Go. It has parenthesis in the subject line.")
	assert.Nil(t, err)

	err = cl.Publish("(another test) test-message", "This is a test message from Go. It has parenthesis in the beginning of the subject line.")
	assert.Nil(t, err)
}
