package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName                         = "reseeding"
	DefaultCodespace sdk.CodespaceType = ModuleName
	StoreKey                           = ModuleName
	RouterKey                          = ModuleName
	QuerierRoute                       = ModuleName
)

var (
	CollectionsKeyPrefix = []byte{0x00} // key for reseeding collections
	OwnersKeyPrefix      = []byte{0x01} // key for balance of NFTs held by an address
)
