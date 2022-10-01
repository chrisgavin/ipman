package types

import (
	"net"

	"github.com/chrisgavin/ipman/internal/validation"
	"go.uber.org/multierr"
)

type Network struct {
	File
	Name  string `yaml:"-"`
	Sites []Site `yaml:"-"`
	Range string `yaml:"range"`
}

func (network *Network) ParseRange() (*net.IPNet, error) {
	_, ipRange, err := net.ParseCIDR(network.Range)
	return ipRange, err
}

func (network *Network) Validate() error {
	var errors error
	networkRange, err := validation.ValidateIPRange(network.Path, network.Range, network)
	if err != nil {
		errors = multierr.Append(errors, err)
	}
	for _, site := range network.Sites {
		errors = multierr.Append(errors, site.Validate())
		siteRange, _ := site.ParseRange()
		err = validation.ValidateIPRangeFullyContains(site.Path, siteRange, networkRange)
		if err != nil {
			errors = multierr.Append(errors, err)
		}
	}
	return errors
}
