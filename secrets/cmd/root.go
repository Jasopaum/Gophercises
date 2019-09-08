package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "secrets",
	Short: "secrets is a CLI secret manager",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var keyphrase string

func init() {
	rootCmd.PersistentFlags().StringVarP(&keyphrase, "key", "k", "", "keyphrase used to encrypt/decrypt secrets")
	rootCmd.MarkFlagRequired("key")
}

func secretsFile() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return filepath.Join(home, ".secrets")
}
