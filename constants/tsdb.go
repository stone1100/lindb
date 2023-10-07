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

package constants

const (
	// MaxSuggestions represents the max number of suggestions count
	MaxSuggestions = 100

	// MetricMaxAheadDuration controls the global max write ahead duration.
	// If current timestamp is 2021-08-19 23:00:00, metric after 2021-08-20 23:00:00 will be dropped.
	MetricMaxAheadDuration    = int64(24 * 60 * 60 * 1000)
	MetricMaxAheadDurationStr = "1d"
	// MetricMaxBehindDuration controls the global max write behind duration.
	// If current timestamp is 2021-08-19 23:00:00, metric before 2021-08-18 23:00:00 will be dropped.
	MetricMaxBehindDuration    = int64(24 * 60 * 60 * 1000)
	MetricMaxBehindDurationStr = "1d"

	// NamespaceID represents global namespace store id in kv.
	NamespaceID = uint32(1)
)
