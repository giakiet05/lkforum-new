# ADD To Code Trace Map

File này dùng để mở source code nhanh khi bị hỏi từng ASR trong ADD.

Main trace format:

```text
BRD User Story -> ASR in ADD -> Code evidence
```

## 1. High-Throughput Content Retrieval

Quality attribute: Performance

BRD sources:

- US-01 Browse posts
- US-02 View posts/comments

Open these files:

- `backend/internal/service/post_service.go`
  - `GetPosts(...)`: checks cache before database lookup.
  - `feedCacheKey(...)`: builds Redis feed cache key.
  - `getCachedFeed(...)`: reads cached feed response.
  - `setCachedFeed(...)`: stores feed response in Redis.
  - `invalidateFeedCache(...)`: clears cached feed data after post mutations.
- `backend/internal/platform/metrics/metrics.go`
  - metrics registry and Prometheus-style renderer.
- `backend/internal/middleware/metrics.go`
  - records HTTP request metrics.

Explain:

- Browse/read-heavy feed requests are cached with Redis.
- Cache hit/miss metrics prove cache behavior at runtime.

Verification:

- `/api/posts?page=1&limit=3` was tested.
- `/metrics` showed feed cache hit/miss metrics.

## 2. Asynchronous Processing

Quality attribute: Performance

BRD sources:

- US-23 Reply notifications
- US-24 Mention notifications
- US-27 Private messages

Open these files:

- `backend/internal/service/reputation_service.go`
  - `processEvents(...)`: event processing loop running separately.
- `backend/internal/platform/ws/hub.go`
  - WebSocket hub manages realtime client delivery.
  - updates active connection metrics.
- `backend/internal/service/notification_service.go`
  - notification logic is isolated in service layer.
- `backend/internal/controller`
  - controllers delegate work to service layer.
- `backend/internal/service`
  - business logic and side effects are kept outside controllers.

Explain:

- Current code separates request handling, WebSocket delivery, notification logic, and event-style processing.
- This is architecture-ready async processing.
- Do not claim that a full external message queue is implemented.

Verification:

- WebSocket and notification paths are separated from normal REST controllers.

## 3. Authentication And Authorization

Quality attribute: Security

BRD sources:

- US-07 Sign up
- US-08 Google sign in
- US-09 Local sign in
- US-32 Assign moderator roles

Open these files:

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
- `backend/internal/route/post_route.go`
  - protected post actions.
- `backend/internal/route/community_route.go`
  - protected community actions.
- `backend/internal/route/message_route.go`
  - protected message actions.
- `backend/internal/route/admin_user_route.go`
  - admin-only user management routes.

Explain:

- Authentication identifies the user through login/JWT.
- Authorization protects routes and role-sensitive actions.

Verification:

- Login flow was tested.
- Protected routes require auth middleware/user context.

## 4. Abuse Prevention With Rate Limiting

Quality attribute: Security

BRD sources:

- US-09 Local sign in
- US-13 Create text post
- US-28 Report inappropriate content

Open these files:

- `backend/internal/middleware/rate_limit.go`
  - `RateLimit(...)`: creates Redis-backed request counters.
  - returns HTTP 429 when limit is exceeded.
- `backend/internal/route/auth_route.go`
  - rate limits login, OTP, verification, forgot-password flows.
- `backend/internal/route/post_route.go`
  - rate limits post creation and report submission.

Explain:

- Login, OTP, post creation, and report endpoints can be abused.
- Redis-backed counters enforce request limits per scope/client.

Verification:

- Middleware is wired into high-risk auth/post routes.

## 5. Data Privacy For Private Messaging

Quality attribute: Security

BRD source:

- US-27 Send private messages

Open these files:

- `frontend/user-web/src/services/e2ee-service.ts`
  - encrypts/decrypts private message content.
  - detects encrypted message payloads.
- `frontend/user-web/src/services/message-service.ts`
  - sends encrypted payloads.
  - decrypts fetched messages.
- `frontend/user-web/src/services/websocket-service.ts`
  - decrypts incoming WebSocket messages.
- `frontend/user-web/src/App.svelte`
  - global WebSocket message flow.
- `backend/internal/route/websocket_route.go`
  - exposes backend WebSocket endpoint.
- `backend/internal/platform/ws/hub.go`
  - manages WebSocket clients and message delivery.
- `backend/internal/model/message.go`
  - encrypted fields: `ciphertext`, `nonce`, `algorithm`, `key_version`.
- `backend/internal/dto/message_dto.go`
  - encrypted message fields in REST DTO.
- `backend/internal/dto/ws_message_dto.go`
  - encrypted message fields in WebSocket DTO.
- `backend/internal/service/message_service.go`
  - persists and returns encrypted message fields.

Explain:

- Frontend encrypts message content before sending.
- Backend stores/routes ciphertext and metadata, not plaintext.
- Receiver frontend decrypts locally.
- Current E2EE is a demo-level architecture validation, not production-grade key management.

Verification:

- Runtime test confirmed stored message content did not contain plaintext.
- Backend stored `ciphertext`, `nonce`, and `algorithm = AES-GCM`.
- REST and WebSocket encrypted flows were tested.

## 6. Audit And User Activity Logging

Quality attribute: Supportability

BRD sources:

- US-28 Report content
- US-33 Approve posts
- US-38 Review reports
- US-39 Ban users
- US-40 Restrict communities

Open these files:

- `backend/internal/model/audit_log.go`
  - audit log data model.
- `backend/internal/repo/audit_log_repo.go`
  - MongoDB persistence.
- `backend/internal/service/audit_log_service.go`
  - records audit entries.
- `backend/internal/middleware/audit.go`
  - `RecordAudit(...)` helper.
- `backend/internal/bootstrap/init.go`
  - wires audit repo/service into app.
- `backend/internal/controller/auth_controller.go`
  - login success/failure, verification events.
- `backend/internal/controller/admin_user_controller.go`
  - user ban/unban audit events.
- `backend/internal/controller/admin_community_controller.go`
  - community ban/unban audit events.
- `backend/internal/controller/community_controller.go`
  - moderation decision audit events.
- `backend/internal/controller/report_controller.go`
  - report delete/batch delete audit events.

Explain:

- Audit log records important actions for accountability.
- Fields include actor, action, target, reason, metadata, timestamp.
- Audit log is different from normal debug/application log.

Verification:

- Audit records were persisted to MongoDB.
- Runtime logs showed `audit_log_record_requested` and `audit_log_persisted`.

## 7. Structured Error Diagnostics With slog

Quality attribute: Supportability

BRD source:

- Cross-cutting operational support for all user flows.
- Related to US-37 Monitor platform performance.

Open these files:

- `backend/internal/logging/logger.go`
  - configures JSON slog logger.
  - sets default logger.
- `backend/internal/middleware/request_logger.go`
  - request logging middleware.
- `backend/main.go`
  - configures default logger and logs startup/fatal errors.
- `backend/internal/config/config.go`
  - logs configuration loading.
- `backend/internal/config/mongo.go`
  - logs MongoDB connection and verification.
- `backend/internal/config/redis.go`
  - logs Redis connection status.
- `backend/internal/service/audit_log_service.go`
  - logs audit persistence success/failure.
- `backend/internal/service/post_service.go`
  - logs feed cache hit/miss/invalidation.

Explain:

- `slog` creates structured JSON logs.
- Logs include machine-readable fields such as method, path, status, latency, request id, and error.
- This improves debugging and production diagnosis.

Verification:

- Structured request logs were observed.
- Audit and cache structured logs were observed.

## 8. Health And Readiness Checks

Quality attribute: Manageability

BRD source:

- US-37 Monitor platform performance

Open these files:

- `backend/internal/bootstrap/init.go`
  - registers `/health`.
  - registers `/ready`.
  - `/ready` checks MongoDB and Redis.
- `backend/internal/config/mongo.go`
  - MongoDB connection dependency.
- `backend/internal/config/redis.go`
  - Redis connection dependency.

Explain:

- `/health` shows whether the backend process is alive.
- `/ready` shows whether dependencies are available so the service can handle traffic.

Verification:

- `/health` returned OK.
- `/ready` returned dependency status for MongoDB and Redis.

## 9. Runtime Metrics For Monitoring

Quality attribute: Manageability

BRD source:

- US-37 Monitor platform performance

Open these files:

- `backend/internal/bootstrap/init.go`
  - registers `/metrics`.
- `backend/internal/platform/metrics/metrics.go`
  - counters, gauges, durations, and Prometheus-style text rendering.
- `backend/internal/middleware/metrics.go`
  - HTTP request count and duration metrics.
- `backend/internal/service/post_service.go`
  - feed cache hit/miss/invalidation metrics.
- `backend/internal/platform/ws/hub.go`
  - active WebSocket connection metrics.
- `backend/internal/service/message_service.go`
  - sent-message metrics.
- `backend/internal/service/community_service.go`
  - moderation decision metrics.
- `monitoring/prometheus.yml`
  - sample Prometheus scrape config.

Explain:

- The backend exposes Prometheus-compatible metrics at `/metrics`.
- Prometheus server can scrape this endpoint using `monitoring/prometheus.yml`.
- The project does not use the official Prometheus Go client library.
- The project does not require running a full Prometheus/Grafana stack for demo.

Verification:

- `/metrics` returned Prometheus-style metrics.
- Observed metrics included HTTP request metrics and feed cache hit/miss.

## 10. Layered Architecture And Interface-Based Design

Quality attribute: Maintainability

BRD source:

- Cross-cutting support for US-01 to US-40.

Open these directories/files:

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
  - dependency wiring.
- `backend/internal/test/mocks/mock_gen.go`
  - mock generation entry point for interface-based testing.
- `backend/internal/repo/mocks`
  - generated repository mocks.

Explain:

- Routes, controllers, services, repositories, DTOs, and models are separated.
- This reduces coupling and keeps changes localized.
- Interfaces/mocks make service-level testing easier.

Verification:

- Codebase structure follows layered backend organization.
- Existing tests use service/repository abstractions and mocks.

