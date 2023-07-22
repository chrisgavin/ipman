package registry

import (
	"github.com/chrisgavin/ipman/internal/providers/dns"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/pkg/errors"
)

func NewDNSProvider(typeName string) (types.DNSProvider, error) {
	switch typeName {
	case "dns.NullProvider":
		return &dns.NullProvider{}, nil
	case "dns.CloudflareProvider":
		return &dns.CloudflareProvider{}, nil
	}
	return nil, errors.New("Unknown type " + typeName + ".")
}
