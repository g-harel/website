#!/bin/sh

# https://github.com/GoogleCloudPlatform/cloud-builders-community#build-the-build-step-from-source.

# Check that the configured project is correct.
PROJECT=$(gcloud config get-value "project")
if [ $PROJECT != "website-222818" ]; then
  echo "Unexpected gcloud project configured: $PROJECT"
  exit 1
fi

# Create gitignored workspace.
TEMPDIR=".temp"
rm -rf $TEMPDIR
mkdir -p $TEMPDIR
cd $TEMPDIR

# Build and publish new terraform release from head.
git clone "https://github.com/GoogleCloudPlatform/cloud-builders-community"
cd cloud-builders-community/terraform
gcloud builds submit --config cloudbuild.yaml .
