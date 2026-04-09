package akcputil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/akcp/sensorProbe"
	"github.com/benhur1999/check-akcp/internal/akcp/sensorProbePlus"
	"github.com/benhur1999/check-akcp/internal/snmputil"
	"github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

type AkcpModel int

const (
	AkcpModelAutoDetect AkcpModel = iota
	AkcpModelSensorProbe
	AkcpModelSensorProbePlus
)

const (
	sysDescr    = ".1.3.6.1.2.1.1.1.0"
	sysName     = ".1.3.6.1.2.1.1.5.0"
	sysLocation = ".1.3.6.1.2.1.1.6.0"
)

func New(snmp *gosnmp.GoSNMP, model AkcpModel) (akcp.Akcp, error) {
	result, err := snmp.Get([]string{
		sysDescr,
		sysName,
		sysLocation,
	})
	if err != nil {
		return nil, fmt.Errorf("SNMP failed: %s", err)
	}

	if len(result.Variables) != 3 {
		return nil, errors.New("SNMP failed: did not get all required variables")
	}

	description, _ := snmputil.GetAsString(&result.Variables[0])
	name, _ := snmputil.GetAsString(&result.Variables[1])
	location, _ := snmputil.GetAsString(&result.Variables[2])

	if model == AkcpModelAutoDetect {
		if description == "" {
			return nil, errors.New("Unable to fetch system description")
		}
		log.Debugf("Auto-detecting model from description: %s", description)
		if strings.Contains(description, "sensorProbe") {
			model = AkcpModelSensorProbe
			if snmp.Version != gosnmp.Version1 {
				return nil, fmt.Errorf("SensorProbe models only suppors SNMP version 1 reliably. Current SNMP version is %v",
					snmp.Version)
			}
		} else if strings.Contains(description, "SPX+") {
			model = AkcpModelSensorProbePlus
		} else {
			return nil, fmt.Errorf("Unsupported Model: %s", description)
		}
	}

	// reset values, if user has not changed meaningless default
	if name == "Sys Name" {
		name = ""
	}
	if location == "Sys Location" {
		location = ""
	}

	switch model {
	case AkcpModelSensorProbe:
		return sensorProbe.New(description, name, location), nil
	case AkcpModelSensorProbePlus:
		return sensorProbePlus.New(description, name, location), nil
	default:
		panic(fmt.Sprintf("invalid model code: %d", model))
	}
}
