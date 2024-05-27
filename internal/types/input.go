package types

import "go.uber.org/multierr"

type Input struct {
	File
	Version       int            `yaml:"version"`
	Providers     []string       `yaml:"providers"`
	Networks      []Network      `yaml:"-"`
	DNSProviders  []DNSProvider  `yaml:"-"`
	DHCPProviders []DHCPProvider `yaml:"-"`
}

func (input *Input) Validate() error {
	var errors error
	for _, network := range input.Networks {
		err := network.Validate()
		if err != nil {
			errors = multierr.Append(errors, err)
		}
	}
	return errors
}
