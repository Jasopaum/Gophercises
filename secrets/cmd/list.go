package cmd

import (
	"fmt"
	"gophercises/secrets"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	Run: func(cmd *cobra.Command, args []string) {
		v := secrets.File(keyphrase, secretsFile())
		res, err := v.List()
		if err != nil {
			fmt.Println("Failed to get secrets (maybe wrong keyphrase):", err)
			return
		}
		for _, s := range res {
			fmt.Println(s)
		}
	},
}
