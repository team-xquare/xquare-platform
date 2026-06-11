# WORK.md
This guide defines the procedure for planning, approval, execution, progress updates, validation, and completion reporting.
Follow `AGENTS.md` and `documents/guide/common.md` first. This guide only defines the detailed workflow for performing tasks.
## 1. Scope
This guide applies to any task that may involve:
1. running commands
2. creating, editing, moving, or deleting files
3. writing or refactoring code
4. editing documentation
5. running tests, builds, linters, formatters, or type checks
6. performing Git operations
7. taking any action with side effects
This guide does not require approval for purely informational responses, explanations, summaries, reviews, evaluations, or draft suggestions that do not run commands or modify files.
## 2. Plan Before Execution
Before running commands, modifying files, or performing any action with side effects, provide a plan to the user.
The plan must include:
1. task summary
2. guide files to read
3. files or areas to inspect
4. intended steps
5. possible commands
6. expected files to change
7. validation method
8. planned commit units
9. stop conditions
The plan must be concise and specific enough for the user to decide whether to approve it.
## 3. Approval
1. Do not start execution until the user explicitly approves the plan.
2. Approval must be a clear instruction to continue.
3. If the user changes the plan, treat the updated plan as the approved scope.
4. If the task scope materially changes during execution, stop and provide a new plan.
5. Do not reuse approval from an earlier scope for a materially different task.
## 4. Plan Format
Use this format by default:
```md
## Plan
- Task: ...
- Guides: `AGENTS.md`, `documents/guide/common.md`, ...
- Inspect: ...
- Steps: ...
- Commands: ...
- Expected changes: ...
- Validation: ...
- Commit units: ...
- Stop conditions: ...
```
Use `None` or `Not expected` when an item does not apply.
## 5. Progress Updates
During approved work, provide short progress updates in small units of completed work.
Progress updates should include relevant information such as:
1. what was completed
2. important findings
3. files being changed
4. commands run and their result
5. command failures
6. next action
7. commits created, including their abbreviated hashes and subjects
Keep progress updates concise and factual.
## 6. Command Reporting
When running commands, make the work traceable to the user.
Report:
1. the command
2. the working directory
3. whether it succeeded or failed
4. important output summary
5. next action after failure, if any
Do not paste long command output unless it is necessary to understand an error or result.
## 7. File Change Reporting
When modifying files, report:
1. files changed
2. purpose of the change
3. summary of the change
4. whether unrelated changes were avoided
If formatting or tooling causes a larger-than-expected diff, stop and report it before continuing.
## 8. Scope Changes
Stop and report to the user if any of the following occurs:
1. the task requires changes beyond the approved scope
2. a new task-specific guide becomes relevant
3. a requirement is ambiguous and materially affects the result
4. existing design or public API changes become necessary
5. a destructive action appears necessary
6. a security or sensitive-data issue is found
When the scope changes, provide a new plan and wait for approval before continuing.
## 9. Validation
After making changes, validate the modified scope.
Choose the smallest reliable validation that covers the change. Consider the following, as applicable:
1. directly affected tests
2. package or module tests
3. type checks
4. linters
5. format checks
6. builds
7. full test suite or full build
Do not run larger validation than necessary unless the task requires it or smaller validation is insufficient.
## 10. Validation Failures
If validation fails, report:
1. failed command
2. failure summary
3. whether the failure appears related to the current change
4. planned fix, if the failure will be addressed
5. reason for not fixing it, if it will not be addressed
Do not claim that a failure is pre-existing unless there is evidence.
## 11. Small Verified Commits
1. Every task that creates, modifies, moves, or deletes repository files must
   commit its completed changes.
2. Before editing, divide the task into the smallest logical units that can be
   understood, reviewed, and validated independently.
3. Record those units in the approved plan.
4. Each commit must represent exactly one purpose.
5. Complete, validate, stage, review, and commit one unit before starting the
   next unit.
6. Keep required tests, documentation, generated output, manifests, and
   lockfiles in the same commit as the change that requires them.
7. Split unrelated behavior, cleanup, formatting, renaming, dependency, and
   documentation changes into separate commits.
8. Do not split a change when doing so would leave an uncompilable, invalid,
   unsafe, or misleading intermediate state.
9. Every commit must leave the affected scope in a valid state.
10. Do not create final-history commits labeled `WIP`, `checkpoint`, `fixup`,
    `temporary`, or similar.
11. Stage only the files or patch sections belonging to the current unit.
12. Never include pre-existing, user-authored, generated-by-another-process, or
    otherwise unrelated changes.
13. Review the staged diff and run the smallest reliable validation before each
    commit.
14. Follow `documents/guide/git.md` for commit boundaries, staging, messages,
    and history safety.
15. After committing, verify the commit and report its abbreviated hash and
    subject in the next progress update.
16. Do not amend, squash, or rewrite an earlier commit merely to avoid creating
    a new logical commit. Follow the Git guide and obtain required approval for
    history rewriting.
17. A file-modifying task is not complete while intended changes remain
    uncommitted.
18. If a required commit cannot be created, stop and report the task as
    incomplete with the reason and current repository state.
19. Purely informational tasks and tasks that produce no repository change do
    not require an empty commit.
20. When the repository has a configured remote, a file-modifying task is not
    complete until its task branch is pushed and represented by a pull request.
21. Follow `documents/guide/git.md` for remote, branch, push, and pull-request
    safety.

Use this sequence for each commit unit:

```sh
git status --short
git diff -- <paths>
<validation-command>
git add <paths>
git diff --cached --check
git diff --cached
git commit
git show --stat --oneline --check HEAD
```

## 12. Push
After all approved commit units are complete:
1. verify the working tree is clean
2. verify the current branch is the intended task branch
3. verify the base branch and remote
4. fetch remote state
5. inspect branch divergence
6. run the final validation required for the complete task
7. push the task branch
8. set upstream tracking on the first push
9. verify the remote branch points to the intended commit

Use an explicit push:

```sh
git fetch origin
git status --short --branch
git log --oneline --left-right --graph origin/<base>...HEAD
git push --set-upstream origin HEAD
```

Rules:
1. Do not push directly to the base branch.
2. Do not push unrelated local commits.
3. Do not use `--force`.
4. Use `--force-with-lease` only when history rewriting was explicitly approved
   and the Git guide permits it.
5. Do not push while intended changes remain uncommitted.
6. Do not report push success until the command completes successfully.
7. If push is rejected, inspect remote divergence before deciding whether to
   rebase, merge, or request instructions.
8. Do not bypass branch protection.
9. Do not delete the remote branch while its pull request is open.
10. If no remote is configured, report that push and pull-request creation are
    not applicable.

## 13. Pull Requests
After the branch is pushed:
1. search for an existing open or closed pull request for the same head branch
2. update an existing open pull request instead of creating a duplicate
3. create a pull request when no suitable open pull request exists
4. target the approved base branch
5. verify the pull request head and base branches
6. report the pull request URL

The pull request title must:
1. summarize one coherent task
2. use direct, specific language
3. match the final scope
4. avoid `WIP`, `changes`, `update`, or other vague standalone descriptions

The pull request body must include:
1. problem or purpose
2. changes made
3. validation commands and results
4. risks, compatibility effects, or rollout notes
5. remaining limitations or follow-up work

Use this format by default:

```md
## Summary
- ...

## Validation
- `command`: passed

## Risks
- None
```

Rules:
1. Do not create duplicate pull requests for the same branch.
2. Do not create a pull request before its branch is pushed.
3. Do not target a branch without verifying it is the intended integration
   branch.
4. Do not omit failed or skipped validation.
5. Do not claim a check passed when it was not run.
6. Do not include secrets, credentials, private logs, or sensitive diagnostic
   data in the title or body.
7. Update the title and body when the final scope changes materially.
8. Verify the created or updated pull request through the hosting service.
9. A file-modifying task is incomplete if required pull-request creation fails.
10. If the repository does not use pull requests, that policy must be explicit
    before omitting this step.

Example with GitHub CLI:

```sh
gh pr list --head <branch> --state all
gh pr create \
  --base <base> \
  --head <branch> \
  --title "<title>" \
  --body-file <body-file>
gh pr view <branch> --json number,title,url,headRefName,baseRefName,state
```

## 14. Completion Report
When the task is complete, use this format:
```md
## Completion Report
- Changed: ...
- Commits: ...
- Push: ...
- Pull request: ...
- Validation: ...
- Notes: ...
- Worktree cleanup: Remove the created worktree?
```
`Changed` should include the files or areas modified.
`Commits` should list each abbreviated commit hash and subject in execution
order.
`Push` should include the remote and branch.
`Pull request` should include the pull request number, title, and URL, or state
why it was not applicable.
`Validation` should include the commands run and their result.
`Notes` should include remaining risks, skipped validation, blockers, or follow-up items.
`Worktree cleanup` must ask whether the created worktree should be removed when a worktree was created.
## 15. Blocked or Incomplete Work
If the task cannot be completed, report:
1. what was completed
2. what blocked completion
3. files changed, if any
4. commands run, if any
5. commits created, if any
6. uncommitted changes that remain, if any
7. push status and remote branch, if any
8. pull-request status and URL, if any
9. what instruction is needed from the user
Do not describe blocked, skipped, or failed work as completed.
