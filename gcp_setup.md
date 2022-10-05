## Environment And Resource Hiearchy

Following GCPs recommended folder hierarchy we have 3 folders for 3 environments

* prod
  * I shortened the name from production
* staging
  *  GCP wanted to call this `non-production` but staging seemed better because it was more distinct
* dev
  *  I shortened the name from development

## GCP Projects 

GCP projects should generally follow the naming convention

`$name-$env-starling`

* The starling suffix is to make them unique


## VPC Network Architecture

Following GCP's recommendation we have a **Multiple Host Projects with Multiple Shared VPCs**
that corresponds to the 3 environments (prod, staging, dev) that we set above.

With this architecture I think we have a VPC that allows projects in each environment to talk to each other
but the environments are isolated from each other. There is a host project associated with each environment.

[Diagram](./images/shared_vpc_architecture_diagram.svg)

References

* [Designing your network architecture](https://cloud.google.com/architecture/framework/system-design/networking?_ga=2.228904924.-1724403252.1664976400)
* [Security Foundations](https://cloud.google.com/architecture/security-foundations?_ga=2.228904924.-1724403252.1664976400)


## GCP Policy Groups

As part of setting up GCP I followed the recommendations and created the following (role?) groups

* gcp-billing-admins
* gcp-developers
* gcp-devops
* gcp-logging-admins
* gcp-logging-viewers
* gcp-monitoring-admins
* gcp-network-admins
* gcp-organization-admins
* gcp-security-admins

These groups should have been granted the relevant permissions.

Permissions should be granted to these groups and we should add users to these groups. We should avoid directly
granting permissions to individuals as that won't scale very well.

# Notes from Setup 2022-10-05

When setting up the organization `starlingai.com` I followed GCPs recommendations




