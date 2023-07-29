package registry

import (
	"github.com/chrisgavin/ipman/internal/providers/dhcp"
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

func NewDHCPProvider(typeName string) (types.DHCPProvider, error) {
	switch typeName {
	case "dhcp.NullProvider":
		return &dhcp.NullProvider{}, nil
	case "dhcp.RouterOSProvider":
		return &dhcp.RouterOSProvider{}, nil
	}
	return nil, errors.New("Unknown type " + typeName + ".")
}
