# Cloud Deploy Vertex AI Pipeline Deployer Quickstart

## Overview

This quickstart demonstrates how to deploy a ML Pipeline to an target environment using a Cloud Deploy custom target.


## 1. Clone Repository

Clone this repository and navigate to the quickstart directory (`cloud-deploy-samples/custom-targets/vertex-ai-pipeline/quickstart`) since the commands provided expect to be executed from that directory.

## 2. Environment variables

To simplify the commands in this quickstart, set the following environment variables with your values:

```shell
export PIPELINE_PROJECT_ID="YOUR_PIPELINE_PROJECT_ID"
export PIPELINE_REGION="YOUR_PIPELINE_REGION"
export PIPELINE_PROJECT_NUMBER=$(gcloud projects list \
        --format="value(projectNumber)" \
        --filter="projectId=${PIPELINE_PROJECT_ID}")

export STAGING_PROJECT_ID="YOUR_STAGING_PROJECT_ID"
export STAGING_REGION="YOUR_STAGING_REGION"
export STAGING_PROJECT_NUMBER=$(gcloud projects list \
        --format="value(projectNumber)" \
        --filter="projectId=${STAGING_PROJECT_ID}")
export STAGING_BUCKET_NAME="YOUR_STAGING_BUCKET_NAME"
export STAGING_PREF_DATA="YOUR_STAGING_PREFERENCE_DATASET"
export STAGING_PROMPT_DATA="YOUR_STAGING_PROMPT_DATASET"

export PROD_PROJECT_ID="YOUR_PROD_PROJECT_ID"
export PROD_REGION="YOUR_PROD_REGION"
export PROD_PROJECT_NUMBER=$(gcloud projects list \
        --format="value(projectNumber)" \
        --filter="projectId=${PROD_PROJECT_ID}")
export PROD_BUCKET_NAME="YOUR_PROD_BUCKET_NAME"
export PROD_PREF_DATA="YOUR_PROD_PREFERENCE_DATASET"
export PROD_PROMPT_DATA="YOUR_PROD_PROMPT_DATASET"

export REPO_ID="YOUR_REPO"
export PACKAGE_ID="YOUR_PACKAGE"
export TAG_OR_VERSION="YOUR_TAG_OR_VERSION"
export LARGE_MODEL_REFERENCE="YOUR_LARGE_MODEL_REFERENCE"
export MODEL_DISPLAY_NAME="YOUR_DISPLAY_NAME"
```

```shell
export PIPELINE_PROJECT_ID="scortabarria-internship"
export PIPELINE_REGION="us-central1"
export PIPELINE_PROJECT_NUMBER=$(gcloud projects list \
        --format="value(projectNumber)" \
        --filter="projectId=${PIPELINE_PROJECT_ID}")

export STAGING_PROJECT_ID="imara-staging"
export STAGING_REGION="us-central1"
export STAGING_PROJECT_NUMBER=$(gcloud projects list \
        --format="value(projectNumber)" \
        --filter="projectId=${STAGING_PROJECT_ID}")
export STAGING_BUCKET_NAME="imara-staging-pipeline-artifacts-scorta"
export STAGING_PREF_DATA="gs://imara-staging-rlhf-artifacts/data/preference/*.jsonl"
export STAGING_PROMPT_DATA="gs://imara-staging-rlhf-artifacts/data/prompt/*.jsonl"

export PROD_PROJECT_ID="imara-prod"
export PROD_REGION="us-central1"
export PROD_PROJECT_NUMBER=$(gcloud projects list \
        --format="value(projectNumber)" \
        --filter="projectId=${PROD_PROJECT_ID}")
export PROD_BUCKET_NAME="imara-prod-pipeline-artifacts-scorta"
export PROD_PREF_DATA="gs://imara-prod-rlhf-artifacts/data/preference/*.jsonl"
export PROD_PROMPT_DATA="gs://imara-prod-rlhf-artifacts/data/prompt/*.jsonl"

export REPO_ID="scortabarria-internship-rlhf-pipelines"
export PACKAGE_ID="rlhf-tune-pipeline"
export TAG_OR_VERSION="sha256:e739c5c310d406f8a6a9133b0c97bf9a249715da0a507505997ced042e3e0f17"
export LARGE_MODEL_REFERENCE="text-bison@001"
export MODEL_DISPLAY_NAME="$PIPELINE_REGION/$PIPELINE_PROJECT_ID"

```
## 3. Prerequisites

[Install](https://cloud.google.com/sdk/docs/install) the latest version of the Google Cloud CLI


### APIs
Enable the Cloud Deploy API, Compute Engine API, and Vertex AI API.

   ```shell
   gcloud services enable clouddeploy.googleapis.com aiplatform.googleapis.com compute.googleapis.com --project $PIPELINE_PROJECT_ID
   ```

    ```shell
   gcloud services enable clouddeploy.googleapis.com aiplatform.googleapis.com compute.googleapis.com --project $STAGING_PROJECT_ID
   ```

   ```shell
   gcloud services enable clouddeploy.googleapis.com aiplatform.googleapis.com compute.googleapis.com --project $PROD_PROJECT_ID
   ```

### Permissions
The default service account, `{project_num}-compute@developer.gserviceaccount.com`, used by Cloud Deploy needs the
   following roles:

1. `roles/clouddeploy.jobRunner` - required by Cloud Deploy

   ```shell
   gcloud projects add-iam-policy-binding $PIPELINE_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $PIPELINE_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/clouddeploy.jobRunner"
   ```
2. `roles/clouddeploy.viewer` - required to access Cloud Deploy resources

   ```shell
   gcloud projects add-iam-policy-binding $PIPELINE_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $PIPELINE_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/clouddeploy.viewer"
   ```
3. `roles/aiplatform.user` - required to access the models and deploy endpoints in the custom target

   ```shell
   gcloud projects add-iam-policy-binding $PIPELINE_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $PIPELINE_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/aiplatform.user"

   ```


   ```shell
   

   gcloud projects add-iam-policy-binding $PROD_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $PROD_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/aiplatform.user"

    gcloud iam service-accounts add-iam-policy-binding \
        ${PROJECT_NUMBER}-compute@developer.gserviceaccount.com \
        --member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
        --role=roles/iam.serviceAccountUser \
        --project=${STAGING_PROJECT_ID}

    gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/clouddeploy.jobRunner"

    gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/clouddeploy.viewer"

    gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/aiplatform.user"

    gcloud projects add-iam-policy-binding $STAGING_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/artifactregistry.writer"



    gcloud iam service-accounts add-iam-policy-binding \
        ${PROJECT_NUMBER}-compute@developer.gserviceaccount.com \
        --member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
        --role=roles/iam.serviceAccountUser \
        --project=${PIPELINE_PROJECT_ID}

    gcloud projects add-iam-policy-binding $PIPELINE_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/clouddeploy.jobRunner"

    gcloud projects add-iam-policy-binding $PIPELINE_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/clouddeploy.viewer"

    gcloud projects add-iam-policy-binding $PIPELINE_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/aiplatform.user"

     gcloud projects add-iam-policy-binding $PIPELINE_PROJECT_ID \
       --member=serviceAccount:$(gcloud projects describe $STAGING_PROJECT_ID \
       --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
       --role="roles/artifactregistry.writer"
   ```


## 4. Create a bucket

From the `quickstart` directory, run this command to create a bucket in Cloud Storage:

```shell
# gsutil mb -l $PIPELINE_REGION -p $PIPELINE_PROJECT_ID gs://$PIPELINE_BUCKET_NAME

gsutil mb -l $STAGING_REGION -p $STAGING_PROJECT_ID gs://$STAGING_BUCKET_NAME

gsutil mb -l $PROD_REGION -p $PROD_PROJECT_ID gs://$PROD_BUCKET_NAME

```


## 5. Build and Register a Custom Target Type for Vertex AI

From within the `quickstart` directory, run this command to build the Vertex AI model deployer image and
install the custom target resources:

```shell
../build_and_register.sh -p $PIPELINE_PROJECT_ID -r $PIPELINE_REGION
```


For information about the `build_and_register.sh` script, see the [README](../README.md#build)


## 6. Create delivery pipeline, target, and skaffold

Within the `quickstart` directory, run this second command to make a temporary copy of `clouddeploy.yaml`, `configuration/skaffold.yaml` and
`configuration/staging/pipelineJob.yaml`, and to replace placeholders in the copies with actual values:

```shell
export TMPDIR=$(mktemp -d)
./replace_variables.sh -s $STAGING_PROJECT_ID -r $STAGING_REGION -p $PROD_PROJECT_ID -o $PROD_REGION -t $TMPDIR -b $STAGING_BUCKET_NAME -c $PROD_BUCKET_NAME -f $STAGING_PREF_DATA -m $STAGING_PROMPT_DATA -l $LARGE_MODEL_REFERENCE -d $MODEL_DISPLAY_NAME -y $PROD_PREF_DATA -z $PROD_PROMPT_DATA -e $STAGING_PROJECT_NUMBER -g $PROD_PROJECT_NUMBER
```

<!-- ```shell
export TMPDIR=$(mktemp -d)
./replace_variables.sh -a $PIPELINE_PROJECT_ID -c $PIPELINE_REGION -s $STAGING_PROJECT_ID -r $STAGING_REGION -p $PROD_PROJECT_ID -o $PROD_REGION -t $TMPDIR -b $PIPELINE_BUCKET_NAME -f $STAGING_PREF_DATA -m $STAGING_PROMPT_DATA -l $LARGE_MODEL_REFERENCE -d $MODEL_DISPLAY_NAME
``` -->

The command does the following:
1. Creates temporary directory $TMPDIR and copies `clouddeploy.yaml` and `configuration` into it.
2. Replaces the placeholders in `$TMPDIR/clouddeploy.yaml`, `configuration/skaffold.yaml`, and `configuration/staging/pipelineJob.yaml`
3. Obtains the URL of the latest version of the custom image, built in step 6, and sets it in `$TMPDIR/configuration/skaffold.yaml`

Lastly, apply the Cloud Deploy configuration defined in `clouddeploy.yaml`:

```shell
gcloud deploy apply --file=$TMPDIR/clouddeploy.yaml --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION
```


## 7. Create a release and rollout

Create a Cloud Deploy release for the configuration defined in the `configuration` directory. This automatically
creates a rollout that deploys the first model version to the target.

```shell
gcloud deploy releases create release-001 \
    --delivery-pipeline=pipeline-cd \
    --project=$PIPELINE_PROJECT_ID \
    --region=$PIPELINE_REGION \
    --source=$TMPDIR/configuration \
    --deploy-parameters="customTarget/vertexAIPipeline=https://$PIPELINE_REGION-kfp.pkg.dev/$PIPELINE_PROJECT_ID/$REPO_ID/$PACKAGE_ID/$TAG_OR_VERSION"
```
```shell
gcloud deploy delete --file=$TMPDIR/clouddeploy.yaml --force --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION
```

### Explanation of command line flags

The `--source` command line flag instructs gcloud where to look for the configuration files relative to the working directory where the command is run.

The `--deploy-parameters` flag is used to provide the custom deployer with additional parameters needed to perform the deployment.

Here, we are providing the custom deployer with deploy parameter `customTarget/vertexAIModel`
which specifies the full resource name of the model to deploy

The remaining flags specify the Cloud Deploy Delivery Pipeline. `--delivery-pipeline` is the name of
the delivery pipeline where the release will be created, and the project and region of the pipeline
is specified by `--project` and `--region` respectively.


### Monitor the release's progress

To check release details, run this command:

```shell
gcloud deploy releases describe release-001 --delivery-pipeline=pipeline-cd --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION
```

Run this command to filter only the render status of the release:

```shell
gcloud deploy releases describe release-001 --delivery-pipeline=pipeline-cd --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION --format "(renderState)"
```


## 8. Monitor rollout status

In the [Cloud Deploy UI](https://cloud.google.com/deploy) for your project click on the
`pipeline-cd` delivery pipeline. Here you can see the release created and the rollout to the target for the release.

You can also describe the rollout created using the following command:

```shell
gcloud deploy rollouts describe release-001-to-staging-environment-0001 --release=release-001 --delivery-pipeline=pipeline-cd --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION
```

 
## 9. Promote a release

This promotes the release, automatically moving it to the next target environment.

```shell
gcloud deploy releases promote \
    --release=release-001 \
    --delivery-pipeline=pipeline-cd \
    --to-target=prod-environment \
    --project=$PIPELINE_PROJECT_ID \
    --region=$PIPELINE_REGION
```


### Monitor the release's progress

To check release details, run this command:

```shell
gcloud deploy releases describe release-001 --delivery-pipeline=pipeline-cd --project=$PROD_PROJECT_ID --region=$PROD_REGION
```

Run this command to filter only the render status of the release:

```shell
gcloud deploy releases describe release-001 --delivery-pipeline=pipeline-cd --project=$PROD_PROJECT_ID --region=$PROD_REGION --format "(renderState)"
```


## 10. Monitor rollout status

In the [Cloud Deploy UI](https://cloud.google.com/deploy) for your project click on the
`pipeline-cd` delivery pipeline. Here you can see the release created and the rollout to the target for the release.

You can also describe the rollout created using the following command:

```shell
gcloud deploy rollouts describe release-001-to-staging-environment-0001 --release=release-001 --delivery-pipeline=pipeline-cd --project=$PROD_PROJECT_ID --region=$PROD_REGION
```


## 11. Clean up

To delete Cloud Deploy resources:

```shell
gcloud deploy delete --file=$TMPDIR/clouddeploy.yaml --force --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION
```

```shell
gcloud deploy delete --file=$TMPDIR/clouddeploy.yaml --force --project=$PROD_PROJECT_ID --region=$PROD_REGION
```
