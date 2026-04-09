# check_akcp

Monitoring plugin (Icinga/Nagios compatible) to monitor AKCP sensorProbe and sensorProbeX+ devices. It
supports temperature, humidity and dry-contact sensors.

## Usage
```
Usage:
  check_akcp [flags]
  check_akcp [command]

Available Commands:
  dry-contact  Checks the dry contants
  humidity     Checks the humidity sensors
  list-sensors List all sensors and dry contacts
  sensors      Checks all sensors and dry contacts
  temperature  Checks the temperature sensors

Flags:
  -H, --hostname string     Host name or IP Address
  -p, --port uint16         Port number (default 161)
  -P, --protocol protocol   SNMP version to use [1|2c|3] (default 2c)
  -c, --community string    SNMPv1/SNMPv2c community string (default "public")
  -U, --username string     SNMPv3 username
  -N, --context string      SNMPv3 context
  -L, --seclevel string     SNMPv3 security level [noAuthNoPriv|authNoPriv|authP                                                 riv]
  -a, --authproto string    SNMPv3 authentication password [md5|sha] (default "s                                                 ha")
  -A, --authpass string     SNMPv3 authentication protocol
  -x, --privproto string    SNMPv3 privacy proto [des|aes] (default "des")
  -X, --privpass string     SNMPv3 privacy password
  -F, --auth-file string    Authentication configuration file
  -d, --debug count         Enable debug mode
  -v, --verbose             Enable verbose mode
      --perf-data           Output performance data
  -t, --timeout int         Abort the check after n seconds (default 30)
  -M, --model model         Model [auto, sensorProbe, sensorProbePlus] (default                                                  auto)
      --virtual             Include virtual sensors (only sensorProbe+ models)
  -h, --help                help for check_akcp
```

The plugin tries to auto-detect the model unless overriden by the command line.

## License

Copyright (c) 2026 Leibniz-Institut für Deutsche Sprache (IDS)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not,
see [gnu.org/licenses](https://www.gnu.org/licenses/).
