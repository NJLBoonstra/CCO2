#!/bin/bash

cd ./sorter-backend
gcloud functions deploy HandleUpload --trigger-event google.storage.object.finalize --trigger-resource boonstra-nieuwenhuijzen.appspot.com --runtime go116 --region europe-west1 --env-vars-file .env.yaml
gcloud functions deploy PartialSort --region europe-west1 --trigger-topic sortJobs --runtime go116 --env-vars-file .env.yaml
gcloud functions deploy MergeSort --region europe-west1 --trigger-topic reduceJobs --runtime go116 --env-vars-file .env.yaml
gcloud functions deploy JobRequest --trigger-http --runtime go116 --region europe-west1 --env-vars-file .env.yaml
gcloud functions deploy FindPalindromes --region europe-west1 --trigger-event google.storage.object.metadataUpdate --trigger-resource sort-chunks --runtime go116 --env-vars-file .env.yaml
