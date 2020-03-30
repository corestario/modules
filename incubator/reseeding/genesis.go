package reseeding

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/corestario/modules/incubator/reseeding/internal/types"
)

// InitGenesis sets reseeding information for genesis.
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	return NewGenesisState()
}
