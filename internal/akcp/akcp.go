package akcp

import (
	"fmt"

	"github.com/gosnmp/gosnmp"
)

type SensorStatus int64

const (
	StatusNoStatus     SensorStatus = 1
	StatusNormal       SensorStatus = 2
	StatusHighWarning  SensorStatus = 3
	StatusHighCritical SensorStatus = 4
	StatusLowWarning   SensorStatus = 5
	StatusLowCritical  SensorStatus = 6
	StatusSensorError  SensorStatus = 7
)

var statusStrings = map[SensorStatus]string{
	StatusNoStatus:     "No Status",
	StatusNormal:       "Normal",
	StatusHighWarning:  "High Warning",
	StatusHighCritical: "High Critical",
	StatusLowWarning:   "Low Warning",
	StatusLowCritical:  "Low Critical",
	StatusSensorError:  "Sensor Error",
}

type TemperatureUnit uint64

const (
	TemperatureUnitCelsius TemperatureUnit = iota
	TemperatureUnitFahrenheit
)

type TemperatureSensor struct {
	Port         string
	Description  string
	Degree       float64
	Unit         TemperatureUnit
	Status       SensorStatus
	Online       bool
	LowCritical  *float64
	LowWarning   *float64
	HighWarning  *float64
	HighCritical *float64
}

func (d *TemperatureSensor) GetUnit() string {
	switch d.Unit {
	case TemperatureUnitCelsius:
		return "C"
	case TemperatureUnitFahrenheit:
		return "F"
	default:
		return ""
	}
}

func (d *TemperatureSensor) GetStatus() string {
	return statusStrings[d.Status]
}

type HumidityUnit uint64

const (
	HumidityUnitRelativeHumidity HumidityUnit = iota
)

type HumiditySensor struct {
	Port         string
	Description  string
	Percent      float64
	Unit         HumidityUnit
	Status       SensorStatus
	Online       bool
	LowCritical  *float64
	LowWarning   *float64
	HighWarning  *float64
	HighCritical *float64
}

func (d *HumiditySensor) GetUnit() string {
	switch d.Unit {
	case HumidityUnitRelativeHumidity:
		return "%"
	default:
		return ""
	}
}

func (d *HumiditySensor) GetStatus() string {
	return statusStrings[d.Status]
}

type DryContactType uint64

const (
	DryContactTypeInput  DryContactType = 7
	DryContactTypeOutput DryContactType = 8
	DryContactTypeArray  DryContactType = 22
)

type DryContactStatus int64

const (
	DryContactStatusNoStatus     DryContactStatus = 1
	DryContactStatusNormal       DryContactStatus = 2
	DryContactStatusHighCritical DryContactStatus = 4
	DryContactStatusLowCritical  DryContactStatus = 6
	DryContactStatusSensorError  DryContactStatus = 7
	DryContactStatusOutputLow    DryContactStatus = 8
	DryContactStatusOutputHigh   DryContactStatus = 9
)

var dryContactStatusStrings = map[DryContactStatus]string{
	DryContactStatusNoStatus:     "No Status",
	DryContactStatusNormal:       "Normal",
	DryContactStatusHighCritical: "High Critical",
	DryContactStatusLowCritical:  "Low Critical",
	DryContactStatusSensorError:  "Sensor Error",
	DryContactStatusOutputLow:    "Output Low",
	DryContactStatusOutputHigh:   "Output High",
}

type DryContactDirection int64

const (
	DryContactDirectionInput  DryContactDirection = 0
	DryContactDirectionOutput DryContactDirection = 1
)

type DryContactNormalState int64

const (
	DryContactNormalStateClosed DryContactNormalState = 0
	DryContactNormalStateOpen   DryContactNormalState = 1
)

type DryContact struct {
	Port                string
	Description         string
	Type                DryContactType
	Status              DryContactStatus
	Online              bool
	Direction           DryContactDirection
	NormalState         DryContactNormalState
	CriticalDescription string
	NormalDescription   string
}

func (d *DryContact) GetStatus() string {
	return dryContactStatusStrings[d.Status]
}

func (d *DryContact) IsOutput() bool {
	return d.Direction == DryContactDirectionOutput
}

func (d *DryContact) GetStateDescription() string {
	switch d.Status {
	case DryContactStatusNormal:
		return d.NormalDescription
	case DryContactStatusHighCritical, DryContactStatusLowCritical:
		return d.CriticalDescription
	default:
		return dryContactStatusStrings[d.Status]
	}
}

func (d *DryContact) GetNormalState() string {
	switch d.NormalState {
	case DryContactNormalStateOpen:
		return "open"
	case DryContactNormalStateClosed:
		return "closed"
	default:
		return ""
	}
}

type SensorType int64

const (
	SensorTypeTemperature     SensorType = 1
	SensorTypeHumidityDual    SensorType = 2
	SensorTypeTemperatureDual SensorType = 3
	SensorTypeDryInOut        SensorType = 7
	SensorTypeDryIn           SensorType = 8
	SensorTypeVirtual         SensorType = 129
	SensorTypeHumidity        SensorType = 256
)

type Sensor struct {
	Port        string
	SensorType  SensorType
	Description string
	Virtual     bool
}

func (s *Sensor) GetType() string {
	switch s.SensorType {
	case SensorTypeTemperature:
		return "Temperature"
	case SensorTypeTemperatureDual:
		return "Temperature (dual)"
	case SensorTypeHumidityDual:
		return "Humidity (dual)"
	case SensorTypeDryInOut:
		return "Dry-Contact (in/out)"
	case SensorTypeDryIn:
		return "Dry-Contact (in)"
	case SensorTypeVirtual:
		return "Virtual"
	case SensorTypeHumidity:
		return "Humidity"
	default:
		return "Unsupported"
	}
}

func IsSensorSupported(sensorType SensorType, virtual bool) bool {
	switch sensorType {
	case SensorTypeTemperature, SensorTypeTemperatureDual, SensorTypeHumidityDual, SensorTypeDryInOut, SensorTypeDryIn:
		return true
	case SensorTypeVirtual:
		return virtual
	default:
		return false
	}
}

type Akcp interface {
	GetName() string
	GetLocation() string
	GetDescription() string
	GetOverallSummaryLine() string
	GetTemperatureSensors(snmp *gosnmp.GoSNMP) ([]TemperatureSensor, error)
	GetTemperatureSensor(snmp *gosnmp.GoSNMP, sensorPort string) (*TemperatureSensor, error)
	GetHumiditySensors(snmp *gosnmp.GoSNMP) ([]HumiditySensor, error)
	GetHumiditySensor(snmp *gosnmp.GoSNMP, sensorPort string) (*HumiditySensor, error)
	GetDryContacts(snmp *gosnmp.GoSNMP) ([]DryContact, error)
	GetDryContact(snmp *gosnmp.GoSNMP, sensorPort string) (*DryContact, error)
	ListSensors(snmp *gosnmp.GoSNMP, includeVirtual bool) ([]Sensor, error)
	GetVirtualTemperatureSensors(snmp *gosnmp.GoSNMP) ([]TemperatureSensor, error)
	GetVirtualHumiditySensors(snmp *gosnmp.GoSNMP) ([]HumiditySensor, error)
	GetVirtualDryContacts(snmp *gosnmp.GoSNMP) ([]DryContact, error)
	ValidatePort(sensorPort string) bool
}

type AkcpBase struct {
	model       string
	description string
	name        string
	location    string
}

func NewAkcpBase(model string, description string, name string, location string) AkcpBase {
	return AkcpBase{
		model:       model,
		description: description,
		name:        name,
		location:    location,
	}
}

func (m *AkcpBase) GetDescription() string {
	return m.description
}

func (m *AkcpBase) GetName() string {
	return m.name
}

func (m *AkcpBase) GetLocation() string {
	return m.location
}

func (m *AkcpBase) GetOverallSummaryLine() string {
	var s string = ""
	if len(m.name) > 0 {
		s = fmt.Sprintf("%s %s", m.model, m.name)
	} else {
		s = m.model
	}
	if len(m.location) > 0 {
		s = fmt.Sprintf("%s at %s", s, m.location)
	}
	return fmt.Sprintf("%s (%s)", s, m.description)
}
