package dataworks

import (
	"fmt"
	"sort"
	"strings"

	dataworks_public20200518 "github.com/alibabacloud-go/dataworks-public-20200518/v6/client"
)

func (client *Client) ListNodes(projectEnv string) ([]*dataworks_public20200518.ListNodesResponseBodyDataNodes, error) {
	nodes := []*dataworks_public20200518.ListNodesResponseBodyDataNodes{}
	var pageNumber int32 = 1
	var pageSize int32 = 100
	for {
		request := &dataworks_public20200518.ListNodesRequest{
			ProjectId:  &client.ProjectId,
			PageNumber: &pageNumber,
			PageSize:   &pageSize,
			ProjectEnv: &projectEnv,
		}

		res, err := client.dwClient.ListNodes(request)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, res.Body.Data.Nodes...)

		if (*res.Body.Data.PageNumber * *res.Body.Data.PageSize) >= *res.Body.Data.TotalCount {
			break
		}

		pageNumber += 1

	}

	sortKey := func(node dataworks_public20200518.ListNodesResponseBodyDataNodes) string {
		keys := []string{}

		if node.BusinessId != nil {
			keys = append(keys, fmt.Sprintf("%d", *node.BusinessId))
		} else {
			keys = append(keys, "a")
		}

		if node.BaselineId != nil {
			keys = append(keys, fmt.Sprintf("%d", *node.BaselineId))
		} else {
			keys = append(keys, "a")
		}

		keys = append(keys, fmt.Sprintf("%s", *node.NodeName))

		key := strings.Join(keys, "/")
		return key
	}

	sort.Slice(nodes, func(i, j int) bool {
		return strings.Compare(sortKey(*nodes[i]), sortKey(*nodes[j])) > 0
	})
	return nodes, nil
}
