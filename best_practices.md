# Best Practices

## Service Account Impersonation and GCP Management

We should probably administer GCP using service accounts rather than
directly granting individuals elevated admin privileges.

See: https://cloud.google.com/iam/docs/impersonating-service-accounts

gcloud has a global flag for this.

```
     --impersonate-service-account=SERVICE_ACCOUNT_EMAILS
        For this gcloud invocation, all API requests will be made as the given
        service account or target service account in an impersonation
        delegation chain instead of the currently selected account. You can
        specify either a single service account as the impersonator, or a
        comma-separated list of service accounts to create an impersonation
        delegation chain. The impersonation is done without needing to create,
        download, and activate a key for the service account or accounts.

        In order to make API requests as a service account, your currently
        selected account must have an IAM role that includes the
        iam.serviceAccounts.getAccessToken permission for the service account
        or accounts.

        The roles/iam.serviceAccountTokenCreator role has the
        iam.serviceAccounts.getAccessToken permission. You can also create a
        custom role.

        You can specify a list of service accounts, separated with commas. This
        creates an impersonation delegation chain in which each service account
        delegates its permissions to the next service account in the chain.
        Each service account in the list must have the
        roles/iam.serviceAccountTokenCreator role on the next service account
        in the list. For example, when --impersonate-service-account=
        SERVICE_ACCOUNT_1,SERVICE_ACCOUNT_2, the active account must have the
        roles/iam.serviceAccountTokenCreator role on SERVICE_ACCOUNT_1, which
        must have the roles/iam.serviceAccountTokenCreator role on
        SERVICE_ACCOUNT_2. SERVICE_ACCOUNT_1 is the impersonated service
        account and SERVICE_ACCOUNT_2 is the delegate.

        Overrides the default auth/impersonate_service_account property value
        for this command invocation.
```

The advantage of using service account impersonation is that for normal daily use
the principal wouldn't be running with elevated GCP priveleges which guards
against doing accidental deletions and the like.

See [Examples](https://github.com/salrashid123/gcp_impersonated_credentials) includes code and gcloud examples.