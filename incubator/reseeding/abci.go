package reseeding

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/reseeding/internal/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	seed := k.GetCurrentSeed(ctx)
	if len(seed) > 0 {
		k.ResetCurrentSeed(ctx)
		k.ResetSeeds(ctx)
	}

	// TODO: change []abci.ValidatorUpdate to a new type that can hold the seed,
	// set the seed.
	return []abci.ValidatorUpdate{}
}
