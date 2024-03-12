package dataworks

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
)

func CreateClient() (_result *dataworks_public20200518.Client, _err error) {
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")

	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}

	endpoint := "dataworks.cn-beijing.aliyuncs.com"
	config.Endpoint = &endpoint
	return dataworks_public20200518.NewClient(config)
}

func GetFiles() error {
	client, err := CreateClient()
	if err != nil {
		return err
	}

	projectIdString := os.Getenv("PROJECT_ID")
	projectId, err := strconv.ParseInt(projectIdString, 10, 64)
	if err != nil {
		return err
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
			return nil
		}
		files = append(files, res.Body.Data.Files...)

		if pageNumber*pageSize >= *res.Body.Data.TotalCount {
			break
		}

		pageNumber += 1

	}

	bytes, err := json.Marshal(files)
	if err != nil {
		return err
	}

	documents := string(bytes)

	fmt.Println(documents)

	return nil

}
