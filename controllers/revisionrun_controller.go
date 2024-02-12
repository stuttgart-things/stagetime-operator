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
	"io"
	"time"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	stagetimev1beta1 "github.com/stuttgart-things/stagetime-operator/api/v1beta1"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	//"google.golang.org/grpc/credentials"
	revisionrun "github.com/stuttgart-things/stageTime-server/revisionrun"

	"github.com/golang/protobuf/jsonpb"
	//google.golang.org/protobuf/encoding/protojson
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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

var (
	address = "stagetime-server-service.stagetime.svc.cluster.local:80"
)

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

	// READ UNSTRUCTURED STRUCT

	// u := &unstructured.Unstructured{}
	// u.SetGroupVersionKind(schema.GroupVersionKind{
	// 	Group:   "stagetime.sthings.tiab.ssc.sva.de",
	// 	Kind:    "Repo",
	// 	Version: "v1beta1",
	// })

	// _ = r.Client.Get(context.Background(), client.ObjectKey{
	// 	Name:      "repo-sample",
	// 	Namespace: "stagetime-operator-system",
	// }, u)

	// fmt.Println(u.UnstructuredContent()["spec"])
	// repo1 := Repo{}

	// // _ = json.Unmarshal(dbByte, &repo)

	// spec := u.UnstructuredContent()["spec"]
	// specByte, _ := json.Marshal(spec)
	// _ = json.Unmarshal(specByte, &repo1)

	// fmt.Println("REPO1", repo1)

	repo2 := Repo{}
	repoSpec := getUnstructuredStructSpec("stagetime.sthings.tiab.ssc.sva.de", "Repo", "v1beta1", "repo-sample", "stagetime-operator-system", r)

	_ = json.Unmarshal(repoSpec, &repo2)
	fmt.Println(repo2)

	fmt.Println("REPO2", repo2)
	fmt.Println("REPOOO-URL", repo2.Url)

	revisionRunJson := ComposeRevisionRun()

	sendRevisionRun(revisionRunJson)

	// fmt.Println(string(revisionRunJson))

	// fmt.Println("CLIENT STARTED CONNECTING TO.. " + address)

	// conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer conn.Close()

	// stsClient := NewClient(conn, time.Second)
	// err = stsClient.CreateRevisionRun(context.Background(), bytes.NewBuffer(revisionRunJson))

	// fmt.Println("ERR:", err)

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

type Client struct {
	stsClient revisionrun.StageTimeApplicationServiceClient
	timeout   time.Duration
}

func NewClient(conn grpc.ClientConnInterface, timeout time.Duration) Client {
	return Client{
		stsClient: revisionrun.NewStageTimeApplicationServiceClient(conn),
		timeout:   timeout,
	}
}

func (c Client) CreateRevisionRun(ctx context.Context, json io.Reader) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(c.timeout))
	defer cancel()

	req := revisionrun.CreateRevisionRunRequest{}
	if err := jsonpb.Unmarshal(json, &req); err != nil {
		return fmt.Errorf("CLIENT CREATE REVISIONRUN: UNMARSHAL: %w", err)
	}

	fmt.Println(req.Pipelineruns)
	res, err := c.stsClient.CreateRevisionRun(ctx, &req)

	fmt.Println(res)

	if err != nil {
		if er, ok := status.FromError(err); ok {
			return fmt.Errorf("CLIENT CREATE REVISIONRUN: CODE: %s - msg: %s", er.Code(), er.Message())
		}
		return fmt.Errorf("CLIENT CREATE REVISIONRUn: %w", err)
	}

	fmt.Println("RESULT:", res.Result)
	fmt.Println("RESPONSE:", res)

	return nil
}

func ComposeRevisionRun() (revisionRun []byte) {

	var prs []Pipelinerun

	pr := Pipelinerun{
		Name:                 "simulate-stagetime",
		Canfail:              true,
		Stage:                0,
		ResolverParams:       "url=https://github.com/stuttgart-things/stuttgart-things.git, revision=main, pathInRepo=stageTime/pipelines/simulate-stagetime-pipelineruns.yaml",
		Params:               "gitRevision=main, gitRepoUrl=https://github.com/stuttgart-things/stageTime-server.git, gitWorkspaceSubdirectory=stageTime, scriptPath=tests/prime.sh, scriptTimeout=25s",
		Listparams:           "",
		VolumeClaimTemplates: "source=openebs-hostpath;ReadWriteOnce;20Mi",
	}

	prs = append(prs, pr)

	rr := RevisionRun{
		RepoName:     "stuttgart-things",
		PushedAt:     "2024-01-13T13:40:36Z",
		Author:       "patrick-hermann-sva",
		RepoUrl:      "https://codehub.sva.de/Lab/stuttgart-things/stuttgart-things.git",
		CommitId:     "135866d3v43453",
		Pipelineruns: prs,
	}

	rrJson, err := json.Marshal(rr)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	return rrJson

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

func sendRevisionRun(revisionRunJson []byte) {

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
