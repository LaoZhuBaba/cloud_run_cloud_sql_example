EXECUTABLE=example-app
rm -f $EXECUTABLE
gcloud builds submit --region=australia-southeast1 --tag australia-southeast1-docker.pkg.dev/qoria-sandbox/sandbox-repo/$EXECUTABLE
