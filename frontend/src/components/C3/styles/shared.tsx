// frontend/src/components/C3/styles/shared.tsx
import { Global, css } from '@emotion/react';

export const c3SharedStyles = css`
  :root {
    --c3-gold: #c8922a;
    --c3-gold-soft: #fbf0d8;
    --c3-gold-border: #e8c87a;
    --c3-text: #1a1a1a;
    --c3-muted: #888;
    --c3-soft: #f5f5f5;
    --c3-border: #e5e5e5;
    --c3-border-soft: #ebebeb;
    --c3-bg: #fff;
    --c3-purple: #7c3aed;
    --c3-blue: #2563eb;
    --c3-green: #059669;
    --c3-red: #dc2626;
  }

  * {
    box-sizing: border-box;
  }

  html,
  body,
  #root {
    height: 100%;
  }

  body {
    margin: 0;
    font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    color: var(--c3-text);
    background: radial-gradient(circle at top, #ffffff 0%, #fafafa 38%, #f4f4f4 100%);
  }

  button,
  input,
  select,
  textarea {
    font: inherit;
  }

  .glass-surface {
    background: rgba(255, 255, 255, 0.72);
    backdrop-filter: blur(22px) saturate(180%);
    -webkit-backdrop-filter: blur(22px) saturate(180%);
    border: 1px solid rgba(255, 255, 255, 0.55);
    box-shadow: 0 18px 60px rgba(0, 0, 0, 0.08), inset 0 1px 0 rgba(255, 255, 255, 0.65);
  }

  .window {
    width: 1320px;
    height: 820px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.88);
    box-shadow: 0 28px 80px rgba(0, 0, 0, 0.18), 0 8px 24px rgba(0, 0, 0, 0.08);
  }

  .titlebar {
    height: 40px;
    display: flex;
    align-items: center;
    padding: 0 16px;
    flex-shrink: 0;
    background: rgba(245, 245, 245, 0.92);
    border-bottom: 1px solid var(--c3-border-soft);
  }

  .traffic-lights {
    display: flex;
    gap: 8px;
  }

  .tl {
    width: 12px;
    height: 12px;
    border-radius: 50%;
  }

  .tl-red { background: #ff5f57; }
  .tl-yellow { background: #ffbd2e; }
  .tl-green { background: #28ca41; }

  .titlebar-name {
    flex: 1;
    text-align: center;
    font-size: 12px;
    font-weight: 500;
    color: #666;
    letter-spacing: 0.04em;
  }

  .topbar {
    height: 52px;
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 0 20px;
    flex-shrink: 0;
    background: rgba(255, 255, 255, 0.82);
  }

  .topbar-logo {
    font-weight: 700;
    font-size: 14px;
    letter-spacing: 0.12em;
    color: var(--c3-text);
  }

  .workspace-pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 10px;
    border-radius: 999px;
    font-size: 12px;
    color: #555;
    background: rgba(245, 245, 245, 0.92);
    border: 1px solid var(--c3-border);
  }

  .topbar-spacer {
    flex: 1;
  }

  .tab-bar {
    display: flex;
    gap: 8px;
    padding: 12px 24px 0;
    flex-shrink: 0;
  }

  .tab {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 7px 12px;
    border-radius: 999px;
    border: 1px solid var(--c3-border);
    background: rgba(255, 255, 255, 0.84);
    color: #666;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
  }

  .tab.active {
    background: var(--c3-gold-soft);
    border-color: var(--c3-gold-border);
    color: var(--c3-gold);
  }

  .btn,
  .ctrl-btn,
  .rha-btn,
  .action-btn,
  .c3-add-btn,
  .tc-start,
  .commit-btn,
  .view-btn {
    border: 1px solid transparent;
    border-radius: 10px;
    transition: transform 160ms ease, border-color 160ms ease, background-color 160ms ease, color 160ms ease, box-shadow 160ms ease;
    cursor: pointer;
  }

  .btn:hover,
  .ctrl-btn:hover,
  .rha-btn:hover,
  .action-btn:hover,
  .c3-add-btn:hover,
  .tc-start:hover,
  .commit-btn:hover,
  .view-btn:hover,
  .tab:hover {
    transform: translateY(-1px);
  }

  .btn-primary {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 14px;
    background: linear-gradient(180deg, #d5a44a 0%, var(--c3-gold) 100%);
    color: #fff;
    font-size: 12px;
    font-weight: 600;
    box-shadow: 0 10px 24px rgba(200, 146, 42, 0.18);
  }

  .btn-ghost {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 14px;
    background: var(--c3-soft);
    border-color: var(--c3-border);
    color: #444;
    font-size: 12px;
    font-weight: 500;
  }

  .page-title {
    padding: 16px 24px 8px;
    font-size: 18px;
    font-weight: 700;
    color: var(--c3-text);
    letter-spacing: 0.04em;
  }

  .page-subtitle {
    padding: 0 24px 14px;
    font-size: 12px;
    color: #9a9a9a;
  }

  .ledger-footer,
  .ledger-footer-hint,
  .inbox-footer {
    padding: 12px 24px;
    border-top: 1px solid var(--c3-border-soft);
    font-size: 12px;
    color: #bbb;
    flex-shrink: 0;
    background: rgba(255, 255, 255, 0.72);
  }




    /* Empty vault sidebar */
    .sidebar {
      background: #fafafa;
      border-right: 1px solid #ebebeb;
      display: flex; flex-direction: column;
      padding: 18px 0 16px;
      flex-shrink: 0;
    }
    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }
    .sidebar-empty-hint {
      margin: 12px 16px;
      padding: 12px;
      border: 1px dashed #e0e0e0;
      border-radius: 8px;
      font-size: 11px; color: #ccc;
      text-align: center;
      line-height: 1.5;
    }
    .add-vault-btn {
      margin: 10px 14px 0;
      padding: 9px 10px;
      border: 1px dashed #ddd;
      border-radius: 7px;
      text-align: center;
      font-size: 12px; color: #C8922A;
      cursor: pointer;
      font-weight: 500;
    }
    .add-vault-btn:hover { background: #FBF0D8; }
`;

export const C3SharedStyles = () => <Global styles={c3SharedStyles} />;