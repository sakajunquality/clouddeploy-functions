# clouddeploy-functions

PoC Automation of Cloud Deploy w/ Cloud Functions

## Functions

### AutoPromote

Automatically promote rollout to next stage in a serialPipeline.

#### How to Deploy

1. Prepare new service account for this function.
2. Add `roles/clouddeploy.releaser` permission at project level.
3. Add `roles/iam.serviceAccountUser` permission at Cloud Deploy execution service account.
4. Deploy Cloud Function with entrypoint `AutoPromote` and subscribe to `clouddeploy-operations` topic.

#### Example Deploy

```bash
gcloud functions deploy cloud-deploy-auto-promote \
    --gen2 \
    --runtime=go119 \
    --project=[PROJECT] \
    --region=[REGION] \
    --source=. \
    --entry-point=AutoPromote \
    --trigger-topic=clouddeploy-operations \
    --service-account=[Service Account] \
    --serve-all-traffic-latest-revision \
    --timeout=540 \
    --memory=128Mi
```

#### Limitations

Currently allowing or disallowing pipelines is not yet implimented. If you need, adding approval is highly recommended. 
