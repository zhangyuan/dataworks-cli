package cmd

import (
	"dataworks-helper/pkg/dataworks"
	"dataworks-helper/pkg/services"
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use: "files",
}

var listFilesCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := services.ListFilesNormalized(fileTypes)

		if err != nil {
			log.Fatalln(err)
		}

		if err := WriteJSON(filesListOutputPath, files); err != nil {
			log.Fatalln(err)
		}
	},
}

var filesDownloadCmd = &cobra.Command{
	Use: "download",
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := os.ReadFile(filesListFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		var files []dataworks.NormalFile
		if err := json.Unmarshal(bytes, &files); err != nil {
			log.Fatalln(err)
		}

		if err := services.DownloadFiles(files, filesOutputDirectoryPath); err != nil {
			log.Fatalln(err)
		}
	},
}

var filesListOutputPath string

var filesListFilePath string
var filesOutputDirectoryPath string
var fileTypes string

func init() {
	rootCmd.AddCommand(filesCmd)

	filesCmd.AddCommand(listFilesCmd)
	listFilesCmd.Flags().StringVarP(&filesListOutputPath, "out", "o", "", "puth to file list output")
	listFilesCmd.Flags().StringVarP(&fileTypes, "file-types", "t", "10", "file types")
	_ = listFilesCmd.MarkFlagRequired("out")

	filesCmd.AddCommand(filesDownloadCmd)
	filesDownloadCmd.Flags().StringVarP(&filesListFilePath, "input", "i", "", "path to file list")
	filesDownloadCmd.Flags().StringVarP(&filesOutputDirectoryPath, "out", "o", "", "path to files directory")
	_ = filesDownloadCmd.MarkFlagRequired("input")
	_ = filesDownloadCmd.MarkFlagRequired("output")

}
