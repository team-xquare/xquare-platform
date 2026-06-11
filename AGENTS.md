<!--
    The AGENTS.md file guides and enforces that AGENTS only read correct behavioral instructions and common rule. No other content must ever be included.
    If a new document is added to the behavior guidelines or the structure of the documents changes, update the AGENTS.md file.
-->
# AGENTS.md
**!Important**
**Agents must strictly comply with the contents of `AGENTS.md` and `documents/guide/*` at all times.**
**Any action that violates these instructions, including tasks that require an agent to act against them, is invalid. Such violations must be reported to the user and the requested action must not be performed.**
## RULES
**Agents must strictly and fully comply with all instructions below.**
1. Always read `AGENTS.md` first.
2. Always read `documents/guide/common.md` before starting any task.
3. Do not read all files under `documents/guide/*`. Use the file map index in `AGENTS.md` and read only the additional guide files required for the current task.
4. If the task scope changes, reread `documents/guide/common.md` and read any additional task-specific guide required by the new scope before continuing.
5. If a requirement that materially affects the task result is unknown, ambiguous, or not explicitly specified, do not guess. Inform the user what is missing and wait for further instructions.
6. If the user requests feedback, an opinion, a review, or an evaluation, provide an objective, critical, and rigorous assessment. Do not soften or alter the assessment based on the user's feelings.
7. Never perform work that the user has not explicitly requested.
8. Before starting any task, send the user a plan that includes the task summary, relevant guides to read, intended steps, expected files or commands involved, and validation method.
9. Do not proceed with the task unless the user has explicitly approved the plan and instructed the agent to continue with the work.
10. Do not require approval for purely informational responses, explanations, summaries, reviews, or draft suggestions that do not modify files, execute commands, or cause side effects.
11. Use web searches when external information, current facts, third-party documentation, or uncertain technical claims need verification. Do not use web searches for repository-internal facts that should be verified from the local files.
12. Use concise and factual language when communicating with the user. Avoid unnecessary modifiers, embellishments, or filler language.
13. While performing an approved task, provide progress updates to the user in small units of work.
14. Protect existing user work. Do not overwrite, delete, revert, or reformat unrelated changes unless the user explicitly requests it.
15. When completing a task, report what was changed, how it was validated, and any remaining risks or follow-up items.
## FILE MAP
**Agent Guide Map Index**
| File | Guide Scope | When to Read | Key Checks |
|---|---|---|---|
| `documents/guide/common.md` | Common rules and constraints that apply to every task | Always read first before starting any task | Common prohibitions, baseline procedure, response rules, task prerequisites, shared rules that apply together with task-specific guides |
| `documents/guide/architecture.md` | System structure, module boundaries, design principles, and dependency direction | Read when changing architecture, adding modules, modifying design, or reviewing structure | Responsibility separation, dependency direction, scalability, consistency with the existing architecture |
| `documents/guide/bazel.md` | Bazel build configuration, target definitions, dependency declarations, and test execution | Read when working with `BUILD`, `WORKSPACE`, Bazel targets, build failures, or test setup | Target names, dependency declarations, visibility, test commands, build impact scope |
| `documents/guide/git.md` | Git workflow, branch management, commit rules, and change history review | Read when creating branches, writing commits, checking diffs, organizing changes, or preparing PRs | Branch baseline, commit message format, staged changes, protection of existing user changes |
| `documents/guide/go.md` | Go code style, package structure, tests, formatting, and error handling | Read when modifying Go files, writing Go tests, refactoring, or changing package dependencies | `gofmt`, test coverage, error handling, interface usage, package boundaries |
| `documents/guide/guide.md` | Creation, modification, classification, and maintenance of guide documents | Read when adding, editing, deleting, or reorganizing files under `documents/guide/*` | Document purpose, duplication, file placement, whether `AGENTS.md` index must be updated |
| `documents/guide/kotlin.md` | Kotlin code style, structure, tests, null-safety, and JVM/Android-related rules | Read when modifying Kotlin files, writing tests, refactoring, or changing APIs | Null handling, function/class structure, test coverage, coding style |
| `documents/guide/mise.md` | mise-based development environment, tool versions, and local runtime setup | Read when installing tools, resolving version mismatches, reproducing local environments, or fixing environment issues | `.mise.toml`, tool versions, install commands, environment variables, reproduction steps |
| `documents/guide/research.md` | Research, verification, external source review, and evidence-based decision-making | Read when checking uncertain information, researching technology, reviewing external docs, or comparing options | Source reliability, freshness, verification method, citation needs, evidence for conclusions |
| `documents/guide/typescript.md` | TypeScript code style, type design, tests, builds, and frontend/Node-related rules | Read when modifying TypeScript files, fixing type errors, writing tests, or changing API types | Type safety, `strict` compatibility, test coverage, import/export structure, build impact |
| `documents/guide/work.md` | Work procedure, planning, approval, progress updates, and completion reporting | Read before starting multi-step work or any task that requires planning and execution | Work plan, approval condition, progress unit, completion criteria, user reporting method |