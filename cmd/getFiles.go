package cmd

import (
	"aliyun-dataworks/pkg/dataworks"
	"log"

	"github.com/spf13/cobra"
)

// getFilesCmd represents the getFiles command
var getFilesCmd = &cobra.Command{
	Use: "getFiles",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dataworks.GetFiles(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getFilesCmd)
}
