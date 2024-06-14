package cmd

import (
	"dataworks-cli/pkg/services"
	"log"

	"github.com/spf13/cobra"
)

var listDataMoelColumnsCmd = &cobra.Command{
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

var dataModelShowTablesCmd = &cobra.Command{
	Use:   "show-tables",
	Short: "Show tables",
	Run: func(cmd *cobra.Command, args []string) {
		models, err := services.ShowTables(dataModelsModelType)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if err := WriteJSON(dataModelsOutputPath, models); err != nil {
			log.Fatalln(err)
		}
	},
}

var dataModelListCmd = &cobra.Command{
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

var dataModelsListCmd = &cobra.Command{
	Use: "data-models",
}

var dataModelsOutputPath string
var dataModelsModelType string

func init() {
	rootCmd.AddCommand(dataModelsListCmd)
	dataModelsListCmd.AddCommand(listDataMoelColumnsCmd)

	listDataMoelColumnsCmd.Flags().StringVarP(&dataModelsOutputPath, "out", "o", "", "path to output")
	_ = listDataMoelColumnsCmd.MarkFlagRequired("out")

	dataModelsListCmd.AddCommand(dataModelListCmd)
	dataModelListCmd.Flags().StringVarP(&dataModelsOutputPath, "out", "o", "", "path to output")
	_ = dataModelListCmd.MarkFlagRequired("out")

	dataModelsListCmd.AddCommand(dataModelShowTablesCmd)
	dataModelShowTablesCmd.Flags().StringVarP(&dataModelsOutputPath, "out", "o", "", "path to output")
	dataModelShowTablesCmd.Flags().StringVarP(&dataModelsModelType, "type", "t", "", "model type")
	_ = dataModelShowTablesCmd.MarkFlagRequired("out")
}
