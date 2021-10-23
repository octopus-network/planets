package keeper

import (
	"github.com/octopus-network/planets/earth/x/earth/types"
)

var _ types.QueryServer = Keeper{}
