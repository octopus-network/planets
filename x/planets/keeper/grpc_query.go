package keeper

import (
	"github.com/octopus-network/planets/x/planets/types"
)

var _ types.QueryServer = Keeper{}
