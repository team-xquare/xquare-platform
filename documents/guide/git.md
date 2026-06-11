# GIT.md
This guide defines Git conventions for this repository.
It follows the official Git documentation and Pro Git recommendations, with
project rules that keep changes understandable, reviewable, and recoverable.

## 1. Core Principles
1. Treat a commit as a reviewed project snapshot, not as a backup of arbitrary
   local work.
2. Keep each branch focused on one explicit task.
3. Keep each commit focused on one coherent change.
4. Preserve existing user work and shared history.
5. Inspect before staging, committing, rebasing, merging, or deleting.
6. Prefer commands whose effects are narrow and visible.
7. Do not use destructive Git commands without explicit approval.

## 2. Repository State
Before modifying files:

```sh
git status --short --branch
```

1. Identify staged, unstaged, untracked, conflicted, and ignored files.
2. Assume changes you did not create belong to another user or process.
3. Do not discard, overwrite, stage, or include unrelated changes.
4. Use `git diff -- path/to/file` before editing a file that already has
   changes.
5. Use `git diff --cached` before every commit.
6. Use `git check-ignore -v path/to/file` when an expected file is missing from
   status.

## 3. Worktrees
1. Follow `documents/guide/common.md`: file-modifying tasks use a separate Git
   worktree unless the user explicitly approves an exception.
2. Create a worktree from the intended base branch:

```sh
git fetch origin
git worktree add -b docs/update-guides ../update-guides origin/main
```

3. Confirm the new worktree's branch and status before editing.
4. Do not attach the same branch to multiple worktrees.
5. Do not remove a worktree until its changes are committed, transferred, or
   explicitly discarded.
6. Ask before removing a task worktree.
7. Use `git worktree list` before cleanup.

## 4. Branches
1. Branch from the current intended integration branch.
2. Fetch before creating a branch when a remote exists.
3. Use lowercase names with slash-separated category and kebab-case detail:

```text
feature/add-token-refresh
fix/reject-expired-token
refactor/split-auth-handler
docs/expand-go-guide
build/update-bazel-rules
```

4. Allowed categories are `feature`, `fix`, `refactor`, `docs`, `test`,
   `build`, `ci`, and `chore`.
5. Include an issue identifier only when the project uses one:

```text
fix/XQ-123-reject-expired-token
```

6. Do not use personal names, vague names, or generic sequence numbers.
7. Avoid `new`, `temp`, `test`, `final`, `v2`, and `changes` as the meaningful
   part of a branch name.
8. Delete merged branches only after verifying that no unique work remains.

## 5. Commit Boundaries
A commit must:
1. represent one logical change
2. include required tests and documentation for that change
3. avoid unrelated formatting or cleanup
4. leave the repository in a valid state for the affected scope
5. be understandable without depending on a later commit

Split commits when:
1. changes have different reasons
2. one part can be reviewed independently
3. generated output can be separated from its generator change without making
   the history invalid
4. a mechanical rename obscures a behavioral change

Keep changes together when splitting would create an uncompilable or misleading
intermediate state.

## 6. Staging
1. Stage only reviewed changes.
2. Prefer path-specific or patch staging:

```sh
git add path/to/file
git add -p
```

3. Do not use `git add .` or `git add -A` when unrelated files are present.
4. Review both summaries and patches:

```sh
git diff --stat
git diff
git diff --cached --stat
git diff --cached
```

5. Check for whitespace and conflict markers:

```sh
git diff --check
git diff --cached --check
```

6. Do not stage secrets, local environment files, editor state, build outputs,
   or credentials.
7. Do not assume `.gitignore` protects a secret after it has been tracked.

## 7. Commit Messages
Git itself does not mandate a universal commit-message taxonomy. This repository
uses the Git project's imperative-message guidance.

### 7.1 Subject
1. Write the subject in the imperative mood.
2. Describe what applying the commit does.
3. Capitalize the first word after an optional area prefix.
4. Do not end the subject with punctuation.
5. Keep the subject near 50 characters when practical.
6. Use a specific area prefix when it improves scanning:

```text
auth: Reject expired access tokens
bazel: Make token tests hermetic
docs: Define TypeScript nullability rules
```

7. Omit the prefix when the change is repository-wide or a prefix would add no
   information.

Good:

```text
Validate refresh token expiration
```

Avoid:

```text
Fixed bug
Update stuff
WIP
Changes
```

### 7.2 Body
1. Separate the body from the subject with one blank line.
2. Explain why the change is needed and how behavior differs.
3. Include constraints, tradeoffs, migration concerns, and non-obvious
   consequences.
4. Do not repeat the diff line by line.
5. Wrap prose near 72 characters where practical.
6. Reference issues or commits in a footer or explanatory paragraph.
7. When referencing a commit, include an abbreviated hash and subject when that
   context matters.

```text
auth: Reject expired access tokens

The validator previously checked the signature but did not compare the
expiration claim with the request clock. Reject expired tokens before
constructing the authenticated principal.

Refs: XQ-123
```

### 7.3 Breaking Changes
1. State a breaking API, schema, build, or operational change explicitly.
2. Describe the migration path.
3. Do not hide a breaking change behind a vague refactor subject.

## 8. Authoring Commits
1. Run the affected validation before committing.
2. Use an editor for commits that need a body.
3. Avoid `git commit -am`; it stages every modified tracked file and can include
   unrelated work.
4. Verify the final commit:

```sh
git show --stat --oneline HEAD
git show --check HEAD
```

5. Amend only an unshared commit or with explicit coordination.
6. Remember that `git commit --amend` replaces the commit with a new object.

## 9. History Rewriting
1. Rewrite local, unpublished history when it improves reviewability.
2. Do not rebase, amend, squash, or reset commits already used by others without
   explicit coordination.
3. Treat pushed shared history as immutable by default.
4. Use interactive rebase only after checking the exact range:

```sh
git log --oneline --decorate --graph
git rebase -i <base>
```

5. After rewriting a published personal branch, use
   `git push --force-with-lease`, never a blind `--force`.
6. Verify the old and new patch series with `git range-diff` when a reviewed
   branch changes substantially.

## 10. Updating a Branch
1. Fetch remote state explicitly:

```sh
git fetch origin
```

2. Inspect divergence before integrating:

```sh
git log --oneline --left-right --graph HEAD...origin/main
```

3. Rebase an unpublished task branch when a linear update improves review.
4. Merge when preserving branch history is required by the repository workflow.
5. Do not use an unconfigured `git pull` when the merge or rebase behavior is
   unclear.
6. Prefer an explicit command:

```sh
git pull --ff-only
git rebase origin/main
git merge origin/main
```

7. Never resolve conflicts by automatically choosing all of one side.
8. Understand each conflict, preserve intended behavior from both sides, and
   run validation afterward.

## 11. Merge and Pull Request History
1. Follow the repository's configured integration method.
2. Do not switch between merge, squash, and rebase policies within one pull
   request without a reason.
3. A pull request must have one coherent purpose.
4. The pull request description must state:
   - problem and intended outcome
   - important implementation decisions
   - validation performed
   - migration, rollout, or risk notes
5. Keep drive-by cleanup out of the pull request.
6. Update the description when the final scope differs from the initial plan.

## 12. Reverting and Restoring
1. Use `git revert` to undo a shared commit because it records a new inverse
   commit.
2. Use `git restore` for explicit working-tree or index restoration only after
   confirming the affected path.
3. Use `git reset` only when its effect on the branch, index, and working tree
   is fully understood and explicitly approved when destructive.
4. Never use `git reset --hard` or `git clean` without explicit user approval.
5. Before deleting untracked files, list exactly what would be affected.
6. Prefer creating a safety branch or tag before a risky local history
   operation.

## 13. File Moves and Renames
1. Keep pure moves separate from behavioral edits when that materially improves
   review.
2. Use `git mv` when convenient, but remember Git detects renames from content
   similarity rather than storing a rename operation.
3. Preserve file history by avoiding unnecessary rewrite during a move.
4. Validate imports, build labels, documentation links, and ownership metadata
   after moving files.

## 14. Generated Files and Dependencies
1. Commit generated files only when repository policy requires them.
2. Generate them from the same source change.
3. State the generator command in the commit or pull request when it is not
   obvious.
4. Do not hand-edit generated output.
5. Keep dependency manifest and lockfile changes together.
6. Explain unexpected transitive dependency changes.

## 15. Binary and Large Files
1. Do not commit build artifacts or replace source assets with generated
   binaries.
2. Verify whether Git LFS or another repository mechanism is required before
   adding a large binary.
3. Avoid repeated binary churn because Git cannot review or delta it as clearly
   as text.
4. Store reproducible artifacts in the designated artifact system when one
   exists.

## 16. Tags
1. Use annotated tags for releases:

```sh
git tag -a v1.2.3 -m "Release v1.2.3"
```

2. Follow the repository's release versioning policy.
3. Do not move or replace a published tag without explicit coordination.
4. Verify the commit a tag points to before pushing it.

## 17. Security
1. Never commit credentials, tokens, private keys, personal data, or production
   configuration.
2. If a secret is committed, removing it in a later commit is insufficient.
3. Stop, report the exposure without repeating the value, rotate the secret,
   and follow the approved history-cleaning procedure.
4. Do not paste sensitive command output into commit messages or pull requests.
5. Review staged `.env`, certificate, key, archive, and database files
   carefully.

## 18. Recommended Inspection Commands

```sh
git status --short --branch
git diff
git diff --cached
git diff --check
git log --oneline --decorate --graph -20
git show --stat --oneline <commit>
git branch --verbose --verbose
git worktree list
```

Use path limits when reviewing a focused change:

```sh
git diff -- documents/guide/go.md
git log --oneline -- path/to/file
```

## 19. Pre-Commit Checklist
- Is the branch based on the intended branch?
- Are unrelated user changes untouched?
- Does the staged diff contain only the intended task?
- Is the commit one coherent change?
- Are tests and documentation included where required?
- Did `git diff --cached --check` pass?
- Does the subject use a specific imperative phrase?
- Does the body explain why when the reason is not obvious?
- Are secrets and generated artifacts excluded?

## 20. Completion Checklist
- Is the working tree state understood?
- Were all required validation commands run?
- Does the commit history tell a clear sequence?
- Has shared history remained intact?
- Are worktree cleanup and branch deletion deferred until approved?

## References
- [Git Reference](https://git-scm.com/docs)
- [Pro Git: Contributing to a Project](https://git-scm.com/book/en/v2/Distributed-Git-Contributing-to-a-Project)
- [Pro Git: Rebasing](https://git-scm.com/book/en/v2/Git-Branching-Rebasing)
- [Pro Git: Rewriting History](https://git-scm.com/book/en/v2/Git-Tools-Rewriting-History)
- [Git SubmittingPatches](https://git-scm.com/docs/SubmittingPatches)
