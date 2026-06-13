# BAZEL.md
This guide defines Bazel conventions for this repository.
It follows the official Bazel BUILD and Starlark style guides, with additional
project rules that favor explicit, readable, and hermetic build definitions.

## 1. Authority and Scope
1. Follow Bazel's language semantics and official documentation first.
2. Follow this guide for choices that Bazel intentionally leaves to projects.
3. Keep BUILD files declarative, simple, and repetitive.
4. Do not hide source structure or dependencies behind clever Starlark.
5. This repository pins Bazel with `.bazelversion`; invoke Bazel through
   Bazelisk or a repository task that uses Bazelisk.
6. Bazel 9 uses Bzlmod. Declare external dependencies in `MODULE.bazel`, not a
   legacy `WORKSPACE` file.

## 2. Readability Principles
1. A reader must be able to identify a target's sources, direct dependencies,
   visibility, and outputs without following unnecessary indirection.
2. Prefer an explicit rule declaration over a macro when the rule is used only
   once or the macro would merely rename attributes.
3. Prefer repeated, conventional BUILD syntax over compact metaprogramming.
4. Keep package boundaries aligned with source ownership and dependency
   boundaries.
5. A BUILD change must describe the real dependency graph, not only make the
   build pass.

## 3. Repository and Version Configuration
1. Keep `.bazelversion` committed.
2. Change the Bazel version in a dedicated change with validation of all
   supported targets.
3. Do not duplicate the Bazel version in `mise.toml`; `mise.toml` pins
   Bazelisk, while `.bazelversion` selects Bazel.
4. Keep repository-wide options in `.bazelrc`.
5. Put stable, shared options in `.bazelrc`; do not add personal paths, machine
   names, credentials, or local-only flags.
6. Name intentional configurations, for example:

```text
build:ci --color=no
test:ci --test_output=errors
```

7. Do not make aliases silently change semantics. A configuration named `ci`
   may adjust output and remote execution settings, but must not skip required
   tests.
8. Use Bzlmod module dependencies with explicit versions.
9. Commit `MODULE.bazel.lock` when Bazel generates it.
10. Do not edit generated lockfiles manually.
11. Bazel owns the compiler, runtime, and build-time linter versions used by
    build and test actions, including Java, Go, Node, and ktlint.
12. Declare shared language toolchains under `build/<language>/` and
    system-specific tools and dependencies in that system's module fragment.
13. Keep the root `MODULE.bazel` focused on composing shared language and
    system module fragments.
14. Do not require developers or CI to install a Bazel-owned tool separately
    through mise or the host operating system.

## 4. File Names and Package Boundaries
1. Prefer `BUILD.bazel` over `BUILD` for new packages.
2. Use one `BUILD.bazel` file per Bazel package.
3. Add a package boundary when a directory has a distinct responsibility,
   ownership, visibility policy, or dependency lifecycle.
4. Do not create a package for every file.
5. Do not allow one large package to recursively absorb unrelated
   subdirectories.
6. Place reusable build logic in focused `.bzl` files.
7. Use lowercase `snake_case` for `.bzl` file names.
8. Keep language-specific build definitions near the corresponding source.

## 5. BUILD File Structure
Order a BUILD file as follows:
1. `load()` statements
2. `package()` and package-level declarations
3. `licenses()` or `exports_files()` when required
4. public libraries and binaries
5. internal implementation targets
6. tests
7. supporting file groups or aliases

Keep one blank line between top-level declarations.

```starlark
load("@rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:private"])

go_library(
    name = "token",
    srcs = [
        "claims.go",
        "token.go",
    ],
    importpath = "example.com/project/token",
    visibility = ["//services/auth:__pkg__"],
    deps = [
        "//internal/clock",
    ],
)

go_test(
    name = "token_test",
    srcs = ["token_test.go"],
    embed = [":token"],
)
```

## 6. Formatting
1. Format all BUILD and `.bzl` files with `buildifier`.
2. Do not manually align values with extra spaces.
3. Use a trailing comma for every item in a multiline list or call.
4. Keep one list item per line when a list has multiple entries.
5. Let `buildifier` determine attribute ordering and line wrapping.
6. Do not disable formatter rewrites without a documented compatibility reason.
7. Validate formatting without modifying files:

```sh
buildifier -mode=check -lint=warn -r .
```

## 7. Target Names
1. Use lowercase `snake_case`.
2. Name a target for the artifact or responsibility it represents.
3. Use the conventional package target name when it improves label clarity.
4. Suffix tests with `_test`.
5. Suffix binaries only when the distinction is useful; avoid names such as
   `main`, `app`, `lib`, `common`, or `misc` without domain context.
6. Name generated or helper targets so their relationship is visible:
   `api_proto`, `api_go_proto`, `api_codegen`.
7. Do not encode implementation details or temporary ticket numbers in target
   names.

Good:

```text
//services/auth:token_validator
//services/auth:token_validator_test
```

Avoid:

```text
//services/auth:utils
//services/auth:new_lib_v2
```

## 8. Labels
1. Use the shortest unambiguous label accepted in the current context.
2. Use `:target` for a target in the current package.
3. Use `//path/to/package:target` across packages.
4. Do not use filesystem-like relative paths to obscure package boundaries.
5. Keep labels stable; renaming a label is an API change for its dependents.
6. Use repository names from Bzlmod's apparent repository mapping in authored
   BUILD files. Do not hard-code canonical repository names generated by Bazel.

## 9. Sources
1. List sources explicitly when the list is short and stable.
2. A non-recursive `glob()` is acceptable for homogeneous source sets where
   adding a matching file should automatically include it.
3. Never use recursive globs such as `glob(["**/*.go"])`.
4. Exclude tests, generated files, fixtures, and tools explicitly when using a
   glob.
5. Do not use a glob when inclusion requires review or when filenames carry
   different build semantics.
6. Generated files must be outputs of declared Bazel targets.
7. Do not read undeclared files from the source tree at execution time.

Good:

```starlark
srcs = glob(
    ["*.ts"],
    exclude = ["*_test.ts"],
)
```

Prefer an explicit list when there are only a few files:

```starlark
srcs = [
    "parser.ts",
    "token.ts",
]
```

## 10. Dependencies
1. Declare every direct dependency.
2. Do not rely on transitive dependencies.
3. Remove dependencies that are no longer used.
4. Keep `deps` sorted according to `buildifier`.
5. Prefer the narrowest target that provides the required API.
6. Do not depend on a broad aggregate target to avoid declaring precise
   dependencies.
7. Avoid dependency cycles. Fix ownership or extract a focused lower-level
   package instead of adding indirection that preserves the cycle.
8. Review dependency changes with `bazel query`.

```sh
bazel query 'deps(//path/to/package:target)'
bazel query 'rdeps(//..., //path/to/package:target)'
```

## 11. Visibility
1. Default packages to private visibility:

```starlark
package(default_visibility = ["//visibility:private"])
```

2. Open visibility only to actual consumers.
3. Prefer `:__pkg__` for one package and `:__subpackages__` only when all
   descendants are intended consumers.
4. Use `//visibility:public` only for deliberate repository-wide APIs.
5. Treat visibility expansion as an API decision requiring review.
6. Do not disable visibility checking to make a build pass.
7. Restrict `.bzl` load visibility when build logic is internal.
8. Keep implementation targets private even when a facade target is public.

## 12. Rule Selection
1. Use the language's standard Bazel rules.
2. Use library rules for reusable code, binary rules for executable entry
   points, and test rules for executable tests.
3. Do not use `filegroup` as a substitute for a typed language rule.
4. Use `alias` only for a stable compatibility name, configuration selection,
   or a deliberate public facade.
5. Use `genrule` only when no purpose-built rule exists.
6. A `genrule` command must use declared tools and inputs and write only declared
   outputs.
7. Do not call host-installed tools directly from build actions.
8. Prefer a custom rule over repeated complex `genrule` commands.

## 13. Macros and Custom Rules
1. Keep BUILD files direct until repetition represents a real policy.
2. Create a macro when it consistently applies required defaults, connects a
   known family of targets, or prevents a meaningful class of mistakes.
3. Do not create a macro solely to shorten a rule name or hide two attributes.
4. Macro names must describe the concept they create.
5. Preserve conventional attributes such as `name`, `srcs`, `deps`,
   `visibility`, `tags`, and `testonly` where applicable.
6. Forward caller visibility deliberately; do not accidentally make generated
   targets public.
7. Prefix private symbols in `.bzl` files with `_`.
8. Document public macros, rules, providers, and non-obvious attributes.
9. Keep implementation targets generated by a macro predictably named and
   private.
10. Use custom rules only when analysis-time providers, toolchains, execution
    groups, or declared build actions are required.
11. Separate rule implementation from repository setup and user-facing macros.

## 14. Starlark Code
1. Write Starlark for readability, not language cleverness.
2. Use lowercase `snake_case` for functions, variables, and attributes.
3. Use uppercase `UPPER_SNAKE_CASE` only for true constants.
4. Keep functions short and single-purpose.
5. Prefer guard clauses with `fail()` for invalid input.
6. Include the invalid attribute value and required condition in failure
   messages.
7. Avoid mutable module-level state.
8. Avoid top-level list comprehensions in BUILD files.
9. Avoid unnecessary variables in BUILD files.
10. Do not perform work in a macro that belongs in rule analysis.
11. Do not depend on iteration order unless the value is deliberately ordered.
12. Prefer named helper functions over deeply nested expressions.

## 15. Configuration and `select()`
1. Use `select()` only for real platform, feature, or build-mode differences.
2. Name `config_setting` targets for the condition they represent, not the flags
   used to implement them.
3. Include `//conditions:default` unless every possible configuration is
   intentionally covered.
4. Keep selected values of the same conceptual kind.
5. Avoid nested or duplicated `select()` expressions.
6. Extract a shared build setting only when multiple packages use the same
   stable concept.
7. Do not use configuration to hide source-level architecture problems.

## 16. Hermeticity and Reproducibility
1. Every action must declare all inputs, tools, environment requirements, and
   outputs.
2. Do not depend on the current directory, user home, local PATH contents,
   system time, locale, hostname, or network access unless a rule explicitly
   models that dependency.
3. Use Bazel toolchains for compilers and platform tools.
4. Run build-time linters through Bazel targets with declared inputs and the
   registered toolchain runtime.
5. Keep repository rules and module extensions deterministic.
6. Verify downloaded artifacts with checksums when the dependency mechanism
   supports them.
7. Do not fetch dependencies from arbitrary URLs during normal build actions.
8. Do not write into the source tree from a build or test.
9. Keep generated outputs deterministic: stable ordering, stable timestamps,
   and no machine-specific paths.

## 17. Tests
1. Declare tests with the language-specific test rule.
2. Keep each test target focused on one package or coherent behavior.
3. Set `size` according to actual resource usage and duration.
4. Use `timeout` only when the default does not represent the test.
5. Use `tags` to describe execution requirements, not to conceal failures.
6. Do not add `manual`, `flaky`, or exclusion tags without a documented reason.
7. Tests must declare fixtures and data through the rule's `data` attribute.
8. Tests must not depend on the invocation directory or undeclared environment
   variables.
9. Prefer deterministic fakes over network services in unit tests.

## 18. Comments
1. Comment why a non-obvious build constraint exists.
2. Do not narrate attributes that are already clear.
3. Put comments immediately above the declaration or attribute they explain.
4. Include an issue or removal condition for temporary compatibility settings.
5. Do not preserve dead targets or dependencies in comments.

## 19. Change Procedure
1. Inspect the affected targets before editing:

```sh
bazel query //path/to/package:all
```

2. Format the changed BUILD and `.bzl` files.
3. Build the narrowest affected production targets.
4. Test the narrowest affected test targets.
5. Expand validation when shared macros, rules, toolchains, or module
   dependencies change.
6. Review reverse dependencies before changing a public target or visibility.
7. For repository-wide build logic, run:

```sh
bazel build //...
bazel test //...
```

only when the full scope is justified and supported by the repository.

## 20. Review Checklist
- Is each target named for a clear responsibility?
- Are all sources and direct dependencies declared?
- Is visibility as narrow as possible?
- Is the BUILD file understandable without opening a macro?
- Is a new macro enforcing real policy rather than hiding ordinary rules?
- Are build actions hermetic and deterministic?
- Are generated and external files pinned and verified?
- Did `buildifier` pass?
- Did affected builds and tests pass?

## References
- [Bazel BUILD Style Guide](https://bazel.build/build/style-guide)
- [Bazel `.bzl` Style Guide](https://bazel.build/rules/bzl-style)
- [Bazel Best Practices](https://bazel.build/configure/best-practices)
- [Bazel Visibility](https://bazel.build/concepts/visibility)
- [Bazel Bzlmod](https://bazel.build/external/module)
- [Buildifier](https://github.com/bazelbuild/buildtools/tree/master/buildifier)
