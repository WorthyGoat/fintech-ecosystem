# Product Roadmap

Strategic plan to evolve this repo into a **production-grade, open-source fintech platform** — a self-hosted alternative to Stripe — with clear phases for quality, growth, and sustainable monetization.

---

## Vision

Provide **developer-first, scalable, open-source financial infrastructure** that any team can run on their own cloud: payments, ledger, and webhooks with a small, clear scope and a path to hosted offerings and paid support.

---

## Versioned Journey

| Version | Focus | Outcome |
|---------|--------|---------|
| **v0.x** | Foundation | Core primitives (payments, ledger, webhooks), docs, community standards. |
| **v1.0** | Quality & Credibility | Unit/integration tests, idempotency, clean layering, contribution rules. Production-ready for self-host. |
| **v1.x** | Growth & Scale | Scale, observability, SDKs, more primitives. |
| **v2.x** | Services & Monetization | Hosted version, paid support, custom integrations for startups. |

---

## Phase: Quality & Credibility (v1.0)
*Goal: Trust and maintainability. Ensure the platform is safe for production use and easy for contributors to join.*

### 🛠 Reliability & Testing
- [ ] **Unit Tests for Services** — Achieve high coverage for `internal/` (payment, ledger, auth). Focus on core business logic.
- [ ] **Table-Driven Tests** — Implement Go table-driven tests for all handlers and domain logic to cover edge cases efficiently.
- [ ] **Mock Interfaces** — Extract interfaces for repositories and external clients (Redis, Kafka) to allow robust unit testing without dependencies.
- [ ] **Idempotency Keys** — Implement `Idempotency-Key` support for payment creation and confirmation to handle retries safely.

### 🏗 Architecture & Integrity
- [ ] **Ledger-Only Balance Updates** — Remove any code path updating balances directly. Enforce "balance = sum(entries)" as the single source of truth.
- [ ] **Layered Separation** — Clearly separate API (HTTP/gRPC), Domain (business logic), and Infrastructure (DB/Messaging). Keep the domain logic pure and framework-agnostic.

---

## Phase: Growth and Long-Term Scale (v1.x)
*Goal: Expand the ecosystem, improve developer experience (DX), and prepare for high-volume traffic.*

### 🚀 Platform & DX
- [ ] **SDKs & API Stability** — Release official SDKs (Node, Python, Go) and maintain a stable, versioned REST/OpenAPI spec.
- [ ] **Advanced Observability** — Implement detailed dashboards and SLOs for latency and error rates across all services.
- [ ] **Wallets as a First-Class Primitive** — Add dedicated APIs for wallet management (top-ups, transfers), still backed by the ledger.

### 📈 Features & Scale
- [ ] **Subscriptions & Billing** — Build recurring payment logic on top of the existing payment and ledger primitives.
- [ ] **Multi-tenancy & Rate Limiting** — Add tenant isolation and per-API-key quotas to support managed hosting environments.

---

## Phase: Turn it into Services (v2.x)
*Goal: Sustainable open source through optional commercial offerings.*

### ☁️ Managed Offerings
- [ ] **Hosted Version (Fintech Cloud)** — Offer a managed deployment path where we handle infrastructure, security, and updates.
- [ ] **Enterprise Compliance** — SOC2/PCI-DSS compliance documentation and hardened security controls for the hosted tier.

### 💼 Commercial Support
- [ ] **Paid Support & SLAs** — Offer tiered support packages for companies requiring guaranteed uptime and priority bugfixes.
- [ ] **Custom Integrations** — Provide professional services for complex migrations (e.g., from Stripe) and bespoke marketplace setups.

---

## Completed (Foundation)
- [x] Core Primitives: Payments, Ledger, Webhooks.
- [x] Community: CONTRIBUTING, Code of Conduct, PR Templates.
- [x] Infrastructure: Kubernetes/Helm, Docker Compose, CI/CD.
- [x] Security: API Key hashing, OAuth2/OIDC, Scopes.
- [x] Advanced Features: Connect/Marketplace, RBAC, Webhook Signing.

---

## Contributing
We welcome contributions. See [CONTRIBUTING.md](CONTRIBUTING.md) for good first issues, commit style, and development setup.
