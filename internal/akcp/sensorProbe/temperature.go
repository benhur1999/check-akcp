package sensorProbe

import (
	"fmt"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
)

const (
	sensorProbeTempTable            = ".1.3.6.1.4.1.3854.1.2.2.1.16"
	sensorProbeTempDescription      = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.1"
	sensorProbeTempDegree           = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.3"
	sensorProbeTempStatus           = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.4"
	sensorProbeTempOnline           = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.5"
	sensorProbePlusTempHighWarning  = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.7"
	sensorProbePlusTempHighCritical = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.8"
	sensorProbePlusTempLowCritical  = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.9"
	sensorProbePlusTempLowWarning   = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.10"
	sensorProbeTempDegreeType       = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.12"
	sensorProbeTempDegreeRaw        = ".1.3.6.1.4.1.3854.1.2.2.1.16.1.14"
)

const (
	sensorProbeTempIsOnline         int64 = 1
	sensorProbeDegreeTypeFahrenheit int64 = 0
	sensorProbeDegreeTypeCelsius    int64 = 1
)

func (m *SensorProbe) GetTemperatureSensors(snmp *gosnmp.GoSNMP) ([]akcp.TemperatureSensor, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbeTempTable, []string{
		sensorProbeTempDescription,
		sensorProbeTempDegree,
		sensorProbeTempStatus,
		sensorProbeTempOnline,
		sensorProbePlusTempHighWarning,
		sensorProbePlusTempHighCritical,
		sensorProbePlusTempLowCritical,
		sensorProbePlusTempLowWarning,
		sensorProbeTempDegreeType,
		sensorProbeTempDegreeRaw,
	})
	if err != nil {
		return nil, err
	}

	var result []akcp.TemperatureSensor
	for _, row := range table {
		desc, _ := row.GetAsString(sensorProbeTempDescription)
		if desc == "" {
			desc = fmt.Sprintf("Temperature %s", row.Index)
		}
		degree, found := row.GetAsFloat64(sensorProbeTempDegree)
		if !found {
			degree = -1
		}
		status, _ := row.GetAsInt64(sensorProbeTempStatus)
		online, _ := row.GetAsInt64(sensorProbeTempOnline)
		lowCritical, _ := row.GetAsFloat64(sensorProbePlusTempLowCritical)
		lowWarning, _ := row.GetAsFloat64(sensorProbePlusTempLowWarning)
		highWarning, _ := row.GetAsFloat64(sensorProbePlusTempHighWarning)
		highCritical, _ := row.GetAsFloat64(sensorProbePlusTempHighCritical)

		var unit akcp.TemperatureUnit
		u, _ := row.GetAsInt64(sensorProbeTempDegreeType)
		switch u {
		case sensorProbeDegreeTypeCelsius:
			unit = akcp.TemperatureUnitCelsius
		case sensorProbeDegreeTypeFahrenheit:
			unit = akcp.TemperatureUnitFahrenheit
		}
		degreeRaw, found := row.GetAsFloat64(sensorProbeTempDegreeRaw)
		if found {
			degree = degreeRaw / 10
		}

		// quirk: override status fpr offline sensor
		s := akcp.SensorStatus(status)
		if online != sensorProbeTempIsOnline {
			s = akcp.StatusNoStatus
		}

		result = append(result, akcp.TemperatureSensor{
			Index:        row.Index,
			Description:  desc,
			Degree:       degree,
			Unit:         unit,
			Status:       s,
			Online:       (online == sensorProbeTempIsOnline),
			LowCritical:  &lowCritical,
			LowWarning:   &lowWarning,
			HighWarning:  &highWarning,
			HighCritical: &highCritical,
		})
	}
	return result, nil
}

func (m *SensorProbe) GetTemperatureSensor(snmp *gosnmp.GoSNMP, sensorPort string) (*akcp.TemperatureSensor, error) {
	result, err := snmp.Get([]string{
		snmputil.AppendOid(sensorProbeTempDescription, sensorPort),
		snmputil.AppendOid(sensorProbeTempDegree, sensorPort),
		snmputil.AppendOid(sensorProbeTempStatus, sensorPort),
		snmputil.AppendOid(sensorProbeTempOnline, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempHighWarning, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempHighCritical, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempLowCritical, sensorPort),
		snmputil.AppendOid(sensorProbePlusTempLowWarning, sensorPort),
		snmputil.AppendOid(sensorProbeTempDegreeType, sensorPort),
		snmputil.AppendOid(sensorProbeTempDegreeRaw, sensorPort),
	})
	if err != nil {
		return nil, err
	}

	if len(result.Variables) != 10 {
		return nil, nil
	}

	desc, found := snmputil.GetAsString(&result.Variables[0])
	if !found {
		return nil, nil
	}
	if desc == "" {
		desc = fmt.Sprintf("Temperature %s", sensorPort)
	}
	degree, found := snmputil.GetAsFloat64(&result.Variables[1])
	if !found {
		degree = -1
	}
	status, _ := snmputil.GetAsInt64(&result.Variables[2])
	online, _ := snmputil.GetAsInt64(&result.Variables[3])
	lowCritical, _ := snmputil.GetAsFloat64(&result.Variables[4])
	lowWarning, _ := snmputil.GetAsFloat64(&result.Variables[5])
	highWarning, _ := snmputil.GetAsFloat64(&result.Variables[6])
	highCritical, _ := snmputil.GetAsFloat64(&result.Variables[7])
	var unit akcp.TemperatureUnit
	u, _ := snmputil.GetAsInt64(&result.Variables[8])
	switch u {
	case sensorProbeDegreeTypeCelsius:
		unit = akcp.TemperatureUnitCelsius
	case sensorProbeDegreeTypeFahrenheit:
		unit = akcp.TemperatureUnitFahrenheit
	}
	degreeRaw, found := snmputil.GetAsFloat64(&result.Variables[9])
	if found {
		degree = degreeRaw / 10
	}

	// quirk: override status fpr offline sensor
	s := akcp.SensorStatus(status)
	if online != sensorProbeTempIsOnline {
		s = akcp.StatusNoStatus
	}

	return &akcp.TemperatureSensor{
		Index:        sensorPort,
		Description:  desc,
		Degree:       degree,
		Unit:         unit,
		Status:       s,
		Online:       (online == sensorProbeTempIsOnline),
		LowCritical:  &lowCritical,
		LowWarning:   &lowWarning,
		HighWarning:  &highWarning,
		HighCritical: &highCritical,
	}, nil
}
