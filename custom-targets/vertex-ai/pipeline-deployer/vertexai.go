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
	"google.golang.org/api/aiplatform/v1"
	"google.golang.org/api/option"
	"os"
	"sigs.k8s.io/yaml"
	// "strings"
)

// deployModelFromManifest loads the file provided in `path` and returns the parsed DeployModelRequest
// from the data.
func pipelineRequestFromManifest(path string) (*aiplatform.GoogleCloudAiplatformV1CreatePipelineJobRequest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest file: %v", err)
	}

	createPipelineRequest := &aiplatform.GoogleCloudAiplatformV1CreatePipelineJobRequest{}

	if err = yaml.Unmarshal(data, createPipelineRequest); err != nil {
		return nil, fmt.Errorf("unable to parse createPipelineRequest from manifest file: %v", err)
	}

	return createPipelineRequest, nil
}



// newAIPlatformService generates a Service that can make API calls in the specified region.
func newAIPlatformService(ctx context.Context, region string) (*aiplatform.Service, error) {
	endPointOption := option.WithEndpoint(fmt.Sprintf("%s-aiplatform.googleapis.com", region))
	regionalService, err := aiplatform.NewService(ctx, endPointOption)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate")
	}

	return regionalService, nil
}

// fetchModel calls the aiplatform API to fetch the Vertex AI model using the given model name.
func fetchPipeline(service *aiplatform.Service, pipelineName string) (*aiplatform.GoogleCloudAiplatformV1PipelineJob, error) {
	pipeline, err := service.Projects.Locations.PipelineJobs.Get(pipelineName).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline job: %v", err)
	}
	return pipeline, nil
}

//NOT SURE IF THIS WORKS!!!!!
// deployModel performs the DeployModel request and awaits the resulting operation until it completes, it times out or an error occurs.
func deployPipeline(ctx context.Context, aiPlatformService *aiplatform.Service, parent string, request *aiplatform.GoogleCloudAiplatformV1CreatePipelineJobRequest) error {
	op, err := aiPlatformService.Projects.Locations.PipelineJobs.Create(parent, request.PipelineJob).Do()

	if err != nil {
		return fmt.Errorf("unable to deploy model: %v", err)
	}

	return poll(ctx, aiPlatformService, op)
}

func regionFromPipeline(pipelineName string) (string, error) {
	matches := pipelineRegex.FindStringSubmatch(pipelineName)
	if len(matches) == 0 {
		return "", fmt.Errorf("unable to parse pipeline name")
	}

	return matches[2], nil
}

