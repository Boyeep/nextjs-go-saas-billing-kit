type BrandMarkProps = {
  className?: string;
};

export function BrandMark({ className = "" }: BrandMarkProps) {
  return (
    <span
      aria-hidden="true"
      className={`inline-flex overflow-hidden rounded-[18px] border border-white/55 bg-[linear-gradient(135deg,rgba(255,255,255,0.9),rgba(255,236,224,0.78))] shadow-[0_10px_24px_rgba(239,125,87,0.18)] ${className}`.trim()}
    >
      <svg
        viewBox="0 0 64 64"
        className="h-full w-full"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
      >
        <rect width="64" height="64" fill="url(#bg)" />
        <rect
          x="12"
          y="12"
          width="40"
          height="40"
          rx="12"
          fill="rgba(255,255,255,0.72)"
          stroke="rgba(31,41,55,0.08)"
        />
        <circle cx="25" cy="25" r="5" fill="#F4A261" />
        <path
          d="M19 43L28 33L35 40L41 31L47 43H19Z"
          fill="#1F6F78"
          fillOpacity="0.92"
        />
        <path
          d="M19 43L28 33L35 40L41 31L47 43"
          stroke="rgba(255,255,255,0.8)"
          strokeWidth="2.5"
          strokeLinejoin="round"
        />
        <defs>
          <linearGradient id="bg" x1="6" y1="8" x2="58" y2="58">
            <stop stopColor="#FFD7BF" />
            <stop offset="0.55" stopColor="#FFF5EC" />
            <stop offset="1" stopColor="#D7EEF0" />
          </linearGradient>
        </defs>
      </svg>
    </span>
  );
}
