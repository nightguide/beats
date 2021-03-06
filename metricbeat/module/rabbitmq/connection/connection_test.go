// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
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

package connection

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
	"github.com/elastic/beats/metricbeat/module/rabbitmq/mtest"

	"github.com/stretchr/testify/assert"
)

func TestFetchEventContents(t *testing.T) {
	server := mtest.Server(t, mtest.DefaultServerConfig)
	defer server.Close()

	config := map[string]interface{}{
		"module":     "rabbitmq",
		"metricsets": []string{"connection"},
		"hosts":      []string{server.URL},
	}

	f := mbtest.NewEventsFetcher(t, config)
	events, err := f.Fetch()
	event := events[0]
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event.StringToPrint())

	assert.EqualValues(t, "[::1]:60938 -> [::1]:5672", event["name"])
	assert.EqualValues(t, "/", event["vhost"])
	assert.EqualValues(t, "guest", event["user"])
	assert.EqualValues(t, "nodename", event["node"])
	assert.EqualValues(t, 8, event["channels"])
	assert.EqualValues(t, 65535, event["channel_max"])
	assert.EqualValues(t, 131072, event["frame_max"])
	assert.EqualValues(t, "network", event["type"])

	packetCount := event["packet_count"].(common.MapStr)
	assert.EqualValues(t, 376, packetCount["sent"])
	assert.EqualValues(t, 376, packetCount["received"])
	assert.EqualValues(t, 0, packetCount["pending"])

	octetCount := event["octet_count"].(common.MapStr)
	assert.EqualValues(t, 3840, octetCount["sent"])
	assert.EqualValues(t, 3764, octetCount["received"])

	assert.EqualValues(t, "::1", event["host"])
	assert.EqualValues(t, 5672, event["port"])

	peer := event["peer"].(common.MapStr)
	assert.EqualValues(t, "::1", peer["host"])
	assert.EqualValues(t, 60938, peer["port"])
}
