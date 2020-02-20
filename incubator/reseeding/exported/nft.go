package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// reseeding non fungible token interface
type reseeding interface {
	GetID() string
	GetOwner() sdk.AccAddress
	SetOwner(address sdk.AccAddress)
	GetTokenURI() string
	EditMetadata(tokenURI string)
	String() string
}
