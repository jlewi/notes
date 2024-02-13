# Git


## Rebasing

Suppose you are doing a rebase and dropping a bunch of changes that have already been merged.
If this produces a lot of conflicts. I think one way to deal with this is to just overwrite
the files with the versions on the commit on your branch  e.g.

git checkout ${COMMIT} <path>