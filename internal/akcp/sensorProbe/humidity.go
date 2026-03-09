package sensorProbe

import (
	"fmt"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
)

const (
	sensorProbeHumTable        = ".1.3.6.1.4.1.3854.1.2.2.1.17"
	sensorProbeHumDescription  = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.1"
	sensorProbeHumPercent      = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.3"
	sensorProbeHumStatus       = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.4"
	sensorProbeHumOnline       = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.5"
	sensorProbeHumLowCritical  = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.7"
	sensorProbeHumLowWarning   = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.8"
	sensorProbeHumHighWarning  = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.9"
	sensorProbeHumHighCritical = ".1.3.6.1.4.1.3854.1.2.2.1.17.1.10"

	sensorProbeHumOnlineIsOnline = 1
)

func (m *SensorProbe) GetHumiditySensors(snmp *gosnmp.GoSNMP) ([]akcp.HumiditySensor, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbeHumTable, []string{
		sensorProbeHumDescription,
		sensorProbeHumPercent,
		sensorProbeHumStatus,
		sensorProbeHumOnline,
		sensorProbeHumLowCritical,
		sensorProbeHumLowWarning,
		sensorProbeHumHighWarning,
		sensorProbeHumHighCritical,
	})
	if err != nil {
		return nil, err
	}

	var result []akcp.HumiditySensor
	for _, row := range table {
		desc, _ := row.GetAsString(sensorProbeHumDescription)
		if desc == "" {
			desc = fmt.Sprintf("Humidity %s", row.Index)
		}
		percent, found := row.GetAsFloat64(sensorProbeHumPercent)
		if !found {
			percent = -1
		}
		status, _ := row.GetAsInt64(sensorProbeHumStatus)
		online, _ := row.GetAsInt64(sensorProbeHumOnline)
		lowCritical, _ := row.GetAsFloat64(sensorProbeHumLowCritical)
		lowWarning, _ := row.GetAsFloat64(sensorProbeHumLowWarning)
		highWarning, _ := row.GetAsFloat64(sensorProbeHumHighWarning)
		highCritical, _ := row.GetAsFloat64(sensorProbeHumHighCritical)

		// quirk: set
		s := akcp.SensorStatus(status)
		if online != sensorProbeTempIsOnline {
			s = akcp.StatusNoStatus
		}

		result = append(result, akcp.HumiditySensor{
			Index:        row.Index,
			Description:  desc,
			Percent:      percent,
			Unit:         akcp.HumidityUnitRelativeHumidity,
			LowCritical:  lowCritical,
			LowWarning:   lowWarning,
			HighWarning:  highWarning,
			HighCritical: highCritical,
			Status:       s,
			Online:       (online == sensorProbeHumOnlineIsOnline),
		})
	}
	return result, nil
}

func (m *SensorProbe) GetHumiditySensor(snmp *gosnmp.GoSNMP, sensorPort string) (*akcp.HumiditySensor, error) {
	result, err := snmp.Get([]string{
		snmputil.AppendOid(sensorProbeHumDescription, sensorPort),
		snmputil.AppendOid(sensorProbeHumPercent, sensorPort),
		snmputil.AppendOid(sensorProbeHumStatus, sensorPort),
		snmputil.AppendOid(sensorProbeHumOnline, sensorPort),
		snmputil.AppendOid(sensorProbeHumLowCritical, sensorPort),
		snmputil.AppendOid(sensorProbeHumLowWarning, sensorPort),
		snmputil.AppendOid(sensorProbeHumHighWarning, sensorPort),
		snmputil.AppendOid(sensorProbeHumHighCritical, sensorPort),
	})
	if err != nil {
		return nil, err
	}

	if len(result.Variables) != 8 {
		return nil, nil
	}

	desc, found := snmputil.GetAsString(&result.Variables[0])
	if !found {
		return nil, nil
	}
	if desc == "" {
		desc = fmt.Sprintf("Humidity %s", sensorPort)
	}
	percent, found := snmputil.GetAsFloat64(&result.Variables[1])
	if !found {
		percent = -1
	}
	status, _ := snmputil.GetAsInt64(&result.Variables[2])
	online, _ := snmputil.GetAsInt64(&result.Variables[3])
	lowCritical, _ := snmputil.GetAsFloat64(&result.Variables[4])
	lowWarning, _ := snmputil.GetAsFloat64(&result.Variables[5])
	highWarning, _ := snmputil.GetAsFloat64(&result.Variables[6])
	highCritical, _ := snmputil.GetAsFloat64(&result.Variables[7])

	// quirk: set
	s := akcp.SensorStatus(status)
	if online != sensorProbeTempIsOnline {
		s = akcp.StatusNoStatus
	}

	return &akcp.HumiditySensor{
		Index:        sensorPort,
		Description:  desc,
		Percent:      percent,
		Unit:         akcp.HumidityUnitRelativeHumidity,
		LowCritical:  lowCritical,
		LowWarning:   lowWarning,
		HighWarning:  highWarning,
		HighCritical: highCritical,
		Status:       s,
		Online:       (online == sensorProbeHumOnlineIsOnline),
	}, nil
}
