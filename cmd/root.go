package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"net/http"
	 "io/ioutil"
	"github.com/darthrevan13/issuspicous/pkg/certParser"
)

var rootCmd = &cobra.Command{
	Use:   "issuspicous",
	Short: "Suspicious site checker",
	Long: "A small utility to verify SSL/TLS certificates of a website and also if it's approved by an online web filter",
	RunE: func(_ *cobra.Command, args []string) error {
		certsChan := certParser.NewCerts(args)
		reg, _ := regexp.Compile(`<h4 class="info_title">Category: (.*)</h4>`)
		//TODO: Get certificates and Frotinet category concurrently
		for c := range certsChan {
			fmt.Println(c.CertificateInfo())
			// TODO: Put in separate package
			// TODO: Handle errors
			resp, _ := http.Get("https://www.fortiguard.com/webfilter?q=" + c.Addr.Addr() + "&version=8")
			if resp.StatusCode == 500 {
				fmt.Println("Fortinet Category:	Response blocked")
			} else if resp.StatusCode == 200 {
				// TODO: Handle errors
				body, _ := ioutil.ReadAll(resp.Body)
				match := reg.FindStringSubmatch(string(body))
				if len(match) > 1 {
					fmt.Println("Fortinet Category:	" + match[1])
				}
			}
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
