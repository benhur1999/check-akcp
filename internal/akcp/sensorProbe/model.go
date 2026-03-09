package sensorProbe

import "github.com/benhur1999/check-akcp/internal/akcp"

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
