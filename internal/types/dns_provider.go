package types

import (
	"context"

	"github.com/chrisgavin/ipman/internal/actions"
)

type DNSProvider interface {
	GetName(ctx context.Context) string
	GetActions(ctx context.Context, network Network, site Site, pool Pool, hosts []Host) ([]actions.DNSAction, error)
	ApplyAction(ctx context.Context, action actions.DNSAction) error
}
