package sns

import "testing"

func TestAlerts_Info(t *testing.T) {
	a := Alerts{
		Client:       NewClient("arn:aws:sns:af-south-1:733171151776:DPS-Alerts", "test"),
		LogReference: "UnitTests",
		Product:      "UnitTests",
		Source:       "UnitTests",
	}
	a.Info("This is a test alert. Please let us know that you received it, and then resolve it.")
}
