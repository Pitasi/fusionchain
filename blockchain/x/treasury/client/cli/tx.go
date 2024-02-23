// Copyright (c) Fusion Laboratories LTD
// SPDX-License-Identifier: BUSL-1.1
package cli

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/qredo/fusionchain/x/treasury/types"
	"github.com/spf13/cobra"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdNewKeyRequest())
	cmd.AddCommand(CmdFulfilKeyRequest())
	cmd.AddCommand(CmdRejectKeyRequest())
	cmd.AddCommand(CmdNewSignatureRequest())
	cmd.AddCommand(CmdFulfilSignatureRequest())
	cmd.AddCommand(CmdNewSignTransactionRequest())
	// this line is used by starport scaffolding # 1

	return cmd
}
