package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/modules/incubator/reseeding/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	reseedingTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Reseeding transactions subcommands",
	}

	reseedingTxCmd.AddCommand(client.PostCommands(
		GetCmdSendSeed(cdc),
	)...)

	return reseedingTxCmd
}

func GetCmdSendSeed(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [sender] [seed]",
		Short: "send a seed",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			sender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			seedStr := args[1]

			msg := types.NewMsgSeed(sender, []byte(seedStr))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}