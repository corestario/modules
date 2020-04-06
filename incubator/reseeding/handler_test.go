package reseeding

import (
	"testing"

	"github.com/corestario/modules/incubator/reseeding/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgSeed(t *testing.T) {
	var (
		req  = require.New(t)
		seed = []byte("test_seed")
	)
	ctx, reseedingKeeper, stakingKeeper, PKs := getTestEnvironment()

	t.Run("sender_is_not_a_validator", func(t *testing.T) {
		msg := types.MsgSeed{}
		_, err := HandleMsgSeed(ctx, msg, reseedingKeeper, stakingKeeper)
		req.Error(err)
	})

	t.Run("seed_successfully_stored", func(t *testing.T) {
		msg := types.MsgSeed{
			Sender: sdk.AccAddress(PKs[0].Address()),
			Seed:   seed,
		}
		_, err := HandleMsgSeed(ctx, msg, reseedingKeeper, stakingKeeper)
		req.NoError(err)

		seeds, err := reseedingKeeper.GetSeeds(ctx)
		req.NoError(err)
		req.Len(seeds, 1)
	})

	t.Run("polka", func(t *testing.T) {
		for i := 1; i < 3; i++ {
			msg := types.MsgSeed{
				Sender: sdk.AccAddress(PKs[i].Address()),
				Seed:   seed,
			}
			_, err := HandleMsgSeed(ctx, msg, reseedingKeeper, stakingKeeper)
			req.NoError(err)

			seeds, err := reseedingKeeper.GetSeeds(ctx)
			req.NoError(err)
			req.Len(seeds[types.GetSeedKey(seed)], i+1)
		}

		currentSeed := reseedingKeeper.GetCurrentSeed(ctx)
		req.Equal(seed, currentSeed)
	})
}
