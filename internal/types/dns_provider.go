package types

import (
	"github.com/chrisgavin/ipman/internal/actions"
)

type DNSProvider interface {
	GetActions(network Network, site Site, pool Pool, hosts []Host) ([]actions.DNSAction, error)
}
