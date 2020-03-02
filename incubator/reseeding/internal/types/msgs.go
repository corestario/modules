package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* --------------------------------------------------------------------------- */
// MsgSeed
/* --------------------------------------------------------------------------- */

// MsgSeed defines a seed obtained by a participant.
type MsgSeed struct {
	Sender sdk.AccAddress
	Seed   []byte
}

// MsgSeed is a constructor function for MsgSeed.
func NewMsgSeed(sender sdk.AccAddress, seed []byte) MsgSeed {
	return MsgSeed{
		Sender: sender,
		Seed:   seed,
	}
}

// Route Implements Msg
func (msg MsgSeed) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgSeed) Type() string { return "seed" }

// ValidateBasic Implements Msg.
func (msg MsgSeed) ValidateBasic() error {
	if msg.Sender.Empty() {
		return errors.New("invalid sender address")
	}
	if len(msg.Seed) == 0 {
		return errors.New("invalid seed")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSeed) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgSeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
