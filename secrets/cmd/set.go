package cmd

import (
	"fmt"
	"gophercises/secrets"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a secret",
	Run: func(cmd *cobra.Command, args []string) {
		v := secrets.File(keyphrase, secretsFile())
		entry, value := args[0], args[1]
		err := v.Set(entry, value)
		if err != nil {
			fmt.Println("Failed to set secret (maybe wrong keyphrase):", err)
		}
	},
}
