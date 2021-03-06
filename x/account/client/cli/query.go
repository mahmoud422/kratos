package cli

import (
	"fmt"

	"github.com/KuChain-io/kuchain/chain/client/flags"
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/account/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the transaction commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the account module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetAccountCmd(cdc),
		GetAuthCmd(cdc),
	)

	return cmd
}

// GetAccountCmd returns a query account that will display the state of the
// account at a given name.
func GetAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name]",
		Short: "Query account data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAccountRetriever(cliCtx)

			key, err := chainTypes.NewAccountIDFromStr(args[0])
			if err != nil {
				fmt.Printf("new account id error %v", err.Error())
				return err
			}

			acc, err := accGetter.GetAccount(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(acc)
		},
	}

	return flags.GetCommands(cmd)[0]
}

// GetAuthCmd returns a query auth that will display the state of the auth
func GetAuthCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth [acc-address]",
		Short: "Query auth data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accGetter := types.NewAccountRetriever(cliCtx)

			key, err := chainTypes.AccAddressFromBech32(args[0])
			if err != nil {
				fmt.Printf("new acc-address error %v", err.Error())
				return err
			}

			data, err := accGetter.GetAddAuth(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(data)
		},
	}

	return flags.GetCommands(cmd)[0]
}
