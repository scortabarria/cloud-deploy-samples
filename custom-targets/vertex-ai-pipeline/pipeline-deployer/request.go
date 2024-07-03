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
	"encoding/json"
)

// Environment variable keys specific to the vertex ai deployer. These are provided via
// deploy parameters in Cloud Deploy.
const (
	pipelineEnvKey = "CLOUD_DEPLOY_customTarget_vertexAIPipeline"
	// parentEnvKey          = "CLOUD_DEPLOY_customTarget_vertexAIParent"
	configPathKey  = "CLOUD_DEPLOY_customTarget_vertexAIPipelineJobConfiguration" //?????
	paramValsKey   = "CLOUD_DEPLOY_customTarget_vertexAIPipelineJobParameterValues"
	locValsKey     = "CLOUD_DEPLOY_customTarget_location"
	projectValsKey = "CLOUD_DEPLOY_customTarget_projectID"
	envValsKey = "CLOUD_DEPLOY_customTarget_environment"
)

// deploy parameters that the custom target requires to be present and provided during render and deploy operations.
// const (
// 	pipelineDPKey = "customTarget/vertexAIPipeline"
// 	//parentDPKey   = "customTarget/vertexAIParent"
// )

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
	project string

	pipeline string

	configPath string

	location string

	pipelineParams map[string]string

	environment string
}

// determineParams returns the supported params provided in the execution environment via environment variables.
func determineParams() (*params, error) {
	location, found := os.LookupEnv(locValsKey)
	if !found {
		fmt.Printf("Required environment variable %s not found. Please verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", locValsKey)
		return nil, fmt.Errorf("environment variable %s not found", locValsKey)
	}
	if location == "" {
		fmt.Printf("Environment variable %s is empty. Please verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", locValsKey)
		return nil, fmt.Errorf("environment variable %s contains empty string", locValsKey)
	}

	project, found := os.LookupEnv(projectValsKey)
	if !found {
		fmt.Printf("Required environment variable %s not found. Please verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", projectValsKey)
		return nil, fmt.Errorf("required environment variable %s not found", projectValsKey)
	}
	if project == "" {
		fmt.Printf("Environment variable %s is empty. lease verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", projectValsKey)
		return nil, fmt.Errorf("environment variable %s contains empty string", projectValsKey)
	}

	

	pipeline, found := os.LookupEnv(pipelineEnvKey)
	if !found {
		fmt.Printf("Required environment variable %s not found. Please verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", pipelineEnvKey)
		return nil, fmt.Errorf("required environment variable %s not found", pipelineEnvKey)
	}
	if pipeline == "" {
		fmt.Printf("environment variable %s is empty. Please verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", pipelineEnvKey)
		return nil, fmt.Errorf("environment variable %s contains empty string", pipelineEnvKey)
	}

	paramString, found := os.LookupEnv(paramValsKey)
	if !found {
		fmt.Printf("Required environment variable %s not found. \n", paramValsKey)
		return nil, fmt.Errorf("required environment variable %s not found", paramValsKey)
	}
	var pipelineParams map[string]string
	err := json.Unmarshal([]byte(paramString), &pipelineParams)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %s", err)
		return nil, fmt.Errorf("unable to unmarshal params json")
	}

	if len(pipelineParams) == 0 {
		fmt.Printf("environment variable %s is empty.\n", paramValsKey)
		return nil, fmt.Errorf("environment variable %s contains empty string", paramValsKey)
	}

	env, found := os.LookupEnv(envValsKey)
	if !found {
		fmt.Printf("Required environment variable %s not found. Please verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", envValsKey)
		return nil, fmt.Errorf("required environment variable %s not found", envValsKey)
	}
	if env == "" {
		fmt.Printf("environment variable %s is empty. Please verify that a valid Vertex AI pipeline resource name was provided through this deploy parameter.\n", envValsKey)
		return nil, fmt.Errorf("environment variable %s contains empty string", envValsKey)
	}

	return &params{
		project:      project,
		pipeline:    pipeline,
		configPath:  os.Getenv(configPathKey),
		location:    location,
		pipelineParams: pipelineParams,
		environment: env,
		
	}, nil
}
