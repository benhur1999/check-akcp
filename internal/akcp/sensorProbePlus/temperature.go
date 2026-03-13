package sensorProbePlus

import (
	"fmt"
	"strings"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
)

const (
	sensorProbePlusTempTable        = ".1.3.6.1.4.1.3854.3.5.2"
	sensorProbePlusTempIndex        = ".1.3.6.1.4.1.3854.3.5.2.1.1"
	sensorProbePlusTempDescription  = ".1.3.6.1.4.1.3854.3.5.2.1.2"
	sensorProbePlusTempDegree       = ".1.3.6.1.4.1.3854.3.5.2.1.4"
	sensorProbePlusTempUnit         = ".1.3.6.1.4.1.3854.3.5.2.1.5"
	sensorProbePlusTempStatus       = ".1.3.6.1.4.1.3854.3.5.2.1.6"
	sensorProbePlusTempOnline       = ".1.3.6.1.4.1.3854.3.5.2.1.8"
	sensorProbePlusTempLowCritical  = ".1.3.6.1.4.1.3854.3.5.2.1.9"
	sensorProbePlusTempLowWarning   = ".1.3.6.1.4.1.3854.3.5.2.1.10"
	sensorProbePlusTempHighWarning  = ".1.3.6.1.4.1.3854.3.5.2.1.11"
	sensorProbePlusTempHighCritical = ".1.3.6.1.4.1.3854.3.5.2.1.12"
	sensorProbePlusTempDegreeRaw    = ".1.3.6.1.4.1.3854.3.5.2.1.20"

	sensorProbePlusTempOnlineIsOnline = 1
)

func (m *SensorProbePlus) GetTemperatureSensors(snmp *gosnmp.GoSNMP) ([]akcp.TemperatureSensor, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbePlusTempTable, []string{
		sensorProbePlusTempIndex,
		sensorProbePlusTempDescription,
		sensorProbePlusTempDegree,
		sensorProbePlusTempUnit,
		sensorProbePlusTempStatus,
		sensorProbePlusTempOnline,
		sensorProbePlusTempLowCritical,
		sensorProbePlusTempLowWarning,
		sensorProbePlusTempHighWarning,
		sensorProbePlusTempHighCritical,
		sensorProbePlusTempDegreeRaw,
	})
	if err != nil {
		return nil, err
	}

	var result []akcp.TemperatureSensor
	for _, row := range table {
		idx, _ := row.GetAsString(sensorProbePlusTempIndex)
		desc, _ := row.GetAsString(sensorProbePlusTempDescription)
		degree, found := row.GetAsFloat64(sensorProbePlusTempDegree)
		if !found {
			degree = -1
		}
		var unit akcp.TemperatureUnit
		u, _ := row.GetAsString(sensorProbePlusTempUnit)
		switch strings.ToUpper(u) {
		case "C":
			unit = akcp.TemperatureUnitCelsius
		case "F":
			unit = akcp.TemperatureUnitFahrenheit
		}
		status, _ := row.GetAsInt64(sensorProbePlusTempStatus)
		online, _ := row.GetAsInt64(sensorProbePlusTempOnline)
		lowCritical := foo(row, sensorProbePlusTempLowCritical)
		lowWarning := foo(row, sensorProbePlusTempLowWarning)
		highWarning := foo(row, sensorProbePlusTempHighWarning)
		highCritical := foo(row, sensorProbePlusTempHighCritical)
		if degreeRaw, ok := row.GetAsFloat64(sensorProbePlusTempDegreeRaw); ok {
			degree = degreeRaw / 10
		}
		result = append(result, akcp.TemperatureSensor{
			Port:         idx,
			Description:  desc,
			Degree:       degree,
			Unit:         unit,
			LowCritical:  lowCritical,
			LowWarning:   lowWarning,
			HighWarning:  highWarning,
			HighCritical: highCritical,
			Status:       akcp.SensorStatus(status),
			Online:       (online == sensorProbePlusTempOnlineIsOnline),
		})
	}
	return result, nil
}

func foo(row *snmputil.Entry, oid string) *float64 {
	if value, ok := row.GetAsFloat64(oid); ok {
		value = value / 10
		return &value
	} else {
		return nil
	}
}

func foo2(pdu *gosnmp.SnmpPDU) *float64 {
	if value, ok := snmputil.GetAsFloat64(pdu); ok {
		value = value / 10
		return &value
	} else {
		return nil
	}
}

func (m *SensorProbePlus) GetTemperatureSensor(snmp *gosnmp.GoSNMP, sensorPort string) (*akcp.TemperatureSensor, error) {
	result, err := snmp.Get([]string{
		snmputil.AppendOid(sensorProbePlusTempIndex, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempDescription, sensorPort),
		snmputil.AppendOid(sensorProbePlusCommonTableType, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempDegree, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempUnit, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempStatus, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempOnline, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempLowCritical, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempLowWarning, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempHighWarning, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempHighCritical, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempDegreeRaw, sensorPort),
	})
	if err != nil {
		return nil, fmt.Errorf("SNMP failed: %s", err)
	}

	port, found := snmputil.GetAsString(&result.Variables[0])
	if !found {
		return nil, nil
	}

	sensor_type, found := snmputil.GetAsInt64(&result.Variables[2])
	if !found {
		return nil, nil
	}
	if sensor_type != int64(akcp.SensorTypeTemperature) && sensor_type != int64(akcp.SensorTypeTemperatureDual) {
		return nil, nil
	}

	desc, _ := snmputil.GetAsString(&result.Variables[1])
	degree, found := snmputil.GetAsFloat64(&result.Variables[2])
	if !found {
		degree = -1
	}

	var unit akcp.TemperatureUnit
	u, _ := snmputil.GetAsString(&result.Variables[4])
	switch strings.ToUpper(u) {
	case "C":
		unit = akcp.TemperatureUnitCelsius
	case "F":
		unit = akcp.TemperatureUnitFahrenheit
	}

	status, _ := snmputil.GetAsInt64(&result.Variables[5])
	online, _ := snmputil.GetAsInt64(&result.Variables[6])
	lowCritical := foo2(&result.Variables[7])
	lowWarning := foo2(&result.Variables[8])
	highWarning := foo2(&result.Variables[9])
	highCritical := foo2(&result.Variables[10])
	if degreeRaw, ok := snmputil.GetAsFloat64(&result.Variables[11]); ok {
		degree = degreeRaw / 10
	}

	return &akcp.TemperatureSensor{
		Port:         port,
		Description:  desc,
		Degree:       degree,
		Unit:         unit,
		LowCritical:  lowCritical,
		LowWarning:   lowWarning,
		HighWarning:  highWarning,
		HighCritical: highCritical,
		Status:       akcp.SensorStatus(status),
		Online:       (online == sensorProbePlusTempOnlineIsOnline),
	}, nil
}
