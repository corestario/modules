package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"
	"github.com/gogo/protobuf/codec"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	rsdTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Reseeding transactions subcommands",
	}

	rsdTxCmd.AddCommand(client.PostCommands(
	//GetCmdTransferNFT(cdc),
	//GetCmdEditNFTMetadata(cdc),
	//GetCmdMintNFT(cdc),
	//GetCmdBurnNFT(cdc),
	)...)

	return rsdTxCmd
}
