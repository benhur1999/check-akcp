package sensorProbePlus

import (
	"fmt"
	"strings"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
)

const (
	sensorProbePlusHumTable        = ".1.3.6.1.4.1.3854.3.5.3"
	sensorProbePlusHumIndex        = ".1.3.6.1.4.1.3854.3.5.3.1.1"
	sensorProbePlusHumDescription  = ".1.3.6.1.4.1.3854.3.5.3.1.2"
	sensorProbePlusHumPercent      = ".1.3.6.1.4.1.3854.3.5.3.1.4"
	sensorProbePlusHumUnit         = ".1.3.6.1.4.1.3854.3.5.3.1.5"
	sensorProbePlusHumStatus       = ".1.3.6.1.4.1.3854.3.5.3.1.6"
	sensorProbePlusHumGoOffline    = ".1.3.6.1.4.1.3854.3.5.3.1.8"
	sensorProbePlusHumLowCritical  = ".1.3.6.1.4.1.3854.3.5.3.1.9"
	sensorProbePlusHumLowWarning   = ".1.3.6.1.4.1.3854.3.5.3.1.10"
	sensorProbePlusHumHighWarning  = ".1.3.6.1.4.1.3854.3.5.3.1.11"
	sensorProbePlusHumHighCritical = ".1.3.6.1.4.1.3854.3.5.3.1.12"

	sensorProbePlusHumGoOfflineOnline = 1
)

func (m *SensorProbePlus) GetHumiditySensors(snmp *gosnmp.GoSNMP) ([]akcp.HumiditySensor, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbePlusHumTable, []string{
		sensorProbePlusHumIndex,
		sensorProbePlusHumDescription,
		sensorProbePlusHumPercent,
		sensorProbePlusHumUnit,
		sensorProbePlusHumStatus,
		sensorProbePlusHumGoOffline,
		sensorProbePlusHumLowCritical,
		sensorProbePlusHumLowWarning,
		sensorProbePlusHumHighWarning,
		sensorProbePlusHumHighCritical,
	})
	if err != nil {
		return nil, err
	}

	var result []akcp.HumiditySensor
	for _, row := range table {
		idx, _ := row.GetAsString(sensorProbePlusHumIndex)
		desc, _ := row.GetAsString(sensorProbePlusHumDescription)
		percent, found := row.GetAsFloat64(sensorProbePlusHumPercent)
		if !found {
			percent = -1
		}
		var unit akcp.HumidityUnit
		u, _ := row.GetAsString(sensorProbePlusHumUnit)
		switch strings.ToUpper(u) {
		case "%":
			unit = akcp.HumidityUnitRelativeHumidity
		}
		status, _ := row.GetAsInt64(sensorProbePlusHumStatus)
		go_offline, _ := row.GetAsInt64(sensorProbePlusHumGoOffline)
		lowCritical, _ := row.GetAsFloat64(sensorProbePlusHumLowCritical)
		lowWarning, _ := row.GetAsFloat64(sensorProbePlusHumLowWarning)
		highWarning, _ := row.GetAsFloat64(sensorProbePlusHumHighWarning)
		highCritical, _ := row.GetAsFloat64(sensorProbePlusHumHighCritical)
		result = append(result, akcp.HumiditySensor{
			Index:        idx,
			Description:  desc,
			Percent:      percent,
			Unit:         unit,
			LowCritical:  &lowCritical,
			LowWarning:   &lowWarning,
			HighWarning:  &highWarning,
			HighCritical: &highCritical,
			Status:       akcp.SensorStatus(status),
			Online:       (go_offline == sensorProbePlusHumGoOfflineOnline),
		})
	}
	return result, nil
}

func (m *SensorProbePlus) GetHumiditySensor(snmp *gosnmp.GoSNMP, sensorPort string) (*akcp.HumiditySensor, error) {
	result, err := snmp.Get([]string{
		snmputil.AppendOid(sensorProbePlusHumIndex, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumDescription, sensorPort),
		snmputil.AppendOid(sensorProbePlusCommonTableType, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumPercent, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumUnit, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumStatus, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumGoOffline, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumLowCritical, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumLowWarning, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumHighWarning, sensorPort),
		snmputil.AppendOid(sensorProbePlusHumHighCritical, sensorPort),
	})
	if err != nil {
		return nil, fmt.Errorf("SNMP failed: %s", err)
	}

	idx, found := snmputil.GetAsString(&result.Variables[0])
	if !found {
		return nil, nil
	}

	desc, _ := snmputil.GetAsString(&result.Variables[1])
	sensor_type, found := snmputil.GetAsInt64(&result.Variables[2])
	if !found {
		return nil, nil
	}

	if sensor_type != int64(akcp.SensorTypeHumidityDual) {
		return nil, nil
	}

	percent, found := snmputil.GetAsFloat64(&result.Variables[3])
	if !found {
		percent = -1
	}

	var unit akcp.HumidityUnit
	u, _ := snmputil.GetAsString(&result.Variables[4])
	switch strings.ToUpper(u) {
	case "%":
		unit = akcp.HumidityUnitRelativeHumidity
	}

	status, _ := snmputil.GetAsInt64(&result.Variables[5])
	go_offline, _ := snmputil.GetAsInt64(&result.Variables[6])
	lowCritical, _ := snmputil.GetAsFloat64(&result.Variables[7])
	lowWarning, _ := snmputil.GetAsFloat64(&result.Variables[8])
	highWarning, _ := snmputil.GetAsFloat64(&result.Variables[9])
	highCritical, _ := snmputil.GetAsFloat64(&result.Variables[10])

	return &akcp.HumiditySensor{
		Index:        idx,
		Description:  desc,
		Percent:      percent,
		Unit:         unit,
		LowCritical:  &lowCritical,
		LowWarning:   &lowWarning,
		HighWarning:  &highWarning,
		HighCritical: &highCritical,
		Status:       akcp.SensorStatus(status),
		Online:       (go_offline == sensorProbePlusHumGoOfflineOnline),
	}, nil

}
