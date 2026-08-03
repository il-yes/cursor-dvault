// frontend/src/components/C3/styles/menu.tsx
import { Global, css } from '@emotion/react';

export const c3MenuStyles = css`
  .sidebar {
    width: 286px;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    padding: 18px 0 16px;
    background: rgba(250, 250, 250, 0.9);
    border-right: 1px solid var(--c3-border-soft);
    overflow-y: auto;
  }

  .sidebar-section-label {
    padding: 0 16px 8px;
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.1em;
    color: #bbb;
    text-transform: uppercase;
  }

  .nav-row,
  .vault-row {
    display: flex;
    align-items: center;
    gap: 9px;
    padding: 7px 16px;
    font-size: 13px;
    color: #555;
    cursor: pointer;
  }

  .nav-row:hover,
  .vault-row:hover {
    background: rgba(240, 240, 240, 0.9);
  }

  .nav-row.active {
    background: var(--c3-gold-soft);
    color: var(--c3-gold);
    font-weight: 600;
  }

  .nav-icon {
    font-size: 13px;
    width: 18px;
    text-align: center;
    flex-shrink: 0;
  }

  .nav-label {
    flex: 1;
  }

  .nav-badge,
  .nav-badge-gray {
    font-size: 10px;
    padding: 1px 6px;
    border-radius: 10px;
    min-width: 18px;
    text-align: center;
  }

  .nav-badge {
    background: var(--c3-gold);
    color: #fff;
    font-weight: 700;
  }

  .nav-badge-gray {
    background: #e8e8e8;
    color: #888;
    font-weight: 600;
  }

  .vault-row.all-vaults {
    font-size: 12px;
    color: #888;
    padding: 5px 16px 10px;
    border-bottom: 1px solid var(--c3-border-soft);
    margin-bottom: 6px;
  }

  .vault-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .unread-pip {
    margin-left: auto;
    width: 6px;
    height: 6px;
    background: var(--c3-gold);
    border-radius: 50%;
  }
`;

export const C3MenuStyles = () => <Global styles={c3MenuStyles} />;