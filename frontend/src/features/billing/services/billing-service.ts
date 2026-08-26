import { fetchJson } from "@/lib/api";

export type Workspace = { id: string; name: string; slug: string };
export type BillingSummary = {
  workspace: Workspace;
  membership_role: "owner" | "admin" | "member";
  subscription: {
    provider_customer_id?: string;
    plan: "free" | "pro" | "team";
    status: "trialing" | "active" | "past_due" | "canceled";
    current_period_end?: string;
    grace_period_end?: string;
    cancel_at_period_end: boolean;
  };
  entitlements: Record<string, number>;
  usage: Record<string, number>;
};

const authorized = (token: string) => ({ Authorization: `Bearer ${token}` });

export function listWorkspaces(token: string) {
  return fetchJson<Workspace[]>("/workspaces", { headers: authorized(token) });
}

export function getBilling(token: string, workspaceId: string) {
  return fetchJson<BillingSummary>(`/workspaces/${workspaceId}/billing`, {
    headers: authorized(token),
  });
}

export function createCheckout(
  token: string,
  workspaceId: string,
  priceId: string,
) {
  return fetchJson<{ url: string }>(
    `/workspaces/${workspaceId}/billing/checkout`,
    {
      method: "POST",
      headers: authorized(token),
      body: JSON.stringify({ price_id: priceId }),
    },
  );
}

export function createPortal(token: string, workspaceId: string) {
  return fetchJson<{ url: string }>(
    `/workspaces/${workspaceId}/billing/portal`,
    {
      method: "POST",
      headers: authorized(token),
      body: JSON.stringify({ return_url: window.location.href }),
    },
  );
}
