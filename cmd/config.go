package cmd

import (
	"fmt"
	"time"

	"github.com/NETWAYS/go-check"
	"github.com/benhur1999/check-akcp/internal/akcp/akcputil"
	"github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
	"github.com/thediveo/enumflag/v2"
)

type ConfigModel enumflag.Flag

const (
	ConfigModelAutoDetect ConfigModel = iota
	ConfogModelSensorProbe
	ConfogModelSensorProbePlus
)

var ConfigModelIds = map[ConfigModel][]string{
	ConfigModelAutoDetect:      {"auto"},
	ConfogModelSensorProbe:     {"sensorProbe"},
	ConfogModelSensorProbePlus: {"sensorProbePlus"},
}

var SnmpVersionIds = map[gosnmp.SnmpVersion][]string{
	gosnmp.Version1:  {"1"},
	gosnmp.Version2c: {"2c"},
	gosnmp.Version3:  {"3"},
}

type Config struct {
	Hostname             string
	Port                 uint16
	SnmpVersion          gosnmp.SnmpVersion
	SnmpCommunity        string
	SnmpV3Username       string
	SnmpV3Context        string
	SnmpV3SecurityLevel  string
	SnmpV3AuthProtocol   string
	SnmpV3AuthPassphrase string
	SnmpV3PrivProtocol   string
	SnmpV3PrivPassphrase string
	Username             string
	Password             string
	AuthFileName         string
	Timeout              int
	Verbose              bool
	Debug                int
	Perfdata             bool
	Model                ConfigModel
}

var config Config = Config{
	Model:       ConfigModelAutoDetect,
	SnmpVersion: gosnmp.Version2c,
}

func validateGlobalArguments() {
	if len(config.Hostname) == 0 {
		check.ExitRaw(check.Unknown, "Hostname is required")
	}

	if config.Port < 1 {
		check.ExitRaw(check.Unknown, "Port must be between 1 and 65535")
	}

	// switch config.SnmpVersion {
	// case "1":
	// 	fallthrough
	// case "2c":
	// 	if len(config.SnmpCommunity) == 0 {
	// 		check.ExitRaw(check.Unknown, "Community string is required for SNMP version 1 and 2c")
	// 	}
	// case "3":
	// 	if len(config.SnmpV3Username) == 0 {
	// 		check.ExitRaw(check.Unknown, "Username required for SNMP version 3")
	// 	}
	// 	if (config.SnmpV3AuthProtocol != "md5") && (config.SnmpV3AuthProtocol != "sha") {
	// 		check.ExitRaw(check.Unknown, "Invalid value for SNMP version 3 authentication protocol")
	// 	}
	// 	if len(config.SnmpV3AuthPassphrase) == 0 {
	// 		check.ExitRaw(check.Unknown, "Authentication passphrase required for SNMP version 3")
	// 	}

	// 	// check if priv proto is set
	// 	if (config.SnmpV3PrivProtocol != "des") && (config.SnmpV3PrivProtocol != "aes") {
	// 		check.ExitRaw(check.Unknown, "Invalid value for SNMP version 3 privacy protocol")
	// 	}
	// 	if len(config.SnmpV3PrivPassphrase) == 0 {
	// 		check.ExitRaw(check.Unknown, "Privacy passphrase required for SNMP version 3")
	// 	}
	// default:
	// 	check.ExitRaw(check.Unknown, "SNMP version must be either \"1\", \"2c\" or \"3\"")
	// }

}

func (c *Config) GetModel() akcputil.AkcpModel {
	switch c.Model {
	case ConfigModelAutoDetect:
		return akcputil.AkcpModelAutoDetect
	case ConfogModelSensorProbe:
		return akcputil.AkcpModelSensorProbe
	case ConfogModelSensorProbePlus:
		return akcputil.AkcpModelSensorProbePlus
	default:
		panic(fmt.Sprintf("invalid config model: %d", c.Model))
	}
}

func NewSnmpClient(connect bool) (*gosnmp.GoSNMP, error) {
	retires := 3
	timeout := config.Timeout / retires

	var snmp *gosnmp.GoSNMP

	switch config.SnmpVersion {
	case gosnmp.Version1:
		fallthrough
	case gosnmp.Version2c:
		snmp = &gosnmp.GoSNMP{
			Port:               config.Port,
			Transport:          "udp",
			Community:          config.SnmpCommunity,
			Version:            config.SnmpVersion,
			Timeout:            time.Duration(timeout) * time.Second,
			Retries:            retires,
			ExponentialTimeout: false,
			MaxOids:            gosnmp.MaxOids,
			Target:             config.Hostname,
		}
	case gosnmp.Version3:
		var securityPareters gosnmp.UsmSecurityParameters

		securityPareters.UserName = config.SnmpV3Username

		switch config.SnmpV3AuthProtocol {
		case "md5":
			securityPareters.AuthenticationProtocol = gosnmp.MD5
		case "sha":
			securityPareters.AuthenticationProtocol = gosnmp.SHA
		}
		securityPareters.AuthenticationPassphrase = config.SnmpV3AuthPassphrase

		switch config.SnmpV3PrivProtocol {
		case "des":
			securityPareters.PrivacyProtocol = gosnmp.DES
		case "aes":
			securityPareters.PrivacyProtocol = gosnmp.AES
		}
		securityPareters.PrivacyPassphrase = config.SnmpV3PrivPassphrase

		snmp = &gosnmp.GoSNMP{
			Port:               config.Port,
			Transport:          "udp",
			Version:            config.SnmpVersion,
			Timeout:            time.Duration(timeout) * time.Second,
			Retries:            retires,
			ExponentialTimeout: false,
			MaxOids:            gosnmp.MaxOids,
			Target:             config.Hostname,
			SecurityModel:      gosnmp.UserSecurityModel,
			MsgFlags:           gosnmp.AuthPriv,
			SecurityParameters: &securityPareters,
			ContextName:        config.SnmpV3Context,
		}
	default:
		return nil, nil
	}

	if connect {
		log.Debugf("connecting to %s using version %v ...", snmp.Target, snmp.Version)
		err := snmp.Connect()
		if err != nil {
			log.Fatalf("connect error: %v", err)
			return nil, err
		}
	}
	return snmp, nil
}
