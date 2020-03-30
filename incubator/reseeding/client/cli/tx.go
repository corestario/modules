package cli

import (
	"bufio"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/corestario/modules/incubator/reseeding/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	reseedingTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Reseeding transactions subcommands",
	}

	reseedingTxCmd.AddCommand(
		flags.PostCommands(
			GetCmdSendSeed(cdc),
		)...,
	)

	return reseedingTxCmd
}

func GetCmdSendSeed(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [sender] [seed]",
		Short: "send a seed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			seedStr := args[0]

			msg := types.NewMsgSeed(cliCtx.GetFromAddress(), []byte(seedStr))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
