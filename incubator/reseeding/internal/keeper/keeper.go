package keeper

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"
)

const seedsKey = "seeds"
const currentSeedKey = "current_seed"

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The amino codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the reseeding Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k Keeper) StoreSeed(ctx sdk.Context, sender sdk.Address, seed []byte) (err error) {
	store := ctx.KVStore(k.storeKey)
	seeds, err := k.GetSeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get seeds: %w", err)
	}

	seeds[sender.String()] = seed
	bz, _ := json.Marshal(seeds)
	store.Set([]byte(seedsKey), bz)

	return nil
}

func (k Keeper) GetSeeds(ctx sdk.Context) (map[string][]byte, error) {
	store := ctx.KVStore(k.storeKey)
	var seeds = make(map[string][]byte)
	if err := json.Unmarshal(store.Get([]byte(seedsKey)), &seeds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal seeds: %w", err)
	}

	return seeds, nil
}

func (k Keeper) ResetSeeds(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	bz, _ := json.Marshal(make(map[string][]byte))
	store.Set([]byte(seedsKey), bz)
}

func (k Keeper) SetCurrentSeed(ctx sdk.Context, seed []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(currentSeedKey), seed)
}

func (k Keeper) ResetCurrentSeed(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(currentSeedKey), []byte{})
}

func (k Keeper) GetCurrentSeed(ctx sdk.Context) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get([]byte(currentSeedKey))
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		return nil, nil
	}
}
