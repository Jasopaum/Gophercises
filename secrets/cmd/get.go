package cmd

import (
	"fmt"
	"gophercises/secrets"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a secret",
	Run: func(cmd *cobra.Command, args []string) {
		v := secrets.File(keyphrase, secretsFile())
		entry := args[0]
		res, err := v.Get(entry)
		if err != nil {
			fmt.Println("Failed to get secret (maybe wrong keyphrase):", err)
			return
		}
		fmt.Printf("%s: %s\n", entry, res)
	},
}
