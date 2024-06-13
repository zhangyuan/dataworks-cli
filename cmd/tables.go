package cmd

import (
	"dataworks-cli/pkg/services"
	"log"

	"github.com/spf13/cobra"
)

// tablesCmd represents the tables command
var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "export tables",
	Run: func(cmd *cobra.Command, args []string) {
		tables, err := services.GetTables(appGuid)
		if err != nil {
			log.Fatalln(err)
		}

		if err := WriteJSON(tablesOutputPath, tables); err != nil {
			log.Fatalln(err)
		}
	},
}

var tablesOutputPath string
var appGuid string

func init() {
	rootCmd.AddCommand(tablesCmd)

	tablesCmd.Flags().StringVarP(&tablesOutputPath, "out", "o", "", "puth to output file")
	_ = tablesCmd.MarkFlagRequired("out")

	tablesCmd.Flags().StringVarP(&appGuid, "app-guid", "a", "", "max compute database name")
	_ = tablesCmd.MarkFlagRequired("app-guid")
}
