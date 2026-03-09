package main

import (
	"github.com/benhur1999/check-akcp/cmd"
)

func main() {
	cmd.Execute()
	// err := flags.Parse(os.Args)
	// if err != nil {
	// 	check.Exitf(check.Unknown, "%e", err)
	// }

	// // setup logging
	// log.SetFormatter(&log.TextFormatter{
	// 	ForceColors: true,
	// })
	// log.SetOutput(os.Stdout)

	// if config.Debug > 0 {
	// 	log.SetLevel(log.DebugLevel)
	// 	log.Debugf("setting debug level to %d", config.Debug)
	// } else if config.Verbose {
	// 	log.SetLevel(log.InfoLevel)
	// } else {
	// 	log.SetLevel(log.WarnLevel)
	// }

	// validateGlobalArguments()

	// log.Debugf("Hostname: %s", config.Hostname)

	// // arm timeout
	// go check.HandleTimeout(config.Timeout)

	// client, err := NewSnmpClient(true)
	// if err != nil {
	// 	check.Exitf(check.Unknown, "error: %s", err)
	// }
	// defer client.Close()

	// log.Infof("Model: %v", config.Model)
	// akcp.DetectModel(client)

	// plugin := check.NewConfig()
	// plugin.Name = "check_akcp"
	// plugin.Version = "0.1.0"
	// plugin.Readme = "Check plugin for the AKCP Sensorprobe X plus"

	// // setup command line flags
	// //flags := pflag.NewFlagSet("check_akcp", pflag.ContinueOnError)
	// flags := plugin.FlagSet
	// flags.StringVarP(&config.Hostname, "hostname", "H", "", "Host name or IP Address")
	// flags.Uint16VarP(&config.Port, "port", "p", 161, "Port number")
	// flags.StringVarP(&config.SnmpProtocol, "protocol", "P", "2c", "SNMP version to use [1|2c|3]")
	// flags.StringVarP(&config.SnmpCommunity, "community", "c", "public", "SNMPv1/SNMPv2c community string")
	// flags.StringVarP(&config.SnmpV3Username, "username", "U", "", "SNMPv3 username")
	// flags.StringVarP(&config.SnmpV3Context, "context", "N", "", "SNMPv3 context")
	// flags.StringVarP(&config.SnmpV3SecurityLevel, "seclevel", "L", "", "SNMPv3 security level [noAuthNoPriv|authNoPriv|authPriv]")
	// flags.StringVarP(&config.SnmpV3AuthProtocol, "authproto", "a", "sha", "SNMPv3 authentication password [md5|sha]")
	// flags.StringVarP(&config.SnmpV3AuthPassphrase, "authpass", "A", "", "SNMPv3 authentication protocol")
	// flags.StringVarP(&config.SnmpV3PrivProtocol, "privproto", "x", "des", "SNMPv3 privacy proto [des|aes]")
	// flags.StringVarP(&config.SnmpV3PrivPassphrase, "privpass", "X", "", "SNMPv3 privacy password")
	// flags.StringVarP(&config.AuthFileName, "auth-file", "F", "", "Authentication configuration file")
	// //flags.CountVarP(&config.Debug, "debug", "d", "Enable debug mode")
	// //flags.BoolVarP(&config.Verbose, "verbose", "v", false, "Enable verbose mode")
	// flags.BoolVar(&config.Perfdata, "perf-data", false, "Output performance data")
	// //flags.IntVarP(&config.Timeout, "timeout", "t", 30, "Abort the check after n seconds")
	// flags.VarP(enumflag.New(&config.Model, "model", ConfigModelIds, enumflag.EnumCaseInsensitive), "model", "M", "Model [auto, ]")

	// plugin.ParseArguments()

}
