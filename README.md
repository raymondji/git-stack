<!--- DO NOT EDIT: this file is generated --->

# git stack

A minimal Git CLI subcommand for managing stacked branches/pull requests. Works with Gitlab and Github.

Core commands:
- `git stack list`: list all stacks
- `git stack show`: view your current stack
- `git stack push`: push branches in the current stack and open MRs/PRs

## What is stacking?

https://graphite.dev/guides/stacked-diffs has a good overview on what it is and why you might want to do it.

## Where does native Git fall short with stacking?

Stacking branches natively with Git is completely doable, but cumbersome.
- While modern Git has made updating stacked branches much easier with [`--update-refs`](https://andrewlock.net/working-with-stacked-branches-in-git-is-easier-with-update-refs/), other tasks like keeping track of your stacks or pushing all branches in a stack are left to the user.
- Moreover, stacking also typically involves additional manual steps on Gitlab/Github/etc, such as setting the correct target branch on each pull request.

## How does `git stack` compare to `<other stacking tool>`?

There are two main areas where `git stack` differs from most existing tools:
- `git stack` is designed to feel like a minimal addition to the Git CLI. It works with existing Git concepts and functionality (like `--update-refs`), and aims to unintrusively fill in the gaps. Unlike most stacking tools, it's also stateless, so there's no state to keep in sync between Git and `git stack`. Instead, it works by automatically inferring stacks from the structure of your commits.
- `git stack` integrates with Gitlab (and Github). I was surprised to find most of the [popular](https://graphite.dev/) [stacking](https://github.com/aviator-co/av) [tools](https://github.com/gitbutlerapp/gitbutler) only support Github. Besides `git stack`, some other projects I've found with Gitlab support include [git-town](https://github.com/git-town/git-town), [git-spice](https://github.com/abhinav/git-spice) and the new [`glab stack`](https://docs.gitlab.com/ee/user/project/merge_requests/stacked_diffs.html) CLI command. They all work pretty differently and have different feature sets.

## Limitations

- `git stack` requires maintaining linear commit histories in your feature branches to be able to infer stacks. Thus it's effectively tied to using `git rebase`, which seemed reasonable given that `git rebase --update-refs` is the native way of updating stacked branches in Git. However, this means `git stack` is not compatible with `git merge` workflows (at least within feature branches, merging into `main` is no problem).

## Installation

Go version >= 1.22 is required. To install Go on macOS:
```
brew install go 
```

To install `git stack`:
```
go install github.com/raymondji/git-stack-cli/cmd/git-stack@0.26.0
```

## Getting started

The `git stack` binary is named `git-stack`. Git offers a handy trick allowing binaries named `git-<foo>` to be invoked as git subcommands, so `git stack` can be invoked as `git stack`.

`git stack` needs a Gitlab/Github personal access token in order to manage MRs/PRs for you. To set this up:
```
cd ~/your/git/repo
git stack init
```

To learn how to use `git stack`, you can access an interactive tutorial built-in to the CLI:
```
git stack learn
```

## Sample usage

This sample output is taken from `git stack learn --chapter=1 --mode=exec`.

```

```

## How does it work?

When working with Git we often think in terms of branches as the unit of work, and Gitlab/Github both tie pull requests to branches. Thus, as much as possible, `git stack` tries to frame stacks in terms of "stacks of branches".

However, branches in Git don't inherently make sense as belonging to a "stack", i.e. where one branch is stacked on top of another branch. Branches in Git are just pointers to commits, so:
- Multiple branches can point to the same commit
- Branches don't inherently have a notion of parent branches or children branches

Under the hood, `git stack` therefore thinks about stacks as "stacks of commits", not "stacks of branches". Commits serve this purpose much better than branches because:
- Each commit is a unique entity
- Commits do inherently have a notion of parent commits and children commits

Merge commits still pose a problem. It's easy to reason about a linear series of commits as a stack, but merge commits have multiple parents. So, `git stack` takes the simple option of being incompatible with merge commits. If it encounters a merge commit, it will print an error message and otherwise ignore the commit.

`git stack` doesn't persist any additional state to keep track of your stacks - it relies purely on parsing the structure of your commits to infer which commits form a stack (and in turn which branches form a stack). If `git stack` encounters a structure it can't parse (e.g. merge commits), it tries to print a helpful error and otherwise ignores the incompatible commit(s).

## Attribution

Some code is adapted from sections of https://github.com/aviator-co/av (MIT license). A copy of av's license is included at `attribution/aviator-co/av/LICENSE`.
- `exec.go` is adapted from [aviator-co/av/internal/git/git.go](https://github.com/aviator-co/av/blob/fbcb5bfc0f19c8a7924e309cb1e86678a9761daa/internal/git/git.go#L178)
