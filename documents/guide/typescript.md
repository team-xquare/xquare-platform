# TYPESCRIPT.md
This guide defines TypeScript conventions for this repository.
It follows the official TypeScript Handbook and TSConfig reference, with
project rules that favor explicit domain models, narrow public APIs, and code
whose runtime behavior is visible from the source.

## 1. Goals
TypeScript code must:
1. preserve JavaScript runtime correctness
2. make invalid states difficult to represent
3. keep control flow direct
4. use inference without hiding contracts
5. validate untrusted runtime data
6. avoid abstractions that are harder to read than the repeated code
7. remain understandable without advanced type-system tricks

TypeScript types are erased at runtime. A type annotation never replaces
runtime validation, authorization, bounds checks, or error handling.

## 2. Compiler Baseline
1. Enable `strict`.
2. Keep code compatible with all enabled strict checks.
3. Enable additional correctness flags for new projects unless the runtime or
   framework has a documented incompatibility:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noImplicitOverride": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "useUnknownInCatchVariables": true,
    "allowUnreachableCode": false,
    "allowUnusedLabels": false
  }
}
```

4. Do not disable strictness for an entire project to fix a local issue.
5. Use the narrowest local suppression only when a third-party or generated
   boundary cannot be corrected.
6. Every suppression must explain why it is safe and when it can be removed.
7. Match `module`, `moduleResolution`, `target`, and library settings to the
   actual runtime and build tool.
8. Do not copy TSConfig settings from another runtime without understanding
   their emitted JavaScript and resolution behavior.

## 3. Formatting
1. Use the repository formatter as the only formatting authority.
2. Do not manually align assignments, object properties, or imports.
3. Keep one statement per line.
4. Use braces for multiline control-flow bodies.
5. Prefer trailing commas in multiline constructs when supported by the
   formatter.
6. Do not mix formatter-only changes with behavioral changes across unrelated
   files.
7. Do not disable formatter rules without a documented technical reason.

## 4. File Names and Organization
1. Use lowercase kebab-case file names: `token-validator.ts`.
2. Use conventional suffixes consistently:
   - `.test.ts` or `.spec.ts` for tests, according to the repository's chosen
     test runner convention
   - `.d.ts` for declarations
   - `.generated.ts` for generated source when generated files are committed
3. Name a file for its primary responsibility.
4. Avoid vague files such as `utils.ts`, `helpers.ts`, `common.ts`, or
   `types.ts` when a domain-specific name is available.
5. Keep implementation, tests, and focused types near their owning feature.
6. Split files by responsibility, not by an arbitrary one-type-per-file rule.
7. Do not create deep directory hierarchies that only repeat namespace words.
8. Keep public package entry points small and deliberate.

## 5. Modules
1. Every source file must be a module.
2. Use ECMAScript `import` and `export` syntax.
3. Do not add values to the global scope.
4. Avoid TypeScript `namespace` in module-based application code.
5. Prefer named exports because import sites remain searchable and rename-safe.
6. Use a default export only when a framework, generated contract, or
   single-entry plugin convention requires it.
7. Export only symbols intended for other modules.
8. Keep implementation details unexported.
9. Avoid wildcard re-export chains.
10. Use barrel files only at a deliberate package or feature boundary.
11. Do not use an internal barrel when direct imports make dependencies clearer
    or avoid cycles.
12. Use `import type` for type-only dependencies when required by compiler
    settings or when it clarifies runtime imports.
13. Avoid side-effect imports except for explicit startup registration,
    polyfills, or styles required by the runtime.
14. Document required side-effect imports.

## 6. Import Ordering
Use one blank line between these groups:
1. intentional side-effect imports
2. platform or third-party imports
3. repository-absolute imports
4. relative imports

Within a group, let the formatter or linter sort consistently.
Do not use deep relative paths when an established repository alias expresses
the same boundary more clearly.
Do not introduce an alias only to shorten a single import.

## 7. Naming
1. Use `PascalCase` for classes, interfaces, type aliases, and enums when an
   enum is justified.
2. Use `camelCase` for functions, methods, variables, and properties.
3. Use `UPPER_SNAKE_CASE` for true module-level constants only.
4. Use descriptive nouns for data and precise verbs for operations.
5. Name booleans as predicates: `isReady`, `hasToken`, `canRetry`, `enabled`.
6. Use `Id`, `Url`, `Http`, and `Json` inside ordinary mixed-case names:
   `userId`, `parseUrl`, `HttpClient`.
7. Preserve externally defined names at protocol boundaries.
8. Do not prefix interfaces with `I`.
9. Do not prefix private fields with `_`.
10. Avoid type information in value names: `users`, not `userArray`.
11. Avoid vague names such as `data`, `item`, `object`, `result`, `manager`,
    `handler`, or `processor` when a domain name is available.
12. Short names are acceptable only in very small scopes with obvious meaning.
13. Name units explicitly when the type does not express them:
    `timeoutMs`, `sizeBytes`.

## 8. Variables
1. Use `const` by default.
2. Use `let` only when reassignment is required.
3. Never use `var`.
4. Declare variables close to first use.
5. Do not reuse one variable for different meanings.
6. Prefer a direct expression over a temporary that adds no meaning.
7. Introduce a named intermediate value when it explains a domain step or
   shortens a complex condition.
8. Avoid module-level mutable state.
9. Make lifecycle and ownership explicit when shared state is unavoidable.

## 9. Type Inference
1. Let TypeScript infer obvious local types.
2. Write explicit types at public boundaries, exported values, object
   construction boundaries, callbacks with subtle contracts, and places where
   widening would be harmful.
3. Prefer an explicit function return type for exported functions.
4. Do not annotate a literal with the same obvious primitive type.
5. Use `satisfies` when validating an expression against a type while
   preserving its inferred narrow type.
6. Use `as const` for intentional literal immutability and narrow literal
   inference.
7. Do not use type annotations to conceal an unsafe value from inference.

```ts
const retryPolicy = {
  attempts: 3,
  mode: "exponential",
} satisfies RetryPolicy;
```

## 10. `any`, `unknown`, and Assertions
1. Do not use `any` in application code unless integrating an untyped boundary
   that cannot be represented otherwise.
2. Prefer `unknown` for values whose type is not yet known.
3. Narrow `unknown` before use.
4. Do not use the global `Function`, `Object`, `String`, `Number`, `Boolean`, or
   `Symbol` types.
5. Use `object`, primitive types, or a precise function signature instead.
6. Treat a type assertion as a claim that requires evidence.
7. Prefer narrowing, runtime validation, or a better API over `as`.
8. Never use a double assertion such as `value as unknown as Target` without a
   documented compatibility boundary.
9. Avoid non-null assertions (`!`).
10. A non-null assertion is allowed only when an invariant is guaranteed by an
    external lifecycle that TypeScript cannot express, and the reason is
    documented.
11. Never use assertions to make untrusted network, storage, or user data
    appear validated.

## 11. Nullability and Optional Values
1. Use `undefined` for an omitted or not-yet-provided value in internal APIs.
2. Use `null` only when the domain or external protocol explicitly distinguishes
   it from omission.
3. Do not use `null` and `undefined` interchangeably for the same state.
4. Use optional properties for absence:

```ts
interface UserPatch {
  displayName?: string;
}
```

5. Use `property: T | undefined` only when the property must exist but its value
   may be undefined.
6. Narrow nullable values close to the boundary.
7. Prefer guard clauses over optional chaining through required domain state.
8. Use `??` for nullish defaults; do not use `||` when `0`, `false`, or `""`
   are valid values.
9. Do not return `undefined` for an unexpected failure; throw or return an
   explicit result.

## 12. Object Types
1. Use an `interface` for an object contract intended for extension or
   implementation.
2. Use a `type` alias for unions, intersections, tuples, mapped types,
   conditional types, and function types.
3. For a closed internal object shape, either form is acceptable; remain
   consistent within the owning module.
4. Do not rely on declaration merging except for intentional third-party or
   platform augmentation.
5. Prefer `extends` for interface composition over deep intersections of object
   types.
6. Avoid index signatures when the set of keys is known.
7. Use `Record<Key, Value>` or a mapped type for a known key domain.
8. Mark properties `readonly` when mutation is not part of the contract.
9. Do not create a second type that duplicates an existing domain shape without
   a semantic reason.

## 13. Domain Modeling
1. Represent a finite set of states with a union of literals or a discriminated
   union.
2. Give each union member a stable discriminant.
3. Keep state-specific fields on the corresponding member.
4. Do not model mutually exclusive states with several independent booleans.

Avoid:

```ts
interface RequestState {
  isLoading: boolean;
  isSuccess: boolean;
  error?: Error;
  data?: User;
}
```

Prefer:

```ts
type RequestState =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "success"; user: User }
  | { status: "failure"; error: Error };
```

5. Use branded types only when confusing structurally identical values creates
   a meaningful correctness risk.
6. Do not create branded primitives for every identifier by default.
7. Keep transport DTOs separate from domain models when validation,
   normalization, or naming differs.

## 14. Enums and Constants
1. Prefer literal unions with `as const` for simple closed values.
2. Use an enum only when runtime enum identity, reverse mapping, or integration
   with an external enum contract is required.
3. Do not use numeric enums for persisted or network values.
4. Assign explicit values to externally visible enums.
5. Keep constants near their owner.
6. Do not create a constants file containing unrelated domains.

```ts
const tokenKinds = ["access", "refresh"] as const;
type TokenKind = (typeof tokenKinds)[number];
```

## 15. Functions
1. A function should perform one understandable operation.
2. Prefer named functions for exported operations and stack-trace clarity.
3. Use arrow functions for concise callbacks and lexical `this`.
4. Keep parameter lists short.
5. Use an object parameter when there are three or more related parameters,
   optional parameters, or any ambiguous booleans.
6. Avoid boolean mode parameters.
7. Prefer separate functions when modes have meaningfully different behavior.
8. Put optional parameters after required parameters.
9. Do not mark callback parameters optional unless the callback may actually be
   invoked without them.
10. Prefer a union parameter over overloads when input variants share one
    return contract.
11. Use overloads only when they materially improve caller types and one
    implementation can honor every signature.
12. Keep side effects visible in names and placement.
13. Do not mix data fetching, validation, transformation, and persistence in one
    function when each is independently meaningful.

## 16. Generics
1. Use a generic only when it expresses a real relationship between types.
2. A type parameter should normally appear at least twice.
3. Use as few type parameters as possible.
4. Push type parameters down to the narrowest value that needs them.
5. Prefer a concrete type when a generic adds no caller flexibility.
6. Keep constraints minimal.
7. Name simple parameters `T`, `K`, and `V`; use a descriptive PascalCase name
   when several parameters would otherwise be confusing.
8. Avoid conditional or recursive types that make common errors unreadable.
9. Export advanced utility types only when multiple public consumers need the
   exact abstraction.
10. Do not implement type-level computation for behavior that should be runtime
    code.

## 17. Control Flow
1. Put the happy path at the lowest indentation level.
2. Use guard clauses for invalid state and errors.
3. Avoid nested ternary expressions.
4. Use a `switch` for a discriminated union or several mutually exclusive
   states.
5. Make switches exhaustive when all states must be handled.
6. Use `never` to detect missing variants:

```ts
function assertNever(value: never): never {
  throw new Error(`Unexpected value: ${String(value)}`);
}
```

7. Always use braces for nested or multiline control flow.
8. Do not rely on truthiness when valid values include `0`, `false`, or an empty
   string.
9. Avoid clever short-circuit expressions for side effects.
10. Use array methods when they express the operation directly; use a loop when
    control flow, early exit, or mutation is clearer.

## 18. Classes
1. Prefer plain functions and data for stateless behavior.
2. Use a class when identity, encapsulated mutable state, lifecycle, or
   polymorphic behavior makes it clearer.
3. Prefer composition over inheritance.
4. Keep constructors cheap and deterministic.
5. Do not perform network or filesystem I/O in a constructor.
6. Make required dependencies constructor parameters.
7. Use parameter properties only when they remain readable and the repository
   formatter supports the style.
8. Mark fields `readonly` when they do not change after construction.
9. Use `override` for overridden members.
10. Avoid static mutable state.
11. Do not create `Manager`, `Helper`, or `Base` classes without a precise
    responsibility.

## 19. Immutability and Ownership
1. Prefer immutable inputs and return values at module boundaries.
2. Use `readonly` and `ReadonlyArray<T>` when callers must not mutate values.
3. Do not expose internal mutable arrays, maps, sets, or objects directly.
4. Clone at a boundary when shared mutation would violate ownership.
5. Do not deep-clone by default; define ownership clearly instead.
6. Avoid mutating function arguments.
7. Local mutation is acceptable when it is contained and clearer than repeated
   allocation.
8. Do not use `Object.freeze` as a substitute for a clear type and ownership
   contract.

## 20. Collections
1. Use arrays for ordered sequences.
2. Use `Set` for uniqueness when membership behavior matters.
3. Use `Map` when keys are dynamic or not naturally object property names.
4. Use objects or `Record` for fixed string-keyed records.
5. Avoid sparse arrays.
6. Do not assume indexed access is in bounds.
7. Sort explicitly before producing deterministic output.
8. Avoid repeated array passes in measured hot paths, but prefer readable
   transformations elsewhere.
9. Do not use `reduce` when a loop or named intermediate steps are clearer.

## 21. Asynchronous Code
1. Return `Promise<T>` from asynchronous public functions.
2. Await promises whose completion or failure matters.
3. Do not create a floating promise.
4. If an operation is intentionally detached, mark it with `void`, handle its
   rejection, and document lifecycle ownership.
5. Use `Promise.all` for independent work that should fail as a group.
6. Do not serialize independent operations with unnecessary sequential awaits.
7. Do not use `new Promise(async (...) => ...)`.
8. Accept an `AbortSignal` for cancellable long-running or I/O operations when
   the platform supports it.
9. Propagate cancellation to downstream operations.
10. Define timeout ownership explicitly.
11. Clean up timers, listeners, streams, and subscriptions.
12. Do not swallow promise rejections.

## 22. Errors
1. Throw `Error` objects, not strings, numbers, or arbitrary objects.
2. Use a custom error class only when callers need stable classification or
   structured fields.
3. Preserve the original error with `cause` when adding context and the
   configured JavaScript runtime target supports `Error.cause`.
4. Catch at a boundary that can recover, translate, add context, or report.
5. Do not catch only to rethrow the same value.
6. Treat caught values as `unknown`.
7. Narrow before reading properties:

```ts
try {
  await saveUser(user);
} catch (error: unknown) {
  if (error instanceof Error) {
    throw new Error("Failed to save user", { cause: error });
  }
  throw new Error("Failed to save user");
}
```

8. Use result unions for expected domain outcomes when failure is part of normal
   control flow.
9. Use exceptions for unexpected inability to complete an operation.
10. Do not expose sensitive payloads in error messages.
11. Do not log and rethrow at every layer.

## 23. Runtime Boundaries
Validate all untrusted data:
1. HTTP requests and responses
2. message queues and events
3. environment variables
4. files and databases
5. browser storage
6. command-line arguments
7. values returned by untyped libraries

Rules:
1. Decode to `unknown`.
2. Validate shape, allowed values, ranges, and cross-field invariants.
3. Normalize once at the boundary.
4. Convert boundary DTOs to domain types.
5. Do not spread an untrusted object into a trusted model.
6. Do not assume `JSON.parse` returns a validated type.
7. Keep validation errors actionable without exposing secrets.

## 24. Comments and Documentation
1. Explain why a constraint or workaround exists.
2. Do not narrate syntax or repeat names.
3. Document exported APIs when their contract is not fully obvious from types.
4. Document runtime behavior that types cannot express:
   - side effects
   - ordering
   - retries
   - cancellation
   - mutation
   - thrown errors
   - units and formats
5. Keep comments synchronized with code.
6. Use `TODO(owner-or-issue): reason` for actionable temporary work.
7. Do not retain commented-out code.
8. Do not use documentation to excuse an unnecessarily complicated API;
   simplify the API first.

## 25. Testing
1. Test observable behavior.
2. Keep tests deterministic and isolated.
3. Use descriptive test names that state the condition and outcome.
4. Follow Arrange, Act, Assert when it improves readability, without requiring
   comments for trivial sections.
5. Prefer real domain values and focused fakes over broad mocks.
6. Mock at external boundaries, not every internal function.
7. Test success, expected failure, boundary values, nullability, and
   cancellation where relevant.
8. Avoid sleeps and real network calls in unit tests.
9. Restore global state, timers, and mocks after each test.
10. Do not weaken production types to simplify tests.
11. Type-check test code.

## 26. Generated Code and Declarations
1. Do not hand-edit generated files.
2. Keep generated files clearly named or located.
3. Commit generated output only when repository policy requires it.
4. Validate generator and generated output in the same change.
5. Keep ambient declarations narrow.
6. Do not use a broad `declare module "*"` to silence missing types.
7. Prefer upstream or precise local declarations over `any` shims.

## 27. Common Prohibitions
Do not:
- disable `strict` to fix a local type error
- use `any` as the default escape hatch
- trust a type assertion at a runtime boundary
- use `var`
- use boxed primitive types
- use `Function` for callbacks
- use non-null assertions routinely
- represent exclusive states with unrelated booleans
- use default exports without a convention-based reason
- create circular barrel imports
- use nested ternaries
- leave promises floating
- throw non-Error values
- catch and ignore failures
- mutate shared data without clear ownership
- add advanced generic types that make callers harder to understand

## 28. Validation
Run the repository's configured equivalents of:

```sh
tsc --noEmit
<formatter> --check .
<linter> .
<test-runner>
```

1. Start with the smallest command covering the changed package.
2. Expand validation for shared types, compiler settings, public exports, and
   build configuration.
3. Do not claim type safety when only transpilation was run.
4. Do not suppress a linter or compiler error without explaining the safety
   argument.

## 29. Review Checklist
- Is the runtime behavior visible from the code?
- Are public and boundary types explicit?
- Is untrusted data validated before use?
- Are nullable states modeled consistently?
- Can a discriminated union replace conflicting booleans?
- Are assertions and `any` absent or narrowly justified?
- Are generics expressing real relationships?
- Are functions and modules focused on one responsibility?
- Are promises awaited, returned, or deliberately owned?
- Are errors classified and contextualized without duplicate logging?
- Are imports and exports narrow and cycle-free?
- Did type checking, formatting, linting, and tests pass?

## References
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Everyday Types](https://www.typescriptlang.org/docs/handbook/2/everyday-types.html)
- [Narrowing](https://www.typescriptlang.org/docs/handbook/2/narrowing.html)
- [More on Functions](https://www.typescriptlang.org/docs/handbook/2/functions.html)
- [Modules](https://www.typescriptlang.org/docs/handbook/2/modules.html)
- [TSConfig Reference](https://www.typescriptlang.org/tsconfig/)
