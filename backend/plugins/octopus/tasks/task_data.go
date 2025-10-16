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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/octopus/models"
)

type OctopusTaskData struct {
	Options       *OctopusOptions
	ApiClient     *api.ApiAsyncClient
	RegexEnricher *api.RegexEnricher
}

type OctopusOptions struct {
	ConnectionId     uint64                     `json:"connectionId"`
	ApplicationName  string                     `json:"applicationName"`
	ScopeConfigId    uint64                     `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId,omitempty"`
	ScopeConfig      *models.OctopusScopeConfig `json:"scopeConfig,omitempty" mapstructure:"scopeConfig,omitempty"`
	CreatedDateAfter *time.Time                 `json:"createdDateAfter"`
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*OctopusOptions, errors.Error) {
	var opts OctopusOptions
	if err := api.Decode(options, &opts, nil); err != nil {
		return nil, err
	}
	if opts.ConnectionId == 0 {
		return nil, errors.Default.New("connectionId is required")
	}
	if opts.ApplicationName == "" {
		return nil, errors.Default.New("applicationName is required")
	}
	if opts.ScopeConfig == nil {
		opts.ScopeConfig = &models.OctopusScopeConfig{}
	}
	return &opts, nil
}
