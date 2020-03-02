package reseeding

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/modules/incubator/reseeding/client/cli"
	"github.com/cosmos/modules/incubator/reseeding/client/rest"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic app module basics object
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name defines module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, types.ModuleCdc, RouterKey)
}

// GetTxCmd gets the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd gets the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule supply app module
type AppModule struct {
	AppModuleBasic

	keeper        Keeper
	stakingKeeper staking.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, stakingKeeper staking.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		stakingKeeper:  stakingKeeper,
	}
}

// Name defines module name
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the reseeding module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route module message route name
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler module handler
func (am AppModule) NewHandler() sdk.Handler {
	return GenericHandler(am.keeper, am.stakingKeeper)
}

// QuerierRoute module querier route name
func (AppModule) QuerierRoute() string {
	return ""
}

// NewQuerierHandler module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier()
}

// InitGenesis module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock module begin-block
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock module end-block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) abci.ModuleEndBlock {
	return EndBlocker(ctx, am.keeper)
}

func NewQuerier() sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		return nil, errors.ErrUnknownRequest
	}
}
