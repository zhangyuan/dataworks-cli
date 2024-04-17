package cmd

import (
	"dataworks-helper/pkg/dataworks"
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// tablesCmd represents the tables command
var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "export tables",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Create(tablesOutputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		tables, err := dataworks.GetTables(appGuid)
		if err != nil {
			log.Fatalln(err)
		}

		bytes, err := json.MarshalIndent(tables, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := file.Write(bytes); err != nil {
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
