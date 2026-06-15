# ARCHITECTURE.md
This guide defines the architecture, module boundaries, dependency direction,
and integration rules for `systems/xquare-application` and the root
`contracts` area.

Follow `AGENTS.md` and `documents/guide/common.md` first. Follow
`documents/guide/bazel.md` when changing Bazel packages, dependencies, or
visibility.

## 1. Architectural Style
`systems/xquare-application` is a modular monolith with hexagonal boundaries.

The application:

1. is built and deployed as one Spring Boot process
2. is divided into domain-oriented source and Bazel modules
3. keeps business policy independent from Spring and infrastructure details
4. uses explicit application contracts between business logic and adapters
5. enforces dependency direction through Bazel targets and visibility

The modules are not independently deployed microservices. A module boundary is
an ownership and dependency boundary inside one deployment unit.

Do not add network calls, serialization, or Proto contracts between modules
solely because they are separate Bazel modules.

## 2. Application Modules
The initial modules are:

| Module | Responsibility |
|---|---|
| `xquare-application-identity` | Authentication identities, credentials, sessions, and identity lifecycle |
| `xquare-application-organization` | Organization membership, roles, groups, and organizational relationships |
| `xquare-application-announcement` | Announcements and their publication lifecycle |
| `xquare-application-content` | General authored content that is not owned by the announcement lifecycle |
| `xquare-application-ticket` | Ticket issuance, ownership, state, and processing lifecycle |
| `xquare-application-infrastructure` | Spring Boot composition root and application-wide technical wiring |

These responsibilities are initial boundaries, not a substitute for domain
analysis. Change a boundary when use cases, invariants, ownership, or change
lifecycle show that the current boundary is incorrect.

The `xquare-application-` prefix exists to make modules recognizable in the
monorepo. Do not repeat that prefix in Kotlin package names or treat it as a
domain concept.

## 3. Deployment Boundary
There is one deployable application:

```text
//systems/xquare-application:xquare_application
```

The root target aliases the binary owned by
`xquare-application-infrastructure`.

Domain modules must not define independent Spring Boot entry points. A future
decision to deploy a module independently requires an explicit architecture
change; Bazel modularity alone does not imply independent deployment.

## 4. Module Structure
Each business module starts with these boundaries:

```text
xquare-application-<module>/
└── src/
    ├── main/
    │   ├── kotlin/kr/dsm/hs/xquare/<module>/
    │   │   ├── domain/
    │   │   ├── application/
    │   │   └── adapter/
    │   └── resources/
    └── test/
        └── kotlin/kr/dsm/hs/xquare/<module>/
```

Create subpackages only when real source responsibilities require them. Do not
pre-create empty package trees to imitate an architecture diagram.

An expected package shape after implementation begins is:

```text
<module>/
├── domain/
│   ├── model/
│   ├── service/
│   ├── event/
│   └── exception/
├── application/
│   ├── usecase/
│   └── port/
└── adapter/
    ├── web/
    ├── persistence/
    ├── messaging/
    └── security/
```

This shape is illustrative. Add only the packages used by the module.

## 5. Domain Boundary
The `domain` package owns business meaning and invariant enforcement.

It may contain:

1. entities and aggregates
2. value objects
3. domain services
4. domain events
5. business exceptions
6. pure business policies

The domain must not depend on:

1. Spring
2. HTTP or RPC types
3. persistence APIs or database drivers
4. JWT libraries or token claims
5. generated Proto classes
6. another module's domain model
7. the `application`, `adapter`, or `infrastructure` packages

Prefer ordinary Kotlin types in the domain. An external library is allowed only
when it represents a stable domain-level primitive and does not introduce a
framework, transport, persistence, or runtime dependency. Such an exception
requires explicit review.

## 6. Application Boundary
The `application` package owns use cases and orchestration.

It may contain:

1. use case interfaces
2. command and query models
3. use case implementations
4. transaction-oriented orchestration
5. authorization decisions based on an application principal
6. outbound port interfaces required by use cases

The application may depend on its own domain. It must not depend on adapter or
infrastructure implementations.

The application must remain independent from Spring, persistence frameworks,
HTTP, JWT, and generated transport types. Framework annotations are not a
replacement for an explicit application contract.

### 6.1 Ports
A port is a source-level abstraction owned by the inside of the hexagon. It
does not mean that the application depends on an external implementation.

An outbound port describes a capability required by a use case:

```kotlin
interface IdentityRepository {
    fun findById(identityId: IdentityId): Identity?
}
```

The interface belongs to `application/port` or, when it is a domain-level
abstraction, to the domain package. A persistence adapter implements it and
depends inward on the interface.

An inbound port is the application contract invoked by a driving adapter. In
this repository, a use case interface under `application/usecase` normally
serves this role.

The dependency direction is:

```text
driving adapter -> application use case -> domain
                             |
                             v
                       outbound port <- driven adapter
```

The arrows point toward source dependencies or abstractions. Runtime control
returns through the adapter implementation bound to the outbound port.

## 7. Adapter Boundary
The `adapter` package translates between the application and technical systems.

Driving adapters invoke use cases. Examples include:

1. HTTP controllers
2. gRPC handlers
3. message consumers
4. schedulers
5. command-line handlers

Driven adapters implement outbound ports. Examples include:

1. database repositories
2. cache clients
3. external service clients
4. message publishers
5. object storage clients

Adapters may depend on external frameworks and libraries. They may depend on
their module's application and domain targets. They must not expose framework,
database, JWT, or generated transport types through application or domain
contracts.

`driving` and `driven`, or `inbound` and `outbound`, are explanatory terms, not
mandatory directory names from hexagonal architecture. Prefer concrete package
names such as `web`, `persistence`, or `messaging`. Introduce directional
grouping only when the number of adapters makes it useful.

## 8. Infrastructure Module
`xquare-application-infrastructure` is the composition root for the single
Spring Boot application.

It owns:

1. the `main` function and `@SpringBootApplication`
2. component scanning and module assembly
3. application-wide configuration
4. security filter-chain wiring
5. technical bean construction
6. runtime resources and deployment configuration
7. context-load and composition tests

Infrastructure may depend on module application and adapter targets to assemble
the runtime graph. Business modules must not depend on infrastructure.

Infrastructure is not a general-purpose dumping ground. A domain-specific
controller, repository, client, mapper, or framework configuration belongs in
the corresponding business module's adapter package when that module owns its
behavior.

Cross-cutting technical code belongs in infrastructure only when it genuinely
applies to the application as a whole and has no single business owner.

## 9. Dependency Rules
The default source dependency direction inside a module is:

```text
adapter -> application -> domain
infrastructure -> adapter/application
```

An adapter may depend directly on domain types when translation requires it,
but use case invocation should pass through the application boundary.

The following dependencies are prohibited:

```text
domain -> application
domain -> adapter
domain -> infrastructure
application -> adapter
application -> infrastructure
business module -> infrastructure
```

Declare every direct dependency in Bazel. Visibility grants permission to
depend on a target; it does not create the dependency. Add `deps` only when
source code actually uses the target.

Default every Bazel package to private visibility. Open visibility only to
actual consumers. Treat visibility expansion as an architecture change.

## 10. Module-to-Module Collaboration
A module must not read or mutate another module's repositories, persistence
entities, adapter implementations, or internal domain objects.

Use one of these collaboration forms:

1. call an explicit application-level contract owned by the target module for
   synchronous behavior
2. publish and consume an internal application or domain event for behavior
   that does not require an immediate result
3. duplicate a small read model when that avoids coupling business ownership

Do not create a broad `common`, `shared`, or `utils` module to bypass ownership.
Extract shared code only when it represents a stable concept with clear
ownership and at least two real consumers.

Cross-module Bazel visibility is not enabled by default. Add it only for the
specific application contract being consumed. Never expose an entire adapter
or domain package to make a cross-module call compile.

Internal module collaboration uses Kotlin contracts and in-process calls. Do
not use Proto, HTTP, or gRPC between modules in the same application process.

## 11. Authentication and JWT
JWT is a transport and security mechanism, not a domain model.

The authentication flow is:

```text
request
  -> security or web adapter
  -> JWT signature, issuer, audience, and expiry validation
  -> normalized application principal
  -> application use case
  -> authorization and business policy
```

Rules:

1. JWT parsing and cryptographic validation belong in a security adapter or
   application-wide infrastructure security configuration.
2. JWT library types and raw claim maps must not enter domain or application
   contracts.
3. Adapters must convert validated token data into a small application-owned
   principal, such as an actor ID and granted authorities.
4. Application and domain policy must authorize actions from the normalized
   principal and owned business data, not by reparsing the token.
5. Identity owns token issuance, refresh, revocation, and identity-session
   semantics when those behaviors are implemented.
6. Other modules may consume the normalized principal but must not depend on
   identity persistence or JWT implementation details.
7. Spring Security filter-chain and bean assembly belong in infrastructure;
   identity-specific authentication behavior belongs in the identity adapter.

The JWT provider, claim schema, signing-key source, and refresh-token strategy
are not established by the current scaffold. Decide and document them before
implementing authentication.

## 12. Proto Contracts
Cross-process service contracts are owned by the repository-level `contracts`
area:

```text
contracts/
├── identity/
├── organization/
├── content/
├── notification/
├── ticket/
└── deployment/
```

Proto contracts model externally consumable process or service boundaries.
They do not model internal Kotlin module boundaries and do not have to map
one-to-one to application modules.

Rules:

1. Manage every `.proto` schema with a Bazel `proto_library` target.
2. Keep the schema target language-neutral.
3. Generate language-specific bindings in targets owned by actual consumers
   or providers.
4. Do not expose generated Proto messages as domain entities or application
   command models.
5. Convert generated transport messages at the adapter boundary.
6. Version published packages and import paths when a contract becomes real,
   for example `contracts/identity/v1/identity_service.proto`.
7. Preserve field numbers, reserve removed fields, and review compatibility
   before changing a published contract.
8. Narrow `//visibility:public` to known consumers when the initial consumer
   set is established.

`//contracts/identity:identity_example_proto` is temporary scaffolding. It
proves that schema parsing and descriptor generation are managed by Bazel; it
does not define production identity semantics.

Proto linting and breaking-change automation are follow-up decisions. Add tools
such as Buf only through a separate reviewed build and architecture change.

## 13. Data and Transactions
Each business module owns its persistence model and repository interfaces.
Do not share persistence entities across module boundaries.

A single physical database may be used by the deployment, but table ownership
must remain explicit. One module must not directly update another module's
tables.

Keep transactions within one module by default. A use case that requires
cross-module atomicity must be reviewed before implementation. Prefer explicit
orchestration and reliable events over hidden repository coupling.

The database technology, migration ownership, and cross-module event delivery
mechanism are not established by the current scaffold.

## 14. Testing Boundaries
Test each boundary at the narrowest useful level:

1. domain tests cover invariants and pure business policy without Spring
2. application tests execute use cases with fake or in-memory ports
3. adapter tests cover serialization, mapping, persistence, and framework
   integration
4. infrastructure tests verify Spring composition and application startup
5. contract builds verify Proto parsing and descriptor generation

Do not use a full Spring context test when a domain or application unit test
covers the behavior.

Bazel test targets must follow the same dependency and visibility rules as
production targets.

## 15. Naming Decisions
Use names that describe business ownership or concrete technical roles.

Current decisions:

1. use `identity`, not `identify`
2. use `organization`, not `team`, because the boundary may include groups,
   membership, and roles beyond teams
3. use `announcement`, not `notice`, for the announcement lifecycle
4. use `content`, not `publication` or `blog`, for general authored content
5. use `ticket` as a separate business module
6. use `infrastructure` for the former bootstrap and platform responsibilities

Do not rename a module based only on preference. A rename must reflect a change
in business responsibility or remove a demonstrated ambiguity.

## 16. Current Decisions and Open Decisions
The following decisions are established:

1. one Spring Boot deployment
2. domain-oriented Bazel modules
3. `domain`, `application`, and `adapter` boundaries in each business module
4. infrastructure as the composition root
5. inward dependency direction
6. no framework dependency in domain code
7. no Proto transport between in-process modules
8. root-level, Bazel-managed Proto contracts for process boundaries

The following decisions remain open:

1. detailed domain models and aggregate boundaries
2. specific cross-module application contracts and events
3. JWT provider, claim schema, and session strategy
4. persistence technology and migration ownership
5. transaction handling and event delivery
6. language-specific Proto and gRPC generation targets
7. Proto linting and compatibility tooling
8. final adapter subpackage structure after real adapters exist

Do not present an open decision as current architecture. Resolve it through a
separate reviewed change when implementation requires it.

## 17. Architecture Review Checklist
Before adding or changing application code, verify:

- Does the code belong to the module that owns the business behavior?
- Does the dependency point inward?
- Is a port owned by the code that requires the capability?
- Are framework and transport types confined to adapters or infrastructure?
- Is cross-module access through an explicit application contract?
- Is Proto used only across a process boundary?
- Is JWT converted to an application principal before use cases run?
- Does Bazel visibility expose only intended consumers?
- Does the change preserve one deployable application?
- Is an open architecture decision being introduced without documentation?
