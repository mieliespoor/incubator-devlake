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

package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
)

type StringFloat64 struct {
	v float64
}

func NewStringFloat64FromAny(f interface{}) *StringFloat64 {
	return &StringFloat64{
		v: cast.ToFloat64(f),
	}
}

func (f *StringFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.v)
}

func (f *StringFloat64) String() string {
	return fmt.Sprintf("%v", f.v)
}

func (f *StringFloat64) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	value, err := cast.ToFloat64E(i)
	if err != nil {
		return err
	}
	f.v = value
	return nil
}

func (f *StringFloat64) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}
	return f.v, nil
}

func (f *StringFloat64) Scan(v interface{}) error {
	switch value := v.(type) {
	case float64:
		*f = StringFloat64{
			v: value,
		}
	case string:
		*f = StringFloat64{
			v: cast.ToFloat64(value),
		}
	default:
		return fmt.Errorf("%+v is an unknown type, with value: %v", v, value)
	}
	return nil
}
