package types

import (
	"fmt"

	"github.com/chrisgavin/ipman/internal/validation"
	"go.uber.org/multierr"
)

type Host struct {
	File
	Name       string      `yaml:"-"`
	Interfaces []Interface `yaml:"interfaces"`
}

func (host *Host) Validate() error {
	var errors error
	for _, networkInterface := range host.Interfaces {
		if _, err := networkInterface.ParseAddress(); err != nil {
			errors = multierr.Append(errors, validation.NewValidationError(host.Path, fmt.Sprintf("Interface %s does not have a valid IP address.", networkInterface.Name), err))
		}
	}
	return errors
}
