package types

import (
	"context"

	"github.com/chrisgavin/ipman/internal/actions"
)

type DHCPProvider interface {
	GetName(ctx context.Context) string
	GetActions(ctx context.Context, network Network, site Site, pool Pool, hosts []Host) ([]actions.DHCPAction, error)
	ApplyAction(ctx context.Context, action actions.DHCPAction) error
}
