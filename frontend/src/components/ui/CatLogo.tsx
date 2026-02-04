// Retro Tokyo neon sign cat logo - minimalist with ^^ eyes
export function CatLogo() {
  return (
    <svg
      className="w-12 h-12"
      viewBox="0 0 100 100"
      xmlns="http://www.w3.org/2000/svg"
    >
      <defs>
        <filter id="neon" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="2.5" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
        <filter id="neon-soft" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="1.5" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

      {/* Main cat outline - single continuous neon tube style */}
      <g fill="none" strokeLinecap="round" strokeLinejoin="round">
        {/* Outer head + ears silhouette */}
        <path
          d="M15 70
             Q15 45 25 40
             L30 15
             L40 38
             Q50 32 60 38
             L70 15
             L75 40
             Q85 45 85 70
             Q85 90 50 90
             Q15 90 15 70"
          stroke="#7dcfff"
          strokeWidth="3"
          filter="url(#neon)"
        />

        {/* Eyes - simple ^^ style */}
        <g stroke="#7dcfff" strokeWidth="2.5" filter="url(#neon-soft)">
          <path d="M32 62 L38 54 L44 62" />
          <path d="M56 62 L62 54 L68 62" />
        </g>
      </g>
    </svg>
  );
}
