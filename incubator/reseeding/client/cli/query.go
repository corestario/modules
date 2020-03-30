package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/corestario/modules/incubator/reseeding/internal/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	reseedingQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the reseeding module",
	}

	reseedingQueryCmd.AddCommand()

	return reseedingQueryCmd
}
