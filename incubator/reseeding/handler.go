package reseeding

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/modules/incubator/reseeding/internal/keeper"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenericHandler routes the messages to the handlers
func GenericHandler(k Keeper, stakingKeeper staking.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSeed:
			return HandleMsgSeed(ctx, msg, k, stakingKeeper)
		default:
			return nil, fmt.Errorf("unrecognized reseeding message type: %T", msg)
		}
	}
}

func HandleMsgSeed(ctx sdk.Context, msg types.MsgSeed, k keeper.Keeper, stakingKeeper staking.Keeper) (*sdk.Result, error) {
	validators := stakingKeeper.GetAllValidators(ctx)

	var isValidator bool
	for _, validator := range validators {
		if bytes.Equal(validator.ConsPubKey.Address().Bytes(), msg.Sender.Bytes()) {
			isValidator = true
			break
		}
	}

	if !isValidator {
		return nil, errors.New("got a seed from a non-validator")
	}

	seeds, err := k.GetSeeds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to GetSeeds: %w", err)
	}

	seeds.Add(msg.Seed, msg.Sender.String())

	// We have all necessary seeds.
	if seeds.GetVotesForSeed(msg.Seed) > len(validators)*(2/3) {
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
		return &sdk.Result{Events: ctx.EventManager().Events()}, nil
	}

	// We do not have enough validators, but the seed is O.K., save it and continue.
	if err := k.StoreSeed(ctx, msg.Sender, msg.Seed); err != nil {
		return nil, fmt.Errorf("failed to StoreSeed: %w", err)
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

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
