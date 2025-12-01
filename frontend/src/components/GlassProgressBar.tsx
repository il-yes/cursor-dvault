import React from "react";

interface GlassProgressBarProps {
  value: number; // 0 to 100
  max?: number;
  label?: string;
}
export function GlassProgressBar00({ value, max = 100, label }: GlassProgressBarProps) {
  const percentage = Math.min(Math.max((value / max) * 100, 0), 100);

  // Don't render if progress is 100%
  if (percentage >= 100) {
    return null;
  }

  return (
    <div className="w-full rounded-3xl bg-white/30 dark:bg-zinc-900/30 backdrop-blur-xl border border-white/40 dark:border-zinc-700/40 shadow-lg overflow-hidden select-none fixed bottom-4 left-0 right-0 mx-auto max-w-xl">
      <div
        className="h-6 rounded-3xl bg-gradient-to-r from-[#C9A44A] via-[#D2AF56] to-[#B8934A] shadow-lg transition-all duration-500"
        style={{ width: `${percentage}%` }}
      />
      {label && (
        <span className="absolute right-4 top-1/2 transform -translate-y-1/2 text-sm font-semibold text-foreground drop-shadow-md pointer-events-none">
          {label}
        </span>
      )}
    </div>
  );
}

export function GlassProgressBar({ value, max = 100, label, visible }: GlassProgressBarProps & { visible: boolean }) {
  if (!visible) return null;
  const percentage = Math.min(Math.max((value / max) * 100, 0), 100);

  // Don't render if progress is 100%
  if (percentage >= 100) {
    return null;
  }

  return (
    <div className="w-full rounded-3xl bg-white/30 dark:bg-zinc-900/30 backdrop-blur-xl border border-white/40 dark:border-zinc-700/40 shadow-lg overflow-hidden select-none fixed bottom-4 left-0 right-0 mx-auto max-w-xl">
      <div
        className="h-6 rounded-3xl bg-gradient-to-r from-[#C9A44A] via-[#D2AF56] to-[#B8934A] shadow-lg transition-all duration-500"
        style={{ width: `${percentage}%` }}
      />
      {label && (
        <span className="absolute right-4 top-1/2 transform -translate-y-1/2 text-sm font-semibold text-foreground drop-shadow-md pointer-events-none">
          {label}
        </span>
      )}
    </div>
  );
}
