#!/bin/bash

export _CT_IMAGE_NAME=vertexai

while getopts "s:r:p:o:t:b:f:m:l:d:" arg; do
  case "${arg}" in
    s)
      STAGING_PROJECT="${OPTARG}"
      ;;
    r)
      STAGING_REGION="${OPTARG}"
      ;;
    p)
      PROD_PROJECT="${OPTARG}"
      ;;
    o)
      PROD_REGION="${OPTARG}"
      ;;
    t)
      TMPDIR="${OPTARG}"
      ;;
    b)
      BUCKET="${OPTARG}"
      ;;
    f)
      PREFERENCE="${OPTARG}"
      ;;
    m)
      PROMPT="${OPTARG}"
      ;;
    l)
      MODEL_REFERENCE="${OPTARG}"
      ;;
    d)
      DISPLAY="${OPTARG}"
      ;;
    *)
      usage
      exit 1
      ;;
  esac
done

if [[ ! -v STAGING_PROJECT || ! -v STAGING_REGION || ! -v PROD_PROJECT || ! -v PROD_REGION || ! -v TMPDIR || ! -v BUCKET || ! -v PREFERENCE || ! -v PROMPT || ! -v MODEL_REFERENCE|| ! -v DISPLAY ]]; then
  usage
  exit 1
fi

# get the location where the custom image was uploaded
AR_REPO=$STAGING_REGION-docker.pkg.dev/$STAGING_PROJECT/cd-custom-targets

# get the image digest of the most recently built image
IMAGE_SHA=$(gcloud -q artifacts docker images describe "${AR_REPO}/${_CT_IMAGE_NAME}:latest" --format 'get(image_summary.digest)')


cp clouddeploy.yaml "$TMPDIR"/clouddeploy.yaml
cp -r configuration "$TMPDIR"/configuration

# replace variables in clouddeploy.yaml with actual values
sed -i "s/\$STAGING_PROJECT_ID/${STAGING_PROJECT}/g" "$TMPDIR"/clouddeploy.yaml
sed -i "s/\$STAGING_REGION/${STAGING_REGION}/g" "$TMPDIR"/clouddeploy.yaml
sed -i "s/\$PROD_PROJECT_ID/${PROD_PROJECT}/g" "$TMPDIR"/clouddeploy.yaml
sed -i "s/\$PROD_REGION/${PROD_REGION}/g" "$TMPDIR"/clouddeploy.yaml
sed -i "s|\$PREFERENCE_DATASET|${PREFERENCE}|g" "$TMPDIR"/clouddeploy.yaml
sed -i "s|\$PROMPT_DATASET|${PROMPT}|g" "$TMPDIR"/clouddeploy.yaml
sed -i "s/\$LARGE_MODEL_REFERENCE/${MODEL_REFERENCE}/g" "$TMPDIR"/clouddeploy.yaml
sed -i "s|\$DISPLAY_NAME|${DISPLAY}|g" "$TMPDIR"/clouddeploy.yaml


# replace variables in configuration/skaffold.yaml with actual values
sed -i "s/\$STAGING_REGION/${STAGING_REGION}/g" "$TMPDIR"/configuration/skaffold.yaml
sed -i "s/\$PROJECT_ID/${STAGING_PROJECT}/g" "$TMPDIR"/configuration/skaffold.yaml
sed -i "s/\$_CT_IMAGE_NAME/${_CT_IMAGE_NAME}/g" "$TMPDIR"/configuration/skaffold.yaml
sed -i "s/\$IMAGE_SHA/${IMAGE_SHA}/g" "$TMPDIR"/configuration/skaffold.yaml

# replace variables in configuration/staging/pipelineJob.yaml
sed -i "s/\$BUCKET_NAME/${BUCKET}/g" "$TMPDIR"/configuration/staging/pipelineJob.yaml




