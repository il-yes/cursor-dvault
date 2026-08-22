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

  .workspace-selector {
    position: relative;
    display: inline-flex;
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
    cursor: pointer;
    user-select: none;
    transition: border-color 160ms ease, background-color 160ms ease;
  }

  .workspace-pill:hover {
    border-color: var(--c3-gold-border);
    background: var(--c3-gold-soft);
    color: var(--c3-gold);
  }

  .workspace-pill.active {
    border-color: var(--c3-gold-border);
    background: var(--c3-gold-soft);
    color: var(--c3-gold);
  }

  .workspace-pill-chevron {
    font-size: 10px;
    opacity: 0.6;
    transition: transform 160ms ease;
  }

  .workspace-pill.active .workspace-pill-chevron {
    transform: rotate(180deg);
  }

  .workspace-dropdown-backdrop {
    position: fixed;
    inset: 0;
    z-index: 199;
  }

  .workspace-dropdown {
    position: absolute;
    top: calc(100% + 6px);
    left: 0;
    min-width: 240px;
    max-width: 320px;
    max-height: 320px;
    background: #ffffff;
    backdrop-filter: none;
    border: 1px solid var(--c3-border);
    border-radius: 10px;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12), 0 4px 12px rgba(0, 0, 0, 0.06);
    z-index: 200;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    animation: workspace-dropdown-in 120ms ease-out;
  }

  @keyframes workspace-dropdown-in {
    from {
      opacity: 0;
      transform: translateY(-4px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* Fixed C3 Sliding View Overlay & Container */
  .c3-sliding-view-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.35);
    backdrop-filter: blur(2px);
    z-index: 999;
    pointer-events: auto;
    animation: c3-fade-in 150ms ease-out;
  }

  .c3-sliding-view-container {
    position: fixed;
    top: 0;
    right: 0;
    width: 460px;
    height: 100vh;
    background: #ffffff;
    border-left: 1px solid #e0e0e0;
    box-shadow: -8px 0 32px rgba(0, 0, 0, 0.10);
    z-index: 1000;
    pointer-events: auto;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    animation: c3-slide-in-right 180ms cubic-bezier(0.16, 1, 0.3, 1);
  }

  @keyframes c3-fade-in {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes c3-slide-in-right {
    from { transform: translateX(100%); }
    to { transform: translateX(0); }
  }

  .workspace-dropdown-header {
    padding: 8px 12px 6px;
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #bbb;
    border-bottom: 1px solid var(--c3-border-soft);
    flex-shrink: 0;
  }

  .workspace-dropdown-list {
    overflow-y: auto;
    flex: 1;
    padding: 4px 0;
  }

  .workspace-option {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    font-size: 13px;
    color: var(--c3-text);
    cursor: pointer;
    transition: background-color 100ms ease;
    border: none;
    background: none;
    width: 100%;
    text-align: left;
  }

  .workspace-option:hover {
    background: var(--c3-gold-soft);
  }

  .workspace-option.selected {
    color: var(--c3-gold);
    font-weight: 600;
  }

  .workspace-option-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--c3-gold);
    flex-shrink: 0;
  }

  .workspace-option.selected .workspace-option-dot {
    box-shadow: 0 0 0 2px var(--c3-gold-soft), 0 0 0 3px var(--c3-gold);
  }

  .workspace-option-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .workspace-option-check {
    font-size: 12px;
    color: var(--c3-gold);
    flex-shrink: 0;
  }

  .workspace-dropdown-divider {
    height: 1px;
    background: var(--c3-border-soft);
    margin: 2px 0;
    flex-shrink: 0;
  }

  .workspace-create-action {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 9px 12px;
    font-size: 12px;
    font-weight: 600;
    color: var(--c3-gold);
    cursor: pointer;
    border: none;
    background: none;
    width: 100%;
    text-align: left;
    flex-shrink: 0;
    transition: background-color 100ms ease;
  }

  .workspace-create-action:hover {
    background: var(--c3-gold-soft);
  }

  .workspace-empty {
    padding: 16px 12px;
    text-align: center;
    font-size: 12px;
    color: #aaa;
    line-height: 1.5;
  }

  .workspace-loading {
    padding: 12px;
    text-align: center;
    font-size: 12px;
    color: #aaa;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
  }

  .workspace-loading-spinner {
    display: inline-block;
    width: 12px;
    height: 12px;
    border: 2px solid var(--c3-border);
    border-top-color: var(--c3-gold);
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
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