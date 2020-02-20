package reseeding

import (
	"github.com/cosmos/modules/incubator/reseeding/internal/keeper"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"
)

type (
	Keeper = keeper.Keeper
)

const (
	DefaultCodespace = types.DefaultCodespace
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	RouterKey        = types.RouterKey
)

var (
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)
