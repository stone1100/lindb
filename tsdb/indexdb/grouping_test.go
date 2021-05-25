// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package indexdb

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/series"
)

func TestGroupingContext_Build(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gCtx := series.NewMockGroupingContext(ctrl)
	gCtx.EXPECT().BuildGroup(gomock.Any(), gomock.Any()).Return(nil)

	ctx := &groupingContext{
		gCtx: gCtx,
	}
	ctx.BuildGroup(uint16(10), nil)
}
