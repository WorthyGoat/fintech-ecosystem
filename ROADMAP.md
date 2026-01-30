# Product Roadmap

Strategic plan to evolve this repo into a **production-grade, open-source fintech platform** — a self-hosted alternative to Stripe — with clear phases for quality, growth, and sustainable monetization.

---

## Vision

Provide **developer-first, scalable, open-source financial infrastructure** that any team can run on their own cloud: payments, ledger, and webhooks with a small, clear scope and a path to hosted offerings and paid support.

---

## Versioned journey

| Version | Focus | Outcome |
|---------|--------|---------|
| **v0.x** | Foundation | Core primitives (payments, ledger, webhooks), docs, community standards. |
| **v1.0** | Quality & credibility | Unit/integration tests, idempotency, clean layering, contribution rules. Production-ready for self-host. |
| **v1.x** | Growth | Scale, observability, SDKs, more primitives. |
| **v2.x** | Monetization path | Hosted version, paid support, custom integrations for startups. |

The roadmap below is grouped by **phase** (Quality, Growth, Turn into services) and aligned with this versioning.

---

## Phase: Quality & credibility (v1.0)

*Goal: Trust and maintainability. Safe for production self-host and for contributors.*

### Testing and reliability

- [ ] **Unit tests for services** — High coverage for `internal/*` (payment, ledger, notification, auth). Focus on business logic first.
- [ ] **Table-driven tests** — Use Go table-driven tests for handlers and domain logic (inputs, expected outputs, edge cases).
- [ ] **Mock interfaces** — Define interfaces for repositories and external clients; inject mocks in tests. Keep domain independent of DB and messaging.
- [ ] **Idempotency keys for payments** — Accept `Idempotency-Key` on payment intent create/confirm; store and return cached response for duplicate requests so clients can safely retry.

### Architecture and data integrity

- [ ] **Never update balance directly** — Enforce “balance = sum of ledger entries” only. Remove any code path that updates a balance column without going through the ledger; document this rule in README and ADRs.
- [ ] **Separate API, domain, and infrastructure** — Clear layers: HTTP/gRPC handlers → application/domain logic → repositories and messaging. Keep `internal/` domain-free of framework and DB details where possible.

### Community and contribution

- [ ] **Contribution rules** — Document in CONTRIBUTING:
  - **Good first issue** label and a short list of starter tasks (docs, tests, small refactors).
  - **Commit style** (e.g. Conventional Commits: `feat:`, `fix:`, `docs:`, `test:`).
  - Link to PR template and code style (e.g. Uber Go Style Guide).

---

## Phase: Growth and long-term scale (v1.x)

*Goal: Scale, DX, and ecosystem. Prepare for wider adoption.*

### Platform and DX

- [ ] **SDK and API stability** — Versioned REST/OpenAPI; semantic versioning for SDKs (e.g. Node, Python). Changelog and upgrade guides.
- [ ] **More test coverage and CI** — Integration tests in CI; optional e2e with Docker Compose. Coverage gates for critical paths.
- [ ] **Performance and observability** — Benchmark critical paths (payment confirm, ledger write); dashboards and SLOs for latency and errors.

### Features and scale

- [ ] **Wallets as first-class primitive** — Explicit “wallet” (ledger account + metadata) and APIs if needed, still backed only by ledger entries.
- [ ] **Subscriptions and recurring payments** — Billing cycles, plans, and usage-based billing backed by payment intents and ledger.
- [ ] **Multi-tenant and rate limiting** — Per-API-key limits and quotas; tenant isolation for hosted offering later.
- [ ] **Documentation and examples** — Runnable examples (e.g. “Checkout flow”, “Connect split”), architecture decision records (ADRs).

### Roadmap hygiene

- [ ] **Living roadmap** — Keep ROADMAP.md updated: move completed items to a “Done” section, add new items from issues and community. Review quarterly.

---

## Phase: Turn it into services (v2.x / monetization)

*Goal: Sustainable open source and optional commercial offerings.*

### Hosted version

- [ ] **Hosted version path** — Document or offer a managed deployment (e.g. “Fintech Cloud”): we run the stack, you get API keys and dashboard. Clear comparison: self-host vs hosted.
- [ ] **Security and compliance** — SOC2-style controls, encryption at rest/transit, and compliance docs for hosted tier.

### Paid support and enterprise

- [ ] **Paid support** — Tiered support (e.g. community vs paid): SLAs, priority bugfixes, and security advisories for paying customers.
- [ ] **Custom integrations for startups** — Offer integration packages (e.g. “Stripe migration”, “marketplace setup”) as paid professional services; keep core open source unchanged.

### Licensing and sustainability

- [ ] **Clarify license and trademarks** — Keep core MIT; define trademark and “Fintech Ecosystem” branding for hosted/paid offerings so community and commercial use stay clear.
- [ ] **Funding and sponsors** — GitHub Sponsors, Open Collective, or sponsor tiers; link from README. Use funds for infra, security, and maintainer time.

---

## Completed (foundation)

- [x] Community standards: CONTRIBUTING, Code of Conduct, PR template.
- [x] CI: lint, unit/integration tests, Docker build.
- [x] Security: dependency scanning, secret scanning, API key hashing.
- [x] Kubernetes and Helm in `deploy/`.
- [x] Observability: tracing, metrics, structured logging.
- [x] Migrations, OAuth2/OIDC, scopes, webhook retries and signing.
- [x] Connect/marketplace: connected accounts, revenue splitting.
- [x] Enterprise: SSO, audit logs, team/RBAC.
- [x] Plugin system, fraud rules, multi-currency.

---

## Contributing

We welcome contributions. See [CONTRIBUTING.md](CONTRIBUTING.md) for good first issues, commit style, and development setup.
