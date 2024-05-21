package dataworks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
)

type Client struct {
	dwClient  *dataworks_public20200518.Client
	ProjectId int64
}

type NormalFile struct {
	LastEditTime   time.Time
	CreateTime     time.Time
	NodeId         *int64
	CreateUser     string
	FileName       string
	LastEditUser   string
	Content        string
	ConnectionName string
	FolderPath     string
	FileFolderId   string
	FileId         int64
	BusinessId     int64
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

func NewClient(accessKeyId, accessKeySecret string, endpoint string, projectId int64) (*Client, error) {
	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}

	config.Endpoint = &endpoint
	dwClient, err := dataworks_public20200518.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		dwClient:  dwClient,
		ProjectId: projectId,
	}, nil
}

func (client *Client) ListFiles(fileTypes string) ([]*dataworks_public20200518.ListFilesResponseBodyDataFiles, error) {
	files := []*dataworks_public20200518.ListFilesResponseBodyDataFiles{}
	var pageNumber int32 = 1
	var pageSize int32 = 100
	for {
		listFilesRequest := &dataworks_public20200518.ListFilesRequest{
			ProjectId:  &client.ProjectId,
			PageNumber: &pageNumber,
			PageSize:   &pageSize,
		}
		if fileTypes != "" {
			listFilesRequest.FileTypes = &fileTypes
		}

		res, err := client.dwClient.ListFiles(listFilesRequest)
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

func (client *Client) ListDIJobs() ([]*dataworks_public20200518.ListDIJobsResponseBodyDIJobPagingDIJobs, error) {
	files := []*dataworks_public20200518.ListDIJobsResponseBodyDIJobPagingDIJobs{}
	var pageNumber int32 = 1
	var pageSize int32 = 100
	for {
		listFilesRequest := &dataworks_public20200518.ListDIJobsRequest{
			ProjectId:  &client.ProjectId,
			PageNumber: &pageNumber,
			PageSize:   &pageSize,
		}
		res, err := client.dwClient.ListDIJobs(listFilesRequest)
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

func (client *Client) GetFolders(projectId int64, folderIds []string) ([]Folder, error) {
	folderIds = uniqueString(folderIds)

	var folders []Folder

	for _, folderId := range folderIds {
		request := dataworks_public20200518.GetFolderRequest{
			ProjectId: &projectId,
			FolderId:  &folderId,
		}

		response, err := client.dwClient.GetFolder(&request)
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

func (client *Client) ListFilesNormalized(fileTypes string) ([]NormalFile, error) {
	rawFiles, err := client.ListFiles(fileTypes)
	if err != nil {
		return nil, err
	}

	var folderIds []string
	for idx := range rawFiles {
		folderIds = append(folderIds, *rawFiles[idx].FileFolderId)
	}

	folders, err := client.GetFolders(client.ProjectId, folderIds)
	if err != nil {
		return nil, err
	}

	var files []NormalFile

	for idx := range rawFiles {
		rawFile := rawFiles[idx]
		normalFile := NormalFile{
			FileId:         *rawFile.FileId,
			CommitStatus:   *rawFile.CommitStatus,
			FileFolderId:   *rawFile.FileFolderId,
			ConnectionName: *rawFile.ConnectionName,
			FileName:       *rawFile.FileName,
			FileType:       *rawFile.FileType,
			LastEditTime:   time.UnixMilli(*rawFile.LastEditTime),
			LastEditUser:   *rawFile.LastEditUser,
			CreateUser:     *rawFile.CreateUser,
			CreateTime:     time.UnixMilli(*rawFile.CreateTime),
			Content:        *rawFile.Content,
			NodeId:         rawFile.NodeId,
			BusinessId:     *rawFile.BusinessId,
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

func (client *Client) GetFileContent(file NormalFile) (string, error) {
	getFilesRequest := dataworks_public20200518.GetFileRequest{
		ProjectId: &client.ProjectId,
		FileId:    &file.FileId,
	}

	res, err := client.dwClient.GetFile(&getFilesRequest)
	if err != nil {
		return "", err
	}

	return *res.Body.Data.File.Content, nil
}

func (client *Client) DownloadFile(file NormalFile, directory string) error {
	folderPath := file.FolderPath
	if folderPath == "" {
		folderPath = file.FileFolderId
	}
	targetFolder := filepath.Join(directory, folderPath)

	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		return err
	}

	content, err := client.GetFileContent(file)
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
