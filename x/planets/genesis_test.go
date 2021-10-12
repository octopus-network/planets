package planets_test

import (
	"testing"

	keepertest "github.com/octopus-network/planets/testutil/keeper"
	"github.com/octopus-network/planets/x/planets"
	"github.com/octopus-network/planets/x/planets/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PlanetsKeeper(t)
	planets.InitGenesis(ctx, *k, genesisState)
	got := planets.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
