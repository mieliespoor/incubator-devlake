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
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cast"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

type OperatedBy struct {
	Raw interface{}
	CoreOperatedBy
}

type CoreOperatedBy struct {
	Type     string
	Account  string
	RealName string
}

func (by *OperatedBy) Value() (driver.Value, error) {
	if by == nil {
		return nil, nil
	}
	b, err := json.Marshal(by.CoreOperatedBy)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (by *OperatedBy) Scan(v interface{}) error {
	switch value := v.(type) {
	case string:
		var coreOperatedBy CoreOperatedBy
		if err := json.Unmarshal([]byte(value), &coreOperatedBy); err != nil {
			return err
		}
		*by = OperatedBy{
			Raw:            nil,
			CoreOperatedBy: coreOperatedBy,
		}
	default:
		return fmt.Errorf("%+v is an unknown type, with value: %v", v, value)
	}
	return nil
}

func (by *OperatedBy) MarshalJSON() ([]byte, error) {
	if by == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(by.Raw)
}

func (by *OperatedBy) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	by.Raw = i
	switch i.(type) {
	case string:
		by.Type = "string"
		by.Account = cast.ToString(by.Raw)
		by.RealName = cast.ToString(by.Raw)
	default:
		by.Type = "struct"
		type ByUser struct {
			ID       int    `json:"id"`
			Account  string `json:"account"`
			Avatar   string `json:"avatar"`
			RealName string `json:"realname"`
		}
		var byUser ByUser
		if err := json.Unmarshal(data, &byUser); err != nil {
			return err
		}
		by.Account = byUser.Account
		by.RealName = byUser.RealName
	}
	return nil
}

func (by *OperatedBy) String() string {
	if by == nil {
		return "<nil>"
	}
	return by.Account
}

type ZentaoProject struct {
	common.Scope  `mapstructure:",squash"`
	Id            int64               `json:"id" mapstructure:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	Project       int64               `json:"project" mapstructure:"project"`
	Model         string              `json:"model" mapstructure:"model"`
	Type          string              `json:"type" mapstructure:"type"`
	ProjectType   string              `json:"projectType" mapstructure:"projectType"`
	Lifetime      string              `json:"lifetime" mapstructure:"lifetime"`
	Budget        string              `json:"budget" mapstructure:"budget"`
	BudgetUnit    string              `json:"budgetUnit" mapstructure:"budgetUnit"`
	Attribute     string              `json:"attribute" mapstructure:"attribute"`
	Percent       int                 `json:"percent" mapstructure:"percent"`
	Milestone     string              `json:"milestone" mapstructure:"milestone"`
	Output        string              `json:"output" mapstructure:"output"`
	Auth          string              `json:"auth" mapstructure:"auth"`
	Parent        int64               `json:"parent" mapstructure:"parent"`
	Path          string              `json:"path" mapstructure:"path"`
	Grade         int                 `json:"grade" mapstructure:"grade"`
	Name          string              `json:"name" mapstructure:"name"`
	Code          string              `json:"code" mapstructure:"code"`
	PlanBegin     *common.Iso8601Time `json:"begin" mapstructure:"begin"`
	PlanEnd       *common.Iso8601Time `json:"end" mapstructure:"end"`
	RealBegan     *common.Iso8601Time `json:"realBegan" mapstructure:"realBegan"`
	RealEnd       *common.Iso8601Time `json:"realEnd" mapstructure:"realEnd"`
	Days          int                 `json:"days" mapstructure:"days"`
	Status        string              `json:"status" mapstructure:"status"`
	SubStatus     string              `json:"subStatus" mapstructure:"subStatus"`
	Pri           string              `json:"pri" mapstructure:"pri"`
	Description   string              `json:"desc" mapstructure:"desc"`
	Version       int                 `json:"version" mapstructure:"version"`
	ParentVersion int                 `json:"parentVersion" mapstructure:"parentVersion"`
	PlanDuration  int                 `json:"planDuration" mapstructure:"planDuration"`
	RealDuration  int                 `json:"realDuration" mapstructure:"realDuration"`
	//OpenedBy       string    `json:"openedBy" mapstructure:"openedBy"`
	OpenedDate    *common.Iso8601Time `json:"openedDate" mapstructure:"openedDate"`
	OpenedVersion string              `json:"openedVersion" mapstructure:"openedVersion"`
	//LastEditedBy   string              `json:"lastEditedBy" mapstructure:"lastEditedBy"`
	LastEditedDate *common.Iso8601Time `json:"lastEditedDate" mapstructure:"lastEditedDate"`
	ClosedBy       *OperatedBy         `json:"closedBy" mapstructure:"closedBy" gorm:"-"`
	ClosedDate     *common.Iso8601Time `json:"closedDate" mapstructure:"closedDate"`
	CanceledBy     *OperatedBy         `json:"canceledBy" mapstructure:"canceledBy" gorm:"-"`
	CanceledDate   *common.Iso8601Time `json:"canceledDate" mapstructure:"canceledDate"`
	SuspendedDate  *common.Iso8601Time `json:"suspendedDate" mapstructure:"suspendedDate"`
	PO             string              `json:"po" mapstructure:"po"`
	PM             `json:"pm" mapstructure:"pm"`
	QD             string `json:"qd" mapstructure:"qd"`
	RD             string `json:"rd" mapstructure:"rd"`
	Team           string `json:"team" mapstructure:"team"`
	Acl            string `json:"acl" mapstructure:"acl"`
	Whitelist      `json:"whitelist" mapstructure:"" gorm:"-"`
	OrderIn        int    `json:"order" mapstructure:"order"`
	Vision         string `json:"vision" mapstructure:"vision"`
	DisplayCards   int    `json:"displayCards" mapstructure:"displayCards"`
	FluidBoard     string `json:"fluidBoard" mapstructure:"fluidBoard"`
	Deleted        bool   `json:"deleted" mapstructure:"deleted"`
	Delay          int    `json:"delay" mapstructure:"delay"`
	Hours          `json:"hours" mapstructure:"hours"`
	TeamCount      int    `json:"teamCount" mapstructure:"teamCount"`
	LeftTasks      string `json:"leftTasks" mapstructure:"leftTasks"`
	//TeamMembers   []interface{} `json:"teamMembers" gorm:"-"`
	TotalEstimate float64               `json:"totalEstimate" mapstructure:"totalEstimate"`
	TotalConsumed float64               `json:"totalConsumed" mapstructure:"totalConsumed"`
	TotalLeft     float64               `json:"totalLeft" mapstructure:"totalLeft"`
	Progress      *common.StringFloat64 `json:"progress" mapstructure:"progress" gorm:"-"`
}

type PM struct {
	PmId       int64  `json:"id" mapstructure:"id"`
	PmAccount  string `json:"account" mapstructure:"account"`
	PmAvatar   string `json:"avatar" mapstructure:"avatar"`
	PmRealname string `json:"realname" mapstructure:"realname"`
}
type Whitelist []struct {
	WhitelistID       int64  `json:"id" mapstructure:"id"`
	WhitelistAccount  string `json:"account" mapstructure:"account"`
	WhitelistAvatar   string `json:"avatar" mapstructure:"avatar"`
	WhitelistRealname string `json:"realname" mapstructure:"realname"`
}
type Hours struct {
	HoursTotalEstimate float64 `json:"totalEstimate" mapstructure:"totalEstimate"`
	HoursTotalConsumed float64 `json:"totalConsumed" mapstructure:"totalConsumed"`
	HoursTotalLeft     float64 `json:"totalLeft" mapstructure:"totalLeft"`
	HoursProgress      float64 `json:"progress" mapstructure:"progress"`
	HoursTotalReal     float64 `json:"totalReal" mapstructure:"totalReal"`
}

func (p ZentaoProject) TableName() string {
	return "_tool_zentao_projects"
}

func (p ZentaoProject) ScopeId() string {
	return strconv.FormatInt(p.Id, 10)
}

func (p ZentaoProject) ScopeName() string {
	return p.Name
}

func (p ZentaoProject) ScopeFullName() string {
	return p.Name
}

func (p ZentaoProject) ScopeParams() interface{} {
	return &ZentaoApiParams{
		ConnectionId: p.ConnectionId,
		ProjectId:    p.Id,
	}
}

func (p ZentaoProject) ConvertApiScope() plugin.ToolLayerScope {
	if p.ProjectType == "" {
		p.ProjectType = p.Type
		p.Type = "project"
	}
	return p
}

type ZentaoApiParams struct {
	ConnectionId uint64
	ProjectId    int64
}
