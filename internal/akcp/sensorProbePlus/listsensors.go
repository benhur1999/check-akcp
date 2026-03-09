package sensorProbePlus

import (
	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
)

const (
	sensorProbePlusCommonTable            = ".1.3.6.1.4.1.3854.3.5.1"
	sensorProbePlusCommonTableIndex       = ".1.3.6.1.4.1.3854.3.5.1.1.1"
	sensorProbePlusCommonTableDescription = ".1.3.6.1.4.1.3854.3.5.1.1.2"
	sensorProbePlusCommonTableType        = ".1.3.6.1.4.1.3854.3.5.1.1.3"
	sensorProbePlusCommonTableGoOffline   = ".1.3.6.1.4.1.3854.3.5.1.1.8"

	sensorProbePlusCommonTableGoOfflineOnline = 1
)

func (m *SensorProbePlus) ListSensors(snmp *gosnmp.GoSNMP) ([]akcp.Sensor, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbePlusCommonTable, []string{
		sensorProbePlusCommonTableIndex,
		sensorProbePlusCommonTableDescription,
		sensorProbePlusCommonTableType,
		sensorProbePlusCommonTableGoOffline,
	})
	if err != nil {
		return nil, err
	}

	var result []akcp.Sensor
	for _, row := range table {
		idx, _ := row.GetAsString(sensorProbePlusCommonTableIndex)
		desc, _ := row.GetAsString(sensorProbePlusCommonTableDescription)
		sensor_type, _ := row.GetAsInt64(sensorProbePlusCommonTableType)
		go_offline, _ := row.GetAsInt64(sensorProbePlusCommonTableGoOffline)

		if !akcp.IsSensorSupported(akcp.SensorType(sensor_type)) {
			continue
		}
		if go_offline != sensorProbePlusCommonTableGoOfflineOnline {
			continue
		}
		result = append(result, akcp.Sensor{
			Index:       idx,
			Description: desc,
			SensorType:  akcp.SensorType(sensor_type),
		})
	}
	return result, nil

}
