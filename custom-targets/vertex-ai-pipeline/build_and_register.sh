#!/bin/bash

SOURCE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

export _CT_SRCDIR="${SOURCE_DIR}/pipeline-deployer"
export _CT_IMAGE_NAME=vertexai
export _CT_TYPE_NAME=vertex-ai-pipeline
export _CT_CUSTOM_ACTION_NAME=vertex-ai-pipeline-deployer
export _CT_GCS_DIRECTORY=vertexai
export _CT_SKAFFOLD_CONFIG_NAME=vertexAiConfig

# gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
#     --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
#     --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
#     --role="roles/clouddeploy.jobRunner"

# gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
#     --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
#     --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
#     --role="roles/clouddeploy.viewer"

# gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
#     --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
#     --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
#     --role="roles/aiplatform.user"

# gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
#     --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
#     --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
#     --role="roles/artifactregistry.writer"

"${SOURCE_DIR}/../util/build_and_register.sh" "$@"

