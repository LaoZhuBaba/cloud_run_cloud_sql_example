URL=$(gcloud run services list --format='value(status.url)' --filter example-app)
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" $URL | jq .