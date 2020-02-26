package reseeding

import (
	"github.com/cosmos/modules/incubator/reseeding/internal/keeper"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)

const (
	DefaultCodespace = types.DefaultCodespace
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	RouterKey        = types.RouterKey
	QuerierRoute     = types.QuerierRoute
)

var (
	NewQuerier          = keeper.NewQuerier
	RegisterInvariants  = keeper.RegisterInvariants
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	RegisterCodec       = types.RegisterCodec
	ModuleCdc           = types.ModuleCdc
)
