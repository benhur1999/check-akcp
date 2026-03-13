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

type DryContactCmdConfig struct {
	SensorPort string
}

var dryContactsCmdConfig DryContactCmdConfig = DryContactCmdConfig{
	SensorPort: "",
}

var dryContactsCmd = &cobra.Command{
	Use:   "dry-contact",
	Short: "Checks the dry contants",
	RunE:  runDryContactsCmd,
}

func runDryContactsCmd(cmd *cobra.Command, args []string) error {
	snmp, err := NewSnmpClient(true)
	if err != nil {
		return err
	}
	defer snmp.Close()

	m, err := akcputil.New(snmp, config.GetModel())
	if err != nil {
		return err
	}

	if len(dryContactsCmdConfig.SensorPort) > 0 {
		sensorPort := dryContactsCmdConfig.SensorPort
		if !m.ValidatePort(sensorPort) {
			return fmt.Errorf("Invalid sensor port: %s", sensorPort)
		}

		contact, err := m.GetDryContact(snmp, sensorPort)
		if err != nil {
			return err
		}
		if contact == nil {
			return fmt.Errorf("No dry contact on port %s found!", sensorPort)
		}
		log.Debugf("Index: %s, Description: %s, Status: %s, Online: %t IsOutput: %t, NormalState: %s",
			contact.Port, contact.Description, contact.GetStatus(), contact.Online, contact.IsOutput(), contact.GetNormalState())
		if !contact.Online {
			return fmt.Errorf("Dry contact on port %s is offline", sensorPort)
		}
		if contact.IsOutput() {
			return fmt.Errorf("Dry contact on port %s is output", sensorPort)
		}
		rc, output, pd := processDryContact(contact)
		if pd != nil {
			output = fmt.Sprintf("%s\n|%s", output, pd.String())
		}
		check.ExitRaw(rc, output)
	} else {
		overall := result.Overall{
			Summary: m.GetOverallSummaryLine(),
		}

		count, err := processDryContacts(m, snmp, &overall)
		if err != nil {
			return err
		}
		if count == 0 {
			sc := result.PartialResult{
				Output: "No dry contacts found.",
			}
			sc.SetState(check.Unknown)
			overall.AddSubcheck(sc)
		}
		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	}
	return nil
}

func processDryContacts(m akcp.Akcp, snmp *gosnmp.GoSNMP, overall *result.Overall) (int, error) {
	contacts, err := m.GetDryContacts(snmp)
	if err != nil {
		return 0, err
	}

	if config.IncludeVirtual {
		vcontacts, err := m.GetVirtualDryContacts(snmp)
		if err != nil {
			return 0, err
		}
		contacts = append(contacts, vcontacts...)
	}

	count := 0
	for _, contact := range contacts {
		log.Debugf("Index: %s, Description: %s, Status: %s, Online: %t IsOutput: %t, NormalState: %s",
			contact.Port, contact.Description, contact.GetStatus(), contact.Online, contact.IsOutput(), contact.GetNormalState())
		if !contact.Online {
			log.Debug("... skipping offline dry contact")
			continue
		}

		if contact.IsOutput() {
			log.Debug("... skipping output dry contact")
			continue
		}
		count++
		rc, output, pd := processDryContact(&contact)
		sc := result.PartialResult{
			Output: output,
		}
		sc.SetState(rc)
		if pd != nil {
			sc.Perfdata.Add(pd)
		}
		overall.AddSubcheck(sc)
	}
	return count, nil
}

func processDryContact(contact *akcp.DryContact) (int, string, *perfdata.Perfdata) {
	var add_pd bool = false
	var rc int = check.Unknown
	var output string = ""
	var pd *perfdata.Perfdata = nil

	switch contact.Status {
	case akcp.DryContactStatusNormal:
		rc = check.OK
		add_pd = true
	case akcp.DryContactStatusHighCritical, akcp.DryContactStatusLowCritical:
		rc = check.Critical
		add_pd = true
	case akcp.DryContactStatusOutputHigh, akcp.DryContactStatusOutputLow:
		rc = check.OK
		add_pd = true
	case akcp.DryContactStatusSensorError:
		rc = check.Critical
		output = fmt.Sprintf("%s: %s", contact.Description, contact.GetStatus())
	case akcp.DryContactStatusNoStatus:
		rc = check.Unknown
		output = fmt.Sprintf("%s: %s", contact.Description, contact.GetStatus())
	}
	if output == "" {
		output = fmt.Sprintf("%s: %s", contact.Description, contact.GetStateDescription())
		if add_pd && config.Perfdata {
			var state int = 0
			if contact.Status == akcp.DryContactStatusNormal {
				if contact.NormalState == akcp.DryContactNormalStateOpen {
					state = 0
				} else {
					state = 1
				}
			} else {
				if contact.NormalState == akcp.DryContactNormalStateOpen {
					state = 1
				} else {
					state = 0
				}
			}
			pd = &perfdata.Perfdata{
				Label: contact.Description,
				Value: state,
			}
		}
	}
	return rc, output, pd

}

func init() {
	dryContactsCmd.DisableFlagsInUseLine = true
	flags := dryContactsCmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(&dryContactsCmdConfig.SensorPort, "sensor-port", "S", "", "ABC")
	rootCmd.AddCommand(dryContactsCmd)
}
