# GitHub BigQuery Data

* [Example Query To Find Issues With Google Docs In Them](https://gist.github.com/jlewi/759527e1bcfa6713a8fc07efbe601103)

Select the distinct types of events that are in the GitHub Archive data.

```bash {"id":"01J167SVNQ909ADRPRYXKAE35X","interactive":"false","mimeType":"application/json"}
cat <<EOF >/tmp/query.sql
SELECT DISTINCT type
FROM \`githubarchive.month.202405\`
order by type desc
EOF
export QUERY=$(cat /tmp/query.sql)
bq query --format=json --use_legacy_sql=false "$QUERY"
```

Count the number of closed PRs per day for the last 28 days

```bash {"id":"01J16A41P5Z0XZJYVDW9242DVD","interactive":"false","mimeType":""}
cat <<EOF > /tmp/query.sql
SELECT
  TIMESTAMP_TRUNC(created_at, DAY) as day,
  COUNT(*) as count
FROM
  \`githubarchive.month.202406\`
where
  type = "PullRequestEvent"
  AND 
  repo.name = "jlewi/foyle"
  AND
  json_value(payload, "$.action") = "closed"
  AND
  DATE(created_at) BETWEEN DATE_SUB(CURRENT_DATE(), INTERVAL 28 DAY) AND CURRENT_DATE()
GROUP BY
  TIMESTAMP_TRUNC(created_at, DAY)
ORDER BY
  day desc
EOF
export QUERY=$(cat /tmp/query.sql)
MAXBILLED=1000000000000
bq query --format=json --maximum_bytes_billed=${MAXBILLED} --use_legacy_sql=false "$QUERY"
```

```sh {"id":"01J167SCXYR2HW15H46G9WTCAE"}
# References

* [Issues Comment Event](https://docs.github.com/en/webhooks/webhook-events-and-payloads#issue_comment)

* IssueCommentEvent - I think this misses issue opened events
* IssuesEvent - I think this is issues when an issue is opened
```