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
	"io"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	stagetimev1beta1 "github.com/stuttgart-things/stagetime-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
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

var (
	address            = "stagetime-server-service.stagetime.svc.cluster.local:80"
	resolverOverwrites = make(map[string]interface{})
	paramsOverwrites   = make(map[string]interface{})
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

	// revisionRunCR := &stagetimev1beta1.RevisionRun{}
	// err := r.Get(ctx, req.NamespacedName, revisionRunCR)
	// fmt.Println(revisionRunCR.Spec.TechnologyConfig)

	// GET RIVISION RUN
	revisionRun := &stagetimev1beta1.RevisionRun{}
	err := r.GetRevisionRun(ctx, req, revisionRun)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("RevisionRun RESOURCE NOT FOUND. IGNORING SINCE OBJECT MUST BE DELETED")

			return ctrl.Result{}, nil
		}

		log.Error(err, "FAILED TO GET RevisionRun")

		return ctrl.Result{}, err
	}

	// TRY TO SET INITIAL CONDITION STATUS
	err = r.SetInitialCondition(ctx, req, revisionRun)
	if err != nil {
		log.Error(err, "FAILED TO SET INITIAL CONDITION")

		return ctrl.Result{}, err
	}

	// READ TECHNOLOGIES
	for _, config := range revisionRun.Spec.TechnologyConfig {
		fmt.Println(config.ID)
		fmt.Println(config.Kind)

		// READ TEMPLATE - GET DEFAULTS
		prExists, pipelineRunTemplate := ReadPipelineRunTemplate(config.Kind, r)
		fmt.Println(prExists, pipelineRunTemplate)

		// CHECK IF TEMPLATE EXISTS - IF NOT SKIP
		if !prExists {
			break
		} else {
			if config.Resolver != "" {
				resolverOverwrites = createOverwriteParams(config.Resolver)
				log.Info("RESOLVER OVERWRITES")
				resolverParams := renderParams(pipelineRunTemplate.Resolver, resolverOverwrites)
				fmt.Println(resolverParams)

			}

			if config.Params != "" {
				fmt.Println(config.Params)
				paramsOverwrites = createOverwriteParams(config.Params)
				log.Info("PARAMS OVERWRITES")
			}

			// CHECK FOR OVERWRITES
			if config.Path != "" {
				fmt.Println(config.Path)
			}

			if config.Stage != 99 {
				fmt.Println(config.Stage)
			}

			if config.Listparams != "" {
				fmt.Println(config.Listparams)
			}

			if config.Vclaims != "" {
				fmt.Println(config.Vclaims)
			}

		}
	}

	// DO ORDERING OF PIPELINERUNS TO STAGES

	// READ CRDS
	// repo2 := Repo{}
	// repoSpec := getUnstructuredStructSpec("stagetime.sthings.tiab.ssc.sva.de", "Repo", "v1beta1", "repo-sample", "stagetime-operator-system", r)

	// _ = json.Unmarshal(repoSpec, &repo2)

	// fmt.Println("REPOOO-URL", repo2.Url)

	pipelineRunTemplate := PipelineRunTemplate{}
	pipelineRunTemplateSpec := getUnstructuredStructSpec("stagetime.sthings.tiab.ssc.sva.de", "PipelineRunTemplate", "v1beta1", "pipelineruntemplate-sample", "stagetime-operator-system", r)

	_ = json.Unmarshal(pipelineRunTemplateSpec, &pipelineRunTemplate)

	fmt.Println("RESOLVER", pipelineRunTemplate.Resolver)
	fmt.Println("STAGE", pipelineRunTemplate.Stage)
	fmt.Println("PARAMS", pipelineRunTemplate.Params)
	fmt.Println("VCLAIMS", pipelineRunTemplate.Vclaims)
	fmt.Println("LISTPARAMS", pipelineRunTemplate.ListParams)

	// var resolverTemplateVars []string
	// var resolverTemplate string

	// for _, data := range pipelineRunTemplate.Resolver {
	// 	resolverKV := strings.Split(data, "=")
	// 	resolverDefaults[resolverKV[0]] = resolverKV[1]
	// 	resolverTemplateVars = append(resolverTemplateVars, resolverKV[0]+"={{ ."+resolverKV[0]+" }}")
	// }
	// resolverTemplate = (strings.Join(resolverTemplateVars, ", "))

	// resolverParams, err := sthingsBase.RenderTemplateInline(resolverTemplate, "missingkey=zero", "{{", "}}", resolverDefaults)

	// if err != nil {
	// 	fmt.Printf("Error: %s", err)
	// }

	// fmt.Println(string(resolverParams))

	// resolverParams := renderParams(pipelineRunTemplate.Resolver)
	// fmt.Println(resolverParams)

	// renderedParams := renderParams(pipelineRunTemplate.Params)
	// fmt.Println(renderedParams)

	// renderedVclaims := renderParams(pipelineRunTemplate.Vclaims)
	// fmt.Println(renderedVclaims)

	revisionRunJson := ComposeRevisionRun()

	sendRevisionRun(revisionRunJson)

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
