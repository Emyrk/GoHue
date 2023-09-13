package cmd

import (
	"fmt"
	"github.com/Emyrk/gohue"
	"github.com/spf13/cobra"
	"log/slog"
)

func authenticate() *cobra.Command {
	return &cobra.Command{
		Use:   "authenticate",
		Short: "Authenticate with the Hue Bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ctx = gohue.WithDebugging(ctx, slog.New(slog.NewTextHandler(cmd.OutOrStdout(), nil)))

			cli, err := gohue.NewClient("l57Ry9PcABEOwWKKvR-UnRCG2CgWejeaNMJxYuwV")
			if err != nil {
				return err
			}

			fmt.Println("Go press the button!")
			resp, err := cli.GenerateAPIKey(ctx)
			if err != nil {
				return err
			}

			fmt.Printf("Username: %s\nClientKey: %s\n", resp[0].Success.Username, resp[0].Success.ClientKey)
			return nil
		},
	}
}
