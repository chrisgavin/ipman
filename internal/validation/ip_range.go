package validation

import (
	"fmt"
	"net"
)

type TypeWithRange interface {
	ParseRange() (*net.IPNet, error)
}

func ValidateIPRange(path string, rangeString string, rangeType TypeWithRange) (*net.IPNet, error) {
	ipRange, err := rangeType.ParseRange()
	if err != nil {
		return nil, NewValidationError(path, "Range is not valid.", err)
	}
	if rangeString != ipRange.String() {
		return nil, NewValidationError(path, fmt.Sprintf("Range is not canonical; %s was given but %s was expected.", rangeString, ipRange), nil)
	}
	return ipRange, nil
}

func ValidateIPRangeFullyContains(childPath string, childRange *net.IPNet, parentRange *net.IPNet) error {
	if parentRange != nil && childRange != nil {
		childSize, _ := childRange.Mask.Size()
		parentSize, _ := parentRange.Mask.Size()
		if childSize < parentSize || !parentRange.Contains(childRange.IP) {
			return NewValidationError(childPath, fmt.Sprintf("Range %s is not included in parent range %s.", childRange, parentRange), nil)
		}
	}
	return nil
}
