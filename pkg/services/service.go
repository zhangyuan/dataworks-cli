package services

import (
	"dataworks-cli/pkg/dataworks"
	"os"
	"strconv"
	"time"

	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
)

const (
	NOMARL_USE_TYPE = "NORMAL"
)

func CreateClient() (*dataworks.Client, error) {
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	endpoint := os.Getenv("DATAWORKS_ENDPOINT")
	throttle := os.Getenv("THROTTLE")

	if throttle == "" {
		throttle = "1s"
	}

	duration, err := time.ParseDuration(throttle)
	if err != nil {
		return nil, err
	}

	projectId, err := GetProjectId()
	if err != nil {
		return nil, err
	}

	return dataworks.NewClient(accessKeyId, accessKeySecret, endpoint, projectId, duration)
}

func GetProjectId() (int64, error) {
	projectIdString := os.Getenv("DATAWORKS_PROJECT_ID")
	projectId, err := strconv.ParseInt(projectIdString, 10, 64)
	if err != nil {
		return 0, err
	}

	return projectId, nil
}

func ListFiles(fileTypes string) ([]*dataworks_public20200518.ListFilesResponseBodyDataFiles, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.ListFiles(fileTypes)
}

func ListDIJobs() ([]*dataworks_public20200518.ListDIJobsResponseBodyDIJobPagingDIJobs, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.ListDIJobs()
}

func GetFolders(projectId int64, folderIds []string) ([]dataworks.Folder, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.GetFolders(projectId, folderIds)
}

func ListFilesNormalized(fileTypes string) ([]dataworks.NormalFile, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.ListFilesNormalized(fileTypes)
}

func GetFileContent(file dataworks.NormalFile) (string, error) {
	client, err := CreateClient()
	if err != nil {
		return "", err
	}

	return client.GetFileContent(file)
}

func DownloadFiles(files []dataworks.NormalFile, directory string) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}
	return client.DownloadFiles(files, directory)
}

func ListDISyncTasks(taskType string, dataSourceName string) ([]*dataworks_public20200518.ListRefDISyncTasksResponseBodyDataDISyncTasks, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.ListDISyncTasks(taskType, dataSourceName)
}

func GetTables(appGuid string) ([]dataworks.Table, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.GetTables(appGuid)
}

func ListNodes(productEnv string) ([]*dataworks_public20200518.ListNodesResponseBodyDataNodes, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.ListNodes(productEnv)
}

func ListDataModels() ([]dataworks.DataModel, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.ListDataModels()
}

func ListDataModelColumns() ([]dataworks.RawDataModelColumn, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.ListDataModelColumns()
}

func ShowTables(modelType string) (interface{}, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.DataModelShowTables(modelType)
}

func ShowTableColumns(tableType string, tableCode string) (interface{}, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	return client.DataModelShowTableColumns(tableType, tableCode)
}
