package lib

import (
	"net"
	"strings"
)

// IPs represent multiple net.IPs
type IPs []*net.IP

// StringWithSep combines all the IPs into a single string with the given seperator
func (ips IPs) StringSep(sep string) string {
	return strings.Join(ips.Strings(), sep)
}

func (ips IPs) Strings() (strings []string) {
	strings = make([]string, len(ips))
	for i, ip := range ips {
		strings[i] = ip.String()
	}
	return
}

func (ips IPs) String() string {
	return ips.StringSep(", ")
}