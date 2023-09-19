# Google Drive API Sandbox

## Drive API Docs
* How you can use client libraries and Application Default Credentials to access Google APIs: https://cloud.google.com/docs/authentication/client-libraries
* Using an API key with client libraries: https://cloud.google.com/docs/authentication/api-keys#using-with-client-libs
* [Set up ADC](https://cloud.google.com/docs/authentication/provide-credentials-adc) for the environment where your application is running
* For a local development environment, you can [set up ADC with your user credentials](https://cloud.google.com/docs/authentication/provide-credentials-adc#local-dev) by using the gcloud CLI
* For production environments, you [set up ADC by attaching a service account](https://cloud.google.com/docs/authentication/provide-credentials-adc#attached-sa)
* [Using OAuth 2.0 to Access Google APIs] (https://developers.google.com/identity/protocols/oauth2)
* [Authentication Use Cases](https://cloud.google.com/docs/authentication/use-cases)
* [Creating authenticated transports to Google's HTTP and gRPC APIs](https://pkg.go.dev/google.golang.org/api/transport)

* [Drive create credentials](https://developers.google.com/workspace/guides/create-credentials)
** Key > Set up domain-wide delegation for a service account
* Drive API Go Quickstart: https://developers.google.com/drive/api/quickstart/go
* Drive Activity Go Quickstart: https://developers.google.com/drive/activity/v2/quickstart/go
* Push Notifications: https://developers.google.com/drive/api/guides/push
* `Watch` Reference: https://developers.google.com/drive/api/v3/reference/files/watch

## Rough Notes
- Looks like I can use the Drive Activity API to list/Get the recent drive activity and lookup Create acctions with activity title name strong.csv and grab the name/item ID e.g. 1Gd4CvxH3iHl9YtPKgz_oaKIUQaKi5PEB

```text
activities {"actions":[{"detail":{"create":{"upload":{}}}},{"detail":{"edit":{}}},{"detail":{"permissionChange":{"addedPermissions":[{"role":"OWNER","user":{"knownUser":{"isCurrentUser":true,"personName":"people/108969388462034462798"}}}]}}},{"detail":{"move":{"addedParents":[{"driveItem":{"driveFolder":{"type":"STANDARD_FOLDER"},"folder":{"type":"STANDARD_FOLDER"},"name":"items/1MGBCtMFKGuFRxIuOx20L6vOmzDfiAbYS","title":"strong_app_workout_logs"}}]}}}],"actors":[{"user":{"knownUser":{"isCurrentUser":true,"personName":"people/108969388462034462798"}}}],"primaryActionDetail":{"create":{"upload":{}}},"targets":[{"driveItem":{"driveFile":{},"file":{},"mimeType":"text/csv","name":"items/1Gd4CvxH3iHl9YtPKgz_oaKIUQaKi5PEB","owner":{"user":{"knownUser":{"isCurrentUser":true,"personName":"people/108969388462034462798"}}},"title":"strong.csv"}}],"timestamp":"2023-08-05T22:44:13.971Z"}
```

- Then I should be able to use the Drive API to download/fetch the csv file by ID to feed into the my strong app