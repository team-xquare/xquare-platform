# RESEARCH.md
This guide defines how to investigate project problems, verify improvement
opportunities, and produce evidence-based conclusions.
It applies to logical correctness, reliability, concurrency, performance,
resource usage, and operational behavior.

## 1. Purpose
Research must answer a concrete question with reproducible evidence.
It must not begin with a preferred solution and search only for supporting
arguments.

The standard investigation sequence is:

```text
symptom
-> impact
-> evidence
-> hypothesis
-> reproduction
-> isolation
-> verified cause
-> controlled change
-> before-and-after validation
-> regression detection
```

The purpose of research is not to produce a long list of possible concerns.
The purpose is to distinguish:
1. behavior that is correct
2. behavior that is suspicious but unverified
3. behavior that is demonstrably incorrect
4. performance that is measurable but acceptable
5. a bottleneck that is verified and worth improving

## 2. Core Principles
1. Measure before optimizing.
2. Reproduce before diagnosing.
3. Separate observations from interpretations.
4. State assumptions explicitly.
5. Prefer primary evidence over intuition.
6. Change one important variable at a time.
7. Compare results under equivalent conditions.
8. Include unsuccessful hypotheses in the investigation record when they
   prevent repeated work.
9. Do not call code a bottleneck without runtime evidence.
10. Do not call a behavior a bug without identifying the violated contract or
    invariant.
11. Do not generalize a local benchmark result to production without explaining
    workload and environment differences.
12. Do not introduce a new observability or profiling platform as part of
    research unless the task explicitly approves it.

## 3. Evidence Levels
Classify every finding with one of these levels.

### 3.1 Observation
An event or measurement was recorded, but its cause is unknown.

Examples:
- request `p99` latency increased
- heap usage grew during a test
- a trace-like request log shows time spent before a database call
- a code path appears complex during review

An observation is not a root cause.

### 3.2 Suspected Cause
A plausible explanation is consistent with some evidence but has not been
isolated or reproduced.

Examples:
- repeated repository calls may be an N+1 query pattern
- lock contention may explain low CPU utilization and high latency
- object retention may explain heap growth

Use language such as `may`, `likely`, or `suspected`.

### 3.3 Reproduced Problem
The undesirable behavior occurs reliably under a documented condition.
The root cause may still be unknown.

### 3.4 Verified Cause
The cause is isolated by a profiler, debugger, execution evidence, controlled
experiment, or a minimal reproduction.

### 3.5 Verified Improvement
The proposed change improves the target result under equivalent conditions
without violating correctness or moving the cost to an unmeasured area.

Only levels 3.4 and 3.5 may be described as confirmed conclusions.

## 4. Research Question
Begin with one narrow question.

Good:

```text
Why does POST /sessions exceed 800 ms at p99 during 200 requests per second?
```

```text
Can two concurrent refresh requests create more than one active token?
```

Avoid:

```text
Why is the service slow?
```

```text
Find every problem in the project.
```

Record:
1. the behavior being investigated
2. the expected behavior
3. the affected user or system
4. the relevant service, endpoint, job, or module
5. the time window
6. the deployed version or commit
7. the required confidence level

## 5. Investigation Scope
Define scope before collecting data.

Include:
1. entry point
2. downstream dependencies
3. data stores
4. asynchronous jobs or messages
5. concurrency level
6. environment
7. deployment version
8. representative inputs

Exclude unrelated components explicitly.
If evidence points outside the approved scope, report the scope change and
follow `documents/guide/work.md`.

## 6. Baseline
Create a baseline before making a change.

A baseline must record:
1. exact code revision
2. configuration relevant to behavior
3. runtime and tool versions
4. machine or container resources
5. dataset size and shape
6. warm-up procedure
7. workload and concurrency
8. test duration
9. number of repetitions
10. measured results
11. observed variance

Do not compare:
- a cold run with a warm run
- a local laptop with a production instance
- different datasets
- different concurrency levels
- different runtime versions
- different logging or profiler settings

unless the difference is the variable under investigation.

## 7. Existing Observability Data
Use the project's existing logs and metrics before adding instrumentation.

Correlate records using available fields such as:
1. timestamp with timezone
2. request or correlation ID
3. service name
4. service version or Git commit
5. instance, pod, or host
6. endpoint, operation, or job name
7. safe domain identifier
8. error classification
9. duration

Rules:
1. Use one time range and timezone across all evidence.
2. Confirm clock synchronization before comparing timestamps across services.
3. Confirm that records belong to the same version and environment.
4. Do not assume adjacent log lines belong to the same request.
5. Prefer stable structured fields over message-text parsing.
6. Do not add high-cardinality metric labels such as raw user IDs, request IDs,
   URLs with identifiers, stack traces, or arbitrary error messages.
7. Do not log secrets, tokens, request bodies, private data, or credentials.
8. If correlation fields are absent, state the limitation instead of claiming a
   cross-service causal chain.

## 8. Signals and Their Questions
Use each signal for the question it can answer.

### 8.1 Logs
Logs answer:
- what event occurred
- what decision was made
- which error was observed
- which safe contextual fields were present

Logs usually do not prove where CPU or memory was consumed.

### 8.2 Metrics
Metrics answer:
- whether behavior changed over time
- how frequently it occurs
- whether a resource is saturated
- whether an objective or threshold is violated

Metrics usually identify where to investigate, not the exact code responsible.

### 8.3 Request Timing
Request timing or existing distributed request records answer:
- which operation was slow
- which dependency contributed latency
- whether delay was local or downstream

Timing without code-level profiling does not identify an expensive function.

### 8.4 Profiles
Profiles answer:
- which code consumed CPU
- which code allocated or retained memory
- where threads or goroutines blocked
- which locks were contended
- how runtime scheduling affected execution

Profiles must be captured during the relevant workload.

## 9. Logical Problem Investigation
Logical research starts from a contract or invariant.

### 9.1 Contracts
Inspect:
1. public API contract
2. database constraints
3. message schema
4. state machine
5. authentication and authorization rules
6. retry and timeout behavior
7. ordering guarantees
8. compatibility requirements

If the expected behavior is undocumented, obtain or define the requirement
before classifying behavior as incorrect.

### 9.2 Invariants
Write important invariants as testable statements.

Examples:

```text
At most one active refresh token exists per session.
```

```text
A completed payment cannot return to a pending state.
```

```text
Every accepted event is either processed or present in a retryable queue.
```

An invariant should identify:
1. subject
2. required condition
3. lifetime or transaction boundary
4. concurrency assumptions

### 9.3 Boundary Conditions
Check:
1. empty input
2. maximum and minimum values
3. missing optional values
4. duplicate requests
5. malformed external data
6. expired data
7. clock boundaries and timezones
8. pagination boundaries
9. partial collections
10. integer overflow or precision
11. Unicode and encoding
12. cancellation and timeout

### 9.4 State Transitions
For stateful behavior:
1. list all states
2. list allowed transitions
3. identify the actor allowed to trigger each transition
4. record required preconditions
5. record side effects
6. define behavior for repeated transitions
7. define behavior under concurrent transitions
8. define recovery after partial failure

Use a state table when prose is ambiguous.

| Current state | Event | Next state | Required effect |
|---|---|---|---|
| `pending` | approve | `approved` | persist approval once |
| `approved` | approve | `approved` | return idempotent result |
| `rejected` | approve | unchanged | reject transition |

### 9.5 Concurrency
Investigate:
1. data races
2. lost updates
3. duplicate work
4. check-then-act races
5. lock ordering
6. deadlock
7. starvation
8. stale reads
9. transaction isolation
10. queue ordering
11. cancellation races
12. shutdown behavior

Reproduce with controlled synchronization rather than repeated random sleeps.
Use barriers, latches, channels, test clocks, or deterministic schedulers where
the language provides them.

### 9.6 Partial Failure
Check behavior when failure occurs:
1. before a write
2. between multiple writes
3. after a write but before acknowledgment
4. during an external call
5. after an external success but before local persistence
6. during retry
7. during process shutdown

Determine whether the system:
- retries safely
- duplicates side effects
- loses work
- leaves inconsistent state
- exposes a misleading success response

### 9.7 Error Handling
Trace each relevant error from origin to boundary.

Check:
1. whether the error is ignored
2. whether identity or classification is lost
3. whether it is incorrectly converted to success
4. whether retryability is preserved
5. whether the caller receives an actionable status
6. whether it is logged repeatedly
7. whether sensitive data is exposed
8. whether cleanup still occurs

### 9.8 Logical Reproduction
A logical reproduction must include:
1. initial state
2. exact input
3. execution order
4. concurrency or timing controls
5. expected result
6. actual result
7. evidence of the violated invariant

Convert a verified logical problem into a deterministic regression test whenever
practical.

## 10. Performance Problem Definition
Performance is a requirement, not a general preference for faster code.

Define the violated target:
1. latency
2. throughput
3. error rate under load
4. CPU
5. memory
6. allocation rate
7. storage or network I/O
8. queue delay
9. startup time
10. cost per operation

Good:

```text
Under 200 requests per second, p99 latency must remain below 500 ms while the
error rate remains below 0.1%.
```

Avoid:

```text
This loop should be optimized.
```

## 11. Performance Metrics
Collect metrics relevant to the suspected limit.

### 11.1 Service-Level Results
- request count
- success and error rate
- throughput
- `p50`, `p90`, `p95`, and `p99` latency
- timeout and cancellation rate
- queue waiting time
- work duration

Do not rely on averages alone. Averages can hide tail latency.

### 11.2 Resource Results
- CPU utilization and throttling
- resident and heap memory
- allocation rate
- garbage collection frequency and pause duration
- goroutine or thread count
- event-loop utilization or delay
- file descriptor count
- network throughput and errors
- disk latency and throughput
- connection-pool use and wait time
- worker-pool queue depth

### 11.3 Dependency Results
- database call count and duration
- query rows read and returned
- query plan
- cache hit and miss rate
- external API duration and error rate
- retry count
- message queue lag
- serialization payload size

## 12. Saturation Model
For each constrained resource, inspect:
1. utilization
2. saturation
3. errors

High utilization without saturation may be healthy.
Low utilization with high latency may indicate:
- lock contention
- connection-pool waiting
- downstream I/O
- queueing
- scheduler delay
- rate limiting

Do not conclude that CPU optimization is needed merely because latency is high.

## 13. Workload Design
A performance result is valid only for the workload it measures.

Record:
1. request mix
2. payload distribution
3. read/write ratio
4. dataset cardinality
5. cache state
6. concurrency
7. arrival pattern
8. test duration
9. warm-up duration
10. external dependency behavior

Prefer a workload based on existing production or staging observations.
When production data is unavailable, state how the synthetic workload differs.

Include steady-state and relevant burst behavior separately.
Do not combine them into one unexplained result.

## 14. Measurement Quality
1. Warm up runtimes and caches when measuring steady-state behavior.
2. Record cold-start behavior separately when startup matters.
3. Repeat measurements.
4. Report distribution or confidence information, not only the best run.
5. Keep background load controlled.
6. Pin or record CPU and memory limits.
7. Record profiler overhead when it may affect results.
8. Do not run unrelated workloads on the same test environment.
9. Preserve raw result files when permitted.
10. Make the test command reproducible.

An improvement smaller than normal variance is not a verified improvement.

## 15. Performance Investigation Sequence
Use this order:
1. confirm the user-visible symptom
2. identify the affected operation and time window
3. compare with a healthy baseline
4. determine whether a resource or dependency is saturated
5. narrow the delay to service, database, network, queue, or runtime
6. capture a profile during the representative workload
7. identify dominant stacks or wait sites
8. form one testable cause hypothesis
9. design a minimal experiment
10. compare before and after
11. check correctness and secondary metrics
12. add regression detection

Do not begin with a microbenchmark when the production symptom has not been
connected to the code being benchmarked.

## 16. Profiling Selection
Select the profile based on the symptom.

| Symptom | First profile or evidence |
|---|---|
| high CPU | CPU profile or execution samples |
| high latency and low CPU | blocking, lock, thread, goroutine, or dependency evidence |
| growing memory | heap retention over time |
| high GC activity | allocation profile and GC data |
| increasing process memory | heap plus native/process memory evidence |
| stalled requests | thread or goroutine dump and pool metrics |
| periodic pauses | runtime trace, GC, scheduler, or safepoint evidence |
| slow browser interaction | main-thread performance recording |
| slow Node.js responses | event-loop delay plus CPU or dependency profile |

Capture the profile while the symptom is present.
A profile from an idle or unrelated workload is not evidence for the incident.

## 17. Go Investigation
Use Go's standard diagnostics first.

### 17.1 Benchmarks
Use benchmarks for a narrow, repeatable unit after connecting it to the
observed problem.

```sh
go test -run '^$' -bench BenchmarkName -benchmem -count 10 ./path/to/package
```

Rules:
1. exclude setup from the measured operation
2. prevent unwanted compiler elimination
3. use representative values
4. report allocations when relevant
5. run multiple iterations
6. compare with a statistical benchmark comparison tool when available
7. do not infer service throughput from a microbenchmark alone

### 17.2 CPU Profiles
Capture CPU profiles through tests, `runtime/pprof`, or an approved protected
`net/http/pprof` endpoint.

```sh
go test -run '^$' -bench BenchmarkName -cpuprofile cpu.out ./path/to/package
go tool pprof -http=:0 cpu.out
```

Inspect:
- cumulative time
- flat time
- hot call paths
- repeated serialization or conversion
- expensive runtime work

### 17.3 Memory Profiles
Distinguish allocation rate from retained heap.

```sh
go test -run '^$' -bench BenchmarkName -benchmem \
  -memprofile memory.out ./path/to/package
go tool pprof -http=:0 memory.out
```

Use:
- `alloc_space` or allocation profiles for allocation pressure
- `inuse_space` or heap profiles for retained live memory

A high total allocation count is not automatically a leak.
A leak requires evidence that retained memory continues to grow because objects
remain reachable or resources remain owned.

### 17.4 Goroutine Profiles
Use goroutine profiles or dumps to identify:
- blocked goroutines
- leaked goroutines
- repeated stack patterns
- waits on channels, network, locks, or shutdown

Compare snapshots over time when investigating leaks.

### 17.5 Block and Mutex Profiles
Use block profiles for time spent waiting on synchronization and mutex profiles
for lock contention.
These profiles require explicit sampling configuration and add overhead.
Enable them only for a bounded investigation.

Inspect total delay and contended call sites, not only event counts.

### 17.6 Execution Trace
Use the execution tracer for scheduler, goroutine, network, synchronization, and
garbage-collection behavior.

```sh
go test -run TestName -trace trace.out ./path/to/package
go tool trace trace.out
```

Keep traces short enough to analyze and capture them during the relevant
operation.

### 17.7 Race Detection
Use the race detector for concurrent correctness.

```sh
go test -race ./path/to/package/...
```

The race detector only detects races exercised by the workload.
A passing race-enabled test does not prove that every path is race-free.

### 17.8 Go Evidence Checklist
- Was the profile captured under the affected workload?
- Was the correct profile type selected?
- Are cumulative and flat costs distinguished?
- Are blocking and dependency waits separated from CPU work?
- Is allocation distinguished from retention?
- Was race-enabled behavior exercised?
- Are profiler endpoints protected and disabled when unnecessary?

## 18. Kotlin and Spring Investigation
Use Spring metrics for symptom localization and JVM diagnostics for code-level
evidence.

### 18.1 Spring Actuator and Existing Metrics
Use approved Actuator endpoints and the project's existing Micrometer registry
to inspect:
- request count, errors, and latency
- JVM heap and non-heap memory
- garbage collection
- thread count
- process CPU
- executor and scheduler pools
- HTTP client behavior where instrumented
- database connection-pool use and wait
- cache behavior where instrumented

Rules:
1. expose only required management endpoints
2. protect management endpoints with network and authentication controls
3. do not expose sensitive configuration or environment values
4. do not add unbounded metric tags
5. use metrics to locate the problem before taking a detailed profile

### 18.2 Java Flight Recorder
Use Java Flight Recorder for bounded, low-overhead JVM event recording.

JFR can provide evidence for:
- CPU execution samples
- object allocation
- garbage collection
- thread parking and blocking
- monitor contention
- socket and file I/O
- exceptions
- class loading

Example with `jcmd`:

```sh
jcmd <pid> JFR.start name=research settings=profile duration=60s \
  filename=research.jfr
```

Analyze recordings with JDK Mission Control or the JDK `jfr` tool.

Rules:
1. capture the affected workload and time window
2. record JDK version and JFR settings
3. use a bounded duration
4. measure profiling overhead when latency is sensitive
5. avoid expensive heap-statistics events during latency measurement unless
   memory investigation requires them
6. retain recordings only according to project security policy

### 18.3 Thread Dumps
Use multiple thread dumps separated by a short interval for stalls, deadlocks,
and pool exhaustion.

```sh
jcmd <pid> Thread.print
```

Inspect:
- repeatedly blocked request threads
- lock ownership
- deadlock reports
- connection-pool waits
- executor queue waits
- blocking calls on event-loop or coroutine threads

One dump is a snapshot and may be misleading.

### 18.4 Heap and Memory
Use GC and allocation evidence before taking a heap dump.
A heap dump can pause the process, consume substantial disk space, and contain
sensitive data.

Investigate:
1. heap after garbage collection
2. allocation rate
3. dominant retained objects
4. paths from GC roots
5. class-loader retention
6. caches and unbounded collections
7. thread-local retention
8. native memory when process memory exceeds heap evidence

Do not declare a leak from high heap usage alone.

### 18.5 Coroutines
Inspect:
- blocked threads caused by synchronous I/O
- dispatcher saturation
- excessive coroutine creation
- lost cancellation
- unbounded flows or channels
- work launched outside owned scopes
- locks held across suspension

Use coroutine debugging facilities only when supported and approved for the
environment. Account for their overhead.

### 18.6 Spring Data and Persistence
Inspect:
- query count per request
- N+1 access
- query duration
- rows scanned versus returned
- missing or ineffective indexes
- transaction duration
- lock waits
- connection-pool saturation
- excessive entity loading
- serialization after session closure

Do not enable broad SQL or parameter logging in production without reviewing
volume and sensitive-data risk.

### 18.7 Optional Profilers
An additional profiler such as async-profiler may be used when it is already
approved and available.
Record:
1. tool version
2. event type
3. sampling interval
4. duration
5. command
6. overhead

Do not require a new profiler dependency when JFR or existing tooling can answer
the question.

### 18.8 JVM Evidence Checklist
- Did metrics identify the affected time and resource?
- Was JFR captured during the symptom?
- Were CPU, allocation, GC, lock, and I/O costs distinguished?
- Were multiple thread dumps used for a stall?
- Was connection-pool waiting measured?
- Was a heap dump treated as sensitive and potentially disruptive?
- Were coroutine dispatchers and blocking calls considered?

## 19. TypeScript and Node.js Investigation
TypeScript source executes as JavaScript.
Investigate emitted runtime behavior and preserve source maps needed to map
profiles back to TypeScript.

### 19.1 Event Loop
For Node.js request latency, determine whether delay comes from:
- synchronous CPU work
- synchronous filesystem or process APIs
- excessive callbacks or microtasks
- garbage collection
- dependency I/O
- worker-pool saturation
- event-loop delay

Use supported `node:perf_hooks` measurements such as event-loop utilization and
event-loop delay when available in the selected Node.js version.

Record measurement interval and process load.
Event-loop delay alone does not identify the responsible function.

### 19.2 CPU Profiles
Use the Node.js or V8 inspector CPU profiler, `--cpu-prof`, or the project's
approved profiler.

```sh
node --cpu-prof dist/server.js
```

Capture the representative operation and inspect:
- self time
- total time
- hot JavaScript stacks
- native or runtime frames
- serialization
- regular-expression work
- synchronous library calls

Confirm that source maps resolve generated JavaScript frames to the correct
TypeScript source.

### 19.3 Heap Profiles and Snapshots
Use allocation profiles for allocation pressure.
Use heap snapshots for retained-object analysis.

Heap snapshots can:
- pause the event loop
- temporarily require substantial memory
- terminate a process under memory pressure
- contain application data

Do not take a production heap snapshot without operational approval and a
capacity plan.

Compare snapshots or retained sizes over time.
Inspect retainers, detached objects, listener accumulation, caches, closures,
timers, and unbounded maps.

### 19.4 Async Behavior
Inspect:
- floating promises
- missing cancellation or timeout
- unbounded concurrency
- sequential awaits for independent work
- retry multiplication
- listener leaks
- unresolved promises
- timer accumulation
- request context retained by closures

Do not infer asynchronous causality from stack order alone.
Use request IDs, timestamps, supported async diagnostics, and controlled
reproduction.

### 19.5 Worker Threads
Worker threads are appropriate for measured CPU-intensive JavaScript work, not
as a general solution for I/O latency.

Before recommending workers, measure:
1. CPU cost of the operation
2. event-loop impact
3. transfer or serialization cost
4. worker startup cost
5. pool size and queueing
6. memory cost

Prefer a bounded worker pool over one worker per request.

### 19.6 Node.js Evidence Checklist
- Is event-loop delay measured?
- Is synchronous CPU work separated from dependency I/O?
- Does the CPU profile map back to TypeScript?
- Is heap allocation distinguished from retained memory?
- Were snapshot pause and memory risks considered?
- Are promise concurrency and retries bounded?
- Would worker transfer overhead erase the expected benefit?

## 20. Browser TypeScript Investigation
Apply this section to browser applications.

Use field or existing user-experience metrics to identify the affected page and
interaction.
Use Chrome DevTools Performance and Memory tools for local reproduction.

Inspect:
1. long main-thread tasks
2. script evaluation
3. layout and style recalculation
4. rendering and paint
5. excessive event handlers
6. repeated framework rendering
7. network waterfalls
8. bundle loading and parsing
9. detached DOM nodes
10. listener and timer retention

Rules:
1. test with a clean browser profile
2. record device, CPU, and network throttling
3. distinguish page load from runtime interaction
4. preserve source maps
5. compare the same interaction
6. do not treat Lighthouse alone as proof of a production bottleneck
7. correlate local findings with available real-user evidence when possible

## 21. Database Investigation
Application profiles often show time waiting for a database without explaining
why the database is slow.

Collect:
1. normalized query
2. query count
3. duration distribution
4. execution plan
5. rows estimated, scanned, and returned
6. index use
7. lock waits
8. transaction duration
9. connection acquisition time
10. database CPU, memory, and I/O where available

Rules:
1. use the database's supported plan analysis tools
2. avoid running expensive analysis against production without approval
3. remove sensitive literal values from captured queries
4. verify index benefit under representative data distribution
5. include write amplification and storage cost
6. do not optimize a query solely because it appears frequently
7. prioritize total user impact

## 22. External Dependency Investigation
For HTTP, RPC, queue, cache, or storage dependencies, distinguish:
1. connection acquisition
2. DNS and connection setup
3. request transmission
4. server processing
5. response transfer
6. retries and backoff
7. local deserialization
8. queue waiting

Check:
- timeout policy
- cancellation propagation
- retry ownership
- retry count
- circuit or rate limiting
- pool saturation
- payload size
- dependency error rate

Do not count nested retries as one attempt.
Do not optimize local code when most latency is verified downstream.

## 23. Improvement Experiments
An experiment must test one primary hypothesis.

Record:
1. hypothesis
2. expected metric change
3. controlled variable
4. unchanged conditions
5. implementation or configuration change
6. before result
7. after result
8. variance
9. correctness validation
10. secondary costs

Examples of secondary costs:
- lower latency but higher memory
- lower CPU but greater I/O
- higher throughput but worse tail latency
- fewer allocations but more lock contention
- improved mean but worse errors

Reject an improvement that moves the bottleneck or violates correctness without
an explicit accepted tradeoff.

## 24. Prioritization
Rank verified findings using:
1. user impact
2. frequency
3. severity
4. affected scope
5. operational risk
6. security or data-integrity risk
7. cost
8. confidence
9. implementation effort
10. regression risk

Do not rank by code ugliness alone.
A simple code cleanup and a verified production bottleneck are different kinds
of work.

Suggested priority:

```text
priority = impact x frequency x confidence / effort
```

Use the formula as a discussion aid, not as a substitute for judgment.

## 25. Regression Detection
A verified problem should produce the smallest durable detector that is
practical.

Possible detectors:
- deterministic unit or integration test
- concurrency regression test
- benchmark
- load-test scenario
- metric and alert
- dashboard threshold
- query-count assertion
- resource budget
- startup or bundle-size budget

Do not add a noisy alert or unstable benchmark merely to claim coverage.
Define owner, threshold, and response for operational detectors.

## 26. Production Safety
Profiling and diagnostics can be disruptive or sensitive.

Before production collection:
1. obtain required operational approval
2. define target instance
3. define start time and duration
4. estimate CPU, memory, disk, and pause overhead
5. confirm rollback or stop procedure
6. protect diagnostic endpoints
7. define artifact storage and deletion
8. redact sensitive fields
9. avoid broad user payload capture
10. monitor the process during collection

Treat as sensitive:
- heap dumps
- goroutine and thread dumps
- profiles with function and path names
- logs with identifiers
- database queries
- environment and configuration output
- JFR recordings
- browser recordings

Never expose `pprof`, Actuator, JMX, JVM diagnostic, Node inspector, or browser
debugging endpoints publicly.

## 27. Research Record
Use this format for a logical investigation:

```md
## Logical Finding
- Question:
- Expected contract:
- Violated invariant:
- Scope:
- Environment and version:
- Initial state:
- Reproduction:
- Expected result:
- Actual result:
- Evidence:
- Root cause:
- Confidence:
- Affected scope:
- Proposed correction:
- Regression test:
- Remaining uncertainty:
```

Use this format for a performance investigation:

```md
## Performance Finding
- Question:
- Affected SLI:
- Target:
- Environment and version:
- Workload:
- Dataset:
- Baseline:
- Time window:
- Existing log or metric evidence:
- Profile type and command:
- Profile evidence:
- Verified bottleneck:
- Experiment:
- Before:
- After:
- Variance:
- Secondary effects:
- Correctness validation:
- Regression detection:
- Remaining uncertainty:
```

## 28. Tool Record
For every diagnostic artifact, record:
1. tool and version
2. exact command
3. target process or test
4. code revision
5. configuration
6. start and end time
7. workload
8. sampling or recording settings
9. known overhead
10. artifact location and retention policy

This information is required to reproduce and interpret the result.

## 29. Common Failure Modes
Do not:
- optimize code before confirming the symptom
- treat static complexity as runtime cost
- use average latency as the only service metric
- benchmark a toy input and generalize it to production
- compare measurements from different environments without qualification
- profile an idle process
- use one thread dump as definitive evidence
- call allocation a memory leak without retention evidence
- blame garbage collection without allocation and pause data
- blame the database without query and pool evidence
- enable verbose production logging without volume and privacy review
- expose diagnostic endpoints publicly
- add unbounded metric dimensions
- collect heap dumps without capacity and data-sensitivity review
- state that a profiler proves causation when the workload is unrelated
- report a suspected cause as confirmed
- hide failed hypotheses that materially affect the conclusion
- recommend a new tool when existing approved tools answer the question

## 30. Completion Criteria
Research is complete when:
1. the question is answered or explicitly remains unresolved
2. observations and interpretations are separated
3. the environment and workload are documented
4. the problem is reproducible or the limitation is stated
5. a verified cause has evidence, or the result remains labeled as a hypothesis
6. performance claims include a baseline and variance
7. improvements include an equivalent before-and-after comparison
8. correctness and secondary effects are checked
9. sensitive artifacts are handled safely
10. remaining uncertainty and follow-up work are explicit

## 31. Review Checklist
- Is the research question narrow and testable?
- Is the expected contract or performance target stated?
- Are the environment, version, dataset, and workload recorded?
- Are observations separated from suspected and verified causes?
- Was existing project telemetry used before adding instrumentation?
- Were records correlated with available request, time, version, and instance
  fields?
- Was the profiler selected for the actual symptom?
- Was the profile captured during the affected workload?
- Are logical problems tied to an invariant?
- Are performance claims based on percentiles and resource evidence?
- Are allocation and retention distinguished?
- Are CPU, blocking, queueing, and dependency latency distinguished?
- Was one primary variable changed in the experiment?
- Are before and after conditions equivalent?
- Is the result larger than normal variance?
- Were correctness and secondary costs validated?
- Is regression detection practical and stable?
- Are diagnostic access and artifacts handled securely?
- Are limitations and unresolved questions explicit?

## References
- [Go Diagnostics](https://go.dev/doc/diagnostics)
- [Go Execution Tracer](https://go.dev/doc/diagnostics#execution-tracer)
- [Go Data Race Detector](https://go.dev/doc/articles/race_detector)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Spring Boot Actuator Endpoints](https://docs.spring.io/spring-boot/reference/actuator/endpoints.html)
- [Spring Boot Metrics](https://docs.spring.io/spring-boot/reference/actuator/metrics.html)
- [Micrometer Documentation](https://docs.micrometer.io/micrometer/reference/)
- [JDK Flight Recorder](https://docs.oracle.com/en/java/javase/25/jfapi/flight-recorder.html)
- [JDK Mission Control](https://docs.oracle.com/en/java/java-components/jdk-mission-control/)
- [The `jcmd` Command](https://docs.oracle.com/en/java/javase/25/docs/specs/man/jcmd.html)
- [Profiling Node.js Applications](https://nodejs.org/learn/getting-started/profiling)
- [Node.js Memory Diagnostics](https://nodejs.org/learn/diagnostics/memory)
- [Node.js Flame Graphs](https://nodejs.org/learn/diagnostics/flame-graphs)
- [Node.js Inspector](https://nodejs.org/api/inspector.html)
- [Node.js Performance Hooks](https://nodejs.org/api/perf_hooks.html)
- [Chrome DevTools Performance](https://developer.chrome.com/docs/devtools/performance/)
- [Chrome DevTools Memory](https://developer.chrome.com/docs/devtools/memory-problems/)
