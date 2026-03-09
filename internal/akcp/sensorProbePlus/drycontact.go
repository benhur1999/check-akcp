package sensorProbePlus

import (
	"fmt"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
)

const (
	sensorProbePlusDryContactTable               = ".1.3.6.1.4.1.3854.3.5.4"
	sensorProbePlusDryContactIndex               = ".1.3.6.1.4.1.3854.3.5.4.1.1"
	sensorProbePlusDryContactDescription         = ".1.3.6.1.4.1.3854.3.5.4.1.2"
	sensorProbePlusDryContactType                = ".1.3.6.1.4.1.3854.3.5.4.1.3"
	sensorProbePlusDryContactStatus              = ".1.3.6.1.4.1.3854.3.5.4.1.6"
	sensorProbePlusDryContactGoOffline           = ".1.3.6.1.4.1.3854.3.5.4.1.8"
	sensorProbePlusDryContactDirection           = ".1.3.6.1.4.1.3854.3.5.4.1.22"
	sensorProbePlusDryContactNormalState         = ".1.3.6.1.4.1.3854.3.5.4.1.23"
	sensorProbePlusDryContactCriticalDescription = ".1.3.6.1.4.1.3854.3.5.4.1.46"
	sensorProbePlusDryContactNormalDescription   = ".1.3.6.1.4.1.3854.3.5.4.1.48"

	sensorProbePlusDryContactGoOfflineOnline = 1
)

func (m *SensorProbePlus) GetDryContacts(snmp *gosnmp.GoSNMP) ([]akcp.DryContact, error) {
	table, err := snmputil.FetchTable(snmp, sensorProbePlusDryContactTable, []string{
		sensorProbePlusDryContactIndex,
		sensorProbePlusDryContactDescription,
		sensorProbePlusDryContactType,
		sensorProbePlusDryContactStatus,
		sensorProbePlusDryContactGoOffline,
		sensorProbePlusDryContactDirection,
		sensorProbePlusDryContactNormalState,
		sensorProbePlusDryContactCriticalDescription,
		sensorProbePlusDryContactNormalDescription,
	})
	if err != nil {
		return nil, err
	}

	var result []akcp.DryContact
	for _, row := range table {
		idx, _ := row.GetAsString(sensorProbePlusDryContactIndex)
		desc, _ := row.GetAsString(sensorProbePlusDryContactDescription)
		t, _ := row.GetAsInt64(sensorProbePlusDryContactType)
		status, _ := row.GetAsInt64(sensorProbePlusDryContactStatus)
		go_offline, _ := row.GetAsInt64(sensorProbePlusDryContactGoOffline)
		direction, _ := row.GetAsInt64(sensorProbePlusDryContactDirection)
		normal_state, _ := row.GetAsFloat64(sensorProbePlusDryContactNormalState)
		critical_decr, _ := row.GetAsString(sensorProbePlusDryContactCriticalDescription)
		normal_decr, _ := row.GetAsString(sensorProbePlusDryContactNormalDescription)
		result = append(result, akcp.DryContact{
			Index:               idx,
			Description:         desc,
			Type:                akcp.DryContactType(t),
			Status:              akcp.DryContactStatus(status),
			Online:              (go_offline == sensorProbePlusDryContactGoOfflineOnline),
			Direction:           akcp.DryContactDirection(direction),
			NormalState:         akcp.DryContactNormalState(normal_state),
			CriticalDescription: critical_decr,
			NormalDescritpion:   normal_decr,
		})
	}
	return result, nil
}

func (m *SensorProbePlus) GetDryContact(snmp *gosnmp.GoSNMP, sensorPort string) (*akcp.DryContact, error) {
	result, err := snmp.Get([]string{
		snmputil.AppendOid(sensorProbePlusDryContactIndex, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactDescription, sensorPort),
		snmputil.AppendOid(sensorProbePlusCommonTableType, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactType, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactStatus, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactGoOffline, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactDirection, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactNormalState, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactCriticalDescription, sensorPort),
		snmputil.AppendOid(sensorProbePlusDryContactNormalDescription, sensorPort),
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
	if sensor_type != int64(akcp.SensorTypeDryInOut) && sensor_type != int64(akcp.SensorTypeDryIn) {
		return nil, nil
	}

	desc, _ := snmputil.GetAsString(&result.Variables[1])
	t, _ := snmputil.GetAsInt64(&result.Variables[3])
	status, _ := snmputil.GetAsInt64(&result.Variables[4])
	go_offline, _ := snmputil.GetAsInt64(&result.Variables[5])
	direction, _ := snmputil.GetAsInt64(&result.Variables[6])
	normal_state, _ := snmputil.GetAsFloat64(&result.Variables[7])
	critical_decr, _ := snmputil.GetAsString(&result.Variables[8])
	normal_decr, _ := snmputil.GetAsString(&result.Variables[9])
	return &akcp.DryContact{
		Index:               idx,
		Description:         desc,
		Type:                akcp.DryContactType(t),
		Status:              akcp.DryContactStatus(status),
		Online:              (go_offline == sensorProbePlusDryContactGoOfflineOnline),
		Direction:           akcp.DryContactDirection(direction),
		NormalState:         akcp.DryContactNormalState(normal_state),
		CriticalDescription: critical_decr,
		NormalDescritpion:   normal_decr,
	}, nil
}
