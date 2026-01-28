import { useState, useRef, useEffect } from "react";
import { ChevronDown, Check } from "lucide-react";

interface Option {
  value: string;
  label: string;
}

interface NeonSelectProps {
  options: Option[];
  value: string | null;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export function NeonSelect({
  options,
  value,
  onChange,
  placeholder = "Select...",
  className = "",
}: NeonSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  const selectedOption = options.find((opt) => opt.value === value);

  // Close on outside click
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  // Close on escape
  useEffect(() => {
    function handleEscape(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setIsOpen(false);
      }
    }

    document.addEventListener("keydown", handleEscape);
    return () => document.removeEventListener("keydown", handleEscape);
  }, []);

  return (
    <div ref={containerRef} className={`relative ${className}`}>
      {/* Trigger Button */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className={`
          flex items-center justify-between gap-3 w-full px-4 py-2
          bg-night-surface border rounded-md
          text-left transition-all cursor-pointer
          ${
            isOpen
              ? "border-neon-cyan shadow-[0_0_12px_rgba(125,207,255,0.3)]"
              : "border-night-border hover:border-neon-cyan/50"
          }
        `}
      >
        <span
          className={`whitespace-nowrap ${selectedOption ? "text-text-primary" : "text-text-muted"}`}
        >
          {selectedOption?.label || placeholder}
        </span>
        <ChevronDown
          className={`w-4 h-4 text-neon-cyan transition-transform duration-200 ${isOpen ? "rotate-180" : ""}`}
        />
      </button>

      {/* Dropdown Menu */}
      {isOpen && (
        <div
          className="absolute top-full left-0 right-0 mt-1 z-50
          bg-night-surface border border-neon-cyan/50 rounded-md
          shadow-[0_4px_20px_rgba(0,0,0,0.5),0_0_15px_rgba(125,207,255,0.15)]
          overflow-hidden animate-in fade-in slide-in-from-top-2 duration-150"
        >
          <div className="max-h-[300px] overflow-y-auto py-1">
            {options.length === 0 ? (
              <div className="px-4 py-3 text-text-muted text-sm">
                No options
              </div>
            ) : (
              options.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => {
                    onChange(option.value);
                    setIsOpen(false);
                  }}
                  className={`
                    flex items-center justify-between w-full px-4 py-2.5 text-left
                    transition-colors duration-100
                    ${
                      option.value === value
                        ? "bg-neon-cyan/10 text-neon-cyan"
                        : "text-text-primary hover:bg-night-bg-highlight hover:text-neon-cyan"
                    }
                  `}
                >
                  <span>{option.label}</span>
                  {option.value === value && (
                    <Check className="w-4 h-4 text-neon-cyan" />
                  )}
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}
