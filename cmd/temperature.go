package cmd

import (
	"fmt"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/akcp/akcputil"
	"github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TemperatureCmdConfig struct {
	SensorPort string
}

var temperatureCmdConfig TemperatureCmdConfig = TemperatureCmdConfig{
	SensorPort: "",
}

var temperatureCmd = &cobra.Command{
	Use:   "temperature",
	Short: "Checks the temperature sensors",
	RunE:  runTemperatureCmd,
}

func runTemperatureCmd(cmd *cobra.Command, args []string) error {
	snmp, err := NewSnmpClient(true)
	if err != nil {
		return err
	}
	defer snmp.Close()

	m, err := akcputil.New(snmp, config.GetModel())
	if err != nil {
		return err
	}

	if len(temperatureCmdConfig.SensorPort) > 0 {
		sensorPort := temperatureCmdConfig.SensorPort
		if !m.ValidatePort(sensorPort) {
			return fmt.Errorf("Invalid sensor port: %s", sensorPort)
		}

		sensor, err := m.GetTemperatureSensor(snmp, sensorPort)
		if err != nil {
			return err
		}
		if sensor == nil {
			return fmt.Errorf("No temperature sensor on port %s found!", sensorPort)
		}

		log.Debugf("Index: %s, Description: %s, Status: %s, Online: %t, Degree: %.1f %s [%.0f, %.0f, %.0f, %.0f]",
			sensor.Index, sensor.Description, sensor.GetStatus(), sensor.Online, sensor.Degree, sensor.GetUnit(),
			sensor.LowCritical, sensor.LowWarning, sensor.HighWarning, sensor.HighCritical)
		if !sensor.Online {
			return fmt.Errorf("Temperature sensor on port %s is offline!", sensorPort)
		}

		rc, output, pd := processTemperatureSensor(sensor)
		if pd != nil {
			output = fmt.Sprintf("%s\n|%s", output, pd.String())
		}
		check.ExitRaw(rc, output)
	} else {
		overall := result.Overall{
			Summary: m.GetOverallSummaryLine(),
		}
		count, err := processTemperatureSensors(m, snmp, &overall)
		if err != nil {
			return err
		}
		if count == 0 {
			sc := result.PartialResult{
				Output: "No temperature sensors found.",
			}
			sc.SetState(check.Unknown)
			overall.AddSubcheck(sc)
		}

		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	}
	return nil
}

func processTemperatureSensors(m akcp.Akcp, snmp *gosnmp.GoSNMP, overall *result.Overall) (int, error) {
	sensors, err := m.GetTemperatureSensors(snmp)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, sensor := range sensors {
		log.Debugf("Index: %s, Description: %s, Status: %s, Online: %t, Degree: %.1f %s [%.0f, %.0f, %.0f, %.0f]",
			sensor.Index, sensor.Description, sensor.GetStatus(), sensor.Online, sensor.Degree, sensor.GetUnit(),
			sensor.LowCritical, sensor.LowWarning, sensor.HighWarning, sensor.HighCritical)
		if !sensor.Online {
			log.Debug("... skipping offline sensor")
			continue
		}

		rc, output, pd := processTemperatureSensor(&sensor)
		sc := result.PartialResult{
			Output: output,
		}
		sc.SetState(rc)
		if pd != nil {
			sc.Perfdata.Add(pd)
		}
		overall.AddSubcheck(sc)
		count++
	}
	return count, nil
}

func processTemperatureSensor(sensor *akcp.TemperatureData) (int, string, *perfdata.Perfdata) {
	var add_pd bool = false
	var rc int = check.Unknown
	var output string = ""
	var pd *perfdata.Perfdata = nil

	switch sensor.Status {
	case akcp.StatusNormal:
		rc = check.OK
		add_pd = true
	case akcp.StatusLowWarning, akcp.StatusHighWarning:
		rc = check.Warning
		add_pd = true
	case akcp.StatusLowCritical, akcp.StatusHighCritical:
		rc = check.Critical
		add_pd = true
	case akcp.StatusSensorError:
		output = fmt.Sprintf("%s: %s", sensor.Description, sensor.GetStatus())
		rc = check.Critical
	case akcp.StatusNoStatus:
		output = fmt.Sprintf("%s: %s", sensor.Description, sensor.GetStatus())
		rc = check.Unknown
	}
	if output == "" {
		output = fmt.Sprintf("%s: %.1f °%s", sensor.Description, sensor.Degree, sensor.GetUnit())
		if add_pd && config.Perfdata {
			pd = &perfdata.Perfdata{
				Label: sensor.Description,
				Value: sensor.Degree,
				Uom:   sensor.GetUnit(),
				Warn:  &check.Threshold{Lower: sensor.LowWarning, Upper: sensor.HighWarning},
				Crit:  &check.Threshold{Lower: sensor.LowCritical, Upper: sensor.HighCritical},
			}
		}
	}
	return rc, output, pd
}

func init() {
	temperatureCmd.DisableFlagsInUseLine = true
	flags := temperatureCmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(&temperatureCmdConfig.SensorPort, "sensor-port", "S", "", "ABC")
	rootCmd.AddCommand(temperatureCmd)
}
