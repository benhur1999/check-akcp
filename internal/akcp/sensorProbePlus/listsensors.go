package sensorProbePlus

import (
	"strings"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

const (
	sensorProbePlusCommonTable                        = ".1.3.6.1.4.1.3854.3.5.1"
	sensorProbePlusCommonTableIndex                   = ".1.3.6.1.4.1.3854.3.5.1.1.1"
	sensorProbePlusCommonTableDescription             = ".1.3.6.1.4.1.3854.3.5.1.1.2"
	sensorProbePlusCommonTableType                    = ".1.3.6.1.4.1.3854.3.5.1.1.3"
	sensorProbePlusCommonTableValue                   = ".1.3.6.1.4.1.3854.3.5.1.1.4"
	sensorProbePlusCommonTableUnit                    = ".1.3.6.1.4.1.3854.3.5.1.1.5"
	sensorProbePlusCommonTableStatus                  = ".1.3.6.1.4.1.3854.3.5.1.1.6"
	sensorProbePlusCommonTableRaw                     = ".1.3.6.1.4.1.3854.3.5.1.1.20"
	sensorProbePlusCommonTableHighCriticalDescription = ".1.3.6.1.4.1.3854.3.5.1.1.46"
	sensorProbePlusCommonTableNormalDescription       = ".1.3.6.1.4.1.3854.3.5.1.1.48"
	sensorProbePlusCommonTableGoOffline               = ".1.3.6.1.4.1.3854.3.5.1.1.8"

	sensorProbePlusCommonTableGoOfflineOnline = 1
)

func (m *SensorProbePlus) ListSensors(snmp *gosnmp.GoSNMP, includeVirtual bool) ([]akcp.Sensor, error) {
	table, err := fetchCommonSensorTable(snmp)
	if err != nil {
		return nil, err
	}

	var result []akcp.Sensor
	for _, row := range table {
		idx, _ := row.GetAsString(sensorProbePlusCommonTableIndex)
		desc, _ := row.GetAsString(sensorProbePlusCommonTableDescription)
		st, _ := row.GetAsInt64(sensorProbePlusCommonTableType)
		sensor_type := akcp.SensorType(st)
		unit, _ := row.GetAsString(sensorProbePlusCommonTableUnit)
		go_offline, _ := row.GetAsInt64(sensorProbePlusCommonTableGoOffline)

		log.Debugf("Port: %s, Type: %d, Description: %s, Unit: %s", idx, sensor_type, desc, unit)
		if !akcp.IsSensorSupported(sensor_type, includeVirtual) {
			continue
		}
		if go_offline != sensorProbePlusCommonTableGoOfflineOnline {
			continue
		}

		virtual := false
		if sensor_type == akcp.SensorTypeVirtual {
			virtual = true
			switch strings.ToUpper(unit) {
			case "":
				sensor_type = akcp.SensorTypeDryIn
			case "°C":
				sensor_type = akcp.SensorTypeTemperature
			case "°F":
				sensor_type = akcp.SensorTypeTemperature
			case "%RH":
				sensor_type = akcp.SensorTypeHumidity
			}
		}
		result = append(result, akcp.Sensor{
			Index:       idx,
			Description: desc,
			SensorType:  sensor_type,
			Virtual:     virtual,
		})
	}
	return result, nil
}

func fetchCommonSensorTable(snmp *gosnmp.GoSNMP) (snmputil.Table, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbePlusCommonTable, []string{
		sensorProbePlusCommonTableIndex,
		sensorProbePlusCommonTableDescription,
		sensorProbePlusCommonTableType,
		sensorProbePlusCommonTableUnit,
		sensorProbePlusCommonTableValue,
		sensorProbePlusCommonTableStatus,
		sensorProbePlusCommonTableRaw,
		sensorProbePlusCommonTableHighCriticalDescription,
		sensorProbePlusCommonTableNormalDescription,
		sensorProbePlusCommonTableGoOffline,
	})
	if err != nil {
		return nil, err
	}
	return table, nil
}
