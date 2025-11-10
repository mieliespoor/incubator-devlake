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

package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.ToolLayerScope = (*OctopusProject)(nil)

type OctopusProject struct {
	common.Scope   `mapstructure:",squash" gorm:"embedded"`
	Id             string    `gorm:"primaryKey;type:VARCHAR(255) NOT NULL" json:"id"`
	ConnectionId   uint64    `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"connectionId"`
	Name           string    `gorm:"type:VARCHAR(255)" json:"name"`
	Description    string    `gorm:"type:TEXT" json:"description"`
	LastModifiedBy string    `gorm:"type:VARCHAR(255)" json:"lastModifiedBy"`
	LastModifiedOn time.Time `gorm:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP" json:"lastModifiedOn"`
	Slug           string    `gorm:"type:VARCHAR(255)" json:"slug"`
	SpaceId        string    `gorm:"type:VARCHAR(255)" json:"spaceId"`
}

func (OctopusProject) TableName() string {
	return "octopus_projects"
}

func (p OctopusProject) ScopeId() string {
	return p.Id
}

func (p OctopusProject) ScopeName() string {
	return p.Name
}

func (p OctopusProject) ScopeFullName() string {
	return p.Name
}

func (p OctopusProject) ScopeParams() interface{} {
	return &OctopusApiParams{
		ConnectionId: p.ConnectionId,
		ProjectId:    p.Id,
	}
}

type OctopusApiParams struct {
	ConnectionId uint64 `json:"connectionId"`
	Name         string `json:"name"`
	ProjectId    string `json:"projectId"`
}
