package snmputil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

type Entry struct {
	Index string
	PDUs  []*gosnmp.SnmpPDU
}

type Table []*Entry

func FetchTable(snmp *gosnmp.GoSNMP, tableOid string, oids []string) (Table, error) {
	log.Debugf("fetching table '%s' (matching %d oids) ...", tableOid, len(oids))

	results := make(Table, 0, 8)
	var walkFn gosnmp.WalkFunc = func(pdu gosnmp.SnmpPDU) error {
		if !strings.HasPrefix(pdu.Name, tableOid) {
			return nil
		}

		for _, oid := range oids {
			m_oid := oid + "."
			if strings.HasPrefix(pdu.Name, m_oid) {
				if len(pdu.Name) < len(m_oid) {
					return errors.New("too short")
				}
				idx := pdu.Name[len(m_oid):]
				if len(idx) == 0 {
					return errors.New("no index")
				}

				var entry *Entry = nil
				for _, i := range results {
					if i.Index == idx {
						entry = i
						break
					}
				}
				if entry == nil {
					entry = &Entry{
						Index: idx,
						PDUs:  make([]*gosnmp.SnmpPDU, 0, len(oids)),
					}
					results = append(results, entry)
				}

				entry.PDUs = append(entry.PDUs, &gosnmp.SnmpPDU{
					Value: pdu.Value,
					Name:  oid,
					Type:  pdu.Type,
				})
				break
			}
		}
		return nil
	}

	var err error
	if snmp.Version == gosnmp.Version1 {
		err = snmp.Walk(tableOid, walkFn)
	} else {
		err = snmp.BulkWalk(tableOid, walkFn)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (entry *Entry) GetAsString(oid string) (string, bool) {
	pdu := entry.findPdu(oid)
	return GetAsString(pdu)
}

func (entry *Entry) GetAsInt64(oid string) (int64, bool) {
	pdu := entry.findPdu(oid)
	return GetAsInt64(pdu)
}

func (entry *Entry) GetAsFloat64(oid string) (float64, bool) {
	pdu := entry.findPdu(oid)
	return GetAsFloat64(pdu)
}

func (entry *Entry) findPdu(oid string) *gosnmp.SnmpPDU {
	for _, pdu := range entry.PDUs {
		if pdu.Name == oid {
			return pdu
		}
	}
	return nil
}

func IsFound(pdu *gosnmp.SnmpPDU) bool {
	if pdu != nil {
		return pdu.Type != gosnmp.NoSuchInstance && pdu.Type != gosnmp.NoSuchObject
	}
	return false
}

func GetAsString(pdu *gosnmp.SnmpPDU) (string, bool) {
	if pdu != nil && IsFound(pdu) && pdu.Type == gosnmp.OctetString {
		b := pdu.Value.([]byte)
		return string(b), true
	}
	return "", false
}

func GetAsInt64(pdu *gosnmp.SnmpPDU) (int64, bool) {
	if pdu != nil && IsFound(pdu) && pdu.Type != gosnmp.OctetString {
		v := gosnmp.ToBigInt(pdu.Value)
		if v.IsInt64() {
			return v.Int64(), true
		}
	}
	return -1, false
}

func GetAsFloat64(pdu *gosnmp.SnmpPDU) (float64, bool) {
	if pdu != nil && IsFound(pdu) && pdu.Type != gosnmp.OctetString {
		v := gosnmp.ToBigInt(pdu.Value)
		if v.IsInt64() {
			return float64(v.Int64()), true
		}
	}
	return -1, false
}

func AppendOid(baseId string, oid string) string {
	return fmt.Sprintf("%s.%s", baseId, oid)
}
