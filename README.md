# Next.js + Go SaaS Billing Kit

A full-stack SaaS starter where subscription state controls real product
entitlements. It combines a Next.js App Router frontend, layered Go API,
PostgreSQL, signed Stripe webhooks, and a fake billing provider that works
without paid credentials.

## Included

- Email/password authentication and personal workspace provisioning
- Owner, admin, and member billing permissions
- Free, Pro, and Team plans with API and member quotas
- Stripe Checkout and customer portal session adapters
- Signed, idempotent webhook processing with stored event IDs
- Trialing, active, past-due, canceled, and grace-period handling
- Billing dashboard with usage and entitlement visibility
- PostgreSQL migrations and repository/service/handler separation
- OpenAPI contract with generated frontend types
- Unit, integration, frontend, and Playwright smoke-test foundations
- CI, secret scanning, workflow linting, and production Dockerfiles

## Architecture

```text
Next.js UI -> generated API contract -> Go handlers
                                      -> billing service
                                      -> fake or Stripe adapter
                                      -> PostgreSQL repositories

Stripe webhook -> signature check -> idempotency table -> subscription state
```

The provider boundary keeps local development deterministic. Entitlements are
calculated server-side and remain available during a past-due grace period;
the frontend only displays the server decision.

## Quick start

Requirements: Node.js 24+, Go 1.25+, and Docker.

```bash
npm install
npm run dev
```

The development command creates missing env files from examples, starts
PostgreSQL, then boots the frontend on `:3000` and API on `:8080`.
`BILLING_PROVIDER=fake` is the default, so no Stripe account is needed.

Demo account:

```text
demo@nextjs-go-kit.local
demo12345
```

After login, open `/dashboard/billing`. Fake checkout URLs return to that page
without charging a card.

## Commands

```bash
npm run dev              # full local stack
npm run dev:down         # stop PostgreSQL
npm run check            # frontend and backend quality gate
npm run check:contract   # verify generated OpenAPI types
npm run e2e              # production-mode Playwright smoke tests
npm run check:workflows  # actionlint
npm run check:secrets    # gitleaks
```

## Stripe mode

Set `BILLING_PROVIDER=stripe`, `STRIPE_SECRET_KEY`, and
`STRIPE_WEBHOOK_SECRET` through your secret manager. Configure Stripe prices
for `price_pro` and `price_team`, or replace those identifiers in the frontend.
Point signed events at `POST /api/v1/billing/webhooks/stripe`.

Never trust checkout redirects as proof of payment. Subscription state changes
only through verified webhooks. The adapter uses API timeouts, verifies the
HMAC and timestamp tolerance, stores event IDs before applying state, and
returns success for duplicate delivery.

## Environment and deployment

Backend variables are documented in `backend/.env.example`; frontend variables
are in `frontend/.env.example`. Production must replace database passwords,
JWT secrets, demo accounts, mail settings, and billing credentials.

`npm run check` runs formatting, linting, type checking, tests, and builds. The
backend and frontend Dockerfiles use production processes. Deploy behind TLS,
keep webhook raw bodies intact, and run migrations before serving traffic.

## License

MIT
