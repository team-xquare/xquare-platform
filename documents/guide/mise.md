# MISE.md
This guide defines mise conventions for reproducible local and CI development
environments. It follows the official mise documentation and adds project rules
that keep tool selection and task execution explicit.

## 1. Responsibilities
Use mise for:
1. selecting development-tool versions
2. installing those tools
3. exposing non-secret project environment defaults
4. providing discoverable developer task entry points
5. running commands with the selected tools in local and CI environments

Do not use mise to:
1. replace Bazel's dependency graph or build rules
2. hide application configuration in a developer-only file
3. store secrets
4. duplicate versions owned by another canonical file
5. create complex shell orchestration that belongs in a tested script or build
   rule

## 2. Source of Truth
1. Use the root `mise.toml` as the shared project configuration.
2. Keep personal overrides out of committed shared configuration.
3. Use `mise.local.toml` only for local values and keep it ignored.
4. Define each version in one canonical place.
5. This repository pins Bazelisk in `mise.toml`.
6. This repository pins Bazel itself in `.bazelversion`.
7. Do not add a second Bazel version to `mise.toml`.
8. Do not add `.tool-versions`, `.node-version`, `.go-version`, or similar files
   when `mise.toml` already owns the same version unless another supported tool
   requires that file.
9. If an idiomatic version file is required, document which file is canonical
   and configure mise's support explicitly.

## 3. Configuration Format
1. Use valid TOML.
2. Keep sections in this order when present:
   - `min_version` or repository-wide settings
   - `[tools]`
   - `[env]`
   - task definitions
3. Use one tool per line.
4. Sort tools alphabetically unless dependency order needs explanation.
5. Quote version strings.
6. Add comments only for non-obvious compatibility or ownership decisions.
7. Keep configuration direct; avoid templates when a literal is sufficient.

```toml
[tools]
bazelisk = "1.29.0"
```

## 4. Tool Versions
1. Never omit a version in committed configuration because mise otherwise
   resolves `latest`.
2. Do not commit floating aliases such as `latest`, `stable`, or `lts` for
   build-critical tools.
3. Use an exact version for build-critical tools when no lockfile is used.
4. A major or minor range is acceptable only when `mise.lock` is committed and
   CI installs from it.
5. Change tool versions in a dedicated, reviewable update.
6. Use the CLI to update configuration:

```sh
mise use --pin tool@version
```

7. Do not hand-edit a resolved lockfile.
8. When updating a tool, verify dependent formatters, linters, builds, and tests.
9. Avoid adding two tools that provide the same command unless the selection is
   intentional and documented.
10. Prefer official or mise-registry backends over an unreviewed custom plugin.

## 5. Lockfiles
1. Use `mise.lock` when reproducible resolution and checksums are required.
2. Once adopted, commit `mise.lock`.
3. Generate or update it with mise:

```sh
touch mise.lock
mise install
mise lock
```

4. Review tool-version, URL, checksum, and platform changes.
5. Keep the lockfile change in the same commit as the corresponding
   `mise.toml` change.
6. Do not delete a lockfile to bypass a resolution failure.
7. CI should fail rather than silently select a different version when locked
   reproducibility is required.

## 6. Tool Installation and Execution
1. Install declared tools with:

```sh
mise install
```

2. Verify selected versions with:

```sh
mise current
mise ls
```

3. In scripts and CI, prefer explicit mise execution when shell activation is
   not guaranteed:

```sh
mise exec -- bazel version
mise run test
```

4. Do not assume a developer has globally installed the required tool.
5. Do not bypass mise with an arbitrary system binary in repository tasks.
6. Keep local and CI entry points equivalent.

## 7. Tasks
1. Define a mise task only when it is a meaningful, repeated developer
   operation.
2. Use Bazel targets for build and test graph semantics; a mise task may provide
   a short entry point to those targets.
3. Use lowercase task names with colon-separated namespaces:

```text
format
lint
test
build
ci:check
deps:update
```

4. Give every non-obvious task a concise `description`.
5. A task name must state its outcome, not its implementation.
6. Keep commands short and visible in `mise.toml`.
7. Move long or multi-platform logic to a versioned script.
8. Do not duplicate the same command across several tasks; use task
   dependencies or references.
9. Declare dependencies with mise's task dependency fields rather than calling
   `mise run` recursively in shell.
10. Use `sources` and `outputs` only when they accurately model freshness.
11. Do not use stale-cache behavior for validation tasks that must always run.
12. Add confirmation to destructive or release tasks.
13. Pass arguments through explicitly and document them with the task `usage`
   field when needed.
14. Do not hide required validation behind an optional task.

```toml
[tasks.format]
description = "Format repository source files"
run = "bazel run //tools:format"

[tasks.test]
description = "Run all repository tests"
run = "bazel test //..."

[tasks."ci:check"]
description = "Run required CI validation"
depends = ["format:check", "lint", "test"]
```

## 8. Task Shell Code
1. Prefer one command for a simple task.
2. Use a multiline script only when sequencing is necessary.
3. Let mise's default fail-fast shell behavior stop on errors.
4. Do not disable fail-fast behavior without handling every exit status.
5. Quote shell variables.
6. Do not parse human-readable command output when a structured output mode is
   available.
7. Avoid platform-specific commands in shared tasks unless the supported
   platforms are explicit.
8. Do not change the user's global configuration from a repository task.
9. Do not install undeclared tools from within a normal task.

## 9. Environment Variables
1. Commit only non-secret defaults that are safe and consistent for all users.
2. Use uppercase `SNAKE_CASE` names.
3. Name variables for their real scope and units.
4. Do not use environment variables as undocumented function parameters.
5. Do not store tokens, passwords, credentials, private keys, or personal data
   in `mise.toml`.
6. Use the approved secret manager or ignored local files for secrets.
7. Provide `.env.example` only when the repository has adopted that pattern,
   and include names and safe placeholders rather than values.
8. Treat environment changes as API changes for tasks and applications that
   consume them.
9. Use `false` to unset an inherited variable only when the removal is
   deliberate and documented.

## 10. Trust and Security
1. Review a repository's mise configuration before trusting it.
2. Use `mise trust` only for a known repository.
3. Diagnose trust problems with:

```sh
mise doctor
```

4. Do not globally trust broad directories merely to bypass prompts.
5. Review custom backends, plugins, URLs, task scripts, and environment-file
   loading as executable code.
6. Prefer checksummed locked downloads when available.
7. Do not print secrets in task output.

## 11. Updates
1. Check outdated tools with:

```sh
mise outdated
```

2. Update one tool or one coherent toolchain at a time.
3. Preserve the configured precision unless intentionally changing policy.
4. Use `mise upgrade tool` for an update within the declared range.
5. Use `mise upgrade --bump tool` only when intentionally changing the declared
   range.
6. Review release notes for behavior, security, and compatibility changes.
7. Run all validation affected by the updated tool.
8. Include regenerated lock data in the same change.

## 12. CI
1. Install mise through the repository's approved pinned method.
2. Run `mise install` before invoking tools.
3. Use `mise exec -- <command>` or `mise run <task>` so selected versions are
   unambiguous.
4. Reuse the same task names locally and in CI.
5. Do not put CI-only semantics into a task with a misleading general name.
6. Do not let CI mutate and commit `mise.toml` or `mise.lock`.
7. Cache installed tools only when cache keys include relevant configuration
   and lockfile content.

## 13. Validation
After changing mise configuration:

```sh
mise fmt --check
mise doctor
mise install
mise current
mise tasks ls
mise tasks validate
```

Also run every changed or newly added task.
When a tool version changes, run that tool's affected formatter, linter, build,
and test commands.

## 14. Review Checklist
- Is `mise.toml` the correct owner of this setting?
- Is the version explicit and reproducible?
- Is another file already the canonical source?
- Is the selected backend trusted and maintained?
- Does each task have one obvious outcome?
- Is complex build logic kept in Bazel or a tested script?
- Are environment values free of secrets?
- Are local and CI commands equivalent?
- Were `mise doctor`, installation, and affected tasks validated?

## References
- [mise Configuration](https://mise.jdx.dev/configuration.html)
- [mise Dev Tools](https://mise.jdx.dev/dev-tools/)
- [mise Lockfile](https://mise.jdx.dev/dev-tools/mise-lock.html)
- [mise Tasks](https://mise.jdx.dev/tasks/task-configuration.html)
- [mise Environments](https://mise.jdx.dev/environments/)
