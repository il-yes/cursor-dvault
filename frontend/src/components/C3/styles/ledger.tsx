// frontend/src/components/C3/styles/ledger.tsx
import { Global, css } from '@emotion/react';

export const c3LedgerStyles = css`
  .layout {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .ledger-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-width: 0;
  }

  .ledger-topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 24px 14px;
    flex-shrink: 0;
    background: rgba(255, 255, 255, 0.68);
  }

  .ledger-title {
    font-size: 17px;
    font-weight: 700;
    letter-spacing: 0.06em;
    color: var(--c3-text);
  }

  .ledger-controls {
    display: flex;
    gap: 8px;
    align-items: center;
  }

  .ctrl-btn {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 5px 12px;
    background: rgba(255, 255, 255, 0.95);
    border-color: var(--c3-border);
    border-radius: 10px;
    font-size: 12px;
    color: #555;
  }

  .table-wrap {
    flex: 1;
    overflow: auto;
    padding: 0 24px;
    scrollbar-gutter: stable;
  }

  table {
    width: 100%;
    border-collapse: separate;
    border-spacing: 0;
  }

  thead th {
    position: sticky;
    top: 0;
    z-index: 1;
    padding: 0 10px 10px;
    text-align: left;
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.09em;
    text-transform: uppercase;
    color: #bbb;
    background: rgba(255, 255, 255, 0.96);
    border-bottom: 1px solid var(--c3-border-soft);
    white-space: nowrap;
  }

  tbody tr {
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  tbody tr:hover {
    background: rgba(250, 250, 250, 0.88);
  }

  tbody tr:hover .row-hover-actions {
    display: flex;
  }

  tbody tr:hover .stellar-val {
    display: none;
  }

  td {
    padding: 11px 10px;
    vertical-align: middle;
    border-bottom: 1px solid #f2f2f2;
    font-size: 13px;
  }

  .sdot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
  .s-ok { background: #22c55e; }
  .s-pend { background: #f59e0b; }
  .s-dispute { background: #dc2626; }

  .th-line1 { font-size: 13px; font-weight: 500; color: var(--c3-text); white-space: nowrap; }
  .th-type { font-size: 9px; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-right: 6px; }
  .th-line2 { font-size: 11px; color: #aaa; margin-top: 2px; }
  .ts { font-size: 12px; color: #bbb; white-space: nowrap; }

  .stellar-val {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 10px;
    color: #bbb;
    white-space: nowrap;
  }

  .flow { display: flex; align-items: center; gap: 4px; }

  .vb {
    width: 22px;
    height: 22px;
    border-radius: 50%;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 8px;
    font-weight: 700;
    color: #fff;
    flex-shrink: 0;
  }

  .fa { color: #ccc; font-size: 10px; font-weight: 400; }

  .pipeline { display: flex; gap: 3px; align-items: center; }
  .pseg { height: 4px; width: 20px; border-radius: 999px; }
  .pseg-done { background: var(--c3-gold); }
  .pseg-wait { background: #e8e8e8; }
  .pseg-reject { background: #ef4444; }

  .c3b {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 3px 8px;
    border-radius: 999px;
    font-size: 11px;
    font-weight: 700;
    white-space: nowrap;
    cursor: pointer;
  }

  .c3-ext {
    background: var(--c3-soft);
    border: 1px solid var(--c3-border);
    color: #999;
  }

  .c3-ext:hover {
    border-color: var(--c3-gold);
    color: var(--c3-gold);
    background: var(--c3-gold-soft);
  }

  .c3-linked {
    background: var(--c3-gold-soft);
    border: 1px solid var(--c3-gold-border);
    color: var(--c3-gold);
  }

  .c3-active {
    background: rgba(124, 58, 237, 0.08);
    border: 1px solid rgba(124, 58, 237, 0.2);
    color: #7c3aed;
  }
    /* Channel extension */


.policy-box {

display:flex;

flex-direction:column;

gap:6px;

}



.channel-table td {

vertical-align:top;

}



.asset-box {

border-bottom:
1px solid #222;

padding:
8px 0;

}



.asset-box:last-child {

border-bottom:none;

}



.event {

font-size:12px;

color:#888;

}



.event span {

font-family:
monospace;

color:#555;

}



.ledger-table tr:hover {

background:
rgba(255,255,255,0.03);

}
.thread-drawer {

    width:520px;
    height:100vh;

    margin-left:auto;

    background:#0b0b0f;

    border-left:1px solid #222;

}
    .thread-drawer {

    width:520px;
    height:100vh;

    margin-left:auto;

    padding:0;

    background:transparent;

    border:none;

}


.detail-panel {

    height:100%;

    background:#fff;

    color:#111;

    display:flex;
    flex-direction:column;

    overflow:hidden;

    border-left:1px solid #ddd;

}
    .drawer-reset {
    width: 520px !important;
    height: 100vh !important;
    margin-left:auto;
    padding:0 !important;
    background:transparent !important;
}

    * { margin: 0; padding: 0; box-sizing: border-box; }

    

    .window {
      width: 1320px;
      height: 820px;
      background: #fff;
      border-radius: 12px;
      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }

    .titlebar {
      height: 40px;
      background: #f5f5f5;
      border-bottom: 1px solid #e0e0e0;
      display: flex; align-items: center;
      padding: 0 16px;
      flex-shrink: 0;
    }
    .traffic-lights { display: flex; gap: 8px; }
    .tl { width: 12px; height: 12px; border-radius: 50%; }
    .tl-red { background: #FF5F57; }
    .tl-yellow { background: #FFBD2E; }
    .tl-green { background: #28CA41; }
    .titlebar-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }

    .topbar {
      height: 52px;
      background: #fff;
      border-bottom: 1px solid #ebebeb;
      display: flex; align-items: center;
      padding: 0 20px;
      flex-shrink: 0; gap: 16px;
    }
    .topbar-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; color: #1a1a1a; }
    .workspace-pill { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }
    .topbar-spacer { flex: 1; }
    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 6px 14px; border-radius: 7px; font-size: 12px; font-weight: 500; cursor: pointer; border: none; }
    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #ccc; pointer-events: none; }
    .btn-primary { background: #C8922A; color: #fff; }

    .layout {
      flex: 1;
      display: flex;
      overflow: hidden;
    }


    /* Ledger area */
    .ledger-area {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }

    .ledger-topbar {
      display: flex; align-items: center; justify-content: space-between;
      padding: 18px 24px 14px;
      flex-shrink: 0;
    }
    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; color: #1a1a1a; }
    .ledger-controls { display: flex; gap: 8px; }
    .ctrl-btn { display: flex; align-items: center; gap: 5px; padding: 5px 12px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; font-size: 12px; color: #ccc; pointer-events: none; }

    /* Empty area */
    .empty-area {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: flex-start;
      padding: 32px 24px;
      overflow-y: auto;
    }

    .empty-headline {
      font-size: 15px; font-weight: 600; color: #1a1a1a;
      margin-bottom: 6px;
      text-align: center;
    }
    .empty-subtext {
      font-size: 13px; color: #aaa;
      margin-bottom: 36px;
      text-align: center;
      max-width: 500px;
      line-height: 1.6;
    }

    /* Template card grid */
    .template-grid {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      justify-content: center;
      max-width: 880px;
      width: 100%;
    }

    .template-card {
      width: 240px;
      background: #fff;
      border: 1px solid #e8e8e8;
      border-radius: 10px;
      padding: 22px 20px 18px;
      cursor: pointer;
      transition: box-shadow 0.15s, border-color 0.15s;
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
    .template-card:hover {
      border-color: #C8922A;
      box-shadow: 0 4px 16px rgba(200, 146, 42, 0.15);
    }
    .template-card:hover .tc-start { background: #C8922A; color: #fff; border-color: #C8922A; }

    .tc-icon {
      width: 38px; height: 38px;
      border-radius: 9px;
      display: flex; align-items: center; justify-content: center;
      font-size: 18px;
      margin-bottom: 4px;
    }
    .tc-type-label {
      font-size: 10px; font-weight: 600;
      letter-spacing: 0.1em; text-transform: uppercase;
      color: #bbb;
    }
    .tc-name {
      font-size: 14px; font-weight: 600; color: #1a1a1a;
      line-height: 1.3;
    }
    .tc-desc {
      font-size: 11px; color: #aaa;
      line-height: 1.5;
      flex: 1;
    }
    .tc-flow {
      display: flex; align-items: center; gap: 4px;
      padding: 6px 0 2px;
      border-top: 1px solid #f0f0f0;
    }
    .tc-vault {
      font-size: 10px; color: #888;
      background: #f5f5f5; padding: 2px 6px;
      border-radius: 4px;
    }
    .tc-arrow { font-size: 10px; color: #ddd; }
    .tc-start {
      display: block;
      width: 100%;
      padding: 8px;
      border: 1px solid #e5e5e5;
      border-radius: 7px;
      text-align: center;
      font-size: 12px; font-weight: 500;
      color: #888;
      cursor: pointer;
      margin-top: 4px;
      background: #fafafa;
    }

    /* Bottom hint */
    .ledger-footer-hint {
      padding: 14px 24px;
      border-top: 1px solid #ebebeb;
      font-size: 12px; color: #ccc;
      text-align: center;
      flex-shrink: 0;
    }
   
`;

export const C3LedgerStyles = () => <Global styles={c3LedgerStyles} />;