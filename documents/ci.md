# Continuous Integration

This repository uses GitHub Actions to validate changes before they enter
`develop` or `main`. CI does not build deployment images, publish artifacts, or
deploy services.

## Branch Flow

```text
feature branch
  -> pull request to develop
  -> merge queue
  -> develop
  -> promotion pull request to main
  -> merge queue
  -> main
  -> external deployment system
```

A merge to `main` triggers deployment outside this repository's CI workflows.
The complete `main` qualification must therefore pass before the merge.

## Validation Policy

| Situation | Bazel validation |
|---|---|
| Pull request to `develop` | Affected targets |
| Merge group targeting `develop` | Affected targets from the temporary merge |
| Pull request to `main` | Full build and test |
| Merge group targeting `main` | Full build and test |
| Repository-wide build or CI change | Full build and test |
| File deletion or rename | Full build and test |
| Documentation-only change | Bazel validation skipped |
| Impact-analysis failure | Full build and test |
| Manual workflow execution | Full build and test |
| Push to `main` | No CI workflow; deployment is external |

Repository-wide changes include:

- `.bazelrc`
- `.bazelversion`
- `MODULE.bazel`
- `MODULE.bazel.lock`
- `mise.toml`
- `build/**`
- `.github/**`
- `*.bzl`
- `*.MODULE.bazel`

The affected-target path uses Target Determinator with a pinned version and
SHA-256 checksum. If affected validation fails for any reason, the workflow
runs the full Bazel build and test. A failing affected test is therefore not
hidden by target selection.

## Required Checks

Configure the branch rules for `develop` and `main` to require:

- `CI / Gate`
- the CodeQL language checks after the repository becomes public

`CI / Gate` aggregates:

- workflow and Bazel formatting policy
- Bazel validation
- Trivy vulnerability, configuration, and secret scanning
- dependency review when GitHub provides it

The gate accepts an explicitly skipped check only when that check is not
applicable, such as dependency review before the repository is public.

## Merge Queue

Both protected branches should require the merge queue. The `merge_group`
event validates GitHub's temporary merge result against the latest target
branch state. Do not replace merge-group validation with the earlier pull
request result.

The repository does not require CODEOWNERS. Branch rules should instead require
at least one pull-request approval, dismiss stale approvals after new commits,
block force pushes and branch deletion, and prevent administrator bypass.

## Security

The CI workflow grants read-only repository access by default. External actions
are pinned to full commit SHAs, and the workflow rejects newly introduced
unpinned actions.

Trivy runs for pull requests and merge groups. Dependency Review and CodeQL are
conditioned on public repository visibility because their availability differs
for private repositories without GitHub Advanced Security.

CodeQL covers:

- Go
- Java and Kotlin
- JavaScript and TypeScript

CodeQL runs for pull requests, merge groups, manual requests, and once per week.
There is no scheduled full build or test.

Fork pull requests receive no repository secrets and no write permissions. The
workflows do not use `pull_request_target`.

## Scaling

Bazel owns target parallelism and dependency scheduling. GitHub jobs should not
be created per service because that model produces excessive runner startup and
workflow-management overhead as the monorepo grows.

When execution time requires shared infrastructure, add a Bazel Remote
Execution API provider through repository-wide Bazel configuration:

1. remote cache for trusted validation results
2. Build Event Protocol collection for timing and cache metrics
3. remote execution with autoscaled workers

Provider selection, credentials, and trust boundaries must be decided before
adding those settings. Untrusted fork pull requests must not receive cache
write credentials or access to private workers.

Track at least:

- pull-request feedback time at `p50` and `p95`
- runner queue time
- affected-target fallback frequency
- remote cache hit rate when enabled
- flaky-test rate
- CI infrastructure failure rate

## Local Validation

Run the same repository checks locally with pinned tools:

```sh
mise exec -- buildifier -mode=check -lint=warn -r .
mise exec -- bazelisk build //...
mise exec -- bazelisk test --test_output=errors //...
```

The tool versions remain owned by `mise.toml` and `.bazelversion`. CI-specific
tools are not added to `mise.toml`.
