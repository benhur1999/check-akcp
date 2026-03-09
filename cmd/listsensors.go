package cmd

import (
	"fmt"
	"os"

	"github.com/benhur1999/check-akcp/internal/akcp/akcputil"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listSensorsCmd = &cobra.Command{
	Use:   "list-sensors",
	Short: "List all sensors and dry contants",
	RunE:  runListSensorsCmd,
}

func runListSensorsCmd(cmd *cobra.Command, args []string) error {
	snmp, err := NewSnmpClient(true)
	if err != nil {
		return err
	}
	defer snmp.Close()

	m, err := akcputil.New(snmp, akcputil.AkcpModelAutoDetect)
	if err != nil {
		return err
	}

	sensors, err := m.ListSensors(snmp)
	if err != nil {
		return err
	}

	if len(sensors) > 0 {
		table := tablewriter.NewTable(os.Stdout)
		table.Header("SensorId", "Type", "Description")
		for _, sensor := range sensors {
			log.Debugf("Index: %s, Description: %s, Sensor Type: %s (%d)",
				sensor.Index, sensor.Description, sensor.GetType(), sensor.SensorType)
			table.Append([]string{
				sensor.Index, sensor.GetType(), sensor.Description,
			})
		}
		table.Render()
	} else {
		fmt.Printf("No supported sensors found")
		os.Exit(1)
	}
	return nil
}

func init() {
	listSensorsCmd.DisableFlagsInUseLine = true
	rootCmd.AddCommand(listSensorsCmd)
}
