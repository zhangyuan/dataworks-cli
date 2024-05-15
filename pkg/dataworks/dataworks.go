package dataworks

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
)

const (
	NOMARL_USE_TYPE = "NORMAL"
)

func CreateClient() (_result *dataworks_public20200518.Client, _err error) {
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")

	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}

	endpoint := os.Getenv("DATAWORKS_ENDPOINT")
	config.Endpoint = &endpoint
	return dataworks_public20200518.NewClient(config)
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

	projectId, err := GetProjectId()
	if err != nil {
		return nil, err
	}

	files := []*dataworks_public20200518.ListFilesResponseBodyDataFiles{}
	var pageNumber int32 = 1
	var pageSize int32 = 100
	for {
		listFilesRequest := &dataworks_public20200518.ListFilesRequest{
			ProjectId:  &projectId,
			PageNumber: &pageNumber,
			PageSize:   &pageSize,
		}
		if fileTypes != "" {
			listFilesRequest.FileTypes = &fileTypes
		}

		res, err := client.ListFiles(listFilesRequest)
		if err != nil {
			return nil, err
		}

		files = append(files, res.Body.Data.Files...)

		if pageNumber*pageSize >= *res.Body.Data.TotalCount {
			break
		}

		pageNumber += 1

	}

	return files, nil
}

func ListDIJobs() ([]*dataworks_public20200518.ListDIJobsResponseBodyDIJobPagingDIJobs, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	projectIdString := os.Getenv("DATAWORKS_PROJECT_ID")
	projectId, err := strconv.ParseInt(projectIdString, 10, 64)
	if err != nil {
		return nil, err
	}

	files := []*dataworks_public20200518.ListDIJobsResponseBodyDIJobPagingDIJobs{}
	var pageNumber int32 = 1
	var pageSize int32 = 100
	for {
		listFilesRequest := &dataworks_public20200518.ListDIJobsRequest{
			ProjectId:  &projectId,
			PageNumber: &pageNumber,
			PageSize:   &pageSize,
		}
		res, err := client.ListDIJobs(listFilesRequest)
		if err != nil {
			return nil, err
		}

		files = append(files, res.Body.DIJobPaging.DIJobs...)

		if pageNumber*pageSize >= *res.Body.DIJobPaging.TotalCount {
			break
		}

		pageNumber += 1

	}

	return files, nil
}

type NormalFile struct {
	LastEditTime   time.Time
	CreateTime     time.Time
	FolderId       string
	Content        string
	FileName       string
	LastEditUser   string
	CreateUser     string
	ConnectionName string
	FolderPath     string
	FileId         int64
	CommitStatus   int32
	FileType       int32
}

type Folder struct {
	FolderId   string
	FolderPath string
}

func uniqueString(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func GetFolders(projectId int64, folderIds []string) ([]Folder, error) {
	folderIds = uniqueString(folderIds)

	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	var folders []Folder

	for _, folderId := range folderIds {
		request := dataworks_public20200518.GetFolderRequest{
			ProjectId: &projectId,
			FolderId:  &folderId,
		}

		response, err := client.GetFolder(&request)
		if err != nil {
			return nil, err
		}

		folders = append(folders, Folder{
			FolderId:   *response.Body.Data.FolderId,
			FolderPath: *response.Body.Data.FolderPath,
		})

		time.Sleep(1 * time.Second)
	}

	return folders, nil
}

func ListFilesNormalized(fileTypes string) ([]NormalFile, error) {
	rawFiles, err := ListFiles(fileTypes)
	if err != nil {
		return nil, err
	}

	var folderIds []string
	for idx := range rawFiles {
		folderIds = append(folderIds, *rawFiles[idx].FileFolderId)
	}

	projectId, err := GetProjectId()
	if err != nil {
		return nil, err
	}

	folders, err := GetFolders(projectId, folderIds)
	if err != nil {
		return nil, err
	}

	var files []NormalFile

	for idx := range rawFiles {
		rawFile := rawFiles[idx]
		normalFile := NormalFile{
			FileId:         *rawFile.FileId,
			CommitStatus:   *rawFile.CommitStatus,
			FolderId:       *rawFile.FileFolderId,
			ConnectionName: *rawFile.ConnectionName,
			FileName:       *rawFile.FileName,
			FileType:       *rawFile.FileType,
			LastEditTime:   time.UnixMilli(*rawFile.LastEditTime),
			LastEditUser:   *rawFile.LastEditUser,
			CreateUser:     *rawFile.CreateUser,
			CreateTime:     time.UnixMilli(*rawFile.CreateTime),
			Content:        *rawFile.Content,
		}

		for folderIdx := range folders {
			folder := folders[folderIdx]
			if folder.FolderId == *rawFile.FileFolderId {
				normalFile.FolderPath = folder.FolderPath
			}
		}

		files = append(files, normalFile)
	}

	return files, nil
}

func GetFileContent(file NormalFile) (string, error) {
	client, err := CreateClient()
	if err != nil {
		return "", err
	}

	projectIdString := os.Getenv("DATAWORKS_PROJECT_ID")
	projectId, err := strconv.ParseInt(projectIdString, 10, 64)
	if err != nil {
		return "", err
	}

	getFilesRequest := dataworks_public20200518.GetFileRequest{
		ProjectId: &projectId,
		FileId:    &file.FileId,
	}

	res, err := client.GetFile(&getFilesRequest)
	if err != nil {
		return "", err
	}

	return *res.Body.Data.File.Content, nil
}

func DownloadFile(file NormalFile, directory string) error {
	folderPath := file.FolderPath
	if folderPath == "" {
		folderPath = file.FolderId
	}
	targetFolder := filepath.Join(directory, folderPath)

	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		return err
	}

	content, err := GetFileContent(file)
	if err != nil {
		return err
	}

	targetFilePath := filepath.Join(targetFolder, file.FileName)
	fileExt := GetFileExt(file.FileType)
	if fileExt != "" {
		targetFilePath = fmt.Sprintf("%s.%s", targetFilePath, fileExt)
	}
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = targetFile.Write([]byte(content))
	return err
}

func GetFileExt(fileType int32) string {
	// ODPS SQL
	if fileType == 10 {
		return "sql"
	}
	// 数据集成
	if fileType == 23 {
		return "json"
	}
	return ""
}

type Table struct {
	Name    string
	Guid    string
	Columns []TableColumn
}

type TableColumn struct {
	Caption           *string
	IsPrimaryKey      *bool
	Name              string
	Guid              string
	Comment           string
	Position          int32
	IsPartitionColumn bool
}

func ListTables(client *dataworks_public20200518.Client, appGuid string) ([]Table, error) {
	var pageNumber int32 = 1
	var pageSize int32 = 100

	dataSourceType := "odps"

	var tables []Table

	for {
		request := dataworks_public20200518.GetMetaDBTableListRequest{
			AppGuid:        &appGuid,
			PageNumber:     &pageNumber,
			PageSize:       &pageSize,
			DataSourceType: &dataSourceType,
		}

		res, err := client.GetMetaDBTableList(&request)
		if err != nil {
			return nil, err
		}

		for _, tableEntity := range res.Body.Data.TableEntityList {
			tables = append(tables, Table{
				Name: *tableEntity.TableName,
				Guid: *tableEntity.TableGuid,
			})
		}

		if int64(pageNumber*pageSize) >= *res.Body.Data.TotalCount {
			break
		}

		pageNumber++
		time.Sleep(500 * time.Millisecond)
	}

	return tables, nil
}

func GetTableFullInfo(client *dataworks_public20200518.Client, table *Table) (*Table, error) {
	dataSourceType := "odps"

	var pageNumber int32 = 1
	var pageSize int32 = 100

	var columns []TableColumn

	for {
		getMetaTableFullInfoRequest := dataworks_public20200518.GetMetaTableFullInfoRequest{
			DataSourceType: &dataSourceType,
			PageNum:        &pageNumber,
			PageSize:       &pageSize,
			TableGuid:      &table.Guid,
			TableName:      &table.Name,
		}

		res, err := client.GetMetaTableFullInfo(&getMetaTableFullInfoRequest)
		if err != nil {
			return nil, err
		}

		for _, c := range res.Body.Data.ColumnList {
			columns = append(columns, TableColumn{
				Caption:           c.Caption,
				Name:              *c.ColumnName,
				Guid:              *c.ColumnGuid,
				Position:          *c.Position,
				Comment:           *c.Comment,
				IsPrimaryKey:      c.IsPrimaryKey,
				IsPartitionColumn: *c.IsPartitionColumn,
			})
		}

		if int64(pageNumber*pageSize) >= *res.Body.Data.TotalColumnCount {
			break
		}

		pageNumber++
	}

	return &Table{
		Name:    table.Name,
		Guid:    table.Guid,
		Columns: columns,
	}, nil
}

func GetTables(appGuid string) ([]Table, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	tableList, err := ListTables(client, appGuid)
	if err != nil {
		return nil, err
	}

	var tables []Table
	for _, t := range tableList {
		table, err := GetTableFullInfo(client, &t)
		if err != nil {
			return nil, err
		}
		tables = append(tables, *table)
		time.Sleep(500 * time.Millisecond)
	}
	return tables, nil
}
