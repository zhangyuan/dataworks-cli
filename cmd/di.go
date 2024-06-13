package cmd

import (
	"dataworks-cli/pkg/services"
	"log"

	"github.com/spf13/cobra"
)

var diCmd = &cobra.Command{
	Use: "di",
}

var diListCmd = &cobra.Command{
	Use: "list-sync-tasks",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := services.ListDISyncTasks(diTaskType, diDataSourceName)
		if err != nil {
			log.Fatalln(err)
		}

		if err := WriteJSON(diTasksOutputPath, files); err != nil {
			log.Fatalln(err)
		}
	},
}

var diTasksOutputPath string

var diTaskType string
var diDataSourceName string

func init() {
	rootCmd.AddCommand(diCmd)
	diCmd.AddCommand(diListCmd)

	diListCmd.Flags().StringVarP(&diTasksOutputPath, "out", "o", "", "path to output")
	_ = diListCmd.MarkFlagRequired("out")
	diListCmd.Flags().StringVarP(&diTaskType, "task-type", "t", "DI_OFFLINE", "DI task type")
	diListCmd.Flags().StringVarP(&diDataSourceName, "data-source-name", "s", "", "DI task data source name")
	_ = diListCmd.MarkFlagRequired("data-source-name")
}
