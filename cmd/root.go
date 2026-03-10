package cmd

import (
	"os"

	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

var rootCmd = &cobra.Command{
	Use:              "check_akcp",
	Short:            "Check plugin for the AKCP SensorProbe series",
	PersistentPreRun: preRunCmd,
	Run:              usage,
}

func Execute() {
	// make sure we catch panics and exit with unknown
	defer check.CatchPanic()

	if err := rootCmd.Execute(); err != nil {
		check.Exitf(check.Unknown, "%s", err)
	}
}

func preRunCmd(cmd *cobra.Command, args []string) {
	// setup logging
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	log.SetOutput(os.Stdout)

	if config.Debug > 0 {
		log.SetLevel(log.DebugLevel)
		log.Debugf("setting debug level to %d", config.Debug)
	} else if config.Verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	// validate flags
	validateGlobalArguments()

	// arm timeout
	go check.HandleTimeout(config.Timeout)
}

func usage(cmd *cobra.Command, args []string) {
	_ = cmd.Usage()

	os.Exit(check.Unknown)
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableAutoGenTag = true
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.DisableSuggestions = true

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	pfs := rootCmd.PersistentFlags()
	pfs.SortFlags = false
	pfs.StringVarP(&config.Hostname, "hostname", "H", "", "Host name or IP Address")
	pfs.Uint16VarP(&config.Port, "port", "p", 161, "Port number")
	pfs.VarP(enumflag.NewWithoutDefault(&config.SnmpVersion, "protocol", SnmpVersionIds, enumflag.EnumCaseSensitive), "protocol", "P", "SNMP version to use [1|2c|3]")
	pfs.StringVarP(&config.SnmpCommunity, "community", "c", "public", "SNMPv1/SNMPv2c community string")
	pfs.StringVarP(&config.SnmpV3Username, "username", "U", "", "SNMPv3 username")
	pfs.StringVarP(&config.SnmpV3Context, "context", "N", "", "SNMPv3 context")
	pfs.StringVarP(&config.SnmpV3SecurityLevel, "seclevel", "L", "", "SNMPv3 security level [noAuthNoPriv|authNoPriv|authPriv]")
	pfs.StringVarP(&config.SnmpV3AuthProtocol, "authproto", "a", "sha", "SNMPv3 authentication password [md5|sha]")
	pfs.StringVarP(&config.SnmpV3AuthPassphrase, "authpass", "A", "", "SNMPv3 authentication protocol")
	pfs.StringVarP(&config.SnmpV3PrivProtocol, "privproto", "x", "des", "SNMPv3 privacy proto [des|aes]")
	pfs.StringVarP(&config.SnmpV3PrivPassphrase, "privpass", "X", "", "SNMPv3 privacy password")
	pfs.StringVarP(&config.AuthFileName, "auth-file", "F", "", "Authentication configuration file")
	pfs.CountVarP(&config.Debug, "debug", "d", "Enable debug mode")
	pfs.BoolVarP(&config.Verbose, "verbose", "v", false, "Enable verbose mode")
	pfs.BoolVar(&config.Perfdata, "perf-data", false, "Output performance data")
	pfs.IntVarP(&config.Timeout, "timeout", "t", 30, "Abort the check after n seconds")
	pfs.VarP(enumflag.New(&config.Model, "model", ConfigModelIds, enumflag.EnumCaseInsensitive), "model", "M",
		"Model [auto, sensorProbe, sensorProbePlus]")
	rootCmd.Flags().SortFlags = false
}
