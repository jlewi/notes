# BigQuery

See [controlling costs](https://cloud.google.com/bigquery/docs/best-practices-costs)

* You can set maxBillingBytes on individual jobs to set limit of bytes to be billed
  if query would exceed this it won't run and you won't get charged

* I think you can also set a per project limit 
  I think in the Google Cloud Console you can set the limit [Query usage per day](https://cloud.google.com/bigquery/quotas)
  * Not sure that's the right value (its per project not query)
  * Default is unlimited