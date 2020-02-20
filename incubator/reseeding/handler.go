package reseeding

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/modules/incubator/reseeding/internal/keeper"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenericHandler routes the messages to the handlers
func GenericHandler(k Keeper, stakingKeeper staking.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSeed:
			HandleMsgSeed(ctx, msg, k, stakingKeeper)
		default:
			errMsg := fmt.Sprintf("unrecognized reseeding message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func HandleMsgSeed(ctx sdk.Context, msg types.MsgSeed, k keeper.Keeper, stakingKeeper staking.Keeper) sdk.Result {
	validators := stakingKeeper.GetAllValidators(ctx)

	var isValidator bool
	for _, validator := range validators {
		if bytes.Equal(validator.ConsPubKey.Address().Bytes(), msg.Sender.Bytes()) {
			isValidator = true
			break
		}
	}

	if !isValidator {
		return sdk.NewError(types.DefaultCodespace, 0, "got a seed from a non-validator").Result()
	}

	seeds, err := k.GetSeeds(ctx)
	if err != nil {
		return sdk.NewError(types.DefaultCodespace, 0, err.Error()).Result()
	}

	if _, ok := seeds[msg.Sender.String()]; ok {
		return sdk.NewError(types.DefaultCodespace, 0, "duplicate seeds from validator %s", msg.Sender.String()).Result()
	}

	// Check that we have identical seeds.
	for _, seed := range seeds {
		if !bytes.Equal(seed, msg.Seed) {
			k.ResetSeeds(ctx)
			return sdk.NewError(types.DefaultCodespace, 0, "mismatching seeds, %+v", seeds).Result()
		}
	}

	// We have all necessary seeds.
	if len(validators) == len(seeds)+1 {
		// All good, set current seed.
		k.SetCurrentSeed(ctx, msg.Seed)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				"EventNewCurrentSeed",
				sdk.NewAttribute("new_current_seed", string(msg.Seed)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
		})
		return sdk.Result{Events: ctx.EventManager().Events()}
	}

	// We do not have enough validators, but the seed is O.K., save it and continue.
	if err := k.StoreSeed(ctx, msg.Sender, msg.Seed); err != nil {
		return sdk.NewError(types.DefaultCodespace, 0, err.Error()).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"EventSeedSaved",
			sdk.NewAttribute("seed_saved", string(msg.Seed)),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
