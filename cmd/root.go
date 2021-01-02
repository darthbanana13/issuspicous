package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"github.com/darthrevan13/issuspicous/pkg/addressParser"

)

var sslPort = "443";

var rootCmd = &cobra.Command{
	Use:   "issuspicous",
	Short: "Suspicious site checker",
	Long: "A small utility to verify SSL/TLS certificates of a website and also if it's approved by an online web filter",
	RunE: func(_ *cobra.Command, args []string) error {
		addressParser.Parse(args)
		return nil;
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
