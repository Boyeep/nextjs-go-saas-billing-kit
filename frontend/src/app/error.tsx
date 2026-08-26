"use client";

import { useEffect } from "react";
import { Button } from "@/components/ui/button";

type GlobalErrorProps = {
  error: Error & { digest?: string };
  reset: () => void;
};

export default function GlobalError({ error, reset }: GlobalErrorProps) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <main className="mx-auto grid min-h-[70vh] w-[min(var(--max-width),calc(100%-2rem))] place-items-center max-[720px]:w-[min(var(--max-width),calc(100%-1.25rem))]">
      <section className="grid max-w-[700px] gap-4 rounded-[36px] border border-red-200 bg-white/90 px-8 py-10 shadow-[var(--shadow)]">
        <span className="section-heading__eyebrow bg-red-100 text-red-700">
          App Error
        </span>
        <h1 className="m-0 font-[family-name:var(--font-display)] text-[clamp(2rem,5vw,3.6rem)] leading-[1]">
          Something broke, but the template is recoverable.
        </h1>
        <p className="m-0 text-[var(--muted-text)]">
          This is a built-in error boundary for unexpected runtime failures. You
          can customize it per project.
        </p>
        <div className="flex flex-wrap gap-3">
          <Button type="button" onClick={reset}>
            Try again
          </Button>
          <Button
            type="button"
            variant="secondary"
            onClick={() => window.location.assign("/")}
          >
            Go home
          </Button>
        </div>
      </section>
    </main>
  );
}
