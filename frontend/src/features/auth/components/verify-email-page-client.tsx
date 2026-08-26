"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { CheckCircle2, MailCheck, ShieldAlert } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AuthSplitShell } from "@/features/auth/components/auth-split-shell";
import { verifyEmail } from "@/features/auth/services/auth-service";

type VerifyEmailPageClientProps = {
  token?: string;
};

export function VerifyEmailPageClient({ token }: VerifyEmailPageClientProps) {
  const router = useRouter();
  const [statusMessage, setStatusMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const hasTriggeredVerification = useRef(false);

  const verifyMutation = useMutation({
    mutationFn: async () => {
      if (!token) {
        throw new Error(
          "Missing verification token. Please reopen the link from your email.",
        );
      }

      return verifyEmail({ token });
    },
    onSuccess: () => {
      setStatusMessage("Your email is verified. Redirecting you to login...");
      setTimeout(() => {
        router.push("/login");
      }, 1200);
    },
    onError: (error) => {
      setErrorMessage(
        error instanceof Error ? error.message : "Failed to verify email.",
      );
    },
  });

  useEffect(() => {
    if (token && !hasTriggeredVerification.current) {
      hasTriggeredVerification.current = true;
      verifyMutation.mutate();
    }
  }, [token, verifyMutation]);

  return (
    <AuthSplitShell
      badge="Email Verification"
      heading="Confirm your email so your account is ready to use."
      description="We only activate new accounts after the owner proves they can open the inbox tied to the email address."
    >
      <Card className="rounded-[32px] border border-white/70 bg-white/85 p-[clamp(1.5rem,4vw,2rem)] shadow-[var(--shadow)]">
        <CardHeader className="mb-6 grid gap-2 p-0">
          <CardTitle>Verifying your email</CardTitle>
          <CardDescription>
            We are confirming your verification token now.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          {statusMessage ? (
            <p className="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
              {statusMessage}
            </p>
          ) : null}
          {errorMessage ? (
            <p className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
              {errorMessage}
            </p>
          ) : null}

          <div className="rounded-[28px] bg-[var(--surface)] p-5">
            <div className="flex items-center gap-3">
              {verifyMutation.isPending ? (
                <MailCheck className="h-5 w-5 text-[var(--accent-brand)]" />
              ) : errorMessage ? (
                <ShieldAlert className="h-5 w-5 text-red-600" />
              ) : (
                <CheckCircle2 className="h-5 w-5 text-emerald-600" />
              )}
              <p className="font-semibold">
                {verifyMutation.isPending
                  ? "Confirming your email..."
                  : "Verification status updated."}
              </p>
            </div>
          </div>
          <div className="flex flex-wrap gap-3">
            <Button asChild>
              <Link href="/login">Back to login</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </AuthSplitShell>
  );
}
