// Tokyo Night themed cat logo SVG
export function CatLogo() {
  return (
    <svg
      className="w-12 h-12"
      viewBox="0 0 100 100"
      xmlns="http://www.w3.org/2000/svg"
    >
      {/* Ears */}
      <path
        d="M20 45 L30 15 L40 40"
        stroke="#565f89"
        strokeWidth="3"
        fill="#24283b"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M60 40 L70 15 L80 45"
        stroke="#565f89"
        strokeWidth="3"
        fill="#24283b"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      {/* Inner ears - neon pink */}
      <path
        d="M25 42 L30 22 L35 40"
        stroke="#f7768e"
        strokeWidth="2"
        fill="#f7768e33"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M65 40 L70 22 L75 42"
        stroke="#f7768e"
        strokeWidth="2"
        fill="#f7768e33"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      {/* Head */}
      <ellipse
        cx="50"
        cy="58"
        rx="32"
        ry="28"
        stroke="#565f89"
        strokeWidth="3"
        fill="#24283b"
      />
      {/* Eyes - neon yellow */}
      <ellipse
        cx="38"
        cy="52"
        rx="6"
        ry="7"
        stroke="#565f89"
        strokeWidth="2"
        fill="#e0af68"
      />
      <ellipse
        cx="62"
        cy="52"
        rx="6"
        ry="7"
        stroke="#565f89"
        strokeWidth="2"
        fill="#e0af68"
      />
      <ellipse cx="38" cy="53" rx="3" ry="4" fill="#1a1b26" />
      <ellipse cx="62" cy="53" rx="3" ry="4" fill="#1a1b26" />
      <circle cx="36" cy="51" r="1.5" fill="#7dcfff" />
      <circle cx="60" cy="51" r="1.5" fill="#7dcfff" />
      {/* Nose - neon pink */}
      <path
        d="M47 62 L50 66 L53 62 Z"
        stroke="#565f89"
        strokeWidth="1.5"
        fill="#f7768e"
      />
      {/* Mouth */}
      <path
        d="M50 66 L50 70"
        stroke="#565f89"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <path
        d="M50 70 Q44 76 38 72"
        stroke="#565f89"
        strokeWidth="2"
        fill="none"
        strokeLinecap="round"
      />
      <path
        d="M50 70 Q56 76 62 72"
        stroke="#565f89"
        strokeWidth="2"
        fill="none"
        strokeLinecap="round"
      />
      {/* Whiskers */}
      <path
        d="M30 65 L12 60"
        stroke="#565f89"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <path
        d="M30 68 L12 68"
        stroke="#565f89"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <path
        d="M30 71 L12 76"
        stroke="#565f89"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <path
        d="M70 65 L88 60"
        stroke="#565f89"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <path
        d="M70 68 L88 68"
        stroke="#565f89"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <path
        d="M70 71 L88 76"
        stroke="#565f89"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  );
}
