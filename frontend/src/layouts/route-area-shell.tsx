import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

type RouteAreaTone = "site" | "auth" | "workspace";

type RouteAreaShellProps = {
  tone: RouteAreaTone;
  children: ReactNode;
};

export function RouteAreaShell({ tone, children }: RouteAreaShellProps) {
  return (
    <section className={cn("relative", tone !== "site" && "isolate py-2")}>
      {tone === "auth" ? (
        <>
          <div
            aria-hidden
            className="pointer-events-none absolute top-16 left-[6%] h-44 w-44 rounded-full bg-[rgba(239,125,87,0.14)] blur-3xl"
          />
          <div
            aria-hidden
            className="pointer-events-none absolute top-28 right-[8%] h-52 w-52 rounded-full bg-[rgba(31,111,120,0.12)] blur-3xl"
          />
        </>
      ) : null}

      {tone === "workspace" ? (
        <>
          <div
            aria-hidden
            className="pointer-events-none absolute top-6 left-1/2 h-56 w-[min(72rem,92vw)] -translate-x-1/2 rounded-full bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.34),rgba(255,255,255,0.04)_68%,transparent_82%)] blur-3xl"
          />
          <div
            aria-hidden
            className="pointer-events-none absolute top-12 right-[10%] h-40 w-40 rounded-full bg-[rgba(31,111,120,0.12)] blur-3xl"
          />
        </>
      ) : null}

      <div className="relative">{children}</div>
    </section>
  );
}
