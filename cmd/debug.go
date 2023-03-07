package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func debugCmd(cfg *settings) {
	fmt.Printf("%#v\n", cfg)

}

func newDebugCmd(cfg *settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "debug",
		Run: func(_ *cobra.Command, args []string) {
			debugCmd(cfg)
		},
	}
	return cmd
}
