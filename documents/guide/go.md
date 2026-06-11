# GO.md
This guide defines Go conventions for this repository.
It combines official Go conventions with the Uber Go Style Guide. Where they
conflict, the Go language specification, standard tooling, and official Go
guidance take precedence.

## 1. Goals and Priority
Write Go code that is:
1. correct
2. simple
3. explicit about ownership and failure
4. idiomatic to Go
5. easy to test
6. safe under concurrency
7. easy to remove or change

Apply rules in this order:
1. Go language and standard-library behavior
2. `gofmt` and standard Go tooling
3. official Go documentation and review comments
4. this guide
5. Uber Go Style Guide where this guide does not specify otherwise
6. consistent local package convention

Do not copy a local pattern that is clearly unsafe or contradicts a higher
priority rule.

## 2. Required Tooling
1. Format every Go file with `gofmt`.
2. Use `goimports` when available to format and group imports.
3. Run focused tests for the changed package.
4. Run `go vet` for affected packages.
5. Run the race detector for changed concurrent behavior when supported.

```sh
gofmt -w path/to/file.go
go test ./path/to/package
go vet ./path/to/package
go test -race ./path/to/package
```

Do not manually format code against `gofmt`.

## 3. Package Design
1. Give each package one clear responsibility.
2. Package names are short, lowercase, singular, and descriptive.
3. Avoid generic package names such as `util`, `utils`, `common`, `base`,
   `helper`, or `misc`.
4. Do not repeat the package name in exported identifiers.
5. Keep package APIs smaller than their implementations.
6. Put code in `internal` when external modules must not depend on it.
7. Use `cmd/<name>` for multiple executable entry points.
8. Keep `main` packages thin: parse configuration, wire dependencies, start the
   process, and report terminal errors.
9. Do not create packages only to mirror architecture layers when the package
   has no independent responsibility.
10. Avoid import cycles by fixing responsibility boundaries.

Good:

```go
package user

type Service struct{}
```

Avoid:

```go
package user

type UserService struct{}
```

## 4. Files
1. Use lowercase file names.
2. Use underscores only when they improve readability or are required by Go
   conventions such as `_test.go` and platform suffixes.
3. Group files by responsibility, not by arbitrary type category.
4. Keep closely related declarations near one another.
5. Do not create one file per small type by default.
6. Split a file when it contains multiple independent responsibilities or
   becomes difficult to navigate.
7. Put tests in files ending with `_test.go`.
8. Keep generated code clearly marked with the standard generated-code header.
9. Do not hand-edit generated files.

## 5. Imports
1. Let `goimports` group imports.
2. Use import groups in this order:
   - standard library
   - third-party packages
   - repository-local packages
3. Avoid dot imports.
4. Use blank imports only for a documented side effect required by the package.
5. Use an import alias only to resolve a collision or clarify a genuinely
   misleading package name.
6. The alias should describe the imported package, not its call site.
7. Remove unused dependencies and imports immediately.

## 6. Naming
1. Use `MixedCaps` or `mixedCaps`; do not use underscores in identifiers.
2. Preserve common initialisms: `ID`, `URL`, `HTTP`, `API`, `JSON`, `SQL`.
3. Use `userID`, `httpClient`, and `parseURL`, not `userId`, `httpclient`, or
   `parseUrl`.
4. Use short names for small scopes and descriptive names for larger scopes.
5. Name variables for meaning, not type: `users`, not `userSlice`.
6. Avoid one-letter names outside small loops, receivers, or mathematical code.
7. Avoid redundant context words already supplied by the package or type.
8. Name booleans as predicates: `enabled`, `hasToken`, `isReady`, `canRetry`.
9. Name functions with verbs when they perform work.
10. Name types and values with nouns when they represent data.
11. Do not use `Get` for a simple accessor unless the operation performs a
    meaningful retrieval.
12. Avoid vague verbs such as `Handle`, `Process`, `Manage`, or `Do` when a
    precise verb exists.

## 7. Receiver Names
1. Use a short, consistent receiver derived from the type name.
2. Use the same receiver name for all methods on a type.
3. Do not use `this`, `self`, or `me`.
4. Use a pointer receiver when the method mutates the receiver, contains a
   synchronization primitive, or the type is large.
5. Do not mix pointer and value receivers without a deliberate reason.
6. Do not copy types containing `sync.Mutex` or other no-copy state.

```go
func (s *Service) Create(ctx context.Context, input CreateInput) error {
    // ...
}
```

## 8. Declarations and Zero Values
1. Prefer the zero value when it is valid and meaningful.
2. Use `var` for a zero-value declaration.
3. Use `:=` for a local value with an obvious type.
4. Use an explicit type when inference would hide an important domain choice.
5. Group declarations only when they belong together.
6. Do not group unrelated variables merely to reduce lines.
7. Prefer typed constants for domain values.
8. Use `iota` only when values are internal and their numeric representation is
   not persisted or transmitted.
9. Never rely on an `iota` ordinal in a database, protocol, or public API.
10. Avoid mutable package-level variables.
11. If global state is unavoidable, make ownership, synchronization, and test
    reset behavior explicit.
12. Avoid `init`; prefer explicit construction and setup.

## 9. Structs
1. A struct should represent one cohesive concept.
2. Keep fields unexported unless callers require direct data access.
3. Order fields by conceptual grouping, with synchronization fields first when
   they protect the remaining state.
4. Do not embed a `sync.Mutex`; use a named field such as `mu sync.Mutex`.
5. Do not embed implementation types merely to inherit methods.
6. Use embedding only when the promoted API is intentionally part of the outer
   type.
7. Avoid struct tags that duplicate a default unless the explicit form prevents
   ambiguity.
8. Treat changes to exported fields and serialization tags as API changes.
9. Use pointer fields only when absence, identity, mutation, or size justifies
   them.
10. Do not use pointers to small scalar values only to avoid copying.

## 10. Constructors
1. Provide a constructor when valid initialization requires dependencies,
   validation, defaults, or unexported fields.
2. Return a concrete type unless callers need an abstraction.
3. Validate required dependencies at construction time.
4. Do not create a constructor that only returns `&T{}` for a valid zero-value
   type.
5. Keep constructor work deterministic and free of hidden network or background
   activity.
6. Name constructors `New`, `NewClient`, or another precise package-level name
   without repeating the package.
7. Use functional options only when there are multiple optional parameters that
   are expected to evolve.
8. Functional options must validate input and produce deterministic results.
9. Prefer a configuration struct when options are data that should be decoded,
   validated, logged safely, or compared.

## 11. Functions
1. A function should perform one understandable operation.
2. Prefer a short straight-line happy path.
3. Use guard clauses to handle invalid input and errors early.
4. Avoid deep nesting.
5. Keep parameter lists small.
6. Use a parameter struct when multiple parameters form one domain input or
   when boolean arguments would be unclear.
7. Do not pass a boolean that makes call sites ambiguous.
8. Return values in an order that follows Go convention: result values followed
   by `error`.
9. Use named results only when they improve documentation or a short deferred
   function needs them.
10. Do not use naked returns in non-trivial functions.
11. Avoid functions that mix parsing, validation, persistence, and transport
    concerns.

Avoid:

```go
createUser(ctx, name, true, false)
```

Prefer:

```go
createUser(ctx, CreateUserInput{
    Name:       name,
    SendInvite: true,
})
```

## 12. Control Flow
1. Put the normal path at the lowest indentation level.
2. Return early on invalid state and errors.
3. Do not add `else` after a branch that returns.
4. Use a `switch` when it expresses mutually exclusive states more clearly than
   repeated `if` statements.
5. Do not use `switch` fallthrough unless the behavior is essential and
   commented.
6. Keep loop termination and mutation visible.
7. Avoid hidden control flow through panic, global callbacks, or side effects.
8. Use `defer` immediately after successfully acquiring a resource.
9. Check close or flush errors when they can affect correctness.

```go
file, err := os.Open(name)
if err != nil {
    return fmt.Errorf("open %q: %w", name, err)
}
defer file.Close()
```

## 13. Errors
1. Use errors for expected failure paths.
2. Do not panic for invalid input, I/O failure, dependency failure, or other
   normal runtime conditions.
3. Add context when returning an error across an abstraction boundary.
4. Wrap an underlying error with `%w` when callers may need `errors.Is` or
   `errors.As`.
5. Use `%v` when intentionally hiding the underlying error identity.
6. Error strings start with lowercase unless they begin with a proper noun or
   acronym.
7. Error strings do not end with punctuation.
8. Describe the failed operation and relevant safe identifier.
9. Do not include secrets or sensitive payloads.
10. Handle an error once: return it, transform it, or log it at a terminal
    boundary. Do not log and return the same error at every layer.
11. Use sentinel errors only when callers need stable identity for control flow.
12. Use a custom error type when callers need structured fields.
13. Compare errors with `errors.Is` and extract types with `errors.As`.
14. Do not compare wrapped errors by string.
15. Do not discard an error with `_` unless failure is impossible by contract
    or intentionally irrelevant and documented.

```go
token, err := parser.Parse(raw)
if err != nil {
    return Token{}, fmt.Errorf("parse access token: %w", err)
}
```

## 14. Interfaces
1. Define an interface where it is consumed, not where it is implemented.
2. Keep interfaces small and focused.
3. Prefer one-method interfaces when one capability is sufficient.
4. Accept interfaces and return concrete types by default.
5. Do not create an interface only to mock a concrete type.
6. Do not create an interface before there is more than one meaningful
   implementation or a consumer boundary that benefits from it.
7. Do not use pointers to interfaces.
8. Make interface satisfaction implicit.
9. Add a compile-time assertion only when it documents an important contract:

```go
var _ io.Reader = (*Reader)(nil)
```

10. Avoid interface pollution that exposes every method of a service to every
    consumer.

## 15. Context
1. Pass `context.Context` as the first parameter.
2. Name it `ctx`.
3. Do not store a context in a struct.
4. Do not pass a nil context.
5. Propagate the caller's context through downstream operations.
6. Derive a timeout or cancellation context only when the current operation
   owns that bound.
7. Call the returned cancel function, normally with `defer cancel()`.
8. Do not use context values for optional function parameters or mutable state.
9. Context values are only for request-scoped data crossing API boundaries.
10. Use private key types for context values.
11. Long-running operations must observe cancellation.
12. Preserve `context.Canceled` and `context.DeadlineExceeded` identity when
    wrapping.

## 16. Concurrency
1. Do not start a goroutine without defining who owns it and how it stops.
2. Every goroutine must terminate predictably or accept a cancellation signal.
3. The owner must have a way to wait for completion when shutdown correctness
   requires it.
4. Do not use fire-and-forget goroutines in production code.
5. Bound concurrency with workers, semaphores, queues, or another explicit
   limit.
6. Prefer synchronous code until concurrency provides a measured benefit or is
   required for correctness.
7. Use channels to communicate ownership or events, not as a default
   replacement for ordinary calls.
8. The sender normally owns and closes a channel.
9. Do not close a channel from the receiving side.
10. Do not close a channel more than once.
11. Use directional channel types in APIs when the direction is known.
12. Protect shared mutable state with a mutex or confine it to one goroutine.
13. Keep mutex critical sections small.
14. Do not hold a lock during slow I/O, callbacks, logging, or unknown external
    code.
15. Document which fields a mutex protects.
16. Use `sync.Once`, atomics, or condition variables only when their semantics
    are clearer than a mutex.
17. Run race-enabled tests for changed concurrent code.

## 17. Slices and Maps
1. Prefer a nil slice for an empty internal value.
2. Use a non-nil empty slice when an external contract distinguishes `[]` from
   `null`, such as JSON.
3. Do not make callers depend on nil-versus-empty unless the API explicitly
   documents it.
4. Preallocate slices and maps when the size is known and material.
5. Do not preallocate based on guesses that obscure code without measured value.
6. Copy slices and maps at API boundaries when retaining the caller's value
   would allow unintended mutation.
7. Do not expose internal mutable slices or maps directly.
8. Check map membership with the comma-ok form when zero values are ambiguous.
9. Remember that map iteration order is unspecified.
10. Sort keys before producing deterministic output.
11. Do not mutate a collection while another goroutine accesses it without
    synchronization.

## 18. Strings, Bytes, and Formatting
1. Use strings for immutable text and byte slices for mutable or binary data.
2. Be explicit about encoding at external boundaries.
3. Use `strings.Builder` or `bytes.Buffer` for repeated construction when it
   improves clarity or performance.
4. Use `%q` when quoted output makes user input or whitespace unambiguous.
5. Do not use formatting where direct conversion or concatenation is clearer.
6. Avoid converting repeatedly between `string` and `[]byte` in hot paths.
7. Validate UTF-8 only when the contract requires valid text.

## 19. Time
1. Represent durations with `time.Duration`, not bare integers.
2. Include units in configuration and serialized field names when the format
   cannot carry a duration type.
3. Use `time.Time` for instants.
4. Store and transmit timestamps in an explicitly defined format and timezone.
5. Use UTC for persistence unless the domain requires a local timezone.
6. Inject a clock or current-time function when deterministic tests require it.
7. Do not scatter `time.Now()` through domain logic.
8. Stop tickers and timers when their lifecycle ends.

## 20. Dependencies and Side Effects
1. Prefer the standard library when it provides a clear, maintained solution.
2. Add a dependency only when its benefit exceeds maintenance and supply-chain
   cost.
3. Inject external systems at a clear boundary.
4. Do not hide network, filesystem, environment, or clock access in innocent
   value helpers.
5. Avoid package initialization side effects.
6. Keep dependency direction visible in constructors and fields.
7. Do not introduce a dependency-injection framework when explicit
   constructors are sufficient.

## 21. Logging
1. Log at process, request, worker, or other terminal operational boundaries.
2. Do not log an error at every layer.
3. Use structured fields where the configured logger supports them.
4. Use stable field names.
5. Do not log secrets, credentials, tokens, private payloads, or personal data.
6. Include identifiers needed to correlate the operation, but only when safe.
7. Do not use logging as control flow.
8. Libraries should return errors rather than decide process logging policy.
9. Only `main` or an equivalent top-level boundary may terminate the process.

## 22. Comments and Documentation
1. Every exported package, type, function, method, variable, and constant must
   have a useful doc comment.
2. Start an exported identifier's comment with the identifier name.
3. Explain purpose, contract, constraints, and non-obvious behavior.
4. Do not narrate syntax.
5. Comment why, not what.
6. Keep comments current with code.
7. Use complete sentences for prose comments.
8. Document concurrency safety, ownership, nil behavior, side effects, and
   error contracts when callers need them.
9. Use `TODO(owner-or-issue): reason` for actionable temporary work.
10. Remove stale and resolved TODO comments.

## 23. Testing
1. Test observable behavior rather than implementation sequence.
2. Place tests near the code in `_test.go` files.
3. Use table-driven tests when multiple cases share the same behavior and
   setup.
4. Give each case a descriptive name.
5. Use subtests for independent cases.
6. Use `t.Helper()` in test helpers.
7. Use `t.Cleanup()` for cleanup tied to test lifetime.
8. Call `t.Parallel()` only when the test and shared fixtures are safe for
   parallel execution.
9. Avoid sleep-based synchronization.
10. Use deadlines, channels, fakes, and polling with bounds for asynchronous
    behavior.
11. Prefer simple fakes or stubs over large mocking frameworks.
12. Assert meaningful outcomes and errors.
13. Use `errors.Is` and `errors.As` in error assertions.
14. Test cancellation, timeout, empty input, boundary values, and failure paths
    when relevant.
15. Tests must be deterministic and independent of execution order.
16. Examples should compile and should be tests when they document public use.
17. Benchmarks must avoid setup in the measured section and report allocations
    when useful.

## 24. Performance
1. Prefer clear code until profiling identifies a real bottleneck.
2. Benchmark before and after an optimization.
3. Do not pool cheap objects without evidence.
4. Avoid reflection and unsafe operations when ordinary typed code is
   practical.
5. Keep allocations visible at API boundaries.
6. Preallocate only with reliable size information.
7. Document non-obvious optimizations and the measurement that justifies them.
8. Do not sacrifice correctness or cancellation for throughput.

## 25. Generics
1. Use generics when one algorithm or data structure genuinely applies to
   multiple types.
2. Do not replace a small interface or ordinary function with generics solely
   to remove a few repeated lines.
3. Keep constraints minimal and capability-based.
4. Use standard constraints or `~` type terms only when underlying types are
   intentionally supported.
5. Give type parameters short, meaningful names such as `T`, `K`, `V`, or
   domain-specific names when needed.
6. Avoid generic APIs that make call sites harder to understand than concrete
   APIs.

## 26. Common Prohibitions
Do not:
- use panic for ordinary failures
- ignore returned errors without justification
- use `interface{}` or `any` when a concrete type is known
- add getters and setters mechanically
- store context in a struct
- start unmanaged goroutines
- expose internal mutable collections
- use global mutable state for convenience
- add interfaces only for mocking
- use package names such as `utils` or `common`
- use `init` for dependency wiring
- log and return the same error at every layer
- use a bare integer for a duration
- depend on map iteration order

## 27. Review Checklist
- Does the package have one clear responsibility?
- Are names idiomatic and free of redundant context?
- Is the happy path shallow and easy to follow?
- Are errors contextual, wrapped deliberately, and handled once?
- Are interfaces defined by consumers and kept small?
- Is context propagated and cancellation respected?
- Does every goroutine have ownership and a stop path?
- Are mutable slices and maps protected at boundaries?
- Are exported APIs documented?
- Are tests deterministic and focused on behavior?
- Did formatting, tests, vet, and relevant race checks pass?

## References
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Organizing a Go Module](https://go.dev/doc/modules/layout)
- [Go Context](https://go.dev/blog/context)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
