package validation

import (
	"fmt"
	"net"
)

type TypeWithAddress interface {
	ParseAddress() (*net.IP, error)
}

func ValidateIP(path string, ipString string, addressType TypeWithAddress) (*net.IP, error) {
	ip, err := addressType.ParseAddress()
	if err != nil {
		return nil, NewValidationError(path, "Range is not valid.", err)
	}
	if ipString != ip.String() {
		return nil, NewValidationError(path, fmt.Sprintf("Range is not canonical; %s was given but %s was expected.", ipString, ip), nil)
	}
	return ip, nil
}

func ValidateIPRangeContains(childPath string, child *net.IP, parentRange *net.IPNet) error {
	if !parentRange.Contains(*child) {
		return NewValidationError(childPath, fmt.Sprintf("Address %s is not included in parent range %s.", child, parentRange), nil)
	}
	return nil
}
