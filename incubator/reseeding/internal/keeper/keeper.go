package keeper

import (
	"encoding/json"
	"fmt"

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

func (k Keeper) StoreSeed(ctx sdk.Context, seeds types.Seeds) (err error) {
	store := ctx.KVStore(k.storeKey)

	bz, _ := json.Marshal(seeds)
	store.Set([]byte(seedsKey), bz)

	return nil
}

func (k Keeper) GetSeeds(ctx sdk.Context) (types.Seeds, error) {
	store := ctx.KVStore(k.storeKey)
	var seeds = make(types.Seeds)
	seedsBytes := store.Get([]byte(seedsKey))
	if seedsBytes == nil {
		return seeds, nil
	}
	if err := json.Unmarshal(seedsBytes, &seeds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal seeds: %w", err)
	}

	return seeds, nil
}

func (k Keeper) ResetSeeds(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	bz, _ := json.Marshal(make(types.Seeds))
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
