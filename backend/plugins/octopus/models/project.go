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

import "github.com/apache/incubator-devlake/core/plugin"

var _ plugin.ToolLayerScope = (*OctopusProject)(nil)

type OctopusProject struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"connectionId"`
	ProjectId    string `gorm:"primaryKey;type:VARCHAR(255) NOT NULL" json:"projectId"`
	Name         string `gorm:"type:VARCHAR(255)" json:"name"`
	Description  string `gorm:"type:TEXT" json:"description"`
	Created      string `gorm:"type:VARCHAR(255)" json:"created"`
	CreatedBy    string `gorm:"type:VARCHAR(255)" json:"createdBy"`
	Updated      string `gorm:"type:VARCHAR(255)" json:"updated"`
	UpdatedBy    string `gorm:"type:VARCHAR(255)" json:"updatedBy"`

	// add more fields if necessary
}

// ScopeConnectionId implements plugin.ToolLayerScope.
func (p *OctopusProject) ScopeConnectionId() uint64 {
	return p.ConnectionId
}

// ScopeScopeConfigId implements plugin.ToolLayerScope.
func (p *OctopusProject) ScopeScopeConfigId() uint64 {
	return 0
}

func (OctopusProject) TableName() string {
	return "octopus_projects"
}

func (p OctopusProject) ScopeId() string {
	return p.ProjectId
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
		ProjectId:    p.ProjectId,
	}
}

type OctopusApiParams struct {
	ConnectionId uint64 `json:"connectionId"`
	Name         string `json:"name"`
	ProjectId    string `json:"projectId"`
}
