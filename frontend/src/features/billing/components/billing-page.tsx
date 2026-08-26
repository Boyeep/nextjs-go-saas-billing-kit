"use client";

import { useEffect, useState } from "react";
import {
  ArrowUpRight,
  Check,
  CircleAlert,
  CreditCard,
  Gauge,
} from "lucide-react";
import { useAuthStore } from "@/features/auth/store/auth-store";
import {
  createCheckout,
  createPortal,
  getBilling,
  listWorkspaces,
  type BillingSummary,
} from "@/features/billing/services/billing-service";

const plans = [
  {
    id: "free",
    name: "Free",
    price: "$0",
    copy: "For validating the workflow locally.",
    quota: "1,000 requests",
  },
  {
    id: "pro",
    name: "Pro",
    price: "$29",
    copy: "For a focused product team shipping weekly.",
    quota: "100,000 requests",
  },
  {
    id: "team",
    name: "Team",
    price: "$99",
    copy: "For larger workspaces with shared ownership.",
    quota: "1,000,000 requests",
  },
] as const;

export function BillingPage() {
  const token = useAuthStore((state) => state.accessToken);
  const [summary, setSummary] = useState<BillingSummary | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState<string | null>(null);

  useEffect(() => {
    if (!token) return;
    void listWorkspaces(token)
      .then((items) => {
        if (!items[0])
          throw new Error("No workspace is available for this account.");
        return getBilling(token, items[0].id);
      })
      .then(setSummary)
      .catch((reason: unknown) =>
        setError(
          reason instanceof Error
            ? reason.message
            : "Billing could not be loaded.",
        ),
      );
  }, [token]);

  async function checkout(plan: string) {
    if (!token || !summary) return;
    setBusy(plan);
    setError(null);
    try {
      const session = await createCheckout(
        token,
        summary.workspace.id,
        `price_${plan}`,
      );
      window.location.assign(session.url);
    } catch (reason) {
      setError(
        reason instanceof Error
          ? reason.message
          : "Checkout could not be created.",
      );
      setBusy(null);
    }
  }

  async function portal() {
    if (!token || !summary) return;
    setBusy("portal");
    try {
      const session = await createPortal(token, summary.workspace.id);
      window.location.assign(session.url);
    } catch (reason) {
      setError(
        reason instanceof Error
          ? reason.message
          : "Portal could not be opened.",
      );
      setBusy(null);
    }
  }

  if (!token)
    return (
      <main className="mx-auto max-w-5xl px-6 py-20">
        <p>Sign in to manage workspace billing.</p>
      </main>
    );
  if (!summary && !error)
    return (
      <main className="mx-auto max-w-5xl px-6 py-20">
        <p>Loading subscription…</p>
      </main>
    );

  return (
    <main className="mx-auto min-h-[calc(100vh-5rem)] max-w-6xl px-6 py-12">
      <header className="flex flex-col gap-6 border-b border-slate-200 pb-9 md:flex-row md:items-end md:justify-between">
        <div>
          <p className="mb-3 text-xs font-semibold tracking-[0.24em] text-teal-700 uppercase">
            Workspace billing
          </p>
          <h1 className="max-w-2xl text-4xl font-semibold tracking-tight text-slate-950 md:text-6xl">
            Plans tied to real entitlements.
          </h1>
        </div>
        {summary?.subscription.provider_customer_id && (
          <button
            onClick={() => void portal()}
            disabled={busy !== null}
            className="inline-flex h-11 items-center gap-2 rounded-full bg-slate-950 px-5 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:opacity-50"
          >
            <CreditCard size={17} /> Customer portal <ArrowUpRight size={16} />
          </button>
        )}
      </header>

      {error && (
        <div
          role="alert"
          className="mt-6 flex items-center gap-3 rounded-2xl bg-red-50 px-4 py-3 text-sm text-red-800"
        >
          <CircleAlert size={18} />
          {error}
        </div>
      )}

      {summary && (
        <section className="grid gap-8 py-9 lg:grid-cols-[1fr_18rem]">
          <div className="grid gap-px overflow-hidden rounded-[2rem] bg-slate-200 md:grid-cols-3">
            {plans.map((plan) => {
              const current = summary.subscription.plan === plan.id;
              return (
                <article
                  key={plan.id}
                  className="flex min-h-80 flex-col bg-white p-7"
                >
                  <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold">{plan.name}</h2>
                    {current && (
                      <span className="rounded-full bg-teal-50 px-3 py-1 text-xs font-semibold text-teal-800">
                        Current
                      </span>
                    )}
                  </div>
                  <p className="mt-5 text-4xl font-semibold">
                    {plan.price}
                    <span className="text-sm font-normal text-slate-500">
                      {" "}
                      / month
                    </span>
                  </p>
                  <p className="mt-5 text-sm leading-6 text-slate-600">
                    {plan.copy}
                  </p>
                  <p className="mt-6 flex items-center gap-2 text-sm">
                    <Check size={16} className="text-teal-700" />
                    {plan.quota}
                  </p>
                  <button
                    disabled={current || plan.id === "free" || busy !== null}
                    onClick={() => void checkout(plan.id)}
                    className="mt-auto h-11 rounded-full border border-slate-900 px-4 text-sm font-semibold transition hover:bg-slate-950 hover:text-white disabled:border-slate-200 disabled:text-slate-400"
                  >
                    {busy === plan.id
                      ? "Opening…"
                      : current
                        ? "Selected"
                        : `Choose ${plan.name}`}
                  </button>
                </article>
              );
            })}
          </div>
          <aside className="space-y-6">
            <div>
              <p className="text-sm text-slate-500">Workspace</p>
              <p className="mt-1 font-semibold">{summary.workspace.name}</p>
              <p className="text-sm text-slate-500 capitalize">
                {summary.membership_role} access
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-500">Subscription</p>
              <p className="mt-1 font-semibold capitalize">
                {summary.subscription.status.replace("_", " ")}
              </p>
              {summary.subscription.grace_period_end && (
                <p className="mt-1 text-xs text-amber-700">
                  Grace period is active
                </p>
              )}
            </div>
            <div>
              <p className="flex items-center gap-2 text-sm text-slate-500">
                <Gauge size={16} /> Monthly usage
              </p>
              <p className="mt-2 text-3xl font-semibold">
                {summary.usage.api_requests.toLocaleString()}
              </p>
              <p className="text-sm text-slate-500">
                of {summary.entitlements.api_requests.toLocaleString()} requests
              </p>
            </div>
          </aside>
        </section>
      )}
    </main>
  );
}
