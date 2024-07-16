```shell
../build_and_register.sh -p $PIPELINE_PROJECT_ID -r $PIPELINE_REGION


export TMPDIR=$(mktemp -d)
./replace_variables.sh -s $STAGING_PROJECT_ID -r $STAGING_REGION -p $PROD_PROJECT_ID -o $PROD_REGION -t $TMPDIR -b $STAGING_BUCKET_NAME -c $PROD_BUCKET_NAME -f $STAGING_PREF_DATA -m $STAGING_PROMPT_DATA -l $LARGE_MODEL_REFERENCE -d $MODEL_DISPLAY_NAME -y $PROD_PREF_DATA -z $PROD_PROMPT_DATA -e $STAGING_PROJECT_NUMBER -g $PROD_PROJECT_NUMBER


gcloud deploy apply --file=$TMPDIR/clouddeploy.yaml --project=$PIPELINE_PROJECT_ID --region=$PIPELINE_REGION


gcloud deploy releases create release-001 \
   --delivery-pipeline=pipeline-cd \
   --project=$PIPELINE_PROJECT_ID \
   --region=$PIPELINE_REGION \
   --source=$TMPDIR/configuration \
   --deploy-parameters="customTarget/vertexAIPipeline=https://$PIPELINE_REGION-kfp.pkg.dev/$PIPELINE_PROJECT_ID/$REPO_ID/$PACKAGE_ID/$TAG_OR_VERSION"


gcloud deploy releases promote \
--release=release-001 \
--delivery-pipeline=pipeline-cd \
--to-target=prod-environment \
--project=$PIPELINE_PROJECT_ID \
--region=$PIPELINE_REGION
    
```