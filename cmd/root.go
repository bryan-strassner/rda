package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rda",
	Short: "RDA analyzes rendered Deckhand documents",
	Long: `Rendered Document Analyzer
	Analyzes rendered Deckhand documents`,
	Run: func(cmd *cobra.Command, args []string) {
		// Root command doesn't do much
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	fmt.Println("Initializing RDA")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
