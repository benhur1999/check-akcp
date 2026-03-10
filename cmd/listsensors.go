package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/benhur1999/check-akcp/internal/akcp/akcputil"
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

	m, err := akcputil.New(snmp, config.GetModel())
	if err != nil {
		return err
	}

	sensors, err := m.ListSensors(snmp)
	if err != nil {
		return err
	}

	if len(sensors) > 0 {
		col_type := 4
		col_port := 4
		col_desc := 4
		for _, sensor := range sensors {
			if len(sensor.GetType()) > col_type {
				col_type = len(sensor.GetType())
			}
			if len(sensor.Index) > col_port {
				col_port = len(sensor.Index)
			}
			if len(sensor.Description) > col_desc {
				col_desc = len(sensor.Description)
			}
		}
		f_str := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds\n", col_type, col_port, col_desc)
		fmt.Printf("The following sensors are found:\n")
		fmt.Printf(f_str, "Type", "Port", "Name")
		fmt.Printf("%s+%s+%s\n", strings.Repeat("-", col_type+1), strings.Repeat("-", col_port+2), strings.Repeat("-", col_desc+1))
		for _, sensor := range sensors {
			fmt.Printf(f_str, sensor.GetType(), sensor.Index, sensor.Description)
		}
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
