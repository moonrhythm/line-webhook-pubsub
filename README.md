# line-webhook-pubsub

Send events from line webhook to pubsub. 

## Deployment

### Cloud Run

```shell script
gcloud beta run deploy line-webhook-pubsub \
  --async \
  --platform=managed \
  --region=asia-northeast1 \
  --concurrency=80 \
  --allow-unauthenticated \
  --timeout=300 \
  --memory=128Mi \
  --image=gcr.io/moonrhythm-containers/line-webhook-pubsub \
  --set-env-vars=LINE_CHANNEL_SECRET=SECRET,PUBSUB_URL=PUBSUB_URL \
  --service-account=SERVICE_ACCOUNT
```
