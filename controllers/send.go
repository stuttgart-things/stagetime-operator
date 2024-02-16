/*
Copyright 2024 patrick.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Pipelinerun struct {
	Name                 string  `json:"name"`
	Canfail              bool    `json:"canfail"`
	Stage                float64 `json:"stage"`
	Params               string  `json:"params"`
	ResolverParams       string  `json:"resolverParams"`
	Listparams           string  `json:"listparams"`
	Workspaces           string  `json:"workspaces"`
	VolumeClaimTemplates string  `json:"volumeClaimTemplates"`
}

type RevisionRun struct {
	RepoName     string        `json:"repo_name"`
	PushedAt     string        `json:"pushed_at"`
	Author       string        `json:"author"`
	RepoUrl      string        `json:"repo_url"`
	CommitId     string        `json:"commit_id"`
	Pipelineruns []Pipelinerun `json:"pipelineruns"`
}

func ComposeRevisionRun(rr RevisionRun) (revisionRun []byte) {

	rrJson, err := json.Marshal(rr)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	return rrJson

}

func sendRevisionRun(revisionRunJson []byte) {

	fmt.Println(string(revisionRunJson))

	fmt.Println("CLIENT STARTED CONNECTING TO.. " + address)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	stsClient := NewClient(conn, time.Second)
	err = stsClient.CreateRevisionRun(context.Background(), bytes.NewBuffer(revisionRunJson))

	fmt.Println("ERR:", err)

}