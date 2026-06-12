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
git worktree add \
  -b docs/123-update-guides \
  ../update-guides \
  origin/main
```

3. Confirm the new worktree's branch and status before editing.
4. Do not attach the same branch to multiple worktrees.
5. Do not remove a worktree until its changes are committed, transferred, or
   explicitly discarded.
6. Ask before removing a task worktree.
7. Use `git worktree list` before cleanup.

## 4. Issues and Branches
Every task branch and pull request must be traceable to an Issue.

### 4.1 Issue Before Branch
1. Search open and closed Issues before creating a new one.
2. Reuse an existing Issue only when its problem, scope, and acceptance criteria
   cover the intended task.
3. Create an Issue before creating the task branch when no suitable Issue
   exists.
4. Do not use a pull request as a substitute for an Issue.
5. Keep one primary Issue for one task branch.
6. Split the work when one branch would implement unrelated Issues.
7. Related secondary Issues may be referenced in the pull request, but the
   branch name contains only the primary Issue number.
8. If the task scope changes beyond the primary Issue, update the Issue or
   create another Issue before continuing.
9. Confirm the Issue number, title, state, and acceptance criteria before branch
   creation.

An Issue must state:
1. problem or purpose
2. scope
3. acceptance criteria
4. relevant constraints, risks, or dependencies

### 4.2 Branch Format
Use:

```text
<type>/<issue-number>-<kebab-case-summary>
```

Examples:

```text
feat/123-add-token-refresh
fix/456-reject-expired-token
perf/789-reduce-session-allocations
refactor/234-split-auth-handler
docs/345-expand-go-guide
build/567-update-bazel-rules
work/678-require-pull-requests
```

Rules:
1. Use the decimal GitHub Issue number without `#`.
2. Use lowercase ASCII characters.
3. Separate the type from the detail with one `/`.
4. Separate detail words with one `-`.
5. Use a short action-oriented summary that describes the intended outcome.
6. Keep the full branch name at or below 80 characters.
7. Use three to eight meaningful summary words when practical.
8. Do not include the author's name, date, environment, or implementation
   sequence.
9. Do not use `new`, `temp`, `final`, `changes`, `stuff`, `misc`, `test`,
   `issue`, or `branch` as the meaningful summary.
10. Do not use multiple Issue numbers in one branch name.

### 4.3 Branch Types
Use exactly one of:

| Type | Use |
|---|---|
| `feat` | user-visible or externally consumable capability |
| `fix` | correction of incorrect behavior |
| `perf` | measured performance improvement |
| `refactor` | internal restructuring without intended behavior change |
| `docs` | documentation-only change |
| `test` | test-only addition or correction |
| `build` | build system, dependency, or packaging change |
| `ci` | continuous-integration or delivery change |
| `chore` | repository maintenance not covered by another type |
| `work` | agent workflow or repository work-procedure change |
| `revert` | explicit reversal of an earlier change |

Choose the type for the task's primary outcome.
Do not choose `chore` when a more specific type applies.

### 4.4 Branch Creation
1. Fetch the remote before creating a branch.
2. Branch from the approved integration branch.
3. Confirm that the base branch is current.
4. Confirm that no local branch already represents the Issue.
5. Create the worktree and branch together when the task modifies files.

```sh
git fetch origin
git worktree add \
  -b docs/345-expand-go-guide \
  ../expand-go-guide \
  origin/main
```

6. Confirm the branch and worktree state before editing.
7. Do not rename a pushed branch with an open pull request without explicit
   coordination.
8. Delete merged branches only after verifying that no unique work remains.

### 4.5 Automated Branches
Repository automation may use a tool-controlled branch name only when the tool
cannot follow the repository format.
The generated pull request must still reference a valid Issue.
Do not use this exception for manually created branches.

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
uses a typed format together with the Git project's imperative-message
guidance.

Use:

```text
<type>(<scope>): <imperative summary>

<body>

Refs #<primary-issue-number>
```

The scope is optional:

```text
docs(git): Define issue-linked conventions
fix(auth): Reject expired access tokens
work: Require pull request creation
```

### 7.1 Subject
1. Use a lowercase commit type followed by `:`.
2. Add a scope in parentheses when it identifies a stable package, module,
   service, tool, or guide.
3. Use lowercase kebab-case for a multiword scope.
4. Do not use an Issue number as the scope.
5. Write the summary in the imperative mood.
6. Describe what applying the commit does.
7. Capitalize the first word after the colon.
8. Do not end the subject with punctuation.
9. Keep the complete subject at or below 72 characters.
10. Keep it near 50 characters when practical.

```text
fix(auth): Reject expired access tokens
build(bazel): Make token tests hermetic
docs(typescript): Define nullability rules
```

Good:

```text
feat(session): Add refresh token rotation
docs: Define project research protocol
```

Avoid:

```text
fixed bug
update stuff
WIP changes
docs(git): Updates
#123 add feature
```

### 7.2 Commit Types
Commit types use the same meanings as branch types:
`feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `build`, `ci`, `chore`,
`work`, and `revert`.

Rules:
1. The commit type must describe the commit, not the entire branch.
2. A branch may contain different commit types when each commit has a distinct
   logical purpose.
3. Use `test` for test-only changes; keep tests required by a behavior change in
   the behavior commit.
4. Use `docs` for documentation-only changes.
5. Use `refactor` only when observable behavior is intentionally unchanged.
6. Use `perf` only when the change targets measured performance.
7. Use `chore` only when no more specific type applies.

### 7.3 Scope
1. Use a scope when it makes the affected area immediately clear.
2. Use a stable domain or repository name such as `auth`, `token`, `bazel`,
   `git`, `research`, or `typescript`.
3. Omit the scope for a genuinely repository-wide change.
4. Do not use file extensions, ticket numbers, temporary labels, or author names
   as scopes.
5. Keep the same spelling for the same area across commits.

### 7.4 Body
1. Separate the body from the subject with one blank line.
2. Explain why the change is needed and how behavior differs.
3. Include constraints, tradeoffs, migration concerns, and non-obvious
   consequences.
4. Do not repeat the diff line by line.
5. Wrap prose near 72 characters where practical.
6. Include a body when:
   - the reason is not obvious from the subject
   - behavior changes
   - a compatibility or migration concern exists
   - a non-obvious implementation decision needs explanation
   - validation is constrained
7. When referencing another commit, include an abbreviated hash and subject when
   context matters.

```text
fix(auth): Reject expired access tokens

The validator previously checked the signature but did not compare the
expiration claim with the request clock. Reject expired tokens before
constructing the authenticated principal.

Refs #456
```

### 7.5 Issue Footer
1. Every manually authored task commit must reference the primary Issue.
2. Use this footer:

```text
Refs #<issue-number>
```

3. Put the footer after the body, separated by one blank line.
4. Use the Issue number from the branch name.
5. Do not put `#<issue-number>` in the subject.
6. Use `Refs`, not a closing keyword, in ordinary commits.
7. Reserve `Closes`, `Fixes`, or `Resolves` for the pull request body so Issue
   completion follows pull request merge.
8. A platform-generated merge commit is exempt from the authored-message format.
9. An automated commit may use a tool-controlled format only when repository
   automation owns it.

### 7.6 Breaking Changes
1. State a breaking API, schema, build, or operational change explicitly.
2. Describe the migration path.
3. Do not hide a breaking change behind a vague refactor subject.
4. Add a footer after the Issue footer:

```text
Refs #123
BREAKING CHANGE: Describe the incompatible behavior and migration.
```

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
7. Verify that the commit type, optional scope, subject, and Issue footer follow
   section 7.
8. Verify that the Issue footer matches the branch's primary Issue.
9. Do not commit with `--no-verify` unless explicitly approved after reporting
   the failing hook and reason.

Example:

```sh
git commit \
  -m "docs(git): Define issue-linked conventions" \
  -m "Document branch, commit, and pull request traceability." \
  -m "Refs #7"
```

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

## 11. Pull Requests and Merge History
1. Follow the repository's configured integration method.
2. Do not switch between merge, squash, and rebase policies within one pull
   request without a reason.
3. A pull request must have one coherent purpose.
4. Every pull request must reference at least one valid Issue in its body.
5. The pull request must use one primary Issue.
6. Use a closing keyword when merging the pull request fully satisfies the
   primary Issue:

```text
Closes #123
```

7. Use a non-closing reference when the pull request contributes to but does
   not complete the Issue:

```text
Refs #123
```

8. Use `Closes`, `Fixes`, or `Resolves` only when the Issue's acceptance
   criteria are fully satisfied.
9. Do not use a closing keyword for a partial implementation.
10. Do not create an Issue only after the pull request is complete to satisfy
    this rule retroactively.
11. The pull request description must state:
   - problem and intended outcome
   - important implementation decisions
   - validation performed
   - migration, rollout, or risk notes
   - Issue reference
12. Use this body structure by default:

```md
## Summary
- ...

## Validation
- `command`: passed

## Risks
- None

## Issue
Closes #123
```

13. Use a title that follows the commit subject format:

```text
<type>(<optional-scope>): <imperative summary>
```

14. Keep drive-by cleanup out of the pull request.
15. Update the title, description, and Issue references when the final scope
    differs from the initial plan.
16. Verify the Issue exists and the reference is rendered correctly after pull
    request creation.
17. Do not create duplicate pull requests for the same branch.
18. If several Issues are required, list each reference explicitly and explain
    why one pull request is still a coherent review unit.

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
- Does the branch include the primary Issue number?
- Does the branch follow `<type>/<issue-number>-<summary>`?
- Is the branch based on the intended branch?
- Are unrelated user changes untouched?
- Does the staged diff contain only the intended task?
- Is the commit one coherent change?
- Are tests and documentation included where required?
- Did `git diff --cached --check` pass?
- Does the subject follow `<type>(<scope>): <imperative summary>`?
- Does the body explain why when the reason is not obvious?
- Does the footer contain `Refs #<primary-issue>`?
- Are secrets and generated artifacts excluded?

## 20. Completion Checklist
- Is the working tree state understood?
- Were all required validation commands run?
- Does the commit history tell a clear sequence?
- Has shared history remained intact?
- Is the task branch pushed to the expected remote?
- Does the pull request reference a valid primary Issue?
- Does the pull request use `Closes` only when it fully satisfies the Issue?
- Are worktree cleanup and branch deletion deferred until approved?

## References
- [Git Reference](https://git-scm.com/docs)
- [Pro Git: Contributing to a Project](https://git-scm.com/book/en/v2/Distributed-Git-Contributing-to-a-Project)
- [Pro Git: Rebasing](https://git-scm.com/book/en/v2/Git-Branching-Rebasing)
- [Pro Git: Rewriting History](https://git-scm.com/book/en/v2/Git-Tools-Rewriting-History)
- [Git SubmittingPatches](https://git-scm.com/docs/SubmittingPatches)
- [GitHub: Linking a Pull Request to an Issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue)
