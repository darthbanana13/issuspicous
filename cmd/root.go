package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"github.com/darthrevan13/issuspicous/pkg/certParser"
	"github.com/darthrevan13/issuspicous/pkg/fortinetCateg"
)

var rootCmd = &cobra.Command{
	Use:   "issuspicous",
	Short: "Suspicious site checker",
	Long: "A small utility to verify SSL/TLS certificates of a website and also if it's approved by an online web filter",
	RunE: func(_ *cobra.Command, args []string) error {
		certsChan := certParser.NewCerts(args)
		//TODO: Get certificates and Frotinet category concurrently
		for c := range certsChan {
			fmt.Println(c.CertificateInfo())
			categ := fortinetCateg.NewCateg(c.Addr.Addr())
			fmt.Println(categ.String())
			fmt.Println()
		}
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
