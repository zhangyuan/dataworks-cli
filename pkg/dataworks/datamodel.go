package dataworks

import (
	"fmt"

	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
	"github.com/mitchellh/mapstructure"
)

type RawDataModelColumn struct {
	ColumnCategory string
	ColumnCode     string
	ColumnName     string
	ColumnType     string
	ColumnUuid     string
	TableCode      string
	TableName      string
	TableUuid      string
}

type DataModel struct {
	TableUuid string
	TableCode string
	TableName string

	Columns []DataModelColumn
}

type DataModelColumn struct {
	ColumnCode     string
	ColumnUuid     string
	ColumnName     string
	ColumnCategory string
	ColumnTye      string
}

func (client *Client) ListDataModelColumns() ([]RawDataModelColumn, error) {
	query := "show full tables"
	projectId := fmt.Sprintf("%d", client.ProjectId)
	request := dataworks_public20200518.QueryPublicModelEngineRequest{
		ProjectId: &projectId,
		Text:      &query,
	}
	res, err := client.dwClient.QueryPublicModelEngine(&request)
	if err != nil {
		return nil, err
	}

	var columns []RawDataModelColumn

	if err := mapstructure.Decode(res.Body.ReturnValue, &columns); err != nil {
		return nil, err
	}

	return columns, nil
}

func (client *Client) ListDataModels() ([]DataModel, error) {
	columns, err := client.ListDataModelColumns()
	if err != nil {
		return nil, err
	}

	var models []DataModel

	for columnIdx := range columns {
		column := columns[columnIdx]

		var model *DataModel
		for modelIdx := range models {
			if column.TableCode == models[modelIdx].TableCode {
				model = &models[modelIdx]
			}
		}
		if model == nil {
			models = append(models, DataModel{
				TableUuid: column.TableUuid,
				TableCode: column.TableCode,
				TableName: column.TableName,
				Columns:   []DataModelColumn{},
			})
			model = &models[len(models)-1]
		}

		model.Columns = append(model.Columns, DataModelColumn{
			ColumnUuid:     column.ColumnUuid,
			ColumnCode:     column.ColumnCode,
			ColumnName:     column.ColumnName,
			ColumnCategory: column.ColumnCategory,
			ColumnTye:      column.ColumnType,
		})
	}

	return models, nil
}
