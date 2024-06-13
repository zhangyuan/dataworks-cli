package cmd

import (
	"dataworks-cli/pkg/services"
	"log"

	"github.com/spf13/cobra"
)

var listDataModelColumnsCmd = &cobra.Command{
	Use:   "list-columns",
	Short: "List data model coloumns",
	Run: func(cmd *cobra.Command, args []string) {
		models, err := services.ListDataModelColumns()
		if err != nil {
			log.Fatalf("%v", err)
		}

		if err := WriteJSON(dataModelsOutputPath, models); err != nil {
			log.Fatalln(err)
		}
	},
}

var listDataModelsCmd = &cobra.Command{
	Use:   "list",
	Short: "List data models",
	Run: func(cmd *cobra.Command, args []string) {
		models, err := services.ListDataModels()
		if err != nil {
			log.Fatalf("%v", err)
		}

		if err := WriteJSON(dataModelsOutputPath, models); err != nil {
			log.Fatalln(err)
		}
	},
}

var dataModelsCmd = &cobra.Command{
	Use: "data-models",
}

var dataModelsOutputPath string

func init() {
	rootCmd.AddCommand(dataModelsCmd)
	dataModelsCmd.AddCommand(listDataModelColumnsCmd)
	dataModelsCmd.AddCommand(listDataModelsCmd)

	listDataModelColumnsCmd.Flags().StringVarP(&dataModelsOutputPath, "out", "o", "", "path to output")
	_ = listDataModelColumnsCmd.MarkFlagRequired("out")

	listDataModelsCmd.Flags().StringVarP(&dataModelsOutputPath, "out", "o", "", "path to output")
	_ = listDataModelsCmd.MarkFlagRequired("out")
}
