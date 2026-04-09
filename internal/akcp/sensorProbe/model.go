package sensorProbe

import (
	"fmt"

	"github.com/benhur1999/check-akcp/internal/akcp"
)

type SensorProbe struct {
	akcp.AkcpBase
}

func New(description string, name string, location string) *SensorProbe {
	return &SensorProbe{
		AkcpBase: akcp.NewAkcpBase(description, name, location),
	}
}

func (m *SensorProbe) GetModel() string {
	return "SensorProbe"
}

func (m *SensorProbe) GetOverallSummaryLine() string {
	name := m.GetName()
	location := m.GetLocation()
	if len(location) > 0 {
		if len(name) > 0 {
			return fmt.Sprintf("%s %s at location %s (%s)", m.GetModel(), m.GetName(), location, m.GetDescription())
		} else {
			return fmt.Sprintf("%s at location %s (%s)", m.GetModel(), location, m.GetDescription())
		}
	} else {
		if len(name) > 0 {
			return fmt.Sprintf("%s %s (%s)", m.GetModel(), m.GetName(), m.GetDescription())
		} else {
			return fmt.Sprintf("%s (%s)", m.GetModel(), m.GetDescription())
		}
	}
}
