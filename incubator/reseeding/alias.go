package reseeding

import (
	"github.com/corestario/modules/incubator/reseeding/internal/keeper"
	"github.com/corestario/modules/incubator/reseeding/internal/types"
)

type (
	Keeper  = keeper.Keeper
	MsgSeed = types.MsgSeed
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

var (
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	RegisterCodec       = types.RegisterCodec
	NewKeeper           = keeper.NewKeeper
)
