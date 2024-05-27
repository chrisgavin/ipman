package types

import (
	"net"

	"github.com/chrisgavin/ipman/internal/validation"
	"go.uber.org/multierr"
)

type Pool struct {
	File
	Providers []string `yaml:"providers"`
	Name      string   `yaml:"-"`
	Hosts     []Host   `yaml:"-"`
	Range     string   `yaml:"range"`
	Site      *Site    `yaml:"-"`
}

func (pool *Pool) ParseRange() (*net.IPNet, error) {
	_, ipRange, err := net.ParseCIDR(pool.Range)
	return ipRange, err
}

func (pool *Pool) Validate() error {
	var errors error
	poolRange, err := validation.ValidateIPRange(pool.Path, pool.Range, pool)
	if err != nil {
		errors = multierr.Append(errors, err)
	}
	for _, host := range pool.Hosts {
		errors = multierr.Append(errors, host.Validate())
		if len(host.Interfaces) > 0 {
			primaryInterface := host.Interfaces[0]
			interfaceAddress, err := primaryInterface.ParseAddress()
			if err == nil {
				err = validation.ValidateIPRangeContains(host.Path, interfaceAddress, poolRange)
				if err != nil {
					errors = multierr.Append(errors, err)
				}
			}
		}
	}
	return errors
}
