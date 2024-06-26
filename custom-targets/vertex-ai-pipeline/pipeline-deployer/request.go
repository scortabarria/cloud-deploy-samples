// Copyright 2023 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/cloud-deploy-samples/custom-targets/util/clouddeploy"
	"google.golang.org/api/aiplatform/v1"
	"os"
	//"strconv"
	"cloud.google.com/go/storage"
)

// Environment variable keys specific to the vertex ai deployer. These are provided via
// deploy parameters in Cloud Deploy.
const (
	pipelineEnvKey = "CLOUD_DEPLOY_customTarget_vertexAIPipeline"
	// parentEnvKey          = "CLOUD_DEPLOY_customTarget_vertexAIParent"
	configPathKey = "CLOUD_DEPLOY_customTarget_vertexAIPipelineJobConfiguration" //?????
	paramValsKey  = "CLOUD_DEPLOY_customTarget_vertexAIPipelineJobParameterValues"
	
)

// deploy parameters that the custom target requires to be present and provided during render and deploy operations.
const (
	pipelineDPKey = "customTarget/vertexAIPipeline"
	//parentDPKey   = "customTarget/vertexAIParent"
)

// requestHandler interface provides methods for handling the Cloud Deploy params.
type requestHandler interface {
	// Process processes the Cloud Deploy params.
	process(ctx context.Context) error
}

// createRequestHandler creates a requestHandler for the provided Cloud Deploy request.
func createRequestHandler(cloudDeployRequest interface{}, params *params, gcsClient *storage.Client, service *aiplatform.Service) (requestHandler, error) {

	switch r := cloudDeployRequest.(type) {
	case *clouddeploy.RenderRequest:
		return &renderer{
			req:               r,
			params:            params,
			gcsClient:         gcsClient,
			aiPlatformService: service,
		}, nil

	case *clouddeploy.DeployRequest:
		return &deployer{
			req:               r,
			params:            params,
			gcsClient:         gcsClient,
			aiPlatformService: service,
		}, nil

	default:
		return nil, fmt.Errorf("received unsupported cloud deploy request type: %q", os.Getenv(clouddeploy.RequestTypeEnvKey))
	}
}

// params contains the deploy parameter values passed into the execution environment.
type params struct {

	// The model to be deployed. May or may not contain a tag or version number.
	// format is "projects/{project}/locations/{location}/models/{modelId}[@versionId|alias].
	pipeline string

	// The endpoint where the model will be deployed
	// format is "projects/{project}/locations/{location}/endpoints/{endpointId}.
	modelParams string

	// directory path where the renderer should look for target-specific configuration
	// for this deployment, if not provided the renderer will check for a deployModel.yaml
	// fie in the root working directory.
	configPath string
}

// determineParams returns the supported params provided in the execution environment via environment variables.
func determineParams() (*params, error) {
	pipeline, found := os.LookupEnv(pipelineEnvKey)
	if !found {
		fmt.Printf("Required environment variable %s not found. This variable is derived from deploy parameter: %s, please verify that a valid Vertex AI model resource name was provided through this deploy parameter.\n", pipelineEnvKey, pipelineDPKey)
		return nil, fmt.Errorf("required environment variable %s not found", pipelineEnvKey)
	}
	if pipeline == "" {
		fmt.Printf("environment variable %s is empty. This variable is derived from deploy parameter: %s, please verify that a valid Vertex AI model resource name was provided through this deploy parameter.\n", pipelineEnvKey, pipelineDPKey)
		return nil, fmt.Errorf("environment variable %s contains empty string", pipelineEnvKey)
	}

	modelParams, found := os.LookupEnv(paramValsKey)
	if !found {
		fmt.Printf("Required environment variable %s not found. \n", paramValsKey)
		return nil, fmt.Errorf("required environment variable %s not found", pipelineEnvKey)
	}
	if len(paramValsKey) == 0 {
		fmt.Printf("environment variable %s is empty.\n", pipelineEnvKey)
		return nil, fmt.Errorf("environment variable %s contains empty string", pipelineEnvKey)
	}

	return &params{
		pipeline:    pipeline,
		modelParams: modelParams,
		configPath:  os.Getenv(configPathKey),
	}, nil
}
