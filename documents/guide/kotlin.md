# KOTLIN.md
This guide defines Kotlin conventions for this repository.
It follows the official Kotlin coding conventions and library API guidelines,
with additional project rules that favor direct control flow, explicit
ownership, safe null handling, and APIs that are easy to understand at the call
site.

## 1. Goals and Priority
Write Kotlin code that is:
1. correct
2. idiomatic
3. explicit about nullability, mutation, failure, and concurrency
4. easy to read without language tricks
5. easy to test
6. safe at platform and process boundaries
7. small enough to change confidently

Apply rules in this order:
1. Kotlin language semantics
2. official Kotlin coding conventions
3. official Kotlin API and coroutine guidelines
4. this guide
5. consistent local module conventions

Do not preserve a local pattern that is unsafe or contradicts a higher-priority
rule.

## 2. Formatting and Tooling
1. Use the repository-configured Kotlin formatter and linter.
2. Follow the official Kotlin style based on the IntelliJ IDEA Kotlin style.
3. Do not manually align declarations, assignments, arguments, or comments.
4. Use four spaces for indentation.
5. Do not use tabs.
6. Do not use semicolons.
7. Use trailing commas in multiline declarations and calls when supported by
   the configured Kotlin version and formatter.
8. Put one statement per line.
9. Keep lines readable; let the formatter control wrapping.
10. Do not suppress formatter or linter rules without a documented technical
    reason.
11. Keep formatting-only changes separate from unrelated behavioral changes.

Run the repository's configured equivalents of:

```sh
<formatter> --check
<linter>
<build-tool> test
```

## 3. Source Layout
1. Follow the build system's standard source-set layout.
2. Keep production and test source sets separate.
3. Match directory structure to package structure.
4. Do not add redundant directory segments that repeat the module name.
5. Keep platform-specific code in the corresponding platform source set.
6. Keep common code free of platform APIs.
7. Use `internal` implementation packages or declarations when an API is not
   intended outside its module.
8. Do not create packages named `util`, `utils`, `common`, `base`, `helper`, or
   `misc` when a domain-specific package name is available.
9. Give each package one clear responsibility.
10. Avoid cyclic dependencies by correcting responsibility boundaries.

## 4. Packages and Imports
1. Use lowercase package names without underscores.
2. Use meaningful domain segments.
3. Do not encode organization charts or temporary architecture names in package
   paths.
4. Use explicit imports.
5. Do not use wildcard imports.
6. Use import aliases only to resolve a collision or clarify a misleading
   external name.
7. Let the formatter or IDE sort imports.
8. Remove unused imports immediately.
9. Do not import a broad facade when a focused dependency makes ownership
   clearer.

Good:

```kotlin
package com.example.auth.token
```

Avoid:

```kotlin
package com.example.common.utils
```

## 5. File Names and Contents
1. Name a file after its primary class, interface, object, or responsibility.
2. Use `PascalCase.kt` for files centered on a type.
3. Use a descriptive `PascalCase.kt` name for files containing related top-level
   declarations.
4. Do not add suffixes such as `Impl`, `Util`, `Helper`, or `Manager` without a
   precise domain meaning.
5. Keep closely related declarations in the same file when that improves
   navigation.
6. Do not require one file per type for small private declarations.
7. Split a file when it contains independent responsibilities.
8. Keep tests in files named after the subject with a consistent `Test` suffix.
9. Mark generated files with the generator's standard header.
10. Do not hand-edit generated files.

## 6. File Organization
Order a source file as follows:
1. copyright or file annotations when required
2. package declaration
3. imports
4. top-level declarations

Inside a class, prefer this order:
1. property declarations and initializer blocks
2. secondary constructors
3. public functions
4. internal functions
5. protected functions
6. private functions
7. nested types
8. companion object

Keep related state and behavior close together. A strict mechanical order may
be relaxed when separating a public function from its small private helper
makes the code harder to follow.

## 7. Naming
1. Use `PascalCase` for classes, interfaces, objects, type aliases, and
   annotations.
2. Use `camelCase` for functions, properties, parameters, and local variables.
3. Use `UPPER_SNAKE_CASE` only for true constants.
4. Preserve conventional acronyms as words in mixed-case names:
   `userId`, `parseUrl`, `HttpClient`.
5. Preserve externally defined spelling at protocol boundaries.
6. Use descriptive nouns for values and precise verbs for operations.
7. Name booleans as predicates: `isReady`, `hasToken`, `canRetry`, `enabled`.
8. Do not repeat package or enclosing-type context in a name.
9. Avoid vague words such as `data`, `item`, `object`, `thing`, `process`,
   `handle`, `manager`, and `helper` when a domain term exists.
10. Use short names only in small, obvious scopes.
11. Include units where the type does not express them:
    `timeoutMillis`, `sizeBytes`.
12. Do not encode type names in variable names:
    `users`, not `userList`.

## 8. Properties and Variables
1. Use `val` by default.
2. Use `var` only when reassignment is required by the model.
3. Keep mutable state private.
4. Expose read-only views or behavior instead of mutable collections.
5. Declare variables close to first use.
6. Do not reuse a variable for a different meaning.
7. Prefer a direct expression when an intermediate variable adds no meaning.
8. Introduce a named value when it explains a domain step or complex condition.
9. Avoid mutable top-level properties.
10. Do not use public mutable properties as an informal API.
11. Give a property a restricted setter when reads may be public but mutation
    must remain controlled.
12. Use delegated properties only when the delegate's lifecycle and side
    effects are clear.

```kotlin
class Account {
    var status: AccountStatus = AccountStatus.PENDING
        private set
}
```

## 9. Type Inference
1. Use inference for obvious local values.
2. Specify types at public API boundaries.
3. Specify a return type for public and protected functions.
4. Specify a type when inference would hide nullability, numeric width, generic
   variance, or an important domain choice.
5. Do not annotate an obvious local literal with the same primitive type.
6. Do not use inference to expose an anonymous or implementation-specific type.
7. Keep inferred expressions simple enough that a reader can identify the type
   without IDE assistance.

## 10. Visibility
1. Use the narrowest visibility that supports intended callers.
2. Remember that Kotlin declarations are `public` by default.
3. Mark implementation declarations `private` or `internal` deliberately.
4. Use `internal` for module APIs that must not become public contracts.
5. Use `protected` only for a deliberate inheritance extension point.
6. Avoid public setters.
7. Treat widening visibility as an API change.
8. Treat changes to public constructors, properties, functions, default
   arguments, and sealed hierarchies as compatibility-sensitive.
9. Enable explicit API mode for published libraries when supported by the build
   configuration.
10. Do not expose third-party implementation types unnecessarily in public APIs.

## 11. Null Safety
1. Use a non-nullable type unless absence is part of the domain.
2. Use `T?` only when `null` has one clear meaning.
3. Do not use `null` for multiple states such as missing, invalid, forbidden,
   and failed.
4. Convert platform types and untrusted nullable values at the boundary.
5. Prefer early returns or explicit branches for required values.
6. Use safe calls when absence should propagate.
7. Use the Elvis operator for a clear fallback or early exit.
8. Do not chain safe calls through state that should be required.
9. Avoid `!!`.
10. `!!` is allowed only when an invariant is guaranteed outside Kotlin's type
    system and the reason is documented at the assertion.
11. Prefer `requireNotNull`, `checkNotNull`, or a domain-specific error when the
    failed invariant should explain itself.
12. Do not use `lateinit` merely to avoid nullable types.
13. Use `lateinit` only when a framework or lifecycle guarantees initialization
    before use and constructor initialization is impractical.
14. Prefer constructor injection over `lateinit` dependency injection.

Avoid:

```kotlin
val userName = request.user!!.profile!!.name
```

Prefer:

```kotlin
val user = requireNotNull(request.user) { "request user is required" }
val profile = requireNotNull(user.profile) { "user profile is required" }
val userName = profile.name
```

## 12. Smart Casts and Casts
1. Use smart casts after explicit type or null checks.
2. Prefer sealed hierarchies and polymorphism over repeated casts.
3. Use `as?` when a failed cast is an expected outcome.
4. Use `as` only when failure indicates a programming or contract error.
5. Do not use unchecked casts to silence generic design problems.
6. Keep an unavoidable unchecked cast at a narrow boundary.
7. Add `@Suppress("UNCHECKED_CAST")` only to the smallest declaration and
   explain the safety invariant.
8. Do not branch on runtime type when a domain discriminator would be clearer.

## 13. Domain Modeling
1. Use types to make invalid states difficult to construct.
2. Use a data class for transparent value-oriented data.
3. Use a regular class for identity, invariants, controlled mutation, or
   lifecycle.
4. Use a value class for a domain value when it prevents meaningful accidental
   mixing and its runtime limitations are understood.
5. Do not wrap every primitive in a value class without a demonstrated
   correctness benefit.
6. Use a sealed class or sealed interface for a closed set of variants.
7. Keep variant-specific data on the corresponding subtype.
8. Use an exhaustive `when` expression for closed variants.
9. Do not model mutually exclusive states with independent booleans.
10. Use an enum for a simple finite set with no variant-specific state.
11. Give externally serialized enum values an explicit stable representation at
    the boundary.
12. Do not persist or transmit enum ordinals.

Avoid:

```kotlin
data class RequestState(
    val isLoading: Boolean,
    val isSuccessful: Boolean,
    val result: User?,
    val error: Throwable?,
)
```

Prefer:

```kotlin
sealed interface RequestState {
    data object Idle : RequestState
    data object Loading : RequestState
    data class Success(val user: User) : RequestState
    data class Failure(val cause: Throwable) : RequestState
}
```

## 14. Data Classes
1. Use a data class when generated structural equality, `copy`, and
   destructuring match the domain semantics.
2. Keep primary constructor properties as the complete structural identity.
3. Do not put significant equality state only in the class body.
4. Prefer immutable `val` properties.
5. Do not use a data class for an entity whose identity is independent of all
   field values.
6. Be cautious with arrays and mutable collections because generated equality
   and copying may not match expectations.
7. Remember that `copy` is shallow.
8. Do not expose mutable collection properties from a data class unless shared
   mutation is intentional.
9. Avoid positional destructuring of domain objects with several properties;
   named property access is clearer and safer under reordering.

## 15. Classes and Constructors
1. Give each class one clear responsibility.
2. Prefer a primary constructor.
3. Make required dependencies constructor parameters.
4. Keep constructors deterministic and free of network, filesystem, coroutine,
   or process-starting side effects.
5. Validate class invariants during construction.
6. Use `require` for invalid caller arguments.
7. Use `check` for invalid object or application state.
8. Prefer factory functions when construction needs a descriptive name,
   validation result, caching, or subtype selection.
9. Use secondary constructors mainly for interoperability or genuinely distinct
   input forms.
10. Delegate secondary constructors to the primary constructor.
11. Do not add an empty companion object only to imitate static methods.
12. Prefer top-level functions for operations that do not need class state.

```kotlin
class RetryPolicy(
    val maxAttempts: Int,
) {
    init {
        require(maxAttempts > 0) { "maxAttempts must be positive" }
    }
}
```

## 16. Functions
1. A function should perform one understandable operation.
2. Name a function for its observable effect or returned value.
3. Prefer a short, straight-line happy path.
4. Use guard clauses for invalid state and failures.
5. Keep parameter lists small.
6. Use a parameter object when several values form one domain input.
7. Avoid boolean mode parameters.
8. Prefer separate functions when modes have meaningfully different behavior.
9. Use default arguments for stable, unsurprising defaults.
10. Do not add overloads that only reproduce default arguments for Kotlin
    callers.
11. Use named arguments when adjacent values have the same type, a boolean is
    passed, or the call is otherwise ambiguous.
12. Do not use named arguments as a substitute for an oversized parameter list.
13. Use expression bodies only for short, obvious expressions.
14. Use block bodies when control flow, validation, logging, or error handling
    is present.
15. Do not use `Unit` as a generic signal when a domain result would be clearer.
16. Keep side effects visible in function names and module boundaries.

Avoid:

```kotlin
createUser(name, true, false)
```

Prefer:

```kotlin
createUser(
    name = name,
    sendInvite = true,
    requirePasswordChange = false,
)
```

Prefer a parameter type when the options grow:

```kotlin
data class CreateUserCommand(
    val name: String,
    val sendInvite: Boolean,
    val requirePasswordChange: Boolean,
)
```

## 17. Return Values
1. Return the narrowest useful type.
2. Return read-only collection interfaces by default.
3. Do not expose mutable internal collections.
4. Use a nullable return only when absence is an expected, unambiguous outcome.
5. Use a sealed result type when callers must distinguish several expected
   outcomes.
6. Throw an exception for an unexpected inability to complete an operation.
7. Do not return `null`, `false`, or an empty collection interchangeably for
   failure.
8. Do not use `Pair` or `Triple` for domain-rich public return values.
9. Use a named data class when fields require meaning.

## 18. Higher-Order Functions and Lambdas
1. Use a lambda when behavior is short and local.
2. Give lambda parameters meaningful names when `it` is not immediately clear.
3. Use `it` only for a short lambda with one obvious parameter.
4. Do not nest lambdas that all rely on implicit receivers or `it`.
5. Extract a named function when a lambda contains branching, error handling, or
   several operations.
6. Keep non-local returns visible and intentional.
7. Avoid labeled returns when a named function or ordinary loop is clearer.
8. Do not use higher-order functions to hide important control flow.
9. Use function references when they are clearer than a forwarding lambda.

## 19. Scope Functions
Choose a scope function for a clear semantic reason:
1. `let` transforms a value or scopes a nullable value.
2. `run` computes a result with an object as receiver.
3. `apply` configures an object and returns that object.
4. `also` performs an additional side effect and returns the object.
5. `with` groups calls on a non-null object when no chaining is needed.

Rules:
1. Do not chain several scope functions when ordinary variables are clearer.
2. Do not nest implicit receivers.
3. Avoid using `let` only to rename a non-null local value.
4. Do not use `also` for essential state mutation that should be visible as a
   normal statement.
5. Prefer explicit code when the reader must remember whether the receiver is
   `this` or `it`.

Avoid:

```kotlin
user?.let {
    repository.find(it.id)?.also {
        audit.record(it)
    }?.run {
        activate()
    }
}
```

Prefer named steps and guard clauses.

## 20. Extension Functions
1. Add an extension when it expresses a natural operation on the receiver.
2. Keep extensions in the package of the owning feature or consuming boundary.
3. Do not use extensions to simulate adding state or privileged access.
4. Keep extension behavior unsurprising from its name.
5. Avoid broad extensions on `Any`, `String`, collections, or other ubiquitous
   types unless the domain meaning is unmistakable.
6. Do not hide I/O or expensive work behind a property-like extension.
7. Use an extension property only for a cheap, deterministic derived value.
8. Keep nullable-receiver extensions rare and make null behavior explicit.
9. Prefer a regular function when multiple arguments are equally central.
10. Do not create competing extensions with the same name in commonly imported
    packages.

## 21. Inheritance, Interfaces, and Delegation
1. Prefer composition over class inheritance.
2. Use an interface for a stable capability or consumer boundary.
3. Keep interfaces focused.
4. Do not create an interface only to mock one implementation.
5. Define an interface where the abstraction is owned, normally near its
   consumer or domain boundary.
6. Use an abstract class only when implementations share protected state or
   invariant-preserving behavior.
7. Keep classes final by default.
8. Mark a class or member `open` only as a deliberate extension point.
9. Document behavioral contracts for overridable members.
10. Do not call open members from constructors or initializer blocks.
11. Use Kotlin delegation when it makes ownership and forwarded behavior
    clearer.
12. Do not delegate a broad interface when the outer type should expose only a
    subset.

## 22. Objects and Companion Objects
1. Use an `object` for a stateless singleton or one intentional process-wide
   identity.
2. Do not use objects as containers for unrelated functions.
3. Avoid mutable singleton state.
4. Use a companion object for constants, factories, or functionality that is
   conceptually tied to the class.
5. Prefer a named factory function over direct construction when the name
   explains validation or representation.
6. Use `const val` only for compile-time constants of supported types.
7. Use a top-level constant when it is not conceptually owned by a class.

## 23. Collections and Sequences
1. Expose `List`, `Set`, and `Map` unless callers must mutate the collection.
2. Use mutable collections only inside a clearly owned mutation boundary.
3. Do not return a mutable internal collection.
4. Copy at a boundary when retaining or exposing a caller-owned mutable
   collection would permit unintended changes.
5. Use collection operations when they directly express the transformation.
6. Use a loop when early exit, stateful accumulation, or performance-sensitive
   mutation is clearer.
7. Avoid long chains of `map`, `filter`, `flatMap`, and `associate` that obscure
   business steps or create unnecessary intermediate collections.
8. Introduce named intermediate values for meaningful stages.
9. Use `Sequence` for lazy multi-stage processing only when the input size or
   measured allocation cost justifies it.
10. Do not use `Sequence` automatically for small collections.
11. Do not rely on unspecified iteration order.
12. Sort explicitly before deterministic serialization, hashing, or tests.
13. Use `firstOrNull`, `singleOrNull`, or another operation whose failure
    semantics match the contract.
14. Do not catch collection lookup exceptions to implement normal absence.

## 24. Strings and Templates
1. Use string templates for simple interpolation.
2. Use `${expression}` for expressions and `$name` for simple names.
3. Avoid complex logic inside a string template.
4. Use a named intermediate value when formatting logic is non-trivial.
5. Use multiline strings for genuinely multiline content.
6. Apply `trimIndent` or `trimMargin` deliberately.
7. Do not build SQL, shell commands, or other injection-sensitive syntax by
   interpolating untrusted values.
8. Keep user-facing text separate from internal exception and log messages when
   localization is required.

## 25. Numbers, Time, and Units
1. Use the type that matches the domain range and external contract.
2. Do not rely on implicit numeric widening; Kotlin does not perform it.
3. Avoid magic numbers.
4. Give domain constants meaningful names.
5. Represent durations with the project's selected duration type, not bare
   integers.
6. Include a unit in names when the type cannot express it.
7. Use an injected clock or time source when deterministic tests require it.
8. Do not scatter direct current-time calls through domain logic.
9. Define timezone and serialization format at external boundaries.
10. Do not use floating-point numbers for exact monetary values.

## 26. Exceptions and Validation
1. Use exceptions for unexpected failure to complete an operation.
2. Use `require` and `requireNotNull` for invalid caller arguments.
3. Use `check` and `checkNotNull` for invalid object or application state.
4. Throw a domain-specific exception only when callers need stable
   classification or structured context.
5. Include the failed operation and safe identifiers in the message.
6. Preserve the original exception as the cause when adding context.
7. Do not catch `Throwable` in ordinary application code.
8. Catch the narrowest exception that can be handled.
9. Do not catch only to log and rethrow unchanged at every layer.
10. Do not use exceptions for expected branching when a nullable result or
    sealed outcome is clearer.
11. Do not expose secrets, credentials, tokens, or private payloads in exception
    messages.
12. Do not ignore an exception with an empty catch block.
13. Cleanup belongs in `finally`, `use`, or a lifecycle-specific cleanup API.
14. Use `runCatching` only when its `Result` value is handled explicitly and
    cancellation semantics remain correct.
15. Do not wrap large blocks in `runCatching` merely to avoid `try` and `catch`.

```kotlin
val configuration = try {
    parser.parse(rawConfiguration)
} catch (exception: ParseException) {
    throw ConfigurationException(
        message = "failed to parse service configuration",
        cause = exception,
    )
}
```

## 27. Resource Management
1. Acquire resources as late as practical and release them deterministically.
2. Use `use` for `Closeable` or `AutoCloseable` resources where supported.
3. Keep the resource scope small.
4. Handle flush or close failures when they affect correctness.
5. Do not rely on finalizers or garbage collection for timely cleanup.
6. Do not start background work in a resource constructor.
7. Make ownership transfer explicit when a function returns an open resource.

## 28. Coroutines
1. Use structured concurrency.
2. Launch coroutines only in a scope with a clear owner and lifecycle.
3. Prefer a suspending function over a function that launches hidden work.
4. Do not use `GlobalScope`.
5. Do not create an unmanaged `CoroutineScope` only to launch fire-and-forget
   work.
6. Every launched coroutine must have a predictable completion or cancellation
   path.
7. The owner must be able to cancel and, when necessary, join its work.
8. Propagate cancellation.
9. Do not catch and swallow `CancellationException`.
10. Rethrow `CancellationException` when a broad boundary catches exceptions.
11. Use `coroutineScope` when child failure should fail siblings and the parent.
12. Use `supervisorScope` only when child failures are intentionally isolated.
13. Document that isolation and handle every child failure.
14. Keep coroutine builders at orchestration boundaries.
15. Keep domain transformations as ordinary or suspending functions.
16. Do not use `runBlocking` in production suspend paths.
17. Use `runBlocking` only at an allowed synchronous boundary such as a process
    entry point or compatible test.

## 29. Dispatchers and Blocking Work
1. Do not hard-code a dispatcher in domain logic when callers should control
   execution.
2. Switch context at the boundary that knows work is blocking or
   CPU-intensive.
3. Use the configured I/O dispatcher for blocking I/O.
4. Do not wrap an already non-blocking suspending API in an I/O context without
   a reason.
5. Inject dispatchers or an execution abstraction when deterministic tests or
   platform variation require it.
6. Do not assume that `suspend` means a function is non-blocking.
7. Keep blocking calls out of event-loop or main-thread contexts.
8. Limit parallelism for expensive or externally constrained operations.

## 30. Async Composition
1. Execute suspending operations sequentially by default.
2. Use `async` only when operations are independent and concurrency provides a
   real benefit.
3. Start `async` in an enclosing structured scope.
4. Await every deferred value.
5. Do not expose `Deferred` from an API when a suspending function communicates
   the contract more clearly.
6. Do not use `async` as an exception wrapper.
7. Use bounded concurrency for collections.
8. Do not create one coroutine per unbounded input element.
9. Define timeout ownership at the operation boundary.
10. Use timeout APIs only when timeout is part of the caller-visible contract.

## 31. Flow
1. Use `Flow` for an asynchronous stream of multiple values.
2. Use a suspending function for one asynchronous result.
3. Keep flows cold unless shared state is a deliberate part of the API.
4. Make the owner and lifetime of `StateFlow` and `SharedFlow` explicit.
5. Expose read-only flow types from mutable internal flows.
6. Do not expose `MutableStateFlow` or `MutableSharedFlow`.
7. Keep context changes and buffering deliberate.
8. Do not apply `flowOn`, `buffer`, `conflate`, or retry operators without
   understanding ordering, cancellation, and backpressure effects.
9. Preserve exception transparency.
10. Catch only upstream exceptions that can be handled.
11. Do not swallow cancellation in `catch`.
12. Collect a flow only in a lifecycle-owned scope.
13. Avoid deeply chained operators when named transformation functions clarify
    business stages.

## 32. Shared Mutable State
1. Prefer immutable values and isolated ownership.
2. Confine mutable state to one coroutine when practical.
3. Use a mutex, atomic primitive, actor-like owner, or platform synchronization
   mechanism when state is shared.
4. Document which state a lock protects.
5. Keep critical sections small.
6. Do not hold a lock across blocking I/O, unknown callbacks, or arbitrary
   suspending calls.
7. Do not assume a collection is thread-safe because its reference is a `val`.
8. Test concurrent behavior under repeated and cancellation-heavy execution.

## 33. Logging
1. Log at process, request, job, or other operational boundaries.
2. Return or propagate failures from reusable code.
3. Do not log and rethrow the same exception at every layer.
4. Use structured fields when supported by the configured logger.
5. Use stable field names.
6. Do not log secrets, credentials, tokens, private payloads, or personal data.
7. Include only safe identifiers required for correlation.
8. Do not use logging as control flow.
9. Do not terminate the process from reusable library or domain code.

## 34. Java Interoperability
Apply this section when Kotlin code is consumed by or consumes Java.

1. Treat Java platform types as untrusted nullability.
2. Normalize platform values into explicit Kotlin nullable or non-nullable types
   at the boundary.
3. Prefer Java declarations with recognized nullability annotations.
4. Design Java-facing Kotlin APIs deliberately.
5. Use `@JvmOverloads` only when Java callers require overloads.
6. Use `@JvmStatic`, `@JvmField`, and `@JvmName` only for a specific
   interoperability need.
7. Do not add JVM annotations by habit.
8. Avoid exposing Kotlin-specific function types, unsigned types, inline
   classes, default arguments, or suspend internals directly to Java when they
   create an unusable API.
9. Provide a Java-friendly facade when required.
10. Consider checked-exception expectations and `@Throws` at foreign-language
    boundaries.
11. Do not use Kotlin keywords or awkward synthetic names in Java-visible APIs.
12. Validate binary and source compatibility for published JVM APIs.

## 35. Multiplatform Code
Apply this section only to Kotlin Multiplatform modules.

1. Keep common source sets platform-independent.
2. Put platform implementations in the narrowest relevant source set.
3. Use `expect` and `actual` only for a real platform capability boundary.
4. Prefer a common interface with platform-provided implementations when it is
   simpler than `expect` and `actual`.
5. Keep expected APIs small and stable.
6. Do not leak one platform's type into common code.
7. Test common behavior in common tests and platform integration in platform
   tests.
8. Treat changes to expected declarations as public compatibility changes.

## 36. Annotations and Suppressions
1. Use annotations only when they change a required contract, tool behavior, or
   interoperability boundary.
2. Place an annotation at the narrowest effective target.
3. Do not suppress a warning before understanding it.
4. Use `@Suppress` on the smallest declaration.
5. Include a comment when the safety or compatibility reason is not obvious.
6. Remove obsolete suppressions.
7. Use opt-in annotations only after reviewing the experimental API's stability
   and migration cost.
8. Do not expose experimental dependencies through a stable public API without
   an explicit decision.

## 37. Comments and KDoc
1. Explain why, constraints, invariants, and non-obvious tradeoffs.
2. Do not narrate syntax.
3. Keep comments synchronized with code.
4. Document public APIs when their contract is not fully obvious from names and
   types.
5. Document null behavior, mutation, thread safety, cancellation, ordering,
   side effects, and thrown exceptions when callers need that information.
6. Use KDoc links for referenced declarations.
7. Avoid redundant `@param` and `@return` text that repeats the signature.
8. Use examples for APIs whose correct composition is not obvious.
9. Use `TODO(owner-or-issue): reason` for actionable temporary work.
10. Remove stale TODOs and commented-out code.
11. Simplify a complicated API before trying to explain it with extensive
    documentation.

## 38. Testing
1. Test observable behavior, not private implementation sequence.
2. Name tests for the condition and expected outcome.
3. Follow the test framework's established naming convention consistently.
4. Use Arrange, Act, Assert when it improves readability.
5. Do not add section comments to trivial tests.
6. Prefer focused fakes or in-memory implementations over broad mocks.
7. Mock external boundaries, clocks, nondeterministic sources, and expensive
   systems.
8. Do not mock every internal collaborator.
9. Test success, expected failure, nullability, boundaries, and invalid
   invariants where relevant.
10. Test sealed variants exhaustively when behavior differs by variant.
11. Keep tests deterministic and independent of execution order.
12. Avoid real sleeps, wall-clock timing, and external network calls in unit
    tests.
13. Restore global state and close resources after tests.
14. Use coroutine test utilities and a test scheduler for coroutine code.
15. Test cancellation, timeout, child failure, and cleanup for concurrent code.
16. Do not use `runBlocking` when the coroutine test framework provides a
    virtual-time test scope.
17. Keep test data builders domain-specific and simple.
18. Do not weaken production visibility or null safety only for tests.

## 39. Performance
1. Prefer clear code until measurement identifies a bottleneck.
2. Benchmark before and after an optimization.
3. Avoid unnecessary allocation in measured hot paths.
4. Do not replace readable collection code with manual mutation without
   evidence.
5. Use sequences only when laziness provides a demonstrated benefit.
6. Avoid reflection when ordinary typed code is practical.
7. Keep boxing implications in mind for nullable primitives, generics, and
   value classes in hot paths.
8. Document non-obvious optimizations and their measurement.
9. Do not sacrifice cancellation, correctness, or API clarity for speculative
   performance.

## 40. Common Prohibitions
Do not:
- use `!!` as routine null handling
- use `lateinit` to avoid constructor design
- expose mutable collections or mutable flows
- use wildcard imports
- create `Util`, `Helper`, `Manager`, or `Base` abstractions without a precise
  responsibility
- use inheritance when composition is clearer
- create interfaces only for mocking
- hide essential work in a scope-function chain
- use `Pair` or `Triple` for domain-rich public values
- use `GlobalScope`
- launch unmanaged coroutines
- swallow `CancellationException`
- use `runBlocking` inside suspend code
- expose `Deferred` for a single asynchronous result without a reason
- catch `Throwable` for ordinary error handling
- use an empty catch block
- log and rethrow the same failure at every layer
- add broad warning suppressions
- use experimental APIs without an explicit opt-in decision

## 41. Validation
Run the smallest repository commands that cover the changed scope:

```sh
<formatter> --check
<linter>
<build-tool> test <affected-target>
```

Also consider:
1. compilation for every affected target or source set
2. API compatibility checks for public libraries
3. coroutine and concurrency tests
4. Java consumer compilation for Java-facing APIs
5. all relevant platform tests for multiplatform changes
6. full module validation when compiler settings or shared APIs change

Do not claim validation passed unless the command was actually run
successfully.

## 42. Review Checklist
- Does each package, file, class, and function have one clear responsibility?
- Are names based on domain meaning rather than generic roles?
- Is state immutable or narrowly owned?
- Is visibility as narrow as possible?
- Does every nullable value have one clear meaning?
- Are `!!`, `lateinit`, casts, and suppressions absent or narrowly justified?
- Are domain states modeled with types rather than conflicting booleans?
- Is the happy path shallow and direct?
- Are scope functions and collection chains easy to follow?
- Are expected outcomes distinct from exceptional failures?
- Does every coroutine belong to an owned structured scope?
- Is cancellation preserved?
- Are public APIs explicit and interoperability-safe?
- Are tests deterministic and focused on observable behavior?
- Did formatting, linting, compilation, and affected tests pass?

## References
- [Kotlin Coding Conventions](https://kotlinlang.org/docs/coding-conventions.html)
- [Kotlin Library API Guidelines](https://kotlinlang.org/docs/api-guidelines-introduction.html)
- [Minimizing Mental Complexity](https://kotlinlang.org/docs/api-guidelines-minimizing-mental-complexity.html)
- [Kotlin Null Safety](https://kotlinlang.org/docs/null-safety.html)
- [Kotlin Exceptions](https://kotlinlang.org/docs/exceptions.html)
- [Kotlin Scope Functions](https://kotlinlang.org/docs/scope-functions.html)
- [Kotlin Coroutines Guide](https://kotlinlang.org/docs/coroutines-guide.html)
- [Kotlin Coroutine Exception Handling](https://kotlinlang.org/docs/exception-handling.html)
- [Calling Java from Kotlin](https://kotlinlang.org/docs/java-interop.html)
- [Calling Kotlin from Java](https://kotlinlang.org/docs/java-to-kotlin-interop.html)
- [KDoc](https://kotlinlang.org/docs/kotlin-doc.html)
