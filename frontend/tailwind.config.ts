import type { Config } from "tailwindcss";

export default {
  darkMode: ["class"],
  content: ["./pages/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}", "./app/**/*.{ts,tsx}", "./src/**/*.{ts,tsx}"],
  prefix: "",
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
          glow: "hsl(var(--primary-glow))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        success: {
          DEFAULT: "hsl(var(--success))",
          foreground: "hsl(var(--success-foreground))",
        },
        warning: {
          DEFAULT: "hsl(var(--warning))",
          foreground: "hsl(var(--warning-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        sidebar: {
          DEFAULT: "hsl(var(--sidebar-background))",
          foreground: "hsl(var(--sidebar-foreground))",
          primary: "hsl(var(--sidebar-primary))",
          "primary-foreground": "hsl(var(--sidebar-primary-foreground))",
          accent: "hsl(var(--sidebar-accent))",
          "accent-foreground": "hsl(var(--sidebar-accent-foreground))",
          border: "hsl(var(--sidebar-border))",
          ring: "hsl(var(--sidebar-ring))",
        },
        // BuildFlow / Construction POC Stitch Design Tokens
        "inverse-surface": "#2d3133",
        "secondary-container": "#6cf8bb",
        "tertiary-fixed-dim": "#ffb95f",
        "surface-variant": "#e0e3e5",
        "on-tertiary-fixed-variant": "#653e00",
        "on-primary-fixed-variant": "#38485a",
        "on-surface-variant": "#44474c",
        "tertiary-container": "#3e2400",
        "inverse-primary": "#b7c8de",
        "surface-container-highest": "#e0e3e5",
        "surface-bright": "#f7fafc",
        "on-secondary": "#ffffff",
        "secondary-fixed-dim": "#4edea3",
        "on-secondary-container": "#00714d",
        "on-tertiary": "#ffffff",
        "primary-container": "#1a2b3c",
        "tertiary-fixed": "#ffddb8",
        "surface-dim": "#d7dadc",
        "outline": "#74777d",
        "on-error-container": "#93000a",
        "on-background": "#181c1e",
        "on-primary-container": "#8192a7",
        "on-tertiary-container": "#ca8100",
        "surface-container-lowest": "#ffffff",
        "primary-fixed-dim": "#b7c8de",
        "surface-container-high": "#e5e9eb",
        "secondary-fixed": "#6ffbbe",
        "primary-fixed": "#d2e4fb",
        "surface-container": "#ebeef0",
        "on-secondary-fixed": "#002113",
        "on-primary-fixed": "#0b1d2d",
        "on-tertiary-fixed": "#2a1700",
        "surface-container-low": "#f1f4f6",
        "error-container": "#ffdad6",
        "inverse-on-surface": "#eef1f3",
        "surface-tint": "#4f6073",
        "on-surface": "#181c1e",
        "outline-variant": "#c4c6cd",
        "on-secondary-fixed-variant": "#005236",
      },
      backgroundImage: {
        'gradient-primary': 'var(--gradient-primary)',
        'gradient-card': 'var(--gradient-card)',
      },
      boxShadow: {
        'glow': 'var(--shadow-glow)',
        'card': 'var(--shadow-card)',
      },
      fontFamily: {
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace'],
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      keyframes: {
        "accordion-down": {
          from: {
            height: "0",
          },
          to: {
            height: "var(--radix-accordion-content-height)",
          },
        },
        "accordion-up": {
          from: {
            height: "var(--radix-accordion-content-height)",
          },
          to: {
            height: "0",
          },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
} satisfies Config;
