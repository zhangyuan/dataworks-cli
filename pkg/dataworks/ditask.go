package dataworks

import (
	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
)

type DITask struct {
	DiDestinationDatasource string
	DiSourceDatasource      string
	NodeName                string
	TaskType                string
	NodeId                  int64
}

func (client *Client) ListDISyncTasks(taskType string, dataSourceName string) ([]*dataworks_public20200518.ListRefDISyncTasksResponseBodyDataDISyncTasks, error) {
	files := []*dataworks_public20200518.ListRefDISyncTasksResponseBodyDataDISyncTasks{}
	var pageNumber int64 = 1
	var pageSize int64 = 100
	var refType = "to"
	for {
		listFilesRequest := &dataworks_public20200518.ListRefDISyncTasksRequest{
			ProjectId:      &client.ProjectId,
			PageNumber:     &pageNumber,
			PageSize:       &pageSize,
			TaskType:       &taskType,
			DatasourceName: &dataSourceName,
			RefType:        &refType,
		}

		res, err := client.dwClient.ListRefDISyncTasks(listFilesRequest)
		if err != nil {
			return nil, err
		}

		files = append(files, res.Body.Data.DISyncTasks...)

		if len(res.Body.Data.DISyncTasks) == 0 {
			break
		}

		pageNumber += 1

	}

	return files, nil
}
