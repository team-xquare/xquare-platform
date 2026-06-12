# GUIDE.md
This guide defines how to create, modify, classify, review, and remove files
under `documents/guide/`.
It exists to keep behavioral instructions accurate, discoverable,
non-duplicative, and consistent with `AGENTS.md`.

## 1. Scope
Read this guide before:
1. creating a guide
2. modifying a guide
3. deleting a guide
4. renaming or moving a guide
5. changing guide ownership or scope
6. changing the guide map in `AGENTS.md`
7. moving a rule between guides

This guide applies to behavioral instructions for agents.
It does not apply to product documentation, architecture decisions, runbooks,
or user documentation stored outside `documents/guide/`.

## 2. Authority and Precedence
Guide changes must preserve this precedence:
1. system and developer instructions
2. `AGENTS.md`
3. `documents/guide/common.md`
4. task-specific guides
5. local repository conventions not documented as guide rules

Rules:
1. A task-specific guide may add stricter requirements.
2. A task-specific guide must not weaken or override `common.md`.
3. A guide must not contradict `AGENTS.md`.
4. A lower-level guide must not create an exception to a higher-level rule
   unless the higher-level rule explicitly permits that exception.
5. When two rules conflict, stop and resolve the conflict instead of choosing
   whichever rule is easier to follow.
6. Do not use wording such as `normally`, `prefer`, or `when practical` to
   weaken an existing mandatory rule.

## 3. Guide Purpose
Every guide must have one primary responsibility.

A guide should answer:
1. what task scope it governs
2. when it must be read
3. which decisions it standardizes
4. which checks are mandatory
5. which actions are prohibited
6. how compliance is validated

Do not create a guide merely because a topic exists.
Create one when the topic has enough stable, reusable behavioral rules to
justify separate ownership.

## 4. Rule Placement
Place a rule in the narrowest guide that owns it.

Use:
- `common.md` for safety and baseline rules that apply to every task
- `work.md` for planning, approval, execution, progress, validation, commits,
  and completion reporting
- `git.md` for Git state, branches, staging, commits, history, and worktrees
- `research.md` for evidence collection, investigation, profiling, and
  evidence-based conclusions
- language guides for language-specific source, API, tooling, and test rules
- `bazel.md` and `mise.md` for their corresponding tools
- `architecture.md` for module boundaries, responsibilities, and dependency
  direction
- `guide.md` for guide maintenance rules

Rules:
1. Do not copy the same normative rule into several guides.
2. Put the authoritative rule in one guide.
3. Refer to the authoritative guide when another guide needs context.
4. A short task-specific restatement is allowed only when omission would create
   a realistic safety risk.
5. A restatement must preserve the same meaning and strength.
6. If a rule affects every task, consider whether it belongs in `common.md`.
7. If a rule describes execution procedure rather than domain behavior, place
   it in `work.md`.
8. If a rule changes Git behavior, `git.md` must be read and considered even
   when the triggering task is documentation work.

## 5. Before Modifying a Guide
Before editing:
1. read `AGENTS.md`
2. read `common.md`
3. read `guide.md`
4. read `work.md` for executable work
5. read `git.md` for commits or repository-state operations
6. read the target guide
7. read only directly related guides whose scope may overlap
8. inspect the current Git state
9. identify user-authored or unrelated changes
10. state whether the guide map requires an update

Do not read every guide without a scope reason.
Use the map in `AGENTS.md` to select related guides.

## 6. Define the Change
Before writing, record:
1. problem with the current guidance
2. desired behavior
3. affected task types
4. target guide
5. related guides
6. required strength of the rule
7. examples or evidence needed
8. validation method
9. compatibility or migration impact

Do not add a rule when the expected behavior is materially ambiguous.
Obtain the missing policy decision first.

## 7. Normative Language
Use language consistently.

### 7.1 Mandatory Rules
Use `must`, `must not`, `always`, `never`, or direct imperative language only
for requirements.

Examples:

```text
Run the affected tests before committing.
Do not expose diagnostic endpoints publicly.
```

### 7.2 Recommended Rules
Use `should` or `prefer` for a default that permits justified exceptions.

State:
1. why it is preferred
2. when an exception is acceptable
3. what evidence is required for the exception

### 7.3 Optional Actions
Use `may` for a genuinely optional action.

### 7.4 Avoid Ambiguous Language
Avoid:
- `try to`
- `usually` without an exception rule
- `as needed`
- `properly`
- `appropriately`
- `best practice`
- `clean`
- `simple`
- `reasonable`

unless the document defines the concrete decision criteria.

## 8. Rule Design
Every rule should be:
1. necessary
2. specific
3. actionable
4. testable or reviewable
5. scoped
6. consistent
7. safe
8. durable

A strong rule states:
1. trigger
2. required action
3. prohibited action when relevant
4. exception condition
5. validation method

Avoid rules that only express taste.
Explain the risk controlled by a non-obvious rule.

## 9. Rule Strength
Choose the weakest rule strength that reliably controls the risk.

Use a mandatory rule when violation may cause:
- data loss
- security exposure
- destruction of user work
- incorrect behavior
- irreproducible builds
- invalid repository history
- unreviewable changes
- violation of an external contract

Use a recommendation when:
- several valid approaches exist
- consistency is useful but exceptions are safe
- the correct choice depends on local context
- measurement should decide

Do not make every preference mandatory.
Excessive mandatory rules make real safety requirements harder to identify.

## 10. Exceptions
1. State exceptions next to the rule they qualify.
2. Define who or what can authorize an exception.
3. Define required evidence or documentation.
4. Keep exceptions narrower than the rule.
5. Do not create an exception through an example or note.
6. Do not use an exception to bypass approval, security, or user-work
   protections.
7. Remove obsolete exceptions.

## 11. Prohibitions
A prohibition must identify the behavior being prevented.

Good:

```text
Do not use `git add .` when unrelated files are present because it can stage
user-authored changes.
```

Avoid:

```text
Be careful when staging.
```

Do not add broad prohibitions that prevent valid work without an escape
condition.
Do not prohibit a tool merely because another tool is preferred.

## 12. Examples
Examples clarify rules but do not replace them.

Rules:
1. Put the normative statement before the example.
2. Label examples as `Good`, `Avoid`, or describe their purpose.
3. Keep examples minimal.
4. Use placeholders for secrets, domains, IDs, and paths.
5. Do not include real credentials or personal data.
6. Ensure commands are syntactically valid.
7. Ensure code follows the guide it demonstrates.
8. Avoid examples tied to temporary versions unless version-specific behavior
   is the subject.
9. Do not let an example imply an undocumented exception.
10. Update or remove examples when the rule changes.

## 13. Commands
1. Commands must identify the intended working directory when it is not obvious.
2. Commands must be non-destructive by default.
3. Use placeholders such as `<target>` for values the reader must replace.
4. Do not include environment-specific absolute paths unless the guide governs
   that exact path.
5. Do not include commands that print secrets or sensitive configuration.
6. Distinguish commands that modify files from read-only checks.
7. State when a command may be expensive or disruptive.
8. State the minimum reliable validation command.
9. Do not claim a command validates behavior it does not cover.

## 14. External Sources
Use external sources when a guide depends on:
1. current tool behavior
2. language semantics
3. official conventions
4. version-specific configuration
5. security guidance
6. profiler or runtime behavior

Rules:
1. Prefer specifications and official documentation.
2. Use current sources.
3. Verify that every reference URL is accessible.
4. Distinguish official requirements from repository policy.
5. Do not turn a third-party opinion into a mandatory repository rule without a
   project reason.
6. Do not copy large passages from a source.
7. Summarize the relevant behavior.
8. Record version scope when behavior differs by version.
9. Remove references that no longer support the rule.

## 15. Repository-Specific Rules
Repository rules must reflect actual repository decisions.

Do not:
- invent a framework, tool, runtime, or deployment assumption
- mandate a tool that is not adopted
- describe an aspirational system as current behavior
- copy conventions from another repository without checking compatibility
- duplicate version ownership across configuration files

When the repository has no implementation or configuration evidence:
1. keep the rule technology-neutral
2. make platform-specific sections conditional
3. state the assumption
4. request a project decision when the result materially depends on it

## 16. Safety and Security
Guide content must not:
1. expose secrets
2. normalize destructive operations
3. weaken approval requirements
4. permit overwriting unrelated work
5. recommend public diagnostic endpoints
6. include production credentials or identifiers
7. encourage collection of sensitive data without controls
8. conceal operational risk

Security and user-work protections must remain explicit even when they make a
procedure longer.

## 17. Duplication
Before adding a rule:
1. search `AGENTS.md`
2. search `documents/guide/`
3. identify the existing owner
4. decide whether to reference, move, or replace the rule

When removing duplication:
1. select one authoritative location
2. preserve the strongest necessary requirement
3. replace other copies with references when useful
4. verify that task readers will still discover the rule through `AGENTS.md`
5. do not silently weaken behavior

## 18. Contradiction Review
Check for:
1. direct conflicts
2. different meanings for the same term
3. mandatory versus optional wording
4. incompatible command sequences
5. conflicting validation scopes
6. conflicting file ownership
7. circular references
8. impossible combinations of rules

If a contradiction exists:
1. identify both rules
2. identify their authority and scope
3. preserve the higher-priority rule
4. revise the lower-priority rule
5. update references and examples
6. report the behavioral change

Do not resolve a conflict only by adding another ambiguous exception.

## 19. Guide Structure
Use this general structure when applicable:
1. purpose and scope
2. authority and precedence
3. core rules
4. detailed procedures
5. exceptions
6. prohibited patterns
7. validation
8. review checklist
9. references

Rules:
1. Use one level-one heading matching the file topic.
2. Use numbered level-two sections for long guides.
3. Use level-three headings only when they improve navigation.
4. Keep related rules together.
5. Put references at the end.
6. Do not add a table of contents unless navigation is otherwise difficult.
7. Keep Markdown compatible with common renderers.

## 20. File Names
1. Use lowercase file names.
2. Use one stable topic per file.
3. Use `.md`.
4. Keep the level-one heading in uppercase with `.md`, matching the existing
   repository convention.
5. Do not rename a guide only for stylistic preference.
6. Treat renaming as a reference and discovery change.
7. Update all links, map entries, and instructions when renaming.

## 21. `AGENTS.md` Guide Map
Update the map when:
1. adding a guide
2. deleting a guide
3. renaming or moving a guide
4. materially changing a guide's scope
5. changing when the guide must be read
6. changing its key checks

Do not update the map when:
- editing detail without changing scope
- adding examples
- correcting wording
- updating references

Every map entry must include:
1. exact file path
2. concise scope
3. objective read trigger
4. key checks

Rules:
1. Keep map entries accurate.
2. Keep paths in sync with the filesystem.
3. Do not list non-guide documents.
4. Do not add behavior outside the map and guide files to `AGENTS.md`.
5. Verify that every guide is listed exactly once.
6. Verify that every listed guide exists.

## 22. Creating a Guide
Before creating a guide:
1. prove the topic is not owned by an existing guide
2. define its scope
3. define its read trigger
4. define its key checks
5. identify related guides
6. decide whether common rules also need changes

Creation requires:
1. the new guide
2. an updated `AGENTS.md` map entry
3. updated references from related guides when required
4. validation that no duplicate owner exists
5. one coherent commit containing the guide and required discovery changes

Do not create an empty placeholder guide unless the repository explicitly
requires the planned file structure.

## 23. Modifying a Guide
For every modification:
1. state the previous behavior
2. state the new behavior
3. identify whether rule strength changed
4. identify affected task scopes
5. search for duplicate or contradictory rules
6. update examples and references
7. update the `AGENTS.md` map only if scope or discovery changed
8. validate the complete guide, not only the edited lines
9. commit the change as one logical unit

Treat these as material behavioral changes:
- changing `must` to `should`
- adding or removing an exception
- broadening or narrowing scope
- changing approval requirements
- changing required validation
- changing file ownership
- changing destructive-operation policy
- changing security requirements

Report material changes explicitly.

## 24. Moving Rules
When moving a rule:
1. identify its authoritative destination
2. add it to the destination
3. replace or remove the old copy in the same change
4. update cross-references
5. preserve rule strength
6. verify that readers still discover it
7. avoid a commit where two guides claim conflicting ownership

## 25. Renaming or Moving a Guide
Renaming or moving requires:
1. a task-level reason
2. updated `AGENTS.md`
3. updated guide references
4. updated commands and links
5. verification that no stale path remains
6. protection of file history where practical
7. one coherent commit

Search for stale references with:

```sh
rg -n 'documents/guide/<old-name>\.md|<old-name>\.md' .
```

## 26. Deleting a Guide
Delete a guide only when:
1. its rules are obsolete
2. its scope is no longer supported
3. its rules have moved to an authoritative destination
4. it duplicates another guide

Before deletion:
1. classify every rule as obsolete, moved, or intentionally removed
2. update `AGENTS.md`
3. update all references
4. verify no required behavior is lost
5. explain any removed requirement
6. search for stale paths

Do not delete a guide merely because it is short.

## 27. Generated and Automated Changes
1. Do not generate guide content without reviewing every rule.
2. Do not bulk-reformat all guides for a targeted change.
3. Do not accept automated wording changes that alter normative strength.
4. Review changes to `must`, `should`, `may`, `never`, and `not`.
5. Keep generated indexes out of guide files unless the repository adopts them
   explicitly.
6. Do not modify generated documentation when the source guide is the
   authoritative file.

## 28. Validation
Validate the modified guide with the smallest reliable checks.

Required checks:
1. Markdown structure
2. balanced code fences
3. trailing whitespace
4. final newline
5. valid internal paths
6. accessible external references
7. no secrets or sensitive values
8. no duplicate or contradictory rules
9. accurate `AGENTS.md` map when affected
10. staged diff contains only the approved change

Suggested commands:

```sh
git diff --check
rg -n '^#{1,6} ' documents/guide/<guide>.md
rg -n 'documents/guide/' AGENTS.md documents/guide/
git diff --cached
```

Use a configured Markdown linter when available.
Do not install a new linter solely for a small documentation change unless the
task approves it.

## 29. Review Checklist
- Does the rule belong in this guide?
- Is there one authoritative owner?
- Does it preserve `AGENTS.md` and `common.md`?
- Is the trigger clear?
- Is the action testable?
- Is the normative strength intentional?
- Are exceptions explicit and narrow?
- Does it protect user work and sensitive data?
- Does it avoid unsupported repository assumptions?
- Are examples correct and non-sensitive?
- Are commands safe and reproducible?
- Are external sources official and current?
- Are duplicate and contradictory rules removed?
- Does the guide map remain accurate?
- Were references, paths, code fences, and whitespace validated?
- Does the staged diff contain only the approved guide change?

## 30. Completion Report
When completing a guide change, report:
1. guide files changed
2. behavioral rules added, removed, or changed
3. whether rule strength changed
4. whether `AGENTS.md` changed and why
5. validation commands and results
6. external references checked
7. commit hash and subject
8. remaining ambiguity or follow-up work
