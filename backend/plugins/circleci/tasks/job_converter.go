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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
)

var ConvertJobsMeta = plugin.SubTaskMeta{
	Name:             "convertJobs",
	EntryPoint:       ConvertJobs,
	EnabledByDefault: true,
	Description:      "convert circleci project",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertJobs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.CircleciJob{}),
		dal.Where("connection_id = ? AND project_slug = ?", data.Options.ConnectionId, data.Options.ProjectSlug),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.CircleciJob{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			userTool := inputRow.(*models.CircleciJob)
			task := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: getJobIdGen().Generate(data.Options.ConnectionId, userTool.WorkflowId, userTool.Id),
				},
				CicdScopeId:  getProjectIdGen().Generate(data.Options.ConnectionId, data.Project.Id),
				Name:         userTool.Name,
				PipelineId:   userTool.PipelineId,
				StartedDate:  userTool.StartedAt.ToTime(),
				FinishedDate: userTool.StoppedAt.ToNullableTime(),
				DurationSec:  userTool.DurationSec,
				// reference: https://circleci.com/docs/api/v2/index.html#operation/getJobDetails
				Status: devops.GetStatus(&devops.StatusRule[string]{
					Done:    []string{"canceled", "failed", "failing", "success", "not_run", "error", "infrastructure_fail", "timedout", "terminated-unknown"},
					Manual:  []string{"on_hold", "blocked"},
					Default: devops.STATUS_IN_PROGRESS,
				}, userTool.Status),
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{"success"},
					Failed:  []string{"failed", "failing", "error"},
					Skipped: []string{"not_run"},
					Abort:   []string{"canceled"},
				}, userTool.Status),
				Type:        data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, userTool.Name),
				Environment: data.RegexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, userTool.Name),
			}

			return []interface{}{
				task,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
