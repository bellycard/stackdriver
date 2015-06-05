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
	"io/ioutil"
	"net/http"
	"strings"
)

type AnnotationEvent struct {
	// Contents of the message, in plain text. Limited to 256 characters.
	Message string `json:"message"`
	// Person or robot who the attribution should be attributed.
	// [Optional]
	AnnotatedBy string `json:"annotated_by,omitempty"`
	// Event status level. Options are: INFO, WARN, ERROR.
	// [Default: INFO]
	// [Optional]
	Level string `json:"level,omitempty"`
	// Events with a defined instance id show up under the defined instances context.
	// [Optional]
	InstanceId string `json:"instance_id,omitempty"`
	// Unix timestamp of where the event should appear in the timeline.
	// [Default: Now]
	// [Optional]
	EventEpoch int64 `json:"event_epoch"`
}

const (
	// Stackdriver event gateway API URL.
	eventGatewayApiUrl = "https://event-gateway.stackdriver.com/v1/annotationevent"
)

func (sdc *StackdriverClient) NewAnnotationEvent(m, ab, l, iid string, ee int64) error {
	var aem string

	// Stackdriver messags cannot be over 256 characters, if larger truncate.
	if len(m) > 256 {
		aem = m[0:252] + "..."
	} else {
		aem = m
	}

	ae := &AnnotationEvent{Message: aem,
		AnnotatedBy: ab,
		Level:       l,
		InstanceId:  iid,
		EventEpoch:  ee,
	}

	body, err := json.Marshal(ae)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", eventGatewayApiUrl, strings.NewReader(string(body)))
	req.Header.Add("user-agent", "Go Stackdriver API Library")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-stackdriver-apikey", sdc.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if (res.StatusCode > 200) || (res.StatusCode < 200) {
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("AnnotationEvent Stackdriver Error: %s", responseBody)
	}
	return nil
}
