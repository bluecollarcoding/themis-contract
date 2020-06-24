package main

import (
	"fmt"
	"os"

	contract "github.com/informalsystems/themis-contract/pkg/themis-contract"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func newCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new [upstream] [output]",
		Short: "Create a new contract",
		Long:  "Create a new contract, using the specified upstream contract effectively as a template",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			contractPath := defaultContractPath
			if len(args) > 1 {
				contractPath = args[1]
			}

			var upstreamLoc string
			if len(args[0]) == 0 {
				_ = fmt.Errorf("when creating a contract with the `new` command, an upstream contract must be supplied as a template")
				os.Exit(1)
			} else {
				upstreamLoc = args[0]
			}

			if _, err := contract.New(contractPath, upstreamLoc, globalCtx); err != nil {
				log.Error().Err(err).Msg("Failed to create new contract")
				os.Exit(1)
			}
		},
	}
}
