package cmd

import (
	"dataworks-helper/pkg/dataworks"
	"dataworks-helper/pkg/services"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/samber/lo"
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

		if err := WriteJSON(listOutputPath, files); err != nil {
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

		lo.ForEach(files, func(file dataworks.NormalFile, _ int) {
			if err := services.DownloadFile(file, filesOutputDirectoryPath); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(500 * time.Millisecond)
		})
	},
}

var listOutputPath string

var filesListFilePath string
var filesOutputDirectoryPath string
var fileTypes string

func init() {
	rootCmd.AddCommand(filesCmd)

	filesCmd.AddCommand(listFilesCmd)
	listFilesCmd.Flags().StringVarP(&listOutputPath, "out", "o", "", "puth to file list output")
	listFilesCmd.Flags().StringVarP(&fileTypes, "file-types", "t", "10", "file types")
	_ = listFilesCmd.MarkFlagRequired("out")

	filesCmd.AddCommand(filesDownloadCmd)
	filesDownloadCmd.Flags().StringVarP(&filesListFilePath, "input", "i", "", "path to file list")
	filesDownloadCmd.Flags().StringVarP(&filesOutputDirectoryPath, "out", "o", "", "path to files directory")
	_ = filesDownloadCmd.MarkFlagRequired("input")
	_ = filesDownloadCmd.MarkFlagRequired("output")

}
