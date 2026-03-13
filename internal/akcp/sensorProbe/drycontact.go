package sensorProbe

import (
	"fmt"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
)

const (
	sensorProbeSwitchTable       = ".1.3.6.1.4.1.3854.1.2.2.1.18"
	sensorProbeSwitchDescription = ".1.3.6.1.4.1.3854.1.2.2.1.18.1.1"
	sensorProbeSwitchStatus      = ".1.3.6.1.4.1.3854.1.2.2.1.18.1.3"
	sensorProbeSwitchOnline      = ".1.3.6.1.4.1.3854.1.2.2.1.18.1.4"
	sensorProbeSwitchDirection   = ".1.3.6.1.4.1.3854.1.2.2.1.18.1.6"
	sensorProbeSwitchNormalState = ".1.3.6.1.4.1.3854.1.2.2.1.18.1.7"
	sensorProbeSwitchSensorType  = ".1.3.6.1.4.1.3854.1.2.2.1.18.1.9"

	sensorProbeOnlineOnline    = 1
	sensorProbeDirectionInput  = 0
	sensorProbeDirectionOutput = 1

	sensorProbeSensorTypeDryContact = 10
)

func (m *SensorProbe) GetDryContacts(snmp *gosnmp.GoSNMP) ([]akcp.DryContact, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbeSwitchTable, []string{
		sensorProbeSwitchDescription,
		sensorProbeSwitchStatus,
		sensorProbeSwitchOnline,
		sensorProbeSwitchDirection,
		sensorProbeSwitchNormalState,
	})
	if err != nil {
		return nil, err
	}

	var result []akcp.DryContact
	for _, row := range table {
		desc, _ := row.GetAsString(sensorProbeSwitchDescription)
		if desc == "" {
			desc = fmt.Sprintf("Switch %s", row.Index)
		}
		status, _ := row.GetAsInt64(sensorProbeSwitchStatus)
		online, _ := row.GetAsInt64(sensorProbeSwitchOnline)
		direction, _ := row.GetAsInt64(sensorProbeSwitchDirection)
		normal_state, _ := row.GetAsFloat64(sensorProbeSwitchNormalState)

		var t akcp.DryContactType
		switch direction {
		case sensorProbeDirectionInput:
			t = akcp.DryContactTypeInput
		case sensorProbeDirectionOutput:
			t = akcp.DryContactTypeOutput
		}
		result = append(result, akcp.DryContact{
			Port:                row.Index,
			Description:         desc,
			Type:                akcp.DryContactType(t),
			Status:              akcp.DryContactStatus(status),
			Online:              (online == sensorProbeOnlineOnline),
			Direction:           akcp.DryContactDirection(direction),
			NormalState:         akcp.DryContactNormalState(normal_state),
			CriticalDescription: "Critical",
			NormalDescription:   "Normal",
		})
	}
	return result, nil
}

func (m *SensorProbe) GetDryContact(snmp *gosnmp.GoSNMP, sensorPort string) (*akcp.DryContact, error) {
	result, err := snmp.Get([]string{
		snmputil.AppendOid(sensorProbeSwitchDescription, sensorPort),
		snmputil.AppendOid(sensorProbeSwitchStatus, sensorPort),
		snmputil.AppendOid(sensorProbeSwitchOnline, sensorPort),
		snmputil.AppendOid(sensorProbeSwitchDirection, sensorPort),
		snmputil.AppendOid(sensorProbeSwitchNormalState, sensorPort),
		snmputil.AppendOid(sensorProbeSwitchSensorType, sensorPort),
	})
	if err != nil {
		return nil, fmt.Errorf("SNMP failed: %s", err)
	}

	if len(result.Variables) != 6 {
		return nil, nil
	}

	sensor_type, found := snmputil.GetAsInt64(&result.Variables[5])
	if !found || sensor_type != sensorProbeSensorTypeDryContact {
		return nil, nil
	}

	desc, found := snmputil.GetAsString(&result.Variables[0])
	if !found {
		return nil, nil
	}
	if desc == "" {
		desc = fmt.Sprintf("Switch %s", sensorPort)
	}
	status, _ := snmputil.GetAsInt64(&result.Variables[1])
	online, _ := snmputil.GetAsInt64(&result.Variables[2])
	direction, _ := snmputil.GetAsInt64(&result.Variables[3])
	normal_state, _ := snmputil.GetAsFloat64(&result.Variables[4])

	var t akcp.DryContactType
	switch direction {
	case sensorProbeDirectionInput:
		t = akcp.DryContactTypeInput
	case sensorProbeDirectionOutput:
		t = akcp.DryContactTypeOutput
	}
	return &akcp.DryContact{
		Port:                sensorPort,
		Description:         desc,
		Type:                akcp.DryContactType(t),
		Status:              akcp.DryContactStatus(status),
		Online:              (online == sensorProbeOnlineOnline),
		Direction:           akcp.DryContactDirection(direction),
		NormalState:         akcp.DryContactNormalState(normal_state),
		CriticalDescription: "Critical",
		NormalDescription:   "Normal",
	}, nil
}
