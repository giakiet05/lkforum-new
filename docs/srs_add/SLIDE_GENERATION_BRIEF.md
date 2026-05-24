# LKForum Slide Generation Brief For Gemini

## Task For Gemini

Create a 12-15 slide presentation for the LKForum Website project. The presentation duration is about 15 minutes.

The deck must show traceability from documents to implementation:

```text
BRD -> ASR -> ADD -> Code
```

Use the slide content below directly. Do not invent extra features. Do not paste long code blocks. For code evidence, show file paths and short implementation notes only.

## Input Documents

Use these PDFs as source material:

- `docs/srs_add/BRD/BRD.pdf`
- `docs/srs_add/ADD/ADD.pdf`
- `docs/srs_add/ADD/Utility Tree.pdf`

Use the code paths listed in this brief as implementation evidence.

## Global Slide Rules

- Slide 1 is only title, group, and members.
- From slide 2 onward, go directly into the technical content.
- The main story is not "we implemented all requirements".
- The main story is: "we selected architecture-significant ASRs from ADD and implemented them in code".
- Keep the deck between 12 and 15 slides.
- Do not copy full ADD tables into slides.
- Summarize ADD tables into compact traceability blocks.
- Each technical slide must include:
  - Requirement Trace
  - ADD Decision
  - Code Evidence
  - Verification, if available
- Code evidence must be file paths, not full code listings.
- Speaker notes must be included for every slide.

## Selected ASRs

These are the ASRs to present:

| Quality Attribute | ASR | BRD User Story Sources |
| --- | --- | --- |
| Performance | High-Throughput Content Retrieval | US-01 Browse posts; US-02 View posts/comments|
| Performance | Asynchronous Processing | US-23 Reply notifications; US-24 Mention notifications; US-27 Private messages|
| Security | Authentication and Authorization | US-07 Sign up; US-08 Google sign in; US-09 Local sign in; US-32 Assign moderator roles|
| Security | Abuse Prevention with Rate Limiting | US-09 Local sign in; US-13 Create text post; US-28 Report inappropriate content|
| Security | Data Privacy for Private Messaging | US-27 Send private messages|
| Supportability | Audit and User Activity Logging | US-28 Report content; US-33 Approve posts; US-38 Review reports; US-39 Ban users; US-40 Restrict communities|
| Supportability | Structured Error Diagnostics with slog | Cross-cutting operational support for all user flows|
| Manageability | Health and Readiness Checks | US-37 Monitor platform performance|
| Manageability | Runtime Metrics for Monitoring | US-37 Monitor platform performance|
| Maintainability | Layered Architecture and Interface-Based Design | Cross-cutting support for US-01 to US-40|

## Slide 1: LKForum Website

Put on slide:

- LKForum Website
- Software Architecture Presentation
- Group: `<fill group name>`
- Members: `<fill member names>`
- Lecturer: Nguyen Trinh Dong

Speaker notes:

- Introduce the team and project name.
- Do not explain architecture yet.

## Slide 2: Presentation Focus

Put on slide:

- Goal: prove traceability from requirements to architecture and code.
- Scope: selected architecture-significant requirements, not all documented features.
- Main chain: `BRD -> ASR -> ADD -> Code`.
- Implementation target: about 30% of important documented features.

Speaker notes:

- LKForum has many documented features, but the presentation focuses on the most important architectural parts.
- The selected ASRs cover performance, security, supportability, manageability, and maintainability.

## Slide 3: Requirement Traceability Flow

Put on slide as a simple vertical diagram:

```text
BRD
Business needs / user stories
        |
        v
ASR
Architecture-significant requirement
        |
        v
ADD
Architecture decision and tactics
        |
        v
Code
Implemented modules and runtime evidence
```

Speaker notes:

- BRD explains what users and business need.
- ASRs identify quality requirements that affect architecture.
- ADD explains the architecture decisions for those ASRs.
- Code evidence proves that the decisions are implemented.

## Slide 4: Selected ASR Overview

Put on slide as a compact table:

| Quality Attribute | ASR | Implementation Status |
| --- | --- | --- |
| Performance | High-Throughput Content Retrieval | Implemented |
| Performance | Asynchronous Processing | Partially implemented / architecture-ready |
| Security | Authentication and Authorization | Implemented |
| Security | Abuse Prevention with Rate Limiting | Implemented |
| Security | Data Privacy for Private Messaging | Implemented as E2EE demo |
| Supportability | Audit and User Activity Logging | Implemented |
| Supportability | Structured Error Diagnostics with slog | Implemented |
| Manageability | Health and Readiness Checks | Implemented |
| Manageability | Runtime Metrics for Monitoring | Implemented |
| Maintainability | Layered Architecture and Interface-Based Design | Implemented |

Also put this traceability matrix on the slide, or split it into a second compact table if the slide is crowded:

| ASR | BRD User Story Sources |
| --- | --- |
| High-Throughput Content Retrieval | US-01 Browse posts; US-02 View posts/comments|
| Asynchronous Processing | US-23 Reply notifications; US-24 Mention notifications; US-27 Private messages|
| Authentication and Authorization | US-07 Sign up; US-08 Google sign in; US-09 Local sign in; US-32 Assign moderator roles|
| Rate Limiting | US-09 Local sign in; US-13 Create text post; US-28 Report inappropriate content|
| Private Message Data Privacy | US-27 Send private messages|
| Audit Log | US-28 Report content; US-33 Approve posts; US-38 Review reports; US-39 Ban users; US-40 Restrict communities|
| Structured Logging | Cross-cutting operational support for all user flows|
| Health/Readiness | US-37 Monitor platform performance|
| Runtime Metrics | US-37 Monitor platform performance|
| Layered Architecture | Cross-cutting support for US-01 to US-40|

Speaker notes:

- This is the list of ASRs selected from ADD.
- Not every ASR needs a deep explanation during the main presentation.
- If the lecturer asks about any ASR, the team can open the corresponding source files listed in the following slides.

## Slide 5: Performance - High-Throughput Content Retrieval

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-01 Browse posts
- US-02 View posts/comments
ASR: High-Throughput Content Retrieval
ADD Response: The system serves feeds directly from Redis RAM using a Read-Through caching strategy, bypassing MongoDB for cached feed endpoints.
Code: Redis feed cache in post service
```

Architecture flow:

```text
Client -> Post Controller -> Post Service -> Redis Cache -> MongoDB fallback
```

Code Evidence:

- `backend/internal/service/post_service.go`
  - `GetPosts(...)` checks cache before querying MongoDB.
  - `feedCacheKey(...)` builds Redis cache keys.
  - `getCachedFeed(...)` reads cached feed responses.
  - `setCachedFeed(...)` writes feed responses to Redis.
  - `invalidateFeedCache(...)` clears cached feed data after mutations.
- `backend/internal/platform/metrics/metrics.go`
  - stores and renders metrics.
- `backend/internal/middleware/metrics.go`
  - records HTTP request metrics.

Verification:

- Feed endpoint was tested.
- Cache hit/miss behavior was observed through `/metrics`.
- Metrics include `lkforum_feed_cache_total`.

Speaker notes:

- Browsing posts is one of the highest-frequency operations in a social platform.
- The architecture uses Redis as a cache layer so repeated feed queries do not always hit MongoDB.
- This directly supports the performance ASR from ADD.

## Slide 6: Performance - Asynchronous Processing

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-23 Reply notifications
- US-24 Mention notifications
- US-27 Private messages
ASR: Asynchronous Processing
ADD Response: The Post Service saves the content and returns success quickly, then publishes an event through Go Channels so background goroutines can process secondary work such as notifications.
Code: Service layer and event-style components isolate side effects
```

Code Evidence:

- `backend/internal/service/reputation_service.go`
  - contains `processEvents(...)`, which runs event processing separately.
- `backend/internal/platform/ws/hub.go`
  - handles WebSocket delivery and active connection metrics outside normal REST controllers.
- `backend/internal/service/notification_service.go`
  - keeps notification behavior in a dedicated service.
- `backend/internal/controller`
  - controllers remain thin and delegate business work to services.
- `backend/internal/service`
  - service layer owns business logic and side effects.

Verification:

- Architecture is async-ready and has separated event/WebSocket/notification paths.
- Do not claim that a full message queue is implemented.

Speaker notes:

- This ASR is about responsiveness and separation of long-running or secondary work.
- The current implementation supports this through service separation, WebSocket hub behavior, and event-style processing.
- Be honest: this is not a full queue-based architecture yet.

## Slide 7: Security - Authentication And Authorization

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-07 Sign up
- US-08 Google sign in
- US-09 Local sign in
- US-32 Assign moderator roles
ASR: Authentication and Authorization
ADD Response: The system validates JWT signatures and enforces role-based access control before allowing restricted actions such as moderation or admin operations.
Code: Auth service, JWT utilities, auth middleware, protected routes
```

Code Evidence:

- `backend/internal/controller/auth_controller.go`
  - handles login, verification, password reset, and audit hooks.
- `backend/internal/service/auth_service.go`
  - validates credentials and issues tokens.
- `backend/internal/auth/jwt.go`
  - generates and validates JWT access/refresh tokens.
- `backend/internal/auth/token_service.go`
  - token service abstraction.
- `backend/internal/middleware/auth.go`
  - authenticates requests and injects auth user context.
- `backend/internal/route/auth_route.go`
  - defines auth endpoints.
- Protected route files:
  - `backend/internal/route/post_route.go`
  - `backend/internal/route/community_route.go`
  - `backend/internal/route/message_route.go`
  - `backend/internal/route/admin_user_route.go`

Verification:

- Login flow was tested.
- Protected routes require authenticated context.

Speaker notes:

- Authentication identifies the user.
- Authorization decides whether the user can perform a protected action.
- LKForum implements this through JWT and route middleware.

## Slide 8: Security - Abuse Prevention With Rate Limiting

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-09 Local sign in
- US-13 Create text post
- US-28 Report inappropriate content
ASR: Abuse Prevention with Rate Limiting
ADD Response: The system increments Redis counters scoped by endpoint and client identifier; requests within limit continue, while excessive requests are rejected before business logic.
Code: Redis-backed rate limit middleware
```

Code Evidence:

- `backend/internal/middleware/rate_limit.go`
  - `RateLimit(...)` builds Redis keys using scope and client identifier.
  - increments request counters.
  - applies TTL windows.
  - returns HTTP 429 when limit is exceeded.
- `backend/internal/route/auth_route.go`
  - applies rate limits to login, OTP, password reset, and verification endpoints.
- `backend/internal/route/post_route.go`
  - applies rate limits to post creation and report endpoints.

Verification:

- Rate limiting middleware is wired into high-risk routes.
- Redis is used as shared state for counters and time windows.

Speaker notes:

- Social platforms need protection against brute-force login, OTP spam, post spam, and report spam.
- Rate limiting is a security tactic that reduces repeated abusive requests.
- Redis is used because counters need fast shared storage.

## Slide 9: Security - Data Privacy For Private Messaging

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-27 Send private messages one-on-one
ASR: Data Privacy for Private Messaging
ADD Response: The sender client encrypts the message before sending it over WebSocket. The backend stores and routes only ciphertext, nonce, algorithm, and key version metadata.
Code: Frontend E2EE demo + backend encrypted message fields + WebSocket delivery
```

Architecture flow:

```text
Sender Frontend
  -> AES-GCM encrypt
  -> send ciphertext + nonce + algorithm
Backend
  -> store encrypted fields
Receiver Frontend
  -> fetch/receive encrypted message
  -> decrypt locally
```

Code Evidence:

- `frontend/user-web/src/services/e2ee-service.ts`
  - encrypts and decrypts private message content.
  - detects encrypted message payloads.
- `frontend/user-web/src/services/message-service.ts`
  - sends encrypted payloads.
  - decrypts fetched messages.
- `frontend/user-web/src/services/websocket-service.ts`
  - decrypts incoming WebSocket messages.
- `frontend/user-web/src/App.svelte`
  - handles global WebSocket message flow.
- `backend/internal/route/websocket_route.go`
  - exposes backend WebSocket endpoint.
- `backend/internal/platform/ws/hub.go`
  - manages WebSocket clients and message delivery.
- `backend/internal/model/message.go`
  - includes `ciphertext`, `nonce`, `algorithm`, and `key_version`.
- `backend/internal/dto/message_dto.go`
  - exposes encrypted message fields in API DTOs.
- `backend/internal/dto/ws_message_dto.go`
  - exposes encrypted message fields for WebSocket payloads.
- `backend/internal/service/message_service.go`
  - persists and returns encrypted message fields.

Verification:

- Encrypted message test passed.
- Stored message content did not contain plaintext.
- Backend stored `ciphertext`, `nonce`, and `algorithm = AES-GCM`.
- REST and WebSocket flows handled encrypted messages.

Limitation:

- This is demo-level E2EE for architecture validation.
- It is not production-grade Signal/WhatsApp-style key management.

Speaker notes:

- The core architectural idea is that the backend does not need plaintext private message content.
- The frontend encrypts before sending and decrypts after receiving.
- The backend stores encrypted fields only.
- Be transparent that key management is simplified for demo purposes.

## Slide 10: Supportability - Audit And User Activity Logging

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-28 Report content
- US-33 Approve posts
- US-38 Review reports
- US-39 Ban users
- US-40 Restrict communities
ASR: Audit and User Activity Logging
ADD Response: The system records persistent audit entries for critical state-changing actions, including actor, role, action, target, reason, IP address, user agent, metadata, and timestamp.
Code: Audit log model, repository, service, middleware, and controller hooks
```

Code Evidence:

- `backend/internal/model/audit_log.go`
  - defines audit log fields.
- `backend/internal/repo/audit_log_repo.go`
  - persists audit logs to MongoDB.
- `backend/internal/service/audit_log_service.go`
  - records audit logs asynchronously from request flow.
- `backend/internal/middleware/audit.go`
  - provides `RecordAudit(...)`.
- `backend/internal/bootstrap/init.go`
  - wires audit log repository/service into the app.
- `backend/internal/controller/auth_controller.go`
  - records login success/failure and verification events.
- `backend/internal/controller/admin_user_controller.go`
  - records admin ban/unban actions.
- `backend/internal/controller/admin_community_controller.go`
  - records community ban/unban actions.
- `backend/internal/controller/community_controller.go`
  - records moderation decisions.
- `backend/internal/controller/report_controller.go`
  - records report deletion/batch deletion actions.

Verification:

- Audit logs were persisted to MongoDB during testing.
- Runtime logs showed `audit_log_record_requested` and `audit_log_persisted`.

Speaker notes:

- Audit log is different from normal application log.
- Normal logs help developers debug runtime behavior.
- Audit logs preserve important user/admin actions for accountability and investigation.

## Slide 11: Supportability - Structured Logging With slog

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-37 View dashboard with system analytics
- Cross-cutting operational support for all user flows
ASR: Structured Error Diagnostics with slog
ADD Response: The backend emits structured diagnostic logs with stable fields such as action, user id, target id, request path, and error details while keeping user-facing errors generic when needed.
Code: Go slog logger and request logging middleware
```

Code Evidence:

- `backend/internal/logging/logger.go`
  - configures JSON slog logger.
  - sets default logger for the application.
- `backend/internal/middleware/request_logger.go`
  - logs HTTP method, path, status, latency, client IP, user agent, and errors.
- `backend/main.go`
  - calls `logging.ConfigureDefault()`.
  - logs server startup and fatal server errors.
- `backend/internal/config/config.go`
  - logs configuration loading.
- `backend/internal/config/mongo.go`
  - logs MongoDB connection and collection verification.
- `backend/internal/config/redis.go`
  - logs Redis connection status.
- `backend/internal/service/audit_log_service.go`
  - logs audit persist success/failure.
- `backend/internal/service/post_service.go`
  - logs feed cache hit/miss/invalidation.

Example log fields:

- `method`
- `path`
- `status`
- `latency`
- `request_id`
- `error`

Verification:

- Structured request logs were observed.
- Audit-related structured logs were observed.
- Cache-related structured logs were observed.

Speaker notes:

- Plain text logs are hard to search and aggregate.
- Structured logs make production diagnosis easier because each field can be queried.
- Go `slog` is a standard library logging package, so it fits well in a Go backend.

## Slide 12: Manageability - Health, Readiness, And Metrics

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- US-37 View dashboard with system analytics
ASR: Health and Readiness Checks + Runtime Metrics
ADD Response: The backend exposes `/health`, `/ready`, and Prometheus-compatible `/metrics`; readiness checks MongoDB and Redis, while metrics expose counters, gauges, and durations for runtime monitoring.
Code: /health, /ready, /metrics
```

Code Evidence:

- `backend/internal/bootstrap/init.go`
  - registers `/health`.
  - registers `/ready`.
  - registers `/metrics`.
  - readiness checks MongoDB and Redis.
- `backend/internal/platform/metrics/metrics.go`
  - stores counters, gauges, and durations.
  - renders Prometheus-style text metrics.
- `backend/internal/middleware/metrics.go`
  - records HTTP request counts and durations.
- `backend/internal/platform/ws/hub.go`
  - updates WebSocket active connection metrics.
- `backend/internal/service/message_service.go`
  - increments sent-message metrics.
- `backend/internal/service/community_service.go`
  - increments moderation decision metrics.
- `monitoring/prometheus.yml`
  - Prometheus scrape configuration for LKForum backend.

Verification:

- `/health` returned OK.
- `/ready` checked MongoDB and Redis.
- `/metrics` exposed request, cache, WebSocket, message, and moderation metrics.

Speaker notes:

- Health check means the process is alive.
- Readiness check means dependencies are available and the service can handle traffic.
- Metrics support monitoring and debugging after deployment.

## Slide 13: Maintainability - Layered Architecture

Put on slide:

Requirement Trace:

```text
BRD User Story Sources:
- Cross-cutting support for US-01 to US-40
ASR: Layered Architecture and Interface-Based Design
ADD Response: The backend separates routing, request handling, business rules, data access, and infrastructure wiring so new capabilities can be added in the proper layer with limited unrelated changes.
Code: controller/service/repository/model/dto package structure
```

Architecture flow:

```text
Routes
  -> Controllers
  -> Services
  -> Repositories
  -> MongoDB / Redis

DTOs define API boundaries.
Models define stored domain data.
Interfaces define contracts between layers.
```

Code Evidence:

- `backend/internal/route`
  - maps HTTP endpoints to controllers.
- `backend/internal/controller`
  - handles HTTP input/output.
- `backend/internal/service`
  - owns business logic.
- `backend/internal/repo`
  - owns database access.
- `backend/internal/model`
  - defines persisted domain models.
- `backend/internal/dto`
  - defines request/response DTOs.
- `backend/internal/bootstrap/init.go`
  - wires dependencies.

Speaker notes:

- Maintainability matters because the system has many modules and user stories.
- Layering keeps HTTP logic, business logic, and persistence logic separate.
- This makes future changes and testing easier.

## Slide 14: Verification Summary

Put on slide as a table:

| Feature / ASR | Verification Result |
| --- | --- |
| Health endpoint | `/health` tested OK |
| Readiness endpoint | `/ready` tested OK, checked MongoDB and Redis |
| Runtime metrics | `/metrics` tested OK |
| Feed cache | Cache hit/miss behavior observed |
| Audit log | Audit records persisted to MongoDB |
| Structured logging | slog request/audit/cache logs observed |
| Rate limiting | Middleware wired to auth and post routes |
| E2EE demo | Encrypted storage and REST/WebSocket flow tested |
| Frontend build | Build passed |

Add note:

- `npm run check` has unrelated old errors in `Settings_OLD.svelte`.

Speaker notes:

- This slide proves that the selected ASRs are not only written in ADD but also implemented and verified.
- The implementation evidence is spread across backend services, middleware, routes, and selected frontend services.

## Slide 15: Conclusion

Put on slide:

- LKForum does not implement every requirement from BRD.
- The project implements selected architecture-significant ASRs from ADD.
- The implemented ASRs cover:
  - Performance
  - Security
  - Supportability
  - Manageability
  - Maintainability
- The selected implementation satisfies the requirement of implementing about 30% of important documented features.
- The trace is clear: `BRD -> ASR -> ADD -> Code`.

Speaker notes:

- End with the main claim: the selected ADD content and codebase match.
- The project focuses on important architecture quality attributes instead of trying to demo every UI feature.
- If asked, open the source files listed in the slide deck as direct implementation evidence.

## Optional Backup Slide Rule

If Gemini needs backup slides, create them after Slide 15 only.

Backup slides should not count as main presentation slides.

Recommended backup slides:

- Backup: E2EE source code evidence
- Backup: Audit log source code evidence
- Backup: Metrics endpoints and Prometheus output
- Backup: Rate limiting routes

Each backup slide should show only file paths and 1-line implementation notes.
