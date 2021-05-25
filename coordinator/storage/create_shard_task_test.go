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

package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/service"
)

func TestCreateShardProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageService := service.NewMockStorageService(ctrl)
	processor := newCreateShardProcessor(storageService)
	assert.Equal(t, 1, processor.Concurrency())
	assert.Equal(t, time.Duration(0), processor.RetryBackOff())
	assert.Equal(t, 0, processor.RetryCount())
	assert.Equal(t, constants.CreateShard, processor.Kind())

	err := processor.Process(context.TODO(), task.Task{Params: []byte{1, 1, 1}})
	assert.NotNil(t, err)
	param := models.CreateShardTask{}
	storageService.EXPECT().CreateShards(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = processor.Process(context.TODO(), task.Task{Params: encoding.JSONMarshal(&param)})
	assert.NotNil(t, err)

	storageService.EXPECT().CreateShards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = processor.Process(context.TODO(), task.Task{Params: encoding.JSONMarshal(&param)})
	assert.Nil(t, err)
}
