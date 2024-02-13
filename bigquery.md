# BigQuery

See [controlling costs](https://cloud.google.com/bigquery/docs/best-practices-costs)

* You can set maxBillingBytes on individual jobs to set limit of bytes to be billed
  if query would exceed this it won't run and you won't get charged

* I think you can also set a per project limit 
  I think in the Google Cloud Console you can set the limit [Query usage per day](https://cloud.google.com/bigquery/quotas)
  * Not sure that's the right value (its per project not query)
  * Default is unlimited


## BigQuery and GoLang

* You can use InferSchema to infer the schema from a struct

InferSchema requires we use nullable types such as bigquery.NullString to indicate that a column is optional; this is required
if columns are added later. Nullable types can be a pain to work with in the code. On the other hand,
if we try to deserialize a null record into field which doesn't have NullableType then we get an error.
This could potentially be addressed by 1) using BigQuery to set a non-null default value and/or 2) automatically
deserializing null values to the corresponding golang default value e.g. empty string, zero, etc...

Examples

[CreateDocumentsTable](https://github.com/starlingai/flock/blob/4feae65ebf62844d7b4f7ef9ce35cbe650da2491/go/pkg/openai/tables.go#L20)