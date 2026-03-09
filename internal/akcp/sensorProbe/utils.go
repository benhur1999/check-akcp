package sensorProbe

import "regexp"

func (m *SensorProbe) ValidatePort(sensorPort string) bool {
	matched, _ := regexp.MatchString("^([0-9]|[1-9][0-9])$", sensorPort)
	return matched
}
