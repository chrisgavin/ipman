package types

import (
	"context"

	"github.com/chrisgavin/ipman/internal/actions"
)

type DNSProvider interface {
	GetActions(ctx context.Context, network Network, site Site, pool Pool, hosts []Host) ([]actions.DNSAction, error)
}
