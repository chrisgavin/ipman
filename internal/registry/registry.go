package registry

import (
	"github.com/chrisgavin/ipman/internal/providers/dhcp"
	"github.com/chrisgavin/ipman/internal/providers/dns"
	"github.com/chrisgavin/ipman/internal/types"
	"github.com/pkg/errors"
)

func NewDNSProvider(typeName string, name string) (types.DNSProvider, error) {
	switch typeName {
	case "dns.NullProvider":
		return &dns.NullProvider{Name: name}, nil
	case "dns.CloudflareProvider":
		return &dns.CloudflareProvider{Name: name}, nil
	}
	return nil, errors.Errorf("Unknown type %s for provider %s.", typeName, name)
}

func NewDHCPProvider(typeName string, name string) (types.DHCPProvider, error) {
	switch typeName {
	case "dhcp.NullProvider":
		return &dhcp.NullProvider{Name: name}, nil
	case "dhcp.RouterOSProvider":
		return &dhcp.RouterOSProvider{Name: name}, nil
	case "dhcp.EdgeOSProvider":
		return &dhcp.EdgeOSProvider{Name: name}, nil
	}
	return nil, errors.Errorf("Unknown type %s for provider %s.", typeName, name)
}
