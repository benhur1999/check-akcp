package sensorProbePlus

import "github.com/benhur1999/check-akcp/internal/akcp"

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
