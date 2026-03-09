package sensorProbe

import (
	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/gosnmp/gosnmp"
)

func (m *SensorProbe) ListSensors(snmp *gosnmp.GoSNMP) ([]akcp.Sensor, error) {
	var result []akcp.Sensor

	tsensors, err := m.GetTemperatureSensors(snmp)
	if err != nil {
		return nil, err
	}

	hsensors, err := m.GetHumiditySensors(snmp)
	if err != nil {
		return nil, err
	}

	dcontacts, err := m.GetDryContacts(snmp)
	if err != nil {
		return nil, err
	}

	if len(tsensors) > 0 {
		for _, tsensor := range tsensors {
			if !tsensor.Online {
				continue
			}
			t := akcp.SensorTypeTemperature
			// maybe need to adjust sensor type to dual, if humidity sensor shares port with temperature sensor
			if len(hsensors) > 0 {
				for _, hsensor := range hsensors {
					if tsensor.Index == hsensor.Index && hsensor.Online {
						t = akcp.SensorTypeTemperatureDual
						break
					}
				}
			}
			result = append(result, akcp.Sensor{
				Index:       tsensor.Index,
				SensorType:  t,
				Description: tsensor.Description,
			})
		}
	}

	if len(hsensors) > 0 {
		for _, sensor := range hsensors {
			if !sensor.Online {
				continue
			}
			result = append(result, akcp.Sensor{
				Index:       sensor.Index,
				SensorType:  akcp.SensorTypeHumidityDual,
				Description: sensor.Description,
			})
		}
	}

	if len(dcontacts) > 0 {
		for _, dcontact := range dcontacts {
			if !dcontact.Online || dcontact.Direction != akcp.DryContactDirectionInput {
				continue
			}
			result = append(result, akcp.Sensor{
				Index:       dcontact.Index,
				SensorType:  akcp.SensorTypeDryIn,
				Description: dcontact.Description,
			})
		}
	}
	return result, nil
}
