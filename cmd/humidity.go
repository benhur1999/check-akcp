package cmd

import (
	"fmt"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/benhur1999/check-akcp/internal/akcp"
	"github.com/benhur1999/check-akcp/internal/akcp/akcputil"
	"github.com/benhur1999/check-akcp/internal/utils"
	"github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type HumidityCmdConfig struct {
	SensorPort string
}

var humidityCmdConfig HumidityCmdConfig = HumidityCmdConfig{
	SensorPort: "",
}

var humidityCmd = &cobra.Command{
	Use:   "humidity",
	Short: "Checks the humidity sensors",
	RunE:  runHumidityCmd,
}

func runHumidityCmd(cmd *cobra.Command, args []string) error {
	snmp, err := NewSnmpClient(true)
	if err != nil {
		return err
	}
	defer snmp.Close()

	m, err := akcputil.New(snmp, config.GetModel())
	if err != nil {
		return err
	}

	if len(humidityCmdConfig.SensorPort) > 0 {
		sensorPort := humidityCmdConfig.SensorPort
		if !m.ValidatePort(sensorPort) {
			return fmt.Errorf("Invalid sensor port: %s", sensorPort)
		}

		sensor, err := m.GetHumiditySensor(snmp, sensorPort)
		if err != nil {
			return err
		}
		if sensor == nil {
			return fmt.Errorf("No humiditiy sensor on port %s found!", sensorPort)
		}

		log.Debugf("Index: %s, Description: %s, Status: %s, Online: %t, Percent: %.0f %s [%.0f, %.0f, %.0f, %.0f]",
			sensor.Index, sensor.Description, sensor.GetStatus(), sensor.Online, sensor.Percent, sensor.GetUnit(),
			sensor.LowCritical, sensor.LowWarning, sensor.HighWarning, sensor.HighCritical)
		if !sensor.Online {
			return fmt.Errorf("Humiditiy sensor on %s is offline!", sensorPort)
		}

		rc, output, pd := processHumiditySensor(sensor)
		if pd != nil {
			output = fmt.Sprintf("%s\n|%s", output, pd.String())
		}
		check.ExitRaw(rc, output)
	} else {
		overall := result.Overall{
			Summary: m.GetOverallSummaryLine(),
		}

		count, err := processHumiditySensors(m, snmp, &overall)
		if err != nil {
			return err
		}
		if count == 0 {
			sc := result.PartialResult{
				Output: "No humidity sensors found.",
			}
			sc.SetState(check.Unknown)
			overall.AddSubcheck(sc)
		}
		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	}
	return nil
}

func processHumiditySensors(m akcp.Akcp, snmp *gosnmp.GoSNMP, overall *result.Overall) (int, error) {
	sensors, err := m.GetHumiditySensors(snmp)
	if err != nil {
		return 0, err
	}

	if config.IncludeVirtual {
		vsensors, err := m.GetVirtualHumiditySensors(snmp)
		if err != nil {
			return 0, err
		}
		sensors = append(sensors, vsensors...)
	}

	count := 0
	for _, sensor := range sensors {
		log.Debugf("Index: %s, Description: %s, Status: %s, Online: %t, Percent: %.0f %s [%s, %s, %s, %s]",
			sensor.Index, sensor.Description, sensor.GetStatus(), sensor.Online, sensor.Percent, sensor.GetUnit(),
			utils.FormatFloat(sensor.LowCritical), utils.FormatFloat(sensor.LowWarning),
			utils.FormatFloat(sensor.HighWarning), utils.FormatFloat(sensor.HighCritical))
		if !sensor.Online {
			log.Debug("... skipping offline sensor")
			continue
		}

		rc, output, pd := processHumiditySensor(&sensor)
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

func processHumiditySensor(sensor *akcp.HumiditySensor) (int, string, *perfdata.Perfdata) {
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
		output = fmt.Sprintf("%s: %.0f %s", sensor.Description, sensor.Percent, sensor.GetUnit())
		if add_pd && config.Perfdata {
			pd = &perfdata.Perfdata{
				Label: sensor.Description,
				Value: sensor.Percent,
				Uom:   sensor.GetUnit(),
				Warn:  utils.MakeThreashold(sensor.LowWarning, sensor.HighWarning),
				Crit:  utils.MakeThreashold(sensor.LowCritical, sensor.HighCritical),
			}
		}
	}
	return rc, output, pd

}

func init() {
	humidityCmd.DisableFlagsInUseLine = true
	flags := humidityCmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(&humidityCmdConfig.SensorPort, "sensor-port", "S", "", "ABC")
	rootCmd.AddCommand(humidityCmd)
}
