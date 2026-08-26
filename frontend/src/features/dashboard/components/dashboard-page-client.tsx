"use client";

import Link from "next/link";
import { ArrowRight, CreditCard, ShieldCheck, Users } from "lucide-react";
import { useAuthSession } from "@/features/auth/hooks/use-auth-session";

export function DashboardPageClient() {
  const { user, isLoading } = useAuthSession();

  if (isLoading) {
    return (
      <main className="mx-auto max-w-5xl px-6 py-20">Loading workspace…</main>
    );
  }

  if (!user) {
    return (
      <main className="mx-auto max-w-5xl px-6 py-20">
        <p className="text-sm tracking-[0.2em] text-teal-700 uppercase">
          Private workspace
        </p>
        <h1 className="mt-4 text-4xl font-semibold">
          Sign in to open your SaaS workspace.
        </h1>
        <Link
          href="/login"
          className="mt-8 inline-flex rounded-full bg-slate-950 px-6 py-3 text-sm font-semibold text-white"
        >
          Continue to login
        </Link>
      </main>
    );
  }

  return (
    <main className="mx-auto min-h-[calc(100vh-5rem)] max-w-6xl px-6 py-14">
      <p className="text-xs font-semibold tracking-[0.24em] text-teal-700 uppercase">
        Workspace overview
      </p>
      <h1 className="mt-4 max-w-3xl text-5xl font-semibold tracking-tight text-slate-950 md:text-6xl">
        Subscription state belongs beside the product decisions it controls.
      </h1>
      <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-600">
        Signed in as {user.full_name || user.email}. Your personal workspace is
        provisioned automatically and every billing mutation is checked against
        membership role.
      </p>
      <section className="mt-12 grid gap-px overflow-hidden rounded-[2rem] bg-slate-200 md:grid-cols-3">
        <article className="bg-white p-7">
          <Users className="text-teal-700" />
          <h2 className="mt-8 text-xl font-semibold">Workspace first</h2>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            Membership and tenant context are resolved before subscription data
            is returned.
          </p>
        </article>
        <article className="bg-white p-7">
          <CreditCard className="text-teal-700" />
          <h2 className="mt-8 text-xl font-semibold">Provider boundary</h2>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            Use the deterministic fake adapter locally, then switch to Stripe
            through environment configuration.
          </p>
        </article>
        <article className="bg-white p-7">
          <ShieldCheck className="text-teal-700" />
          <h2 className="mt-8 text-xl font-semibold">Webhook authority</h2>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            Signed, idempotent events—not redirect parameters—drive entitlement
            changes.
          </p>
        </article>
      </section>
      <Link
        href="/dashboard/billing"
        className="mt-10 inline-flex items-center gap-2 rounded-full bg-slate-950 px-6 py-3 text-sm font-semibold text-white transition hover:bg-slate-800"
      >
        Manage billing <ArrowRight size={17} />
      </Link>
    </main>
  );
}
