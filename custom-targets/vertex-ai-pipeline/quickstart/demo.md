```shell
gcloud deploy apply --file=$TMPDIR/clouddeploy.yaml --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION


gcloud deploy releases create release-001 \
   --delivery-pipeline=pipeline-cd \
   --project=$PIPELINE_PROJECT_ID \
   --region=$PIPELINE_REGION \
   --source=$TMPDIR/configuration \
   --deploy-parameters="customTarget/vertexAIPipeline=https://$PIPELINE_REGION-kfp.pkg.dev/$PIPELINE_PROJECT_ID/$REPO_ID/$PACKAGE_ID/$TAG_OR_VERSION"

```