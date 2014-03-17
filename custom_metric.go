// Copyright 2014, Belly, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stackdriver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Metric struct {
	// Timestamp in Unix time notation representing when the custom metric was collected.
	// If value is over an hour before current time the measurement will be rejected.
	CollectedAt int64 `json:"collected_at"`
	// Name of the custom metric as represented in the Stackdriver API.
	Name string
	// Measurement to record for the data point.
	Value interface{}
	// Metrics with a defined instance id show up under the defined instances resources.
	// One metric name can be shared across a number of instances to include on a single graph or for alerting.
	// Custom metrics not associated with an instance id will be found under the Custom resource type when creating
	// charts or alerting policies.
	// [Optional]
	InstanceId string `json:",omitempty"`
}

type GatewayMessage struct {
	// Timestamp the gateway message is created.
	Timestamp int64
	// Protocol version defining the schema of the gateway message.
	ProtocolVersion int64 `json:"proto_version"`
	// Stackdriver assigned Customer Id.
	// [Optional]
	CustomerId string `json:"customer_id,omitempty"`
	// Customer metrics to be sent to Stackdriver API.
	// Each data point must have its own (not necessarily unique) name, value, and collected_at.
	Data []Metric `json:"data"`
}

// Wrapper struct to properly marshal JSON with 'gateway_msg' root value.
type GatewayMessageObject struct {
	Message GatewayMessage `json:"gateway_msg"`
}

const (
	// Stackdriver custom metrics API protocol schema version.
	apiProtocolVersion = 1
	// Stackdriver custom metrics API URL.
	metricApiUrl = "https://custom-gateway.stackdriver.com/v1/custom"
)

// Factory function to create a new gateway message.
func NewGatewayMessage() *GatewayMessage {
	timestamp := time.Now().Unix()
	return &GatewayMessage{Timestamp: timestamp, ProtocolVersion: apiProtocolVersion}
}

// CustomMetric takes a name, instance id, collected-at and value to populates the data slice.
func (gwm *GatewayMessage) CustomMetric(n, id string, ca int64, v interface{}) error {
	if ca-time.Now().Unix() > 3600 {
		return fmt.Errorf("Metric created_at value is older than one hour.")
	}
	gwm.Data = append(gwm.Data, Metric{CollectedAt: ca, Name: n, Value: v, InstanceId: id})
	return nil
}

// Send utilizes HTTP POST to send all currently collected metrics to the Stackdriver API.
func (sdc *StackdriverClient) Send(gwm GatewayMessage) error {
	m := &GatewayMessageObject{Message: gwm}
	m.Message.CustomerId = sdc.CustomerId

	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", metricApiUrl, strings.NewReader(string(body)))
	req.Header.Add("user-agent", "Go Stackdriver API Library")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-stackdriver-apikey", sdc.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if (res.StatusCode > 200) || (res.StatusCode < 200) {
		return fmt.Errorf("Unable to send to Stackdriver API. HTTP response code: %d", res.StatusCode)
	}
	return nil
}
