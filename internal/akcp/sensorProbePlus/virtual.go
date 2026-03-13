package sensorProbePlus

import (
	"strings"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/gosnmp/gosnmp"
)

func (m *SensorProbePlus) GetVirtualTemperatureSensors(snmp *gosnmp.GoSNMP) ([]akcp.TemperatureSensor, error) {
	table, err := fetchCommonSensorTable(snmp)
	if err != nil {
		return nil, err
	}

	var result []akcp.TemperatureSensor
	for _, row := range table {
		st, _ := row.GetAsInt64(sensorProbePlusCommonTableType)
		sensor_type := akcp.SensorType(st)

		if sensor_type != akcp.SensorTypeVirtual {
			continue
		}

		port, _ := row.GetAsString(sensorProbePlusCommonTableIndex)
		desc, _ := row.GetAsString(sensorProbePlusCommonTableDescription)
		unit, _ := row.GetAsString(sensorProbePlusCommonTableUnit)
		value, _ := row.GetAsFloat64(sensorProbePlusCommonTableValue)
		status, _ := row.GetAsInt64(sensorProbePlusCommonTableStatus)
		if value_raw, found := row.GetAsFloat64(sensorProbePlusCommonTableRaw); found {
			value = value_raw / 10
		}
		var u akcp.TemperatureUnit
		switch strings.ToUpper(unit) {
		case "°C":
			u = akcp.TemperatureUnitCelsius
		case "°F":
			u = akcp.TemperatureUnitFahrenheit
		default:
			continue
		}

		result = append(result, akcp.TemperatureSensor{
			Port:        port,
			Description: desc,
			Degree:      value,
			Unit:        u,
			Status:      akcp.SensorStatus(status),
			Online:      true,
		})
	}
	return result, nil
}

func (m *SensorProbePlus) GetVirtualHumiditySensors(snmp *gosnmp.GoSNMP) ([]akcp.HumiditySensor, error) {
	table, err := fetchCommonSensorTable(snmp)
	if err != nil {
		return nil, err
	}
	var result []akcp.HumiditySensor
	for _, row := range table {
		st, _ := row.GetAsInt64(sensorProbePlusCommonTableType)
		sensor_type := akcp.SensorType(st)

		if sensor_type != akcp.SensorTypeVirtual {
			continue
		}

		port, _ := row.GetAsString(sensorProbePlusCommonTableIndex)
		desc, _ := row.GetAsString(sensorProbePlusCommonTableDescription)
		unit, _ := row.GetAsString(sensorProbePlusCommonTableUnit)
		value, _ := row.GetAsFloat64(sensorProbePlusCommonTableValue)
		status, _ := row.GetAsInt64(sensorProbePlusCommonTableStatus)
		if value_raw, found := row.GetAsFloat64(sensorProbePlusCommonTableRaw); found {
			value = value_raw
		}
		var u akcp.HumidityUnit
		switch strings.ToUpper(unit) {
		case "%RH":
			u = akcp.HumidityUnitRelativeHumidity
		default:
			continue
		}

		result = append(result, akcp.HumiditySensor{
			Port:        port,
			Description: desc,
			Percent:     value,
			Unit:        u,
			Status:      akcp.SensorStatus(status),
			Online:      true,
		})
	}
	return result, nil
}

func (m *SensorProbePlus) GetVirtualDryContacts(snmp *gosnmp.GoSNMP) ([]akcp.DryContact, error) {
	table, err := fetchCommonSensorTable(snmp)
	if err != nil {
		return nil, err
	}
	var result []akcp.DryContact
	for _, row := range table {
		st, _ := row.GetAsInt64(sensorProbePlusCommonTableType)
		sensor_type := akcp.SensorType(st)

		if sensor_type != akcp.SensorTypeVirtual {
			continue
		}

		port, _ := row.GetAsString(sensorProbePlusCommonTableIndex)
		desc, _ := row.GetAsString(sensorProbePlusCommonTableDescription)
		unit, _ := row.GetAsString(sensorProbePlusCommonTableUnit)
		status, _ := row.GetAsInt64(sensorProbePlusCommonTableStatus)
		critical_desc, _ := row.GetAsString(sensorProbePlusCommonTableHighCriticalDescription)
		normal_desc, _ := row.GetAsString(sensorProbePlusCommonTableNormalDescription)

		var u akcp.DryContactDirection
		switch strings.ToUpper(unit) {
		case "":
			u = akcp.DryContactDirectionInput
		default:
			continue
		}

		result = append(result, akcp.DryContact{
			Port:                port,
			Description:         desc,
			Type:                akcp.DryContactTypeInput,
			Direction:           u,
			Status:              akcp.DryContactStatus(status),
			Online:              true,
			CriticalDescription: critical_desc,
			NormalDescritpion:   normal_desc,
		})
	}
	return result, nil

}
