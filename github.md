# GitHub Policies and Settings

Require linear history 

* see [explanation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)

* Merge commits must use either squash merge or rebase merge

I think we want to use [squash merge commits](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)

* Configure the default message to be PR title + description.

* If you do rebase merge you get multiple commits which isn't what we want.

* Configure branches to be automatically deleted after PR merge
  * Intermediary commits on a PR still appear to be linkable even if the branch is deleted.


## Creating new repositories

1. Create the repository from the template repository [starlingai/template](https://github.com/starlingai/template)