package sensorProbe

import (
	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

func (m *SensorProbe) GetVirtualTemperatureSensors(snmp *gosnmp.GoSNMP) ([]akcp.TemperatureSensor, error) {
	log.Warn("Virtual sensors are nor supported on sensorProbe models")
	return nil, nil
}

func (m *SensorProbe) GetVirtualHumiditySensors(snmp *gosnmp.GoSNMP) ([]akcp.HumiditySensor, error) {
	log.Warn("Virtual sensors are nor supported on sensorProbe models")
	return nil, nil
}

func (m *SensorProbe) GetVirtualDryContacts(snmp *gosnmp.GoSNMP) ([]akcp.DryContact, error) {
	log.Warn("Virtual sensors are nor supported on sensorProbe models")
	return nil, nil
}
