package dataworks

import (
	"sort"
	"strings"

	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
)

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

func (client *Client) ListTables(appGuid string) ([]Table, error) {
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

		res, err := client.dwClient.GetMetaDBTableList(&request)
		if err != nil {
			return nil, err
		}

		for _, tableEntity := range res.Body.Data.TableEntityList {
			tables = append(tables, Table{
				Name: *tableEntity.TableName,
				Guid: *tableEntity.TableGuid,
			})
		}

		if int64(pageNumber**res.Body.Data.PageSize) >= *res.Body.Data.TotalCount {
			break
		}

		pageNumber++
		client.Wait()
	}

	return tables, nil
}

func (client *Client) GetTableFullInfo(table *Table) (*Table, error) {
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

		res, err := client.dwClient.GetMetaTableFullInfo(&getMetaTableFullInfoRequest)
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

func (client *Client) GetTables(appGuid string) ([]Table, error) {
	tableList, err := client.ListTables(appGuid)
	if err != nil {
		return nil, err
	}

	var tables []Table
	for _, t := range tableList {
		table, err := client.GetTableFullInfo(&t)
		if err != nil {
			return nil, err
		}
		tables = append(tables, *table)
		client.Wait()
	}

	sort.Slice(tables, func(i, j int) bool {
		return strings.Compare(tables[i].Guid, tables[j].Guid) > 0
	})
	return tables, nil
}
