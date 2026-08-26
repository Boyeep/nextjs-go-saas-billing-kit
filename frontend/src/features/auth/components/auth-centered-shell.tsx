import type { ReactNode } from "react";

type AuthCenteredShellProps = {
  children: ReactNode;
};

export function AuthCenteredShell({ children }: AuthCenteredShellProps) {
  return (
    <main className="mx-auto flex min-h-[calc(100vh-12rem)] w-[min(100%,32rem)] items-center justify-center px-5 py-10">
      {children}
    </main>
  );
}
