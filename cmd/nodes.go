package cmd

import (
	"dataworks-helper/pkg/services"
	"log"

	"github.com/spf13/cobra"
)

var nodesCmd = &cobra.Command{
	Use: "nodes",
}

var nodesListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		nodes, err := services.ListNodes(projectEnv)

		if err != nil {
			log.Fatalln(err)
		}

		if err := WriteJSON(nodesListOutputPath, nodes); err != nil {
			log.Fatalln(err)
		}
	},
}

var nodesListOutputPath string

var projectEnv string

func init() {
	rootCmd.AddCommand(nodesCmd)

	nodesCmd.AddCommand(nodesListCmd)
	nodesListCmd.Flags().StringVarP(&nodesListOutputPath, "out", "o", "", "puth to file list output")
	_ = nodesListCmd.MarkFlagRequired("out")

	nodesListCmd.Flags().StringVarP(&projectEnv, "env", "e", "PROD", "product env")
}
