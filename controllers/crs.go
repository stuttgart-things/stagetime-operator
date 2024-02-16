/*
Copyright 2024 PATRICK HERMANN patrick.hermann@sva.de

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
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/samber/lo"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Repo struct {
	Url string `json:"url"`
}

type PipelineRunTemplate struct {
	Stage      int      `json:"stage"`
	Resolver   []string `json:"resolver"`
	Params     []string `json:"params"`
	ListParams []string `json:"listparams,omitempty"`
	Vclaims    []string `json:"vclaims,omitempty"`
}

var (
	stageTimeApi             = "stagetime.sthings.tiab.ssc.sva.de"
	stageTimeApiVersion      = "v1beta1"
	projectTemplateNamespace = "stagetime-operator-system"
	prTemplateResource       = "PipelineRunTemplate"
)

func ReadPipelineRunTemplate(crName string, r *RevisionRunReconciler) (bool, PipelineRunTemplate) {

	pipelineRunTemplate := PipelineRunTemplate{}
	repoSpec := getUnstructuredStructSpec(stageTimeApi, prTemplateResource, stageTimeApiVersion, crName, projectTemplateNamespace, r)

	fmt.Println(crName)
	fmt.Println(len(repoSpec))

	if len(repoSpec) <= 4 {
		fmt.Println(crName, " NOT FOUND")
		return false, pipelineRunTemplate
	} else {
		fmt.Println(crName, " FOUND")

		_ = json.Unmarshal(repoSpec, &pipelineRunTemplate)
		return true, pipelineRunTemplate
	}

}

func getUnstructuredStructSpec(group, kind, version, name, namespace string, r *RevisionRunReconciler) []byte {

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Kind:    kind,
		Version: version,
	})

	_ = r.Client.Get(context.Background(), client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, u)

	specContent := u.UnstructuredContent()["spec"]
	fmt.Println(specContent)

	interfaceBytes, _ := json.Marshal(specContent)

	return interfaceBytes
}

func renderParams(params []string, values map[string]interface{}) string {

	defaults := make(map[string]interface{})
	var remplateVars []string
	var paramsTemplate string

	// CREATE DEFAULTS AND TEMPLATE
	for _, data := range params {
		keyValues := strings.Split(data, "=")
		defaults[keyValues[0]] = keyValues[1]
		remplateVars = append(remplateVars, keyValues[0]+"={{ ."+keyValues[0]+" }}")
	}
	paramsTemplate = (strings.Join(remplateVars, ", "))

	// HERE THE DEFAULTS SHOULD BE OVERWRITTEN
	defaults = lo.Assign(defaults, values)

	fmt.Println("RESOLVER OVERWRITES")

	// RENDER TEMPLATE
	renderedParams, err := sthingsBase.RenderTemplateInline(paramsTemplate, "missingkey=zero", "{{", "}}", defaults)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	return string(renderedParams)
}

func createOverwriteParams(params string) (defaults map[string]interface{}) {

	defaults = make(map[string]interface{})

	paramsList := strings.Split(params, ";")
	fmt.Println(paramsList)

	// if string ends with ; we should cut if off here
	// CR working resolver: revision=main
	// CR not working resolver: revision=main; INDEX OUT OF LOOP

	for _, data := range paramsList {
		keyValues := strings.Split(data, "=")
		defaults[strings.TrimSpace(keyValues[0])] = strings.TrimSpace(keyValues[1])
	}

	return
}

func generateRandomRevisionRunID(n int, pool []rune) string {

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = pool[rand.Intn(len(pool))]
	}
	return string(b)
}
