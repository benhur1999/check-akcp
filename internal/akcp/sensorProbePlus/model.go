package sensorProbePlus

import (
	"fmt"

	"github.com/benhur1999/check-akcp/internal/akcp"
)

type SensorProbePlus struct {
	akcp.AkcpBase
}

func New(description string, name string, location string) *SensorProbePlus {
	return &SensorProbePlus{
		AkcpBase: akcp.NewAkcpBase(description, name, location),
	}
}

func (m *SensorProbePlus) GetModel() string {
	return "SensorProbe+"
}

func (m *SensorProbePlus) GetOverallSummaryLine() string {
	location := m.GetLocation()
	if len(location) > 0 {
		return fmt.Sprintf("%s %s at location %s (%s)", m.GetModel(), m.GetName(), location, m.GetDescription())
	} else {
		return fmt.Sprintf("%s %s (%s)", m.GetModel(), m.GetName(), m.GetDescription())
	}
}
