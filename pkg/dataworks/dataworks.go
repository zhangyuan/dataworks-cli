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
