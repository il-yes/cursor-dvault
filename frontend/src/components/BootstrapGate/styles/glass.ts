// styles/glass.ts
import { css } from "@emotion/react";

export const glassPerfect = css`
  backdrop-filter: blur(30px) saturate(180%);
  -webkit-backdrop-filter: blur(30px) saturate(180%);
  background: rgba(255, 255, 255, 0.06);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.25),
    inset 0 0 40px rgba(255, 255, 255, 0.05);
`;

export const goldAccent = `
  background: linear-gradient(135deg, #D4AF37, #FFB84D);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
`;