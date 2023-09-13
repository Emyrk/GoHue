package cmd

import (
	"github.com/Emyrk/gohue"
	"github.com/spf13/cobra"
	"log/slog"
)

func devices() *cobra.Command {
	return &cobra.Command{
		Use:   "devices",
		Short: "List devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ctx = gohue.WithDebugging(ctx, slog.New(slog.NewTextHandler(cmd.OutOrStdout(), nil)))

			cli, err := gohue.NewClient("l57Ry9PcABEOwWKKvR-UnRCG2CgWejeaNMJxYuwV")
			if err != nil {
				return err
			}

			cli.Devices(ctx)

			return nil
		},
	}
}
