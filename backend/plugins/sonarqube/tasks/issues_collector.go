/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_ISSUES_TABLE = "sonarqube_api_issues"

var _ plugin.SubTaskEntryPoint = CollectIssues

type SonarqubeIssueIteratorNode struct {
	Severity      string
	Status        string
	Type          string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
}

func CollectIssues(taskCtx plugin.SubTaskContext) (err errors.Error) {
	logger := taskCtx.GetLogger()
	logger.Info("collect issues")

	iterator := helper.NewQueueIterator()
	severities := []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"}
	statuses := []string{"OPEN", "CONFIRMED", "REOPENED", "RESOLVED", "CLOSED"}
	types := []string{"BUG", "VULNERABILITY", "CODE_SMELL"}
	for _, severity := range severities {
		for _, status := range statuses {
			for _, typ := range types {
				iterator.Push(
					&SonarqubeIssueIteratorNode{
						Severity:      severity,
						Status:        status,
						Type:          typ,
						CreatedAfter:  nil,
						CreatedBefore: nil,
					},
				)
			}
		}
	}
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUES_TABLE)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		UrlTemplate:        "issues/search",
		Input:              iterator,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("componentKeys", fmt.Sprintf("%v", data.Options.ProjectKey))
			input, ok := reqData.Input.(*SonarqubeIssueIteratorNode)
			if !ok {
				return nil, errors.Default.New(fmt.Sprintf("Input to SonarqubeIssueIteratorNode failed:%+v", reqData.Input))
			}
			query.Set("severities", reqData.Input.(*SonarqubeIssueIteratorNode).Severity)
			query.Set("statuses", reqData.Input.(*SonarqubeIssueIteratorNode).Status)
			query.Set("types", reqData.Input.(*SonarqubeIssueIteratorNode).Type)
			if input.CreatedAfter != nil {
				query.Set("createdAfter", GetFormatTime(input.CreatedAfter))
			}
			if input.CreatedBefore != nil {
				query.Set("createdBefore", GetFormatTime(input.CreatedBefore))
			}
			query.Set("p", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("ps", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Encode()
			return query, nil
		},
		GetTotalPages: func(res *http.Response, args *helper.ApiCollectorArgs) (int, errors.Error) {
			body := &SonarqubePagination{}
			err := helper.UnmarshalResponse(res, body)
			if err != nil {
				return 0, err
			}

			pages := body.Paging.Total / args.PageSize
			if body.Paging.Total%args.PageSize > 0 {
				pages++
			}

			query := res.Request.URL.Query()

			// if get more than 10000 data, that need split it
			if pages > 100 {
				var createdAfterUnix int64
				var createdBeforeUnix int64

				createdAfter, err := getTimeFromFormatTime(query.Get("createdAfter"))
				if err != nil {
					return 0, err
				}
				createdBefore, err := getTimeFromFormatTime(query.Get("createdBefore"))
				if err != nil {
					return 0, err
				}

				if createdAfter == nil {
					createdAfterUnix = 0
				} else {
					createdAfterUnix = createdAfter.Unix()
				}

				if createdBefore == nil {
					createdBeforeUnix = time.Now().Unix()
				} else {
					createdBeforeUnix = createdBefore.Unix()
				}

				// can not split it any more
				if createdBeforeUnix-createdAfterUnix <= 10 {
					return 100, nil
				}

				// split it
				MidTime := time.Unix((createdAfterUnix+createdBeforeUnix)/2+1, 0)

				severity := query.Get("severities")
				status := query.Get("statuses")
				typ := query.Get("types")
				// left part
				iterator.Push(&SonarqubeIssueIteratorNode{
					Severity:      severity,
					Status:        status,
					Type:          typ,
					CreatedAfter:  createdAfter,
					CreatedBefore: &MidTime,
				})

				// right part
				iterator.Push(&SonarqubeIssueIteratorNode{
					Severity:      severity,
					Status:        status,
					Type:          typ,
					CreatedAfter:  &MidTime,
					CreatedBefore: createdBefore,
				})

				logger.Info("split [%s][%s] by mid [%s] for it has pages:[%d] and total:[%d]",
					query.Get("createdAfter"), query.Get("createdBefore"), GetFormatTime(&MidTime), pages, body.Paging.Total)

				return 0, nil
			} else {
				logger.Info("[%s][%s] has pages:[%d] and total:[%d]",
					query.Get("createdAfter"), query.Get("createdBefore"), pages, body.Paging.Total)

				return pages, nil
			}
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resData struct {
				Data []json.RawMessage `json:"issues"`
			}
			err = helper.UnmarshalResponse(res, &resData)
			if err != nil {
				return nil, err
			}

			// check if sonar report updated during collecting
			var issue struct {
				UpdateDate *common.Iso8601Time `json:"updateDate"`
			}
			for _, v := range resData.Data {
				err = errors.Convert(json.Unmarshal(v, &issue))
				if err != nil {
					return nil, err
				}
				if issue.UpdateDate.ToTime().After(data.TaskStartTime) {
					return nil, errors.Default.New(fmt.Sprintf(`Your data is affected by the latest analysis\n
						Please recollect this project: %s`, data.Options.ProjectKey))
				}
			}

			return resData.Data, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()

}

var CollectIssuesMeta = plugin.SubTaskMeta{
	Name:             "CollectIssues",
	EntryPoint:       CollectIssues,
	EnabledByDefault: true,
	Description:      "Collect issues data from Sonarqube api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func GetFormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05-0700")
}

func getTimeFromFormatTime(formatTime string) (*time.Time, errors.Error) {
	if formatTime == "" {
		return nil, nil
	}

	if len(formatTime) < 20 {
		return nil, errors.Default.New(fmt.Sprintf("formatTime [%s] is too short ", formatTime))
	}

	t, err := time.Parse("2006-01-02T15:04:05-0700", formatTime)

	if err != nil {
		return nil, errors.Default.New(fmt.Sprintf("Failed to Parse the time from [%s]:%s", formatTime, err.Error()))
	}

	return &t, nil
}
