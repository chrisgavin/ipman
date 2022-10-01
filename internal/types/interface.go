package types

import (
	"net"

	"github.com/pkg/errors"
)

const AddressTypeLeased = "leased"
const AddressTypeCNAME = "leased"

type Interface struct {
	File
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	MAC     string `yaml:"mac"`
}

func (networkInterface *Interface) ParseAddress() (*net.IP, error) {
	ip := net.ParseIP(networkInterface.Address)
	if ip != nil {
		return &ip, nil
	}
	return nil, errors.New("IP address cannot be parsed.")
}
