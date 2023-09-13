package cmd

import (
	"encoding/json"
	"github.com/Emyrk/gohue"
	"github.com/spf13/cobra"
)

func discover() *cobra.Command {
	return &cobra.Command{
		Use:   "discover",
		Short: "Discover hue bridges",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			bridges, err := gohue.DiscoverBridges(ctx)
			if err != nil {
				return err
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(bridges)
		},
	}
}
