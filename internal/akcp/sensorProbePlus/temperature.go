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

func (m *SensorProbePlus) GetTemperatureSensors(snmp *gosnmp.GoSNMP) ([]akcp.TemperatureData, error) {
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

	var result []akcp.TemperatureData
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
		lowCritical, _ := row.GetAsFloat64(sensorProbePlusTempLowCritical)
		lowWarning, _ := row.GetAsFloat64(sensorProbePlusTempLowWarning)
		highWarning, _ := row.GetAsFloat64(sensorProbePlusTempHighWarning)
		highCritical, _ := row.GetAsFloat64(sensorProbePlusTempHighCritical)
		degreeRaw, found := row.GetAsFloat64(sensorProbePlusTempDegreeRaw)
		if found {
			degree = degreeRaw / 10
		}
		result = append(result, akcp.TemperatureData{
			Index:        idx,
			Description:  desc,
			Degree:       degree,
			Unit:         unit,
			LowCritical:  lowCritical / 10,
			LowWarning:   lowWarning / 10,
			HighWarning:  highWarning / 10,
			HighCritical: highCritical / 10,
			Status:       akcp.SensorStatus(status),
			Online:       (online == sensorProbePlusTempOnlineIsOnline),
		})
	}
	return result, nil
}

func (m *SensorProbePlus) GetTemperatureSensor(snmp *gosnmp.GoSNMP, sensorPort string) (*akcp.TemperatureData, error) {
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

	idx, found := snmputil.GetAsString(&result.Variables[0])
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
	lowCritical, _ := snmputil.GetAsFloat64(&result.Variables[7])
	lowWarning, _ := snmputil.GetAsFloat64(&result.Variables[8])
	highWarning, _ := snmputil.GetAsFloat64(&result.Variables[9])
	highCritical, _ := snmputil.GetAsFloat64(&result.Variables[10])
	degreeRaw, found := snmputil.GetAsFloat64(&result.Variables[11])
	if found {
		degree = degreeRaw / 10
	}

	return &akcp.TemperatureData{
		Index:        idx,
		Description:  desc,
		Degree:       degree,
		Unit:         unit,
		LowCritical:  lowCritical / 10,
		LowWarning:   lowWarning / 10,
		HighWarning:  highWarning / 10,
		HighCritical: highCritical / 10,
		Status:       akcp.SensorStatus(status),
		Online:       (online == sensorProbePlusTempOnlineIsOnline),
	}, nil
}
