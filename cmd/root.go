package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"github.com/darthrevan13/issuspicous/pkg/addrParser"
	"github.com/darthrevan13/issuspicous/pkg/certParser"

)

var rootCmd = &cobra.Command{
	Use:   "issuspicous",
	Short: "Suspicious site checker",
	Long: "A small utility to verify SSL/TLS certificates of a website and also if it's approved by an online web filter",
	RunE: func(_ *cobra.Command, args []string) error {
		addrParser.Parse(args)
		certParser.Get("google.com", "443")
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
