# GoGit

## Configuration

If you get the repo config

```
r.Config()
```

I think that is only loading `$REPO/.git/config` its not merging it with user config e.g.
`$HOME/.gitconfig`.

Per [go-git/go-git#395](https://github.com/go-git/go-git/issues/395) it looks like there
are some potential workarounds.

A workaround is to use `ConfigScoped`

```
c, err := a.repo.ConfigScoped(gogitconfig.GlobalScope)
```


## Examples and Snippets


### Delete and recreate a branch from a specific commit

```
    log := zapr.NewLogger(zap.L())
	log.Info("Running as jlewi; resetting the git branch of the generated code")
	gitRoot, err := gitutil.LocateRoot(pkgPath)
	if err != nil {
		return err
	}
	gitRepo, err := git.PlainOpenWithOptions(gitRoot, &git.PlainOpenOptions{})
	if err != nil {
		return errors.Wrapf(err, "Error opening git repo %v", gitRoot)
	}

	localBranch := "jlewi/main"
	// Try to get a reference to the local branch. We use this to determine whether the branch already exists
	branchRef, err := gitRepo.Reference(plumbing.NewBranchReferenceName(localBranch), false)

	if err != nil && err.Error() != "reference not found" {
		return err
	}

	if err == nil {
		// Branch already exists
		log.Info("Branch already exists; deleting it.", "branch", localBranch, "local", branchRef.Hash())

		if err := gitRepo.Storer.RemoveReference(branchRef.Name()); err != nil {
			return errors.Wrapf(err, "Failed to delete existing local branch %v", branchRef.Name())
		}
	}

	// Commit to create the branch from
	hash := "9adc4a9d8f4e819d5a9f40e5dd40a3f06af7ed2c"

	checkoutOptions := &git.CheckoutOptions{
		// Branch to create
		Branch: plumbing.NewBranchReferenceName(localBranch),
		Force:  true,
		Create: true,
		Hash:   plumbing.NewHash(hash),
	}

	log.Info("Checking out branch", "name", checkoutOptions.Branch, "hash", hash)

	w, err := gitRepo.Worktree()
	if err != nil {
		return err
	}
	err = w.Checkout(checkoutOptions)

	if err != nil {
		return err
	}
	return nil
```