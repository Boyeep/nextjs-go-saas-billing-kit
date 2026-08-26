import * as React from "react";
import { cn } from "@/lib/utils";

const Textarea = React.forwardRef<
  HTMLTextAreaElement,
  React.ComponentProps<"textarea">
>(({ className, ...props }, ref) => {
  return (
    <textarea
      className={cn(
        "border-border placeholder:text-muted-foreground focus-visible:ring-ring flex min-h-[120px] w-full rounded-2xl border bg-white/90 px-4 py-3 text-sm transition outline-none focus-visible:ring-2",
        className,
      )}
      ref={ref}
      {...props}
    />
  );
});
Textarea.displayName = "Textarea";

export { Textarea };
