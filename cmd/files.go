package cmd

import (
	"dataworks-helper/pkg/dataworks"
	"encoding/json"
	"log"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use: "files",
}

var listFilesCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Create(listOutputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		files, err := dataworks.GetScriptsWithContent()
		if err != nil {
			log.Fatalln(err)
		}

		bytes, err := json.MarshalIndent(files, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := file.Write(bytes); err != nil {
			log.Fatalln(err)
		}
	},
}

var listAllFilesCmd = &cobra.Command{
	Use: "list-all",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Create(listOutputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		files, err := dataworks.ListFiles()
		if err != nil {
			log.Fatalln(err)
		}

		bytes, err := json.MarshalIndent(files, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := file.Write(bytes); err != nil {
			log.Fatalln(err)
		}
	},
}

var fetchFilesCmd = &cobra.Command{
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
			if err := dataworks.DownloadFile(file, filesOutputDirectoryPath); err != nil {
				log.Fatalln(err)
			}
		})
	},
}

var listOutputPath string

var filesListFilePath string
var filesOutputDirectoryPath string

func init() {
	rootCmd.AddCommand(filesCmd)

	filesCmd.AddCommand(listAllFilesCmd)
	listAllFilesCmd.Flags().StringVarP(&listOutputPath, "out", "o", "", "puth to file list output")
	_ = listAllFilesCmd.MarkFlagRequired("out")

	filesCmd.AddCommand(listFilesCmd)
	listFilesCmd.Flags().StringVarP(&listOutputPath, "out", "o", "", "puth to file list output")
	_ = listFilesCmd.MarkFlagRequired("out")

	filesCmd.AddCommand(fetchFilesCmd)
	fetchFilesCmd.Flags().StringVarP(&filesListFilePath, "input", "i", "", "path to file list")
	fetchFilesCmd.Flags().StringVarP(&filesOutputDirectoryPath, "out", "o", "", "path to files directory")
	_ = fetchFilesCmd.MarkFlagRequired("input")
	_ = fetchFilesCmd.MarkFlagRequired("output")
}
