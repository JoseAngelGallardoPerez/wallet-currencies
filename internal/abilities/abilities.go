// package for defining and controlling abilities
package abilities

import (
	"github.com/Confialink/wallet-users/rpc/proto/users"

	"github.com/Confialink/wallet-currencies/internal/abilities/actions"
	"github.com/Confialink/wallet-currencies/internal/abilities/resolvers"
	"github.com/Confialink/wallet-currencies/internal/abilities/resources"
	"github.com/Confialink/wallet-currencies/internal/abilities/roles"
)

type abilityFunc func(user *users.User) bool

var abilitiesList = map[string]map[string]map[string]abilityFunc{
	roles.Root: {
		resources.Currency: {
			actions.Read:   resolvers.Allow,
			actions.Index:  resolvers.Allow,
			actions.Create: resolvers.Allow,
		},
		resources.FullCurrency: {
			actions.Read:   resolvers.Allow,
			actions.Index:  resolvers.Allow,
			actions.Update: resolvers.Allow,
		},
		resources.Settings: {
			actions.Read:   resolvers.Allow,
			actions.Update: resolvers.Allow,
		},
		resources.Rates: {
			actions.Index:  resolvers.Allow,
			actions.Update: resolvers.Allow,
			actions.Read:   resolvers.Allow,
		},
	},
	roles.Admin: {
		resources.Currency: {
			actions.Read:  resolvers.Allow,
			actions.Index: resolvers.Allow,
		},
		resources.FullCurrency: {
			actions.Read:   resolvers.CanViewSettings,
			actions.Index:  resolvers.CanViewSettings,
			actions.Update: resolvers.CanEditSettings,
		},
		resources.Settings: {
			actions.Read:   resolvers.CanViewSettings,
			actions.Update: resolvers.CanEditSettings,
		},
		resources.Rates: {
			actions.Index:  resolvers.CanViewSettings,
			actions.Update: resolvers.CanEditSettings,
			actions.Read:   resolvers.Allow,
		},
	},
	roles.User: {
		resources.Currency: {
			actions.Read:  resolvers.Allow,
			actions.Index: resolvers.Allow,
		},
		resources.Rates: {
			actions.Read:  resolvers.Allow,
			actions.Index: resolvers.Allow,
		},
	},
}
