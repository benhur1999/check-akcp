package cmd

import (
	"fmt"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/benhur1999/check-akcp/internal/akcp/akcputil"
	"github.com/spf13/cobra"
)

const (
	selectedSensorsTemperature = 1
	selectedSensorsHumidity    = 2
	selectedSensorsDryContacts = 4
	selectedSensorsAll         = selectedSensorsTemperature | selectedSensorsHumidity | selectedSensorsDryContacts
)

type SensorsConfig struct {
	SelectedSensorsArgs []string
	SelectedSensors     int
}

var sensorCmdConfig SensorsConfig = SensorsConfig{}

var selectedTypesIds = map[string]int{
	"all":         selectedSensorsAll,
	"temperature": selectedSensorsTemperature,
	"humidity":    selectedSensorsHumidity,
	"dry-contact": selectedSensorsDryContacts,
}

var sensorsCmd = &cobra.Command{
	Use:     "sensors",
	Short:   "Checks all sensors and dry contants",
	PreRunE: preRunAllSensorsCmd,
	RunE:    runAllSensorsCmd,
}

func preRunAllSensorsCmd(cmd *cobra.Command, args []string) error {
	var selectedSensors int = 0
	for _, valueStr := range sensorCmdConfig.SelectedSensorsArgs {
		value, ok := selectedTypesIds[valueStr]
		if !ok {
			return fmt.Errorf("Invalid value for option \"-T\": %s", valueStr)
		}
		selectedSensors |= value
	}
	sensorCmdConfig.SelectedSensors = selectedSensors
	return nil
}

func runAllSensorsCmd(cmd *cobra.Command, args []string) error {
	snmp, err := NewSnmpClient(true)
	if err != nil {
		return err
	}
	defer snmp.Close()

	m, err := akcputil.New(snmp, config.GetModel())
	if err != nil {
		return err
	}

	overall := result.Overall{
		Summary: m.GetOverallSummaryLine(),
	}

	count := 0
	if sensorCmdConfig.SelectedSensors&selectedSensorsTemperature > 0 {
		c, err := processTemperatureSensors(m, snmp, &overall)
		if err != nil {
			return nil
		}
		count += c
	}

	if sensorCmdConfig.SelectedSensors&selectedSensorsHumidity > 0 {
		c, err := processHumiditySensors(m, snmp, &overall)
		if err != nil {
			return err
		}
		count += c
	}

	if sensorCmdConfig.SelectedSensors&selectedSensorsDryContacts > 0 {
		c, err := processDryContacts(m, snmp, &overall)
		if err != nil {
			return err
		}
		count += c
	}
	if count == 0 {
		sc := result.PartialResult{
			Output: "No sensors found.",
		}
		sc.SetState(check.Unknown)
		overall.AddSubcheck(sc)
	}

	check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	return nil
}

func init() {
	sensorsCmd.DisableFlagsInUseLine = true
	flags := sensorsCmd.Flags()
	flags.SortFlags = false
	flags.StringSliceVarP(&sensorCmdConfig.SelectedSensorsArgs, "sensor-type", "T", []string{"all"},
		"Selected sensor types [all, temperature, humidity, dry-contact]")
	rootCmd.AddCommand(sensorsCmd)
}
