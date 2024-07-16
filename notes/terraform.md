# Terraform

## Refresh vs. import

* Terraform will only resources defined in the state file
* Resources get added to the state file when they are created by Terraform
* Refresh updates resources in the terraform file with the current state of the infrastructure
   * If a resource owned by TF has been deleted it will be marked as being deleted

* However if a resource exists but was created outside of TF but is defined in your .tf files (not .tfstate) then
   a refresh will not add it to the tfstate and Terraform will not manage it
   * You need to import it to bring it under Terraform management

## Resources

* [Blog describing pain points in Terraform](https://www.qovery.com/blog/terraform-not-the-golden-hammer/)
   * "Dependency is too strong problem" - I think this is basically a version of the you don't really want a DAG on Day 2 problem
      * If you structure your Terraform code as a DAG then it won't run successive steps if there is an error on earlier steps
      * This might be what you want on Day 0 where you don't want to create the resources downstream if the resources upstream fail
         * e.g. in the blog they don't want to create the EKS cluster if DNS setup fails

      * However on Day 2 Cloud Flare (DNS) might be down so you might not be able to update edit it. That could then block
         updates to EKS even if you don't need to make any changes to EKS

   * "Automatic Import" - If you want Terraform to take control of a resource it doesn't know about you have to manually import it

* [What does terraform refresh do](https://stackoverflow.com/questions/42628660/what-does-terraform-refresh-really-do)

* [Article explaining limits of Terraform Plan and Apply](https://medium.com/@bgrant0607/infrastructure-as-code-reminds-me-of-make-run-all-15eb6628f306)

* [Terraform vs. Crossplane](https://blog.crossplane.io/crossplane-vs-terraform/)