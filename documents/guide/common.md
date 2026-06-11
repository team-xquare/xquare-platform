# COMMON.md
This guide defines common safety rules that apply to every task.
Task-specific guides may add stricter rules, but they must not weaken or override this guide.
## 1. Protect Existing Work
1. Preserve all existing user changes.
2. Do not overwrite, delete, revert, or discard unrelated changes.
3. Do not reformat unrelated files.
4. Do not modify generated files unless the task explicitly requires it.
5. If unrelated issues are found, report them separately instead of fixing them.
## 2. Git Worktree
1. Any task that modifies files must be performed in a separate Git worktree.
2. Do not modify files directly in the user's current working tree unless the user explicitly instructs otherwise.
3. Before modifying files, inspect the Git state of the relevant worktree.
4. After completing a file-modifying task, ask the user whether the created worktree should be removed.
5. Do not remove the worktree without explicit user approval.
## 3. Scope Safety
1. Make only the changes required for the requested task.
2. Prefer minimal, targeted changes.
3. Do not make opportunistic refactors, cleanup changes, dependency changes, or architecture changes.
4. Do not rename, move, delete, or reorganize files unless required by the task.
5. Preserve existing behavior unless the user explicitly requests a behavior change.
## 4. Destructive Actions
1. Do not use destructive commands unless the user explicitly requests them.
2. Destructive commands include `git reset --hard`, `git clean`, force push, mass deletion, and broad overwrite operations.
3. If a destructive action appears necessary, stop and explain why explicit approval is required.
## 5. File Editing
1. Follow the existing style, structure, naming, and conventions of the repository.
2. Keep changes small and reviewable.
3. Do not introduce unnecessary abstraction.
4. Do not add new dependencies, tools, frameworks, or architectural patterns unless required by the task.
5. Keep comments and documentation in English unless the surrounding file or repository convention uses another language.
## 6. Validation
1. Validate changes with the smallest reliable command that covers the modified scope.
2. Run relevant tests, formatters, linters, type checks, or builds when applicable.
3. Do not claim validation passed unless the validation command was actually run and completed successfully.
4. If validation cannot be performed, report why.
5. If validation is partial, state what was and was not validated.
## 7. Security and Sensitive Data
1. Do not expose, print, copy, commit, or document secrets.
2. Secrets include tokens, passwords, private keys, credentials, session values, and personal data.
3. If sensitive data is found, report its presence without revealing the value.
4. Do not add secrets to source code, tests, examples, logs, or documentation.
5. Prefer existing environment-variable or secret-management patterns when credentials are required.