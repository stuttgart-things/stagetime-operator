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
	"context"
	"encoding/json"
	"fmt"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	stagetimev1beta1 "github.com/stuttgart-things/stagetime-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RevisionRunReconciler reconciles a RevisionRun object
type RevisionRunReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

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

type Repo struct {
	Url string `json:"url"`
}

//+kubebuilder:rbac:groups=stagetime.sthings.tiab.ssc.sva.de,resources=revisionruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stagetime.sthings.tiab.ssc.sva.de,resources=revisionruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stagetime.sthings.tiab.ssc.sva.de,resources=revisionruns/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RevisionRun object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *RevisionRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := ctrllog.FromContext(ctx)
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "stagetime.sthings.tiab.ssc.sva.de",
		Kind:    "Repo",
		Version: "v1beta1",
	})

	_ = r.Client.Get(context.Background(), client.ObjectKey{
		Name:      "repo-sample",
		Namespace: "stagetime-operator-system",
	}, u)

	fmt.Println(u.GetObjectKind())
	fmt.Println(u.UnstructuredContent())
	fmt.Println(u.GetName())
	fmt.Println(u.GetCreationTimestamp())

	fmt.Println(u.UnstructuredContent()["spec"])
	spec := u.UnstructuredContent()["spec"]

	fmt.Println("SPEC!", spec)

	repo := Repo{}
	dbByte, _ := json.Marshal(spec)
	_ = json.Unmarshal(dbByte, &repo)

	fmt.Println(repo.Url)

	for k, v := range u.Object {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}

	// maps := make(map[string]interface{})

	// u.Object, _ = runtime.DefaultUnstructuredConverter.ToUnstructured(maps)

	// fmt.Println(maps)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RevisionRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stagetimev1beta1.RevisionRun{}).
		Complete(r)
}
