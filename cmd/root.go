package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aliyun-dataworks",
	Short: "dataworks utils",
}

func WriteJSON(outputPath string, data interface{}) error {
	file, err := os.Create(listOutputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(bytes); err != nil {
		return err
	}

	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
