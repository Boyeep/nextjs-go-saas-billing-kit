import * as React from "react";
import { cn } from "@/lib/utils";

const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<"input">>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "border-border placeholder:text-muted-foreground focus-visible:ring-ring flex h-12 w-full rounded-2xl border bg-white/90 px-4 py-2 text-sm transition outline-none focus-visible:ring-2",
          className,
        )}
        ref={ref}
        {...props}
      />
    );
  },
);
Input.displayName = "Input";

export { Input };
