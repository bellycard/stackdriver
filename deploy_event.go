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
)

type DeployEvent struct {
	// Revision of the code that was deployed.
	RevisionId string `json:"revision_id"`
	// Person or robot responsible for deploying the code.
	// [Optional]
	DeployedBy string `json:"deployed_by,omitempty"`
	// Environment code was deployed to. (ie: development, staging, production)
	// [Optional]
	DeployedTo string `json:"deployed_to,omitempty"`
	// Repository (or project) deployed.
	// [Optional]
	Repository string `json:"repository,omitempty"`
}

const (
	// Stackdriver deploy event gateway API URL.
	deployEventGatewayApiUrl = "https://event-gateway.stackdriver.com/v1/deployevent"
)

func (sdc *StackdriverClient) NewDeployEvent(rid, db, dt, r string) error {

	de := &DeployEvent{RevisionId: rid,
		DeployedBy: db,
		DeployedTo: dt,
		Repository: r,
	}

	body, err := json.Marshal(de)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", deployEventGatewayApiUrl, strings.NewReader(string(body)))
	req.Header.Add("user-agent", "Go Stackdriver API Library")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-stackdriver-apikey", sdc.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if (res.StatusCode > 200) || (res.StatusCode < 200) {
		return fmt.Errorf("Unable to send to Stackdriver deploy event gateway API. HTTP response code: %d", res.StatusCode)
	}
	return nil
}
