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
	"fmt"
	"time"
)

type OctopusDeployment struct {
	Id                 string     `gorm:"primaryKey;type:VARCHAR(255) NOT NULL" json:"id"`
	ConnectionId       uint64     `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"connectionId"`
	ProjectId          string     `gorm:"type:VARCHAR(255)" json:"projectId"`
	EnvironmentId      string     `gorm:"type:VARCHAR(255)" json:"environmentId"`
	ReleaseId          string     `gorm:"type:VARCHAR(255)" json:"releaseId"`
	Created            time.Time  `gorm:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP" json:"created"`
	DeployedBy         string     `gorm:"type:VARCHAR(255)" json:"deployedBy"`
	FailureEncountered bool       `gorm:"type:BOOLEAN" json:"failureEncountered"`
	LastModifiedBy     string     `gorm:"type:VARCHAR(255)" json:"lastModifiedBy"`
	LastModifiedOn     time.Time  `gorm:"type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP" json:"lastModifiedOn"`
	Name               string     `gorm:"type:VARCHAR(255)" json:"name"`
	QueueTime          *time.Time `gorm:"type:BIGINT" json:"queueTime,omitempty"`
	SpaceId            string     `gorm:"type:VARCHAR(255)" json:"spaceId"`
	TaskId             string     `gorm:"type:VARCHAR(255)" json:"taskId"`
	TenantId           string     `gorm:"type:VARCHAR(255)" json:"tenantId"`
	UseGuidedFailure   bool       `gorm:"type:BOOLEAN" json:"useGuidedFailure"`
}

func (OctopusDeployment) TableName() string {
	return "octopus_deployments"
}

func (d OctopusDeployment) GetDeploymentKey() string {
	return fmt.Sprintf("%d-%s", d.ConnectionId, d.Id)
}
