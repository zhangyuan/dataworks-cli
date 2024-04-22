package dataworks

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
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

func ListFiles() ([]*dataworks_public20200518.ListFilesResponseBodyDataFiles, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, err
	}

	projectIdString := os.Getenv("DATAWORKS_PROJECT_ID")
	projectId, err := strconv.ParseInt(projectIdString, 10, 64)
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

type NormalFile struct {
	FileId         int64
	CommitStatus   int32
	FolderId       string
	Content        string
	FileName       string
	FileType       int32
	LastEditTime   time.Time
	LastEditUser   string
	CreateTime     time.Time
	CreateUser     string
	ConnectionName string
}

func GetScriptsWithContent() ([]NormalFile, error) {
	allFiles, err := ListFiles()
	if err != nil {
		return nil, err
	}

	normalFiles := lo.Filter(allFiles, func(item *dataworks_public20200518.ListFilesResponseBodyDataFiles, index int) bool {
		return *item.UseType == NOMARL_USE_TYPE && !lo.Contains([]int{23, 99, 1119}, int(*item.FileType))
	})

	files := lop.Map(normalFiles, func(x *dataworks_public20200518.ListFilesResponseBodyDataFiles, _ int) NormalFile {
		return NormalFile{
			FileId:         *x.FileId,
			CommitStatus:   *x.CommitStatus,
			FolderId:       *x.FileFolderId,
			ConnectionName: *x.ConnectionName,
			FileName:       *x.FileName,
			FileType:       *x.FileType,
			LastEditTime:   time.UnixMilli(*x.LastEditTime),
			LastEditUser:   *x.LastEditUser,
			CreateUser:     *x.CreateUser,
			CreateTime:     time.UnixMilli(*x.CreateTime),
			Content:        *x.Content,
		}
	})
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
	targetFolder := filepath.Join(directory, fmt.Sprintf("%s.%s", file.ConnectionName, file.FolderId))
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
	return ""
}

type Table struct {
	Name    string
	Guid    string
	Columns []TableColumn
}

type TableColumn struct {
	Caption           *string
	Name              string
	Guid              string
	Comment           string
	IsPartitionColumn bool
	IsPrimaryKey      *bool
	Position          int32
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
