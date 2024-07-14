```shell
   
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