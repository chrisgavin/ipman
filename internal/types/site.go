package types

import (
	"net"

	"github.com/chrisgavin/ipman/internal/validation"
	"go.uber.org/multierr"
)

type Site struct {
	File
	Providers []string `yaml:"providers"`
	Name      string   `yaml:"-"`
	Pools     []Pool   `yaml:"-"`
	Range     string   `yaml:"range"`
	Network   *Network `yaml:"-"`
}

func (site *Site) ParseRange() (*net.IPNet, error) {
	_, ipRange, err := net.ParseCIDR(site.Range)
	return ipRange, err
}

func (site *Site) Validate() error {
	var errors error
	siteRange, err := validation.ValidateIPRange(site.Path, site.Range, site)
	if err != nil {
		errors = multierr.Append(errors, err)
	}
	for _, pool := range site.Pools {
		errors = multierr.Append(errors, pool.Validate())
		poolRange, _ := pool.ParseRange()
		err = validation.ValidateIPRangeFullyContains(pool.Path, poolRange, siteRange)
		if err != nil {
			errors = multierr.Append(errors, err)
		}
	}
	return errors
}
