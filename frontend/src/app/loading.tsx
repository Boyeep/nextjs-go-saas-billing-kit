export default function Loading() {
  return (
    <main className="mx-auto grid min-h-[60vh] w-[min(var(--max-width),calc(100%-2rem))] place-items-center max-[720px]:w-[min(var(--max-width),calc(100%-1.25rem))]">
      <div className="grid gap-4 rounded-[32px] border border-white/70 bg-white/80 px-8 py-10 text-center shadow-[var(--shadow)]">
        <div className="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-[var(--accent-soft)] border-t-[var(--accent-brand)]" />
        <div className="grid gap-2">
          <h1 className="font-[family-name:var(--font-display)] text-3xl">
            Loading your workspace
          </h1>
          <p className="m-0 max-w-[520px] text-[var(--muted-text)]">
            The template is pulling together page data and shared app chrome.
          </p>
        </div>
      </div>
    </main>
  );
}
