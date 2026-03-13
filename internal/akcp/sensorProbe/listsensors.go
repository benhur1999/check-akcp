package sensorProbe

import (
	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/gosnmp/gosnmp"
)

func (m *SensorProbe) ListSensors(snmp *gosnmp.GoSNMP, includeVirtual bool) ([]akcp.Sensor, error) {
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
					if tsensor.Port == hsensor.Port && hsensor.Online {
						t = akcp.SensorTypeTemperatureDual
						break
					}
				}
			}
			result = append(result, akcp.Sensor{
				Port:        tsensor.Port,
				SensorType:  t,
				Description: tsensor.Description,
				Virtual:     false,
			})
		}
	}

	if len(hsensors) > 0 {
		for _, sensor := range hsensors {
			if !sensor.Online {
				continue
			}
			result = append(result, akcp.Sensor{
				Port:        sensor.Port,
				SensorType:  akcp.SensorTypeHumidityDual,
				Description: sensor.Description,
				Virtual:     false,
			})
		}
	}

	if len(dcontacts) > 0 {
		for _, dcontact := range dcontacts {
			if !dcontact.Online || dcontact.Direction != akcp.DryContactDirectionInput {
				continue
			}
			result = append(result, akcp.Sensor{
				Port:        dcontact.Port,
				SensorType:  akcp.SensorTypeDryIn,
				Description: dcontact.Description,
				Virtual:     false,
			})
		}
	}
	return result, nil
}
