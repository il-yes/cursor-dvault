import { Global,css } from '@emotion/react';

export const C3styles = () => {
  return (
    <>

      <link
        href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap"
        rel="stylesheet"
      />
      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n\n    \n\n    .window {\n      width: 1320px;\n      height: 820px;\n      background: #fff;\n      border-radius: 12px;\n      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);\n      overflow: hidden;\n      display: flex;\n      flex-direction: column;\n    }\n\n    .titlebar {\n      height: 40px;\n      background: #f5f5f5;\n      border-bottom: 1px solid #e0e0e0;\n      display: flex;\n      align-items: center;\n      padding: 0 16px;\n      flex-shrink: 0;\n    }\n    .traffic-lights { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-red { background: #FF5F57; }\n    .tl-yellow { background: #FFBD2E; }\n    .tl-green { background: #28CA41; }\n    .titlebar-name {\n      flex: 1; text-align: center;\n      font-size: 12px; font-weight: 500; color: #666;\n      letter-spacing: 0.04em;\n    }\n\n    /* App topbar */\n    .topbar {\n      height: 52px;\n      background: #fff;\n      border-bottom: 1px solid #ebebeb;\n      display: flex;\n      align-items: center;\n      padding: 0 20px;\n      flex-shrink: 0;\n      gap: 16px;\n    }\n    .topbar-logo {\n      font-weight: 700;\n      font-size: 14px;\n      letter-spacing: 0.12em;\n      color: #1a1a1a;\n    }\n    .workspace-pill {\n      display: flex; align-items: center; gap: 6px;\n      background: #f5f5f5;\n      border: 1px solid #e8e8e8;\n      border-radius: 6px;\n      padding: 4px 10px;\n      font-size: 12px; color: #555;\n      cursor: pointer;\n    }\n    .workspace-pill svg { width: 12px; height: 12px; color: #999; }\n    .topbar-spacer { flex: 1; }\n    .topbar-actions { display: flex; gap: 8px; }\n\n    .btn {\n      display: inline-flex; align-items: center; gap: 6px;\n      padding: 6px 14px;\n      border-radius: 7px;\n      font-size: 12px; font-weight: 500;\n      cursor: pointer; border: none;\n    }\n    .btn-ghost {\n      background: #f5f5f5;\n      border: 1px solid #e5e5e5;\n      color: #444;\n    }\n    .btn-primary {\n      background: #C8922A;\n      color: #fff;\n    }\n\n    /* Main layout */\n    .layout {\n      flex: 1;\n      display: flex;\n      overflow: hidden;\n    }\n\n    /* Vault sidebar */\n    .sidebar {\n      \n      background: #fafafa;\n      border-right: 1px solid #ebebeb;\n      display: flex;\n      flex-direction: column;\n      padding: 18px 0 16px;\n      flex-shrink: 0;\n      overflow-y: auto;\n    }\n    .sidebar-section-label {\n      padding: 0 16px 8px;\n      font-size: 10px;\n      font-weight: 600;\n      letter-spacing: 0.1em;\n      color: #bbb;\n      text-transform: uppercase;\n    }\n    .vault-row {\n      display: flex; align-items: center; gap: 9px;\n      padding: 7px 16px;\n      cursor: pointer;\n      font-size: 13px;\n      color: #333;\n      position: relative;\n    }\n    .vault-row:hover { background: #f0f0f0; }\n    .vault-row.all-vaults {\n      font-size: 12px; color: #888;\n      padding: 5px 16px 10px;\n      border-bottom: 1px solid #ebebeb;\n      margin-bottom: 6px;\n    }\n    .vault-dot {\n      width: 8px; height: 8px;\n      border-radius: 50%;\n      flex-shrink: 0;\n    }\n    .unread-pip {\n      margin-left: auto;\n      width: 6px; height: 6px;\n      background: #C8922A;\n      border-radius: 50%;\n      flex-shrink: 0;\n    }\n    .add-vault-btn {\n      margin: 14px 14px 0;\n      padding: 8px 10px;\n      border: 1px dashed #ddd;\n      border-radius: 7px;\n      text-align: center;\n      font-size: 12px;\n      color: #aaa;\n      cursor: pointer;\n    }\n    .add-vault-btn:hover { border-color: #C8922A; color: #C8922A; }\n\n    /* Ledger area */\n    .ledger-area {\n      flex: 1;\n      display: flex;\n      flex-direction: column;\n      overflow: hidden;\n    }\n\n    .ledger-topbar {\n      display: flex;\n      align-items: center;\n      justify-content: space-between;\n      padding: 18px 24px 14px;\n      flex-shrink: 0;\n    }\n    .ledger-title {\n      font-size: 17px;\n      font-weight: 700;\n      letter-spacing: 0.06em;\n      color: #1a1a1a;\n    }\n    .ledger-controls { display: flex; gap: 8px; }\n    .ctrl-btn {\n      display: flex; align-items: center; gap: 5px;\n      padding: 5px 12px;\n      border: 1px solid #e5e5e5;\n      border-radius: 6px;\n      background: #fff;\n      font-size: 12px; color: #555;\n      cursor: pointer;\n    }\n    .ctrl-btn:hover { border-color: #ccc; }\n\n    /* Table */\n    .table-wrap {\n      flex: 1;\n      overflow-y: auto;\n      padding: 0 24px;\n    }\n\n    table {\n      width: 100%;\n      border-collapse: separate;\n      border-spacing: 0;\n    }\n\n    thead th {\n      padding: 0 10px 10px;\n      text-align: left;\n      font-size: 10px;\n      font-weight: 600;\n      letter-spacing: 0.09em;\n      text-transform: uppercase;\n      color: #bbb;\n      border-bottom: 1px solid #ebebeb;\n      white-space: nowrap;\n    }\n\n    tbody tr {\n      cursor: pointer;\n      border-bottom: 1px solid #f2f2f2;\n    }\n    tbody tr:hover { background: #fafafa; }\n    tbody tr:hover .row-hover-actions { display: flex; }\n    tbody tr:hover .stellar-val { display: none; }\n\n    td {\n      padding: 11px 10px;\n      vertical-align: middle;\n      border-bottom: 1px solid #f2f2f2;\n    }\n\n    /* Status */\n    .sdot {\n      width: 8px; height: 8px;\n      border-radius: 50%;\n      display: inline-block;\n    }\n    .s-ok { background: #22C55E; }\n    .s-pend { background: #F59E0B; }\n    .s-flag { background: #EF4444; }\n\n    /* Thread name */\n    .th-line1 {\n      font-size: 13px; font-weight: 500; color: #1a1a1a;\n      white-space: nowrap;\n    }\n    .th-type {\n      font-size: 9px; font-weight: 600;\n      letter-spacing: 0.1em;\n      text-transform: uppercase;\n      color: #bbb;\n      margin-right: 6px;\n    }\n    .th-line2 {\n      font-size: 11px; color: #aaa;\n      margin-top: 2px;\n    }\n\n    /* Flow */\n    .flow {\n      display: flex; align-items: center; gap: 4px;\n    }\n    .vb {\n      width: 22px; height: 22px;\n      border-radius: 50%;\n      display: inline-flex; align-items: center; justify-content: center;\n      font-size: 8px; font-weight: 700; color: #fff;\n      flex-shrink: 0;\n    }\n    .vb-sm { width: 18px; height: 18px; font-size: 7px; }\n    .fa { color: #ccc; font-size: 10px; font-weight: 400; }\n\n    /* Opacity states */\n    .vb-pending { opacity: 0.4; }\n\n    /* Pipeline mini */\n    .pipeline {\n      display: flex; gap: 3px; align-items: center;\n    }\n    .pseg {\n      height: 4px; width: 20px;\n      border-radius: 2px;\n    }\n    .pseg-done { background: #C8922A; }\n    .pseg-wait { background: #e8e8e8; }\n\n    /* Timestamp */\n    .ts { font-size: 12px; color: #bbb; white-space: nowrap; }\n\n    /* Stellar */\n    .stellar-val {\n      font-family: 'SF Mono', 'Fira Code', monospace;\n      font-size: 10px; color: #ccc;\n      white-space: nowrap;\n    }\n\n    /* Row hover actions */\n    .row-hover-actions {\n      display: none; gap: 5px;\n    }\n    .rha-btn {\n      padding: 3px 9px;\n      border: 1px solid #e0e0e0;\n      border-radius: 4px;\n      font-size: 11px; color: #555;\n      background: #fff;\n      cursor: pointer;\n      white-space: nowrap;\n    }\n    .rha-btn:hover { border-color: #C8922A; color: #C8922A; }\n\n    /* C3 badge */\n    .c3b {\n      display: inline-flex; align-items: center;\n      padding: 3px 8px;\n      border-radius: 5px;\n      font-size: 11px; font-weight: 600;\n      white-space: nowrap; cursor: pointer;\n    }\n    .c3-ext { background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; }\n    .c3-ext:hover { border-color: #C8922A; color: #C8922A; background: #FBF0D8; }\n    .c3-linked { background: #FBF0D8; border: 1px solid #E8C87A; color: #C8922A; }\n\n    /* Footer */\n    .ledger-footer {\n      padding: 12px 24px;\n      border-top: 1px solid #ebebeb;\n      font-size: 12px; color: #bbb;\n      flex-shrink: 0;\n    }\n\n    /* Screen label */\n    .screen-label {\n      position: absolute;\n      top: 8px; right: 16px;\n      font-size: 10px; color: #ccc;\n      letter-spacing: 0.06em;\n      font-family: 'SF Mono', monospace;\n    }\n  * { margin: 0; padding: 0; box-sizing: border-box; }\n\n    \n\n    .window {\n      width: 1320px;\n      height: 820px;\n      background: #fff;\n      border-radius: 12px;\n      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);\n      overflow: hidden;\n      display: flex;\n      flex-direction: column;\n    }\n\n    .titlebar {\n      height: 40px;\n      background: #f5f5f5;\n      border-bottom: 1px solid #e0e0e0;\n      display: flex; align-items: center;\n      padding: 0 16px;\n      flex-shrink: 0;\n    }\n    .traffic-lights { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-red { background: #FF5F57; }\n    .tl-yellow { background: #FFBD2E; }\n    .tl-green { background: #28CA41; }\n    .titlebar-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    .topbar {\n      height: 52px;\n      background: #fff;\n      border-bottom: 1px solid #ebebeb;\n      display: flex; align-items: center;\n      padding: 0 20px;\n      flex-shrink: 0; gap: 16px;\n    }\n    .topbar-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; color: #1a1a1a; }\n    .workspace-pill { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .topbar-spacer { flex: 1; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 6px 14px; border-radius: 7px; font-size: 12px; font-weight: 500; cursor: pointer; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #ccc;  }\n    .btn-primary { background: #C8922A; color: #fff; }\n\n    .layout {\n      flex: 1;\n      display: flex;\n      overflow: hidden;\n    }\n\n    /* Empty vault sidebar */\n    .sidebar {\n      \n      background: #fafafa;\n      border-right: 1px solid #ebebeb;\n      display: flex; flex-direction: column;\n      padding: 18px 0 16px;\n      flex-shrink: 0;\n    }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .sidebar-empty-hint {\n      margin: 12px 16px;\n      padding: 12px;\n      border: 1px dashed #e0e0e0;\n      border-radius: 8px;\n      font-size: 11px; color: #ccc;\n      text-align: center;\n      line-height: 1.5;\n    }\n    .add-vault-btn {\n      margin: 10px 14px 0;\n      padding: 9px 10px;\n      border: 1px dashed #ddd;\n      border-radius: 7px;\n      text-align: center;\n      font-size: 12px; color: #C8922A;\n      cursor: pointer;\n      font-weight: 500;\n    }\n    .add-vault-btn:hover { background: #FBF0D8; }\n\n    /* Ledger area */\n    .ledger-area {\n      flex: 1;\n      display: flex;\n      flex-direction: column;\n      overflow: hidden;\n    }\n\n    .ledger-topbar {\n      display: flex; align-items: center; justify-content: space-between;\n      padding: 18px 24px 14px;\n      flex-shrink: 0;\n    }\n    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; color: #1a1a1a; }\n    .ledger-controls { display: flex; gap: 8px; }\n    .ctrl-btn { display: flex; align-items: center; gap: 5px; padding: 5px 12px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; font-size: 12px; color: #ccc;  }\n\n    /* Empty area */\n    .empty-area {\n      flex: 1;\n      display: flex;\n      flex-direction: column;\n      align-items: center;\n      justify-content: flex-start;\n      padding: 32px 24px;\n      overflow-y: auto;\n    }\n\n    .empty-headline {\n      font-size: 15px; font-weight: 600; color: #1a1a1a;\n      margin-bottom: 6px;\n      text-align: center;\n    }\n    .empty-subtext {\n      font-size: 13px; color: #aaa;\n      margin-bottom: 36px;\n      text-align: center;\n      max-width: 500px;\n      line-height: 1.6;\n    }\n\n    /* Template card grid */\n    .template-grid {\n      display: flex;\n      flex-wrap: wrap;\n      gap: 16px;\n      justify-content: center;\n      max-width: 880px;\n      width: 100%;\n    }\n\n    .template-card {\n      width: 240px;\n      background: #fff;\n      border: 1px solid #e8e8e8;\n      border-radius: 10px;\n      padding: 22px 20px 18px;\n      cursor: pointer;\n      transition: box-shadow 0.15s, border-color 0.15s;\n      display: flex;\n      flex-direction: column;\n      gap: 10px;\n    }\n    .template-card:hover {\n      border-color: #C8922A;\n      box-shadow: 0 4px 16px rgba(200, 146, 42, 0.15);\n    }\n    .template-card:hover .tc-start { background: #C8922A; color: #fff; border-color: #C8922A; }\n\n    .tc-icon {\n      width: 38px; height: 38px;\n      border-radius: 9px;\n      display: flex; align-items: center; justify-content: center;\n      font-size: 18px;\n      margin-bottom: 4px;\n    }\n    .tc-type-label {\n      font-size: 10px; font-weight: 600;\n      letter-spacing: 0.1em; text-transform: uppercase;\n      color: #bbb;\n    }\n    .tc-name {\n      font-size: 14px; font-weight: 600; color: #1a1a1a;\n      line-height: 1.3;\n    }\n    .tc-desc {\n      font-size: 11px; color: #aaa;\n      line-height: 1.5;\n      flex: 1;\n    }\n    .tc-flow {\n      display: flex; align-items: center; gap: 4px;\n      padding: 6px 0 2px;\n      border-top: 1px solid #f0f0f0;\n    }\n    .tc-vault {\n      font-size: 10px; color: #888;\n      background: #f5f5f5; padding: 2px 6px;\n      border-radius: 4px;\n    }\n    .tc-arrow { font-size: 10px; color: #ddd; }\n    .tc-start {\n      display: block;\n      width: 100%;\n      padding: 8px;\n      border: 1px solid #e5e5e5;\n      border-radius: 7px;\n      text-align: center;\n      font-size: 12px; font-weight: 500;\n      color: #888;\n      cursor: pointer;\n      margin-top: 4px;\n      background: #fafafa;\n    }\n\n    /* Bottom hint */\n    .ledger-footer-hint {\n      padding: 14px 24px;\n      border-top: 1px solid #ebebeb;\n      font-size: 12px; color: #ccc;\n      text-align: center;\n      flex-shrink: 0;\n    }\n  "
        }}
      />
      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    .topbar { height: 52px; background: #fff; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; flex-shrink: 0; }\n    .tb-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .tb-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .tb-spacer { flex: 1; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 6px 14px; border-radius: 7px; font-size: 12px; font-weight: 500; border: none; cursor: pointer; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #444; }\n    .btn-primary { background: #C8922A; color: #fff; }\n\n    .layout { flex: 1; display: flex; overflow: hidden; }\n\n    /* Sidebar */\n    .sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; display: flex; flex-direction: column; padding: 18px 0 16px; flex-shrink: 0; }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .nav-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #555; cursor: pointer; border-radius: 0; }\n    .nav-row:hover { background: #f0f0f0; }\n    .nav-row.active { background: #FBF0D8; color: #C8922A; font-weight: 600; }\n    .nav-icon { font-size: 13px; width: 18px; text-align: center; flex-shrink: 0; }\n    .nav-label { flex: 1; }\n    .nav-badge { background: #C8922A; color: #fff; font-size: 10px; font-weight: 700; padding: 1px 6px; border-radius: 10px; min-width: 18px; text-align: center; }\n    .nav-badge-gray { background: #e8e8e8; color: #888; font-size: 10px; font-weight: 600; padding: 1px 6px; border-radius: 10px; }\n    .sidebar-divider { border: none; border-top: 1px solid #ebebeb; margin: 10px 0; }\n\n    .vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #333; cursor: pointer; }\n    .vault-row:hover { background: #f0f0f0; }\n    .vault-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .unread-pip { margin-left: auto; width: 6px; height: 6px; background: #C8922A; border-radius: 50%; }\n\n    /* Main inbox area */\n    .inbox-main { flex: 1; display: flex; flex-direction: column; overflow: hidden; }\n\n    .inbox-topbar { display: flex; align-items: center; justify-content: space-between; padding: 20px 28px 16px; flex-shrink: 0; border-bottom: 1px solid #f0f0f0; }\n    .inbox-title-row { display: flex; align-items: center; gap: 12px; }\n    .inbox-title { font-size: 17px; font-weight: 700; letter-spacing: 0.04em; color: #1a1a1a; }\n    .inbox-subtitle { font-size: 12px; color: #bbb; }\n    .inbox-count-pill { background: #FBF0D8; border: 1px solid #E8C87A; color: #C8922A; font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 10px; }\n    .inbox-controls { display: flex; gap: 8px; align-items: center; }\n    .ctrl-btn { display: flex; align-items: center; gap: 5px; padding: 5px 12px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; font-size: 12px; color: #555; cursor: pointer; }\n    .ctrl-btn:hover { border-color: #C8922A; color: #C8922A; }\n\n    /* Inbox list */\n    .inbox-list { flex: 1; overflow-y: auto; padding: 0; }\n\n    /* Vault group */\n    .vault-group { margin-bottom: 0; }\n    .vg-header { display: flex; align-items: center; gap: 10px; padding: 14px 28px 10px; background: #fafafa; border-bottom: 1px solid #f0f0f0; position: sticky; top: 0; z-index: 1; }\n    .vg-dot { width: 9px; height: 9px; border-radius: 50%; flex-shrink: 0; }\n    .vg-name { font-size: 12px; font-weight: 700; letter-spacing: 0.07em; text-transform: uppercase; color: #444; }\n    .vg-count { font-size: 11px; color: #bbb; }\n    .vg-dot-spacer { flex: 1; }\n\n    /* Inbox item */\n    .inbox-item { display: flex; align-items: center; gap: 0; padding: 0 28px; border-bottom: 1px solid #f5f5f5; cursor: pointer; min-height: 68px; }\n    .inbox-item:hover { background: #fafafa; }\n    .inbox-item.urgent { border-left: 3px solid #EF4444; padding-left: 25px; }\n    .inbox-item.pending { border-left: 3px solid #C8922A; padding-left: 25px; }\n    .inbox-item.c3 { border-left: 3px solid #7C3AED; padding-left: 25px; }\n\n    .ii-left { display: flex; flex-direction: column; justify-content: center; margin-right: 16px; flex-shrink: 0; }\n    .ii-priority { width: 8px; height: 8px; border-radius: 50%; }\n    .pri-urgent { background: #EF4444; }\n    .pri-normal { background: #C8922A; }\n    .pri-low { background: #bbb; }\n\n    .ii-main { flex: 1; padding: 14px 0; }\n    .ii-slot { display: flex; align-items: center; gap: 7px; margin-bottom: 4px; }\n    .slot-tag { font-family: 'SF Mono', monospace; font-size: 12px; font-weight: 600; color: #1a1a1a; }\n    .gate-ok-badge { font-size: 9px; color: #059669; font-weight: 600; background: #F0FDF4; border: 1px solid #86EFAC; padding: 1px 5px; border-radius: 3px; }\n    .gate-blocked-badge { font-size: 9px; color: #D97706; font-weight: 600; background: #FFF7ED; border: 1px solid #FED7AA; padding: 1px 5px; border-radius: 3px; }\n    .c3-badge { font-size: 9px; color: #7C3AED; font-weight: 600; background: #EEF2FF; border: 1px solid #C7D2FE; padding: 1px 5px; border-radius: 3px; }\n    .ii-thread { display: flex; align-items: center; gap: 8px; }\n    .ii-thread-type { font-size: 9px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #C8922A; background: #FBF0D8; padding: 1px 5px; border-radius: 3px; }\n    .ii-thread-name { font-size: 12px; color: #555; }\n    .ii-since { font-size: 11px; color: #bbb; margin-top: 4px; }\n\n    .ii-right { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; padding: 14px 0; flex-shrink: 0; }\n    .ii-ts { font-size: 11px; color: #bbb; white-space: nowrap; }\n    .commit-btn { display: inline-flex; align-items: center; gap: 5px; padding: 6px 14px; background: #C8922A; color: #fff; border-radius: 6px; font-size: 11px; font-weight: 600; border: none; cursor: pointer; white-space: nowrap; }\n    .commit-btn:hover { background: #b8821a; }\n    .view-btn { display: inline-flex; align-items: center; gap: 5px; padding: 6px 12px; background: #f5f5f5; border: 1px solid #e5e5e5; color: #555; border-radius: 6px; font-size: 11px; font-weight: 500; cursor: pointer; }\n    .view-btn:hover { border-color: #C8922A; color: #C8922A; }\n    .blocked-text { font-size: 11px; color: #D97706; background: #FFF7ED; padding: 5px 10px; border-radius: 5px; }\n\n    /* Empty state for vault group */\n    .vg-empty { padding: 16px 28px; font-size: 12px; color: #ccc; font-style: italic; }\n\n    /* Footer */\n    .inbox-footer { padding: 12px 28px; border-top: 1px solid #ebebeb; font-size: 12px; color: #bbb; flex-shrink: 0; display: flex; align-items: center; justify-content: space-between; }\n  "
        }}
      />
      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; position: relative; }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    /* Dimmed bg */\n    .bg-app { position: absolute; inset: 40px 0 0 0; opacity: 0.14;  display: flex; flex-direction: column; }\n    .bg-topbar { height: 52px; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; }\n    .bg-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .bg-sp { flex: 1; }\n    .bg-body { flex: 1; display: flex; }\n    .bg-sidebar {background: #fafafa; border-right: 1px solid #ebebeb; padding: 18px 0; }\n    .bg-sidebar-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .bg-vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; }\n    .bg-dot { width: 8px; height: 8px; border-radius: 50%; }\n    .bg-main { flex: 1; padding: 18px 24px; }\n    .bg-ledger-title { font-size: 17px; font-weight: 700; margin-bottom: 20px; }\n    .bg-row { display: flex; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f2f2f2; }\n    .bg-sdot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .bg-tname { font-size: 13px; font-weight: 500; flex: 1; }\n    .bg-ts { font-size: 11px; color: #aaa; }\n\n    .scrim { position: absolute; inset: 40px 0 0 0; background: rgba(0,0,0,0.32); z-index: 2; }\n    .modal-wrap { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; z-index: 3; }\n    .modal { width: 560px; background: #fff; border-radius: 12px; box-shadow: 0 24px 80px rgba(0,0,0,0.30), 0 8px 24px rgba(0,0,0,0.16); overflow: hidden; display: flex; flex-direction: column; }\n\n    /* Header */\n    .modal-header { padding: 20px 24px 16px; border-bottom: 1px solid #ebebeb; }\n    .mh-top { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 10px; }\n    .vault-badge { display: inline-flex; align-items: center; gap: 5px; padding: 3px 10px; border-radius: 5px; font-size: 11px; font-weight: 600; }\n    .vb-dir { background: #F5F5F5; border: 1px solid #e0e0e0; color: #333; }\n    .vbdot { width: 7px; height: 7px; border-radius: 50%; }\n    .mh-close { width: 24px; height: 24px; background: #f5f5f5; border: 1px solid #e5e5e5; border-radius: 6px; display: flex; align-items: center; justify-content: center; font-size: 12px; color: #888; cursor: pointer; }\n    .mh-title { font-size: 16px; font-weight: 700; color: #1a1a1a; }\n    .mh-sub { font-size: 12px; color: #aaa; margin-top: 3px; }\n\n    /* Thread context */\n    .thread-ctx { display: flex; align-items: center; gap: 10px; padding: 9px 12px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 7px; margin-top: 12px; }\n    .tc-type { font-size: 9px; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: #C8922A; background: #FBF0D8; padding: 2px 6px; border-radius: 3px; flex-shrink: 0; }\n    .tc-name { font-size: 12px; font-weight: 500; color: #333; flex: 1; }\n    .tc-tid { font-family: 'SF Mono', monospace; font-size: 10px; color: #bbb; }\n\n    /* Modal body */\n    .modal-body { padding: 20px 24px; display: flex; flex-direction: column; gap: 16px; }\n\n    /* Received commit box */\n    .commit-box { background: #fafafa; border: 1px solid #ebebeb; border-radius: 9px; padding: 14px 16px; }\n    .cb-header { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; }\n    .cb-from-badge { display: inline-flex; align-items: center; gap: 5px; background: #EEF2FF; border: 1px solid #C7D2FE; padding: 3px 9px; border-radius: 5px; font-size: 11px; font-weight: 600; color: #3730A3; }\n    .cb-from-dot { width: 7px; height: 7px; border-radius: 50%; background: #2563EB; }\n    .cb-arrow { font-size: 11px; color: #ccc; }\n    .cb-slot { font-family: 'SF Mono', monospace; font-size: 12px; font-weight: 600; color: #1a1a1a; }\n    .cb-ts { font-size: 11px; color: #bbb; margin-left: auto; }\n    .cb-value { font-size: 12px; color: #444; line-height: 1.5; padding: 10px 12px; background: #fff; border: 1px solid #e8e8e8; border-radius: 6px; margin-bottom: 10px; font-style: italic; }\n    .cb-value::before { content: '\"'; color: #bbb; } .cb-value::after { content: '\"'; color: #bbb; }\n    .cb-meta { display: flex; align-items: center; gap: 12px; }\n    .cb-cid { font-family: 'SF Mono', monospace; font-size: 10px; color: #bbb; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }\n    .cb-verify { font-size: 10px; color: #C8922A; cursor: pointer; white-space: nowrap; }\n\n    /* Action options */\n    .action-options { display: flex; flex-direction: column; gap: 8px; }\n    .action-opt { display: flex; align-items: flex-start; gap: 12px; padding: 12px 14px; border: 1.5px solid #e5e5e5; border-radius: 9px; cursor: pointer; transition: border-color 0.15s; }\n    .action-opt:hover { border-color: #C8922A; background: #FFFDF8; }\n    .action-opt.selected-ok { border-color: #22C55E; background: #F0FDF4; }\n    .action-opt.selected-reject { border-color: #EF4444; background: #FEF2F2; }\n    .ao-radio { width: 16px; height: 16px; border-radius: 50%; border: 2px solid #d0d0d0; flex-shrink: 0; margin-top: 1px; display: flex; align-items: center; justify-content: center; }\n    .ao-radio.ok-active { border-color: #22C55E; background: #22C55E; }\n    .ao-radio.ok-active::after { content: ''; width: 6px; height: 6px; border-radius: 50%; background: #fff; }\n    .ao-radio.reject-active { border-color: #EF4444; background: #EF4444; }\n    .ao-radio.reject-active::after { content: ''; width: 6px; height: 6px; border-radius: 50%; background: #fff; }\n    .ao-content { flex: 1; }\n    .ao-label { font-size: 13px; font-weight: 600; color: #1a1a1a; }\n    .ao-label.ok { color: #059669; }\n    .ao-label.reject { color: #DC2626; }\n    .ao-desc { font-size: 11px; color: #aaa; margin-top: 2px; line-height: 1.4; }\n\n    /* Rejection reason */\n    .rejection-reason { margin-top: 10px; }\n    .rr-label { font-size: 10px; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; color: #bbb; margin-bottom: 5px; }\n    .rr-code-select { width: 100%; padding: 8px 12px; border: 1.5px solid #EF4444; border-radius: 7px; font-size: 12px; color: #1a1a1a; outline: none; font-family: 'Inter', sans-serif; background: #fff; appearance: none; background-image: url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8' viewBox='0 0 12 8'%3E%3Cpath d='M1 1l5 5 5-5' stroke='%23EF4444' stroke-width='1.5' fill='none'/%3E%3C/svg%3E\"); background-repeat: no-repeat; background-position: right 10px center; cursor: pointer; }\n    .rr-code-select:focus { background-color: #FEF2F2; }\n    .rr-code-hint { font-size: 10px; color: #EF4444; opacity: 0.7; margin-top: 3px; margin-bottom: 8px; }\n    textarea.rr-input { width: 100%; min-height: 64px; padding: 9px 12px; border: 1.5px solid #EF4444; border-radius: 7px; font-size: 12px; color: #1a1a1a; outline: none; resize: none; font-family: 'Inter', sans-serif; line-height: 1.5; background: #fff; }\n    textarea.rr-input:focus { background: #FEF2F2; }\n    .rr-note { font-size: 10px; color: #bbb; margin-top: 5px; }\n    /* Irreversible warning */\n    .irreversible-warn { display: flex; align-items: center; gap: 7px; padding: 8px 12px; background: #FFF7ED; border: 1px solid #FED7AA; border-radius: 6px; }\n    .iw-icon { font-size: 12px; flex-shrink: 0; }\n    .iw-text { font-size: 11px; color: #92400E; line-height: 1.4; }\n\n    /* Stellar preview */\n    .stellar-row { display: flex; align-items: center; gap: 8px; padding: 9px 12px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 7px; }\n    .sr-icon { font-size: 11px; }\n    .sr-text { font-size: 11px; color: #888; flex: 1; }\n    .sr-text strong { color: #1a1a1a; }\n    .sr-dot { width: 5px; height: 5px; border-radius: 50%; background: #22C55E; }\n\n    /* Footer */\n    .modal-footer { padding: 14px 24px; border-top: 1px solid #ebebeb; display: flex; align-items: center; justify-content: space-between; gap: 10px; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 18px; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #555; }\n    .btn-reject { background: #EF4444; color: #fff; font-size: 13px; font-weight: 600; padding: 10px 22px; border-radius: 8px; border: none; cursor: pointer; display: flex; align-items: center; gap: 7px; }\n    .btn-reject:hover { background: #DC2626; }\n  "
        }}
      />
      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n   .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    .topbar { height: 52px; background: #fff; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; flex-shrink: 0; }\n    .tb-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .tb-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .tb-spacer { flex: 1; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 6px 14px; border-radius: 7px; font-size: 12px; font-weight: 500; border: none; cursor: pointer; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #444; }\n    .btn-primary { background: #C8922A; color: #fff; }\n\n    /* Toast notification */\n    .toast-bar { background: #1a1a1a; color: #fff; display: flex; align-items: center; gap: 12px; padding: 10px 20px; font-size: 12px; flex-shrink: 0; }\n    .toast-icon { font-size: 14px; }\n    .toast-text { flex: 1; }\n    .toast-text strong { color: #C8922A; }\n    .toast-action { color: #C8922A; font-weight: 600; cursor: pointer; white-space: nowrap; }\n    .toast-close { color: #666; cursor: pointer; font-size: 14px; }\n\n    .layout { flex: 1; display: flex; overflow: hidden; }\n\n    .sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; display: flex; flex-direction: column; padding: 18px 0 16px; flex-shrink: 0; }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #333; }\n    .vault-row.all-vaults { font-size: 12px; color: #888; padding: 5px 16px 10px; border-bottom: 1px solid #ebebeb; margin-bottom: 6px; }\n    .vault-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .unread-pip { margin-left: auto; width: 6px; height: 6px; background: #C8922A; border-radius: 50%; }\n    .new-channel-btn { margin: 10px 14px 0; padding: 7px; border: 1.5px dashed #C8922A; border-radius: 7px; text-align: center; font-size: 12px; color: #C8922A; font-weight: 500; cursor: pointer; }\n    .new-channel-btn:hover { background: #FBF0D8; }\n\n    .ledger-area { flex: 1; display: flex; flex-direction: column; overflow: hidden; }\n    .ledger-topbar { display: flex; align-items: center; justify-content: space-between; padding: 18px 24px 14px; flex-shrink: 0; }\n    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; }\n    .ledger-controls { display: flex; gap: 8px; }\n    .ctrl-btn { display: flex; align-items: center; gap: 5px; padding: 5px 12px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; font-size: 12px; color: #555; }\n\n    .table-wrap { flex: 1; overflow-y: auto; padding: 0 24px; }\n    table { width: 100%; border-collapse: separate; border-spacing: 0; }\n    thead th { padding: 0 10px 10px; text-align: left; font-size: 10px; font-weight: 600; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; border-bottom: 1px solid #ebebeb; position: sticky; top: 0; background: #fff; z-index: 1; }\n    td { padding: 11px 10px; vertical-align: middle; border-bottom: 1px solid #f2f2f2; font-size: 13px; }\n\n    /* New row styles */\n    .row-new td { background: #FFFDF8; }\n    .row-new-left-border { position: relative; }\n    .row-new-left-border td:first-child { border-left: 3px solid #C8922A; }\n    .new-label { display: inline-block; font-size: 9px; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: #fff; background: #C8922A; padding: 1px 6px; border-radius: 3px; vertical-align: middle; margin-left: 4px; animation: fadeIn 0.3s ease; }\n    @keyframes fadeIn { from { opacity: 0; transform: translateY(-4px); } to { opacity: 1; transform: translateY(0); } }\n\n    .sdot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }\n    .s-ok { background: #22C55E; } .s-pend { background: #F59E0B; } .s-new { background: #F59E0B; animation: pulse 1.5s infinite; }\n    @keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.4; } }\n\n    .th-line1 { font-size: 13px; font-weight: 500; }\n    .th-type { font-size: 9px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-right: 6px; }\n    .th-line1.new-thread { color: #C8922A; }\n    .th-line2 { font-size: 11px; color: #aaa; margin-top: 2px; }\n\n    .flow { display: flex; align-items: center; gap: 4px; }\n    .vb { width: 22px; height: 22px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 8px; font-weight: 700; color: #fff; flex-shrink: 0; }\n    .fa { color: #ccc; font-size: 10px; }\n\n    .pipeline { display: flex; gap: 3px; }\n    .pseg { height: 4px; width: 20px; border-radius: 2px; }\n    .pseg-done { background: #C8922A; }\n    .pseg-wait { background: #e8e8e8; }\n    .pseg-new { background: #F5E6C8; border: 1px dashed #C8922A; }\n\n    .ts { font-size: 12px; color: #bbb; }\n    .stellar-val { font-family: 'SF Mono', monospace; font-size: 10px; color: #ccc; }\n    .stellar-genesis { font-family: 'SF Mono', monospace; font-size: 10px; color: #C8922A; }\n\n    .c3b { display: inline-flex; align-items: center; padding: 3px 8px; border-radius: 5px; font-size: 11px; font-weight: 600; }\n    .c3-ext { background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; }\n    .c3-linked { background: #FBF0D8; border: 1px solid #E8C87A; color: #C8922A; }\n\n    /* First slot hint on new row */\n    .first-slot-hint { display: inline-flex; align-items: center; gap: 5px; font-size: 10px; color: #C8922A; font-weight: 500; margin-top: 3px; }\n    .fsh-dot { width: 5px; height: 5px; border-radius: 50%; background: #C8922A; animation: pulse 1.5s infinite; }\n\n    .ledger-footer { padding: 12px 24px; border-top: 1px solid #ebebeb; font-size: 12px; color: #bbb; flex-shrink: 0; }\n  "
        }}
      />
      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; position: relative; }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    /* Topbar dimmed */\n    .topbar { height: 52px; background: #fff; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; flex-shrink: 0; opacity: 0.45;  }\n    .tb-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .tb-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .tb-spacer { flex: 1; }\n    .btn-g { background: #f5f5f5; border: 1px solid #e5e5e5; color: #444; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n    .btn-p { background: #C8922A; color: #fff; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n\n    .layout { flex: 1; display: flex; overflow: hidden; }\n\n    /* Dimmed ledger wrapper */\n    .ledger-wrapper { flex: 1; display: flex; overflow: hidden; opacity: 0.35;  }\n    .sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; display: flex; flex-direction: column; padding: 18px 0 16px; flex-shrink: 0; }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #333; }\n    .vault-row.all-vaults { font-size: 12px; color: #888; padding: 5px 16px 10px; border-bottom: 1px solid #ebebeb; margin-bottom: 6px; }\n    .vault-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .unread-pip { margin-left: auto; width: 6px; height: 6px; background: #C8922A; border-radius: 50%; }\n    .ledger-area { flex: 1; display: flex; flex-direction: column; overflow: hidden; }\n    .ledger-topbar { display: flex; align-items: center; justify-content: space-between; padding: 18px 24px 14px; flex-shrink: 0; }\n    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; }\n    .table-wrap { flex: 1; overflow: hidden; padding: 0 24px; }\n    table { width: 100%; border-collapse: separate; border-spacing: 0; }\n    thead th { padding: 0 10px 10px; text-align: left; font-size: 10px; font-weight: 600; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; border-bottom: 1px solid #ebebeb; }\n    td { padding: 11px 10px; vertical-align: middle; border-bottom: 1px solid #f2f2f2; font-size: 13px; }\n    .sdot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }\n    .s-ok { background: #22C55E; } .s-pend { background: #F59E0B; }\n    .th-type { font-size: 9px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-right: 6px; }\n    .flow { display: flex; align-items: center; gap: 4px; }\n    .vb { width: 22px; height: 22px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 8px; font-weight: 700; color: #fff; }\n    .fa { color: #ccc; font-size: 10px; }\n    .pipeline { display: flex; gap: 3px; }\n    .pseg { height: 4px; width: 20px; border-radius: 2px; }\n    .pseg-done { background: #C8922A; } .pseg-wait { background: #e8e8e8; }\n    .ts { font-size: 12px; color: #bbb; }\n    .stellar-val { font-family: 'SF Mono', monospace; font-size: 10px; color: #ccc; }\n    .c3b { display: inline-flex; align-items: center; padding: 3px 8px; border-radius: 5px; font-size: 11px; font-weight: 600; }\n    .c3-ext { background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; }\n    .c3-linked { background: #FBF0D8; border: 1px solid #E8C87A; color: #C8922A; }\n    .ledger-footer { padding: 12px 24px; border-top: 1px solid #ebebeb; font-size: 12px; color: #bbb; flex-shrink: 0; }\n\n    /* ===== SLIDE PANEL ===== */\n    .slide-panel { margin-top: -1rem ; height: 100%; background: #fff; border-left: 1px solid #e0e0e0; display: flex; flex-direction: column; flex-shrink: 0; box-shadow: -8px 0 32px rgba(0,0,0,0.10); z-index: 10; overflow: hidden; }\n\n    /* Panel header */\n    .sp-header { padding: 20px 24px 16px; border-bottom: 1px solid #ebebeb; flex-shrink: 0; }\n    .sp-header-row { display: flex; align-items: center; justify-content: space-between; }\n    .sp-title { font-size: 15px; font-weight: 700; color: #1a1a1a; letter-spacing: 0.01em; }\n    .sp-subtitle { font-size: 12px; color: #bbb; margin-top: 3px; }\n    .sp-close { width: 26px; height: 26px; border-radius: 7px; background: #f5f5f5; border: 1px solid #e5e5e5; display: flex; align-items: center; justify-content: center; font-size: 13px; color: #888; cursor: pointer; flex-shrink: 0; }\n\n    /* Panel body */\n    .sp-body { flex: 1; overflow-y: auto; padding: 20px 24px; display: flex; flex-direction: column; gap: 20px; }\n\n    /* Field label style */\n    .fl { font-size: 10px; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; margin-bottom: 6px; }\n    .fl-hint { font-size: 11px; color: #bbb; font-weight: 400; margin-left: 6px; text-transform: none; letter-spacing: 0; }\n\n    /* Channel selector */\n    .channel-selected { display: flex; align-items: center; gap: 10px; padding: 12px 14px; border: 1.5px solid #C8922A; border-radius: 9px; background: #FFFDF8; cursor: pointer; }\n    .cs-icon-wrap { width: 34px; height: 34px; border-radius: 8px; background: #FBF0D8; display: flex; align-items: center; justify-content: center; font-size: 15px; flex-shrink: 0; }\n    .cs-content { flex: 1; }\n    .cs-name { font-size: 13px; font-weight: 600; color: #1a1a1a; }\n    .cs-desc { font-size: 11px; color: #bbb; margin-top: 2px; }\n    .cs-arrow { font-size: 11px; color: #C8922A; }\n\n    /* Channel flow preview */\n    .channel-flow-box { background: #fafafa; border: 1px solid #ebebeb; border-radius: 8px; padding: 14px 16px; }\n    .cfb-row { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }\n    .cfb-vault { display: flex; align-items: center; gap: 5px; padding: 5px 9px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; font-size: 12px; font-weight: 500; color: #333; }\n    .cfb-dot { width: 7px; height: 7px; border-radius: 50%; }\n    .cfb-arrow { color: #ccc; font-size: 12px; }\n    .cfb-meta { display: flex; gap: 12px; margin-top: 10px; padding-top: 10px; border-top: 1px solid #f0f0f0; }\n    .cfb-metaitem { display: flex; align-items: center; gap: 5px; font-size: 11px; color: #bbb; }\n    .cfb-metaitem strong { color: #555; }\n\n    /* Thread name input */\n    .thread-name-wrap { display: flex; align-items: center; border: 1.5px solid #e5e5e5; border-radius: 8px; overflow: hidden; }\n    .thread-name-wrap:focus-within { border-color: #C8922A; background: #FFFDF8; }\n    .thread-name-prefix { padding: 10px 12px; background: #fafafa; border-right: 1px solid #e5e5e5; font-size: 12px; color: #bbb; white-space: nowrap; font-family: 'SF Mono', monospace; }\n    .thread-name-input { flex: 1; padding: 10px 12px; border: none; outline: none; font-size: 13px; color: #1a1a1a; font-family: 'Inter', sans-serif; background: transparent; }\n\n    /* Property fields */\n    .props-grid { display: flex; flex-direction: column; gap: 10px; }\n    .prop-row { display: flex; gap: 8px; align-items: flex-end; }\n    .prop-key-wrap { flex: 1; }\n    .prop-val-wrap { flex: 2; }\n    .prop-label { font-size: 10px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: #bbb; margin-bottom: 4px; }\n    .prop-input { width: 100%; padding: 8px 10px; border: 1px solid #e5e5e5; border-radius: 7px; font-size: 12px; color: #1a1a1a; outline: none; font-family: 'Inter', sans-serif; background: #fff; }\n    .prop-input:focus { border-color: #C8922A; background: #FFFDF8; }\n    .prop-input.prefilled { background: #fafafa; color: #888; }\n    .prop-remove { width: 26px; height: 32px; display: flex; align-items: center; justify-content: center; font-size: 14px; color: #ccc; cursor: pointer; flex-shrink: 0; border-radius: 5px; }\n    .prop-remove:hover { color: #EF4444; background: #FEE2E2; }\n    .add-prop-btn { display: inline-flex; align-items: center; gap: 5px; font-size: 11px; color: #C8922A; cursor: pointer; padding: 4px 0; }\n    .add-prop-btn:hover { opacity: 0.8; }\n\n    /* Vault overrides */\n    .vault-overrides { display: flex; flex-direction: column; gap: 8px; }\n    .vault-override-row { display: flex; align-items: center; gap: 10px; }\n    .vor-role { font-size: 12px; color: #888; width: 68px; flex-shrink: 0; }\n    .vor-select { flex: 1; display: flex; align-items: center; gap: 7px; padding: 7px 10px; border: 1px solid #e5e5e5; border-radius: 7px; background: #fafafa; cursor: pointer; }\n    .vor-select:hover { border-color: #C8922A; }\n    .vor-dot { width: 7px; height: 7px; border-radius: 50%; }\n    .vor-name { font-size: 12px; font-weight: 500; color: #333; flex: 1; }\n    .vor-arrow { font-size: 10px; color: #bbb; }\n\n    /* Stellar info bar */\n    .stellar-info { display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 7px; }\n    .si-icon { font-size: 12px; }\n    .si-text { font-size: 11px; color: #888; flex: 1; line-height: 1.4; }\n    .si-text strong { color: #1a1a1a; }\n    .si-status { display: flex; align-items: center; gap: 5px; }\n    .si-dot { width: 6px; height: 6px; border-radius: 50%; background: #22C55E; }\n    .si-label { font-size: 10px; color: #22C55E; font-weight: 600; }\n\n    /* Footer actions */\n    .sp-footer { padding: 16px 24px; border-top: 1px solid #ebebeb; flex-shrink: 0; display: flex; flex-direction: column; gap: 10px; }\n    .start-btn { width: 100%; padding: 12px; background: #C8922A; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 700; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 8px; letter-spacing: 0.02em; }\n    .start-btn:hover { background: #b8821a; }\n    .footer-note { font-size: 11px; color: #bbb; text-align: center; line-height: 1.4; }\n    .footer-note strong { color: #888; }\n  "
        }}
      />

      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window {\n      width: 1320px; height: 820px;\n      background: #fff; border-radius: 12px;\n      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);\n      overflow: hidden; display: flex; flex-direction: column;\n      position: relative;\n    }\n\n    /* macOS chrome */\n    .titlebar {\n      height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0;\n      display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1;\n    }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    /* Dimmed background app */\n    .bg-app { position: absolute; inset: 40px 0 0 0; opacity: 0.18;  display: flex; flex-direction: column; }\n    .bg-topbar { height: 52px; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; }\n    .bg-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .bg-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .bg-sp { flex: 1; }\n    .bg-body { flex: 1; display: flex; }\n    .bg-sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; padding: 18px 0; }\n    .bg-sidebar-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .bg-vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; }\n    .bg-dot { width: 8px; height: 8px; border-radius: 50%; }\n    .bg-main { flex: 1; padding: 18px 24px; }\n    .bg-ledger-title { font-size: 17px; font-weight: 700; margin-bottom: 20px; }\n    .bg-row { display: flex; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f2f2f2; }\n    .bg-sdot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .bg-tname { font-size: 13px; font-weight: 500; flex: 1; }\n    .bg-ts { font-size: 11px; color: #aaa; }\n\n    /* Dark scrim */\n    .scrim { position: absolute; inset: 40px 0 0 0; background: rgba(0,0,0,0.28); z-index: 2; }\n\n    /* Modal */\n    .modal-wrap { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; z-index: 3; }\n    .modal {\n      width: 640px;\n      background: #fff;\n      border-radius: 12px;\n      box-shadow: 0 24px 80px rgba(0,0,0,0.28), 0 8px 24px rgba(0,0,0,0.14);\n      overflow: hidden;\n      display: flex; flex-direction: column;\n    }\n\n    /* Modal header */\n    .modal-header {\n      padding: 22px 26px 18px;\n      border-bottom: 1px solid #ebebeb;\n      display: flex; align-items: flex-start; justify-content: space-between;\n    }\n    .modal-title { font-size: 13px; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; color: #1a1a1a; }\n    .step-indicator { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }\n    .step-label { font-size: 11px; color: #bbb; }\n    .step-dots { display: flex; gap: 5px; }\n    .sdot-i { width: 7px; height: 7px; border-radius: 50%; background: #e5e5e5; }\n    .sdot-i.active { background: #C8922A; }\n\n    /* Modal body */\n    .modal-body { padding: 22px 26px; flex: 1; overflow-y: auto; }\n\n    /* Input */\n    .field-label { font-size: 11px; font-weight: 600; color: #888; letter-spacing: 0.06em; text-transform: uppercase; margin-bottom: 8px; }\n    .name-input-wrap {\n      display: flex; align-items: center;\n    padding: 10px 14px;\n      gap: 10px;\n      margin-bottom: 22px;\n      background: #FFFDF8;\n    }\n    .name-input {\n      flex: 1; font-size: 14px; font-weight: 500; color: #1a1a1a;\n      border: none; outline: none; background: transparent;\n      font-family: 'Inter', sans-serif;\n    }\n    .name-dropdown-arrow { color: #ccc; font-size: 11px; }\n\n    /* Template grid */\n    .templates-label {\n      font-size: 11px; font-weight: 600; color: #888;\n      letter-spacing: 0.06em; text-transform: uppercase;\n      margin-bottom: 12px;\n    }\n       .tpl-card {\n      border: 1.5px solid #e8e8e8;\n      border-radius: 8px;\n      padding: 12px 14px;\n      cursor: pointer;\n      transition: border-color 0.15s;\n    }\n    .tpl-card:hover { border-color: #C8922A; }\n    .tpl-card.selected {\n      border-color: #C8922A;\n      background: #FFFDF8;\n    }\n    .tpl-card-name {\n      font-size: 12px; font-weight: 600; color: #1a1a1a;\n      margin-bottom: 4px;\n      display: flex; align-items: center; gap: 6px;\n    }\n    .tpl-check { color: #C8922A; font-size: 11px; }\n    .tpl-flow { font-size: 10px; color: #aaa; }\n    .tpl-card-single { grid-column: 1; }\n\n    /* Blank hint */\n    .blank-hint { font-size: 11px; color: #bbb; margin-top: 4px; }\n\n    /* Modal footer */\n    .modal-footer {\n      padding: 16px 26px;\n      border-top: 1px solid #ebebeb;\n      display: flex; align-items: center; justify-content: flex-end;\n      gap: 10px;\n    }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 18px; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #555; }\n    .btn-primary { background: #C8922A; color: #fff; }\n  "
        }}
      />

      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n     .window {\n      width: 1320px; height: 820px;\n      background: #fff; border-radius: 12px;\n      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);\n      overflow: hidden; display: flex; flex-direction: column;\n      position: relative;\n    }\n\n    /* macOS chrome */\n    .titlebar {\n      height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0;\n      display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1;\n    }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    /* Dimmed background app */\n    .bg-app { position: absolute; inset: 40px 0 0 0; opacity: 0.18;  display: flex; flex-direction: column; }\n    .bg-topbar { height: 52px; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; }\n    .bg-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .bg-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .bg-sp { flex: 1; }\n    .bg-body { flex: 1; display: flex; }\n    .bg-sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; padding: 18px 0; }\n    .bg-sidebar-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .bg-vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; }\n    .bg-dot { width: 8px; height: 8px; border-radius: 50%; }\n    .bg-main { flex: 1; padding: 18px 24px; }\n    .bg-ledger-title { font-size: 17px; font-weight: 700; margin-bottom: 20px; }\n    .bg-row { display: flex; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f2f2f2; }\n    .bg-sdot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .bg-tname { font-size: 13px; font-weight: 500; flex: 1; }\n    .bg-ts { font-size: 11px; color: #aaa; }\n\n    /* Dark scrim */\n    .scrim { position: absolute; inset: 40px 0 0 0; background: rgba(0,0,0,0.28); z-index: 2; }\n\n    /* Modal */\n    .modal-wrap { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; z-index: 3; }\n    .modal {\n      width: 640px;\n      background: #fff;\n      border-radius: 12px;\n      box-shadow: 0 24px 80px rgba(0,0,0,0.28), 0 8px 24px rgba(0,0,0,0.14);\n      overflow: hidden;\n      display: flex; flex-direction: column;\n    }\n\n    /* Modal header */\n    .modal-header {\n      padding: 22px 26px 18px;\n      border-bottom: 1px solid #ebebeb;\n      display: flex; align-items: flex-start; justify-content: space-between;\n    }\n    .modal-title { font-size: 13px; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; color: #1a1a1a; }\n    .step-indicator { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }\n    .step-label { font-size: 11px; color: #bbb; }\n    .step-dots { display: flex; gap: 5px; }\n    .sdot-i { width: 7px; height: 7px; border-radius: 50%; background: #e5e5e5; }\n    .sdot-i.active { background: #C8922A; }\n\n    /* Modal body */\n    .modal-body { padding: 22px 26px; flex: 1; overflow-y: auto; }\n\n    /* Input */\n    .field-label { font-size: 11px; font-weight: 600; color: #888; letter-spacing: 0.06em; text-transform: uppercase; margin-bottom: 8px; }\n    .name-input-wrap {\n      display: flex; align-items: center;\n    padding: 10px 14px;\n      gap: 10px;\n      margin-bottom: 22px;\n      background: #FFFDF8;\n    }\n    .name-input {\n      flex: 1; font-size: 14px; font-weight: 500; color: #1a1a1a;\n      border: none; outline: none; background: transparent;\n      font-family: 'Inter', sans-serif;\n    }\n    .name-dropdown-arrow { color: #ccc; font-size: 11px; }\n\n    /* Template grid */\n    .templates-label {\n      font-size: 11px; font-weight: 600; color: #888;\n      letter-spacing: 0.06em; text-transform: uppercase;\n      margin-bottom: 12px;\n    }\n    .tpl-card {\n      border: 1.5px solid #e8e8e8;\n      border-radius: 8px;\n      padding: 12px 14px;\n      cursor: pointer;\n      transition: border-color 0.15s;\n    }\n    .tpl-card:hover { border-color: #C8922A; }\n    .tpl-card.selected {\n      border-color: #C8922A;\n      background: #FFFDF8;\n    }\n    .tpl-card-name {\n      font-size: 12px; font-weight: 600; color: #1a1a1a;\n      margin-bottom: 4px;\n      display: flex; align-items: center; gap: 6px;\n    }\n    .tpl-check { color: #C8922A; font-size: 11px; }\n    .tpl-flow { font-size: 10px; color: #aaa; }\n    .tpl-card-single { grid-column: 1; }\n\n    /* Blank hint */\n    .blank-hint { font-size: 11px; color: #bbb; margin-top: 4px; }\n\n    /* Modal footer */\n    .modal-footer {\n      padding: 16px 26px;\n      border-top: 1px solid #ebebeb;\n      display: flex; align-items: center; justify-content: flex-end;\n      gap: 10px;\n    }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 18px; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #555; }\n    .btn-primary { background: #C8922A; color: #fff; }\n  "
        }}
      />

      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window {\n      width: 1320px; height: 820px; background: #fff; border-radius: 12px;\n      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);\n      overflow: hidden; display: flex; flex-direction: column; position: relative;\n    }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    .bg-app { position: absolute; inset: 40px 0 0 0; opacity: 0.18;  display: flex; flex-direction: column; }\n    .bg-topbar { height: 52px; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; }\n    .bg-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .bg-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .bg-sp { flex: 1; }\n    .bg-body { flex: 1; display: flex; }\n    .bg-sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; padding: 18px 0; }\n    .bg-sidebar-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .bg-vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; }\n    .bg-dot { width: 8px; height: 8px; border-radius: 50%; }\n    .bg-main { flex: 1; padding: 18px 24px; }\n    .bg-ledger-title { font-size: 17px; font-weight: 700; margin-bottom: 20px; }\n    .bg-row { display: flex; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f2f2f2; }\n    .bg-sdot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .bg-tname { font-size: 13px; font-weight: 500; flex: 1; }\n    .bg-ts { font-size: 11px; color: #aaa; }\n\n    .scrim { position: absolute; inset: 40px 0 0 0; background: rgba(0,0,0,0.28); z-index: 2; }\n    .modal-wrap { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; z-index: 3; }\n    .modal { width: 640px; background: #fff; border-radius: 12px; box-shadow: 0 24px 80px rgba(0,0,0,0.28), 0 8px 24px rgba(0,0,0,0.14); overflow: hidden; display: flex; flex-direction: column; }\n\n    .modal-header { padding: 22px 26px 18px; border-bottom: 1px solid #ebebeb; display: flex; align-items: flex-start; justify-content: space-between; }\n    .modal-title { font-size: 13px; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; color: #1a1a1a; }\n    .modal-subtitle { font-size: 12px; color: #aaa; margin-top: 3px; }\n    .step-indicator { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }\n    .step-label { font-size: 11px; color: #bbb; }\n    .step-dots { display: flex; gap: 5px; }\n    .sdot-i { width: 7px; height: 7px; border-radius: 50%; background: #e5e5e5; }\n    .sdot-i.active { background: #C8922A; }\n    .sdot-i.done { background: #C8922A; opacity: 0.4; }\n\n    .modal-body { padding: 22px 26px; overflow-y: auto; }\n\n    .section-label { font-size: 10px; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-bottom: 12px; }\n    .section-hint { font-size: 11px; color: #bbb; margin-bottom: 14px; }\n\n    /* Role table */\n    .role-table { width: 100%; border-collapse: separate; border-spacing: 0; margin-bottom: 6px; }\n    .role-table thead th {\n      font-size: 10px; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase;\n      color: #bbb; padding: 0 10px 8px; text-align: left; border-bottom: 1px solid #ebebeb;\n    }\n    .role-table td { padding: 10px; border-bottom: 1px solid #f5f5f5; vertical-align: middle; }\n\n    .role-name { font-size: 13px; font-weight: 500; color: #1a1a1a; }\n    .role-access { font-size: 12px; color: #aaa; }\n\n    /* Vault dropdown */\n    .vault-select-wrap {\n      display: flex; align-items: center;\n      border: 1px solid #e5e5e5; border-radius: 7px;\n      padding: 6px 10px; gap: 8px;\n      background: #fafafa; cursor: pointer;\n      min-width: 180px;\n    }\n    .vault-select-wrap:hover { border-color: #C8922A; }\n    .vs-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .vs-name { font-size: 12px; font-weight: 500; color: #333; flex: 1; }\n    .vs-arrow { font-size: 10px; color: #bbb; }\n\n    /* Divider */\n    .section-divider { border: none; border-top: 1px solid #ebebeb; margin: 18px 0; }\n\n    /* C3 section */\n    .c3-optional-box {\n      background: #fafafa;\n      border: 1px solid #ebebeb;\n      border-radius: 9px;\n      padding: 18px;\n    }\n    .c3-box-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }\n    .c3-box-title { font-size: 12px; font-weight: 600; color: #1a1a1a; display: flex; align-items: center; gap: 7px; }\n    .c3-chain-icon { font-size: 14px; }\n    .c3-optional-badge {\n      font-size: 10px; font-weight: 600;\n      letter-spacing: 0.08em; text-transform: uppercase;\n      color: #bbb; background: #f0f0f0;\n      padding: 2px 7px; border-radius: 4px;\n    }\n    .c3-box-desc { font-size: 11px; color: #aaa; margin-bottom: 14px; line-height: 1.5; }\n    .c3-add-btn {\n      display: flex; align-items: center; gap: 8px;\n      padding: 9px 14px;\n      border: 1.5px dashed #C8922A;\n      border-radius: 7px;\n      background: #FFFDF8;\n      font-size: 12px; font-weight: 500; color: #C8922A;\n      cursor: pointer; margin-bottom: 14px;\n      width: fit-content;\n    }\n    .c3-add-btn:hover { background: #FBF0D8; }\n\n    /* Invite input */\n    .invite-wrap { display: flex; gap: 8px; align-items: center; }\n    .invite-input {\n      flex: 1;\n      padding: 8px 12px;\n      border: 1px solid #e5e5e5;\n      border-radius: 7px;\n      font-size: 12px; color: #555;\n      outline: none; background: #fff;\n      font-family: 'Inter', sans-serif;\n    }\n    .invite-input:focus { border-color: #C8922A; }\n    .invite-or { font-size: 11px; color: #bbb; white-space: nowrap; }\n    .invite-link-btn {\n      padding: 8px 12px;\n      border: 1px solid #e5e5e5;\n      border-radius: 7px;\n      font-size: 11px; color: #555;\n      background: #fff; cursor: pointer; white-space: nowrap;\n    }\n    .invite-link-btn:hover { border-color: #C8922A; color: #C8922A; }\n    .after-activation-note { font-size: 11px; color: #ccc; margin-top: 10px; font-style: italic; }\n\n    .modal-footer { padding: 16px 26px; border-top: 1px solid #ebebeb; display: flex; align-items: center; justify-content: space-between; gap: 10px; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 18px; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #555; }\n    .btn-primary { background: #C8922A; color: #fff; }\n  "
        }}
      />
      <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window {\n      width: 1320px; height: 820px; background: #fff; border-radius: 12px;\n      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);\n      overflow: hidden; display: flex; flex-direction: column; position: relative;\n    }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    .bg-app { position: absolute; inset: 40px 0 0 0; opacity: 0.18;  display: flex; flex-direction: column; }\n    .bg-topbar { height: 52px; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; }\n    .bg-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .bg-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .bg-sp { flex: 1; }\n    .bg-body { flex: 1; display: flex; }\n    .bg-sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; padding: 18px 0; }\n    .bg-sidebar-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .bg-vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; }\n    .bg-dot { width: 8px; height: 8px; border-radius: 50%; }\n    .bg-main { flex: 1; padding: 18px 24px; }\n    .bg-ledger-title { font-size: 17px; font-weight: 700; margin-bottom: 20px; }\n    .bg-row { display: flex; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f2f2f2; }\n    .bg-sdot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .bg-tname { font-size: 13px; font-weight: 500; flex: 1; }\n    .bg-ts { font-size: 11px; color: #aaa; }\n\n    .scrim { position: absolute; inset: 40px 0 0 0; background: rgba(0,0,0,0.28); z-index: 2; }\n    .modal-wrap { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; z-index: 3; }\n    .modal { width: 640px; background: #fff; border-radius: 12px; box-shadow: 0 24px 80px rgba(0,0,0,0.28), 0 8px 24px rgba(0,0,0,0.14); overflow: hidden; display: flex; flex-direction: column; }\n\n    .modal-header { padding: 22px 26px 18px; border-bottom: 1px solid #ebebeb; display: flex; align-items: flex-start; justify-content: space-between; }\n    .modal-title { font-size: 13px; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; color: #1a1a1a; }\n    .modal-subtitle { font-size: 12px; color: #aaa; margin-top: 3px; }\n    .step-indicator { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }\n    .step-label { font-size: 11px; color: #bbb; }\n    .step-dots { display: flex; gap: 5px; }\n    .sdot-i { width: 7px; height: 7px; border-radius: 50%; background: #e5e5e5; }\n    .sdot-i.active { background: #C8922A; }\n    .sdot-i.done { background: #C8922A; opacity: 0.4; }\n\n    .modal-body { padding: 22px 26px; overflow-y: auto; }\n\n    /* Summary block */\n    .summary-block {\n      background: #fafafa; border: 1px solid #ebebeb;\n      border-radius: 9px; padding: 16px 18px;\n      margin-bottom: 16px;\n    }\n    .summary-row { display: flex; gap: 0; padding: 5px 0; border-bottom: 1px solid #f0f0f0; }\n    .summary-row:last-child { border-bottom: none; }\n    .sum-key { font-size: 11px; color: #bbb; width: 140px; flex-shrink: 0; padding-top: 1px; }\n    .sum-val { font-size: 12px; font-weight: 500; color: #1a1a1a; }\n    .sum-val.mono { font-family: 'SF Mono', 'Fira Code', monospace; font-size: 11px; color: #555; }\n\n    /* Flow visual */\n    .flow-section { margin-bottom: 16px; }\n    .section-label { font-size: 10px; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb;  }\n\n    .flow-visual { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }\n    .fv-vault {\n      display: flex; align-items: center; gap: 6px;\n      padding: 5px 10px;\n      border: 1px solid #e5e5e5;\n      border-radius: 6px;\n      background: #fff;\n      font-size: 12px; font-weight: 500; color: #333;\n    }\n    .fv-dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; }\n    .fv-arrow { color: #ccc; font-size: 14px; }\n\n    /* Slots summary table */\n    .slots-summary { width: 100%; border-collapse: separate; border-spacing: 0; margin-bottom: 16px; }\n    .slots-summary thead th { font-size: 10px; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase; color: #bbb; padding: 0 10px 8px; text-align: left; border-bottom: 1px solid #ebebeb; }\n    .slots-summary td { padding: 8px 10px; border-bottom: 1px solid #f5f5f5; font-size: 12px; vertical-align: middle; }\n    .slots-summary td:first-child { font-weight: 500; color: #1a1a1a; font-family: 'SF Mono', monospace; font-size: 11px; }\n    .vault-chip { display: inline-flex; align-items: center; gap: 5px; }\n    .vc-dot { width: 7px; height: 7px; border-radius: 50%; }\n    .gate-on { display: inline-flex; align-items: center; gap: 4px; font-size: 10px; color: #C8922A; font-weight: 600; }\n\n    /* Status rows */\n    .status-section { margin-bottom: 16px; }\n    .status-row { display: flex; align-items: center; justify-content: space-between; padding: 10px 14px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 8px; margin-bottom: 8px; }\n    .status-row:last-child { margin-bottom: 0; }\n    .sr-left { display: flex; align-items: center; gap: 8px; }\n    .sr-icon { font-size: 13px; }\n    .sr-label { font-size: 12px; font-weight: 500; color: #333; }\n    .sr-val { font-size: 12px; color: #888; display: flex; align-items: center; gap: 6px; }\n    .sr-val.positive { color: #059669; font-weight: 600; }\n    .c3-ext-badge { background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; padding: 2px 7px; border-radius: 4px; font-size: 10px; font-weight: 600; }\n    .stellar-on { display: flex; align-items: center; gap: 6px; }\n    .stellar-dot { width: 7px; height: 7px; border-radius: 50%; background: #22C55E; }\n\n    /* Post-activation preview */\n    .post-preview {\n      margin-top: 4px;\n      padding: 12px 14px;\n      background: #F8F8F8;\n      border: 1px solid #ebebeb;\n      border-radius: 8px;\n      border-left: 3px solid #C8922A;\n    }\n    .pp-label { font-size: 10px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.08em; color: #bbb; margin-bottom: 6px; }\n    .pp-row { display: flex; align-items: center; gap: 10px; }\n    .pp-sdot { width: 8px; height: 8px; border-radius: 50%; background: #F59E0B; flex-shrink: 0; }\n    .pp-name { font-size: 12px; font-weight: 500; color: #555; flex: 1; }\n    .pp-meta { font-size: 11px; color: #bbb; display: flex; gap: 8px; }\n    .pp-pipe { display: flex; gap: 3px; }\n    .pp-seg { width: 16px; height: 3px; border-radius: 2px; background: #e8e8e8; }\n    .pp-c3 { font-size: 10px; color: #999; background: #f0f0f0; padding: 2px 6px; border-radius: 4px; }\n\n    .modal-footer { padding: 16px 26px; border-top: 1px solid #ebebeb; display: flex; align-items: center; justify-content: space-between; gap: 10px; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 18px; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #555; }\n    .btn-activate { background: #C8922A; color: #fff; font-size: 13px; font-weight: 600; padding: 10px 24px; border-radius: 8px; border: none; cursor: pointer; letter-spacing: 0.02em; }\n    .btn-activate:hover { background: #b8821a; }\n  "
        }}
      />
      <style
    dangerouslySetInnerHTML={{
      __html:
        "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; position: relative; }\n\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    /* Topbar dimmed */\n    .topbar { height: 52px; background: #fff; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; flex-shrink: 0; gap: 16px; opacity: 0.45;  }\n    .topbar-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .workspace-pill { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .topbar-spacer { flex: 1; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 6px 14px; border-radius: 7px; font-size: 12px; font-weight: 500; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #444; }\n    .btn-primary { background: #C8922A; color: #fff; }\n\n    .layout { flex: 1; display: flex; overflow: hidden; }\n\n    /* Dimmed ledger side */\n    .ledger-wrapper { flex: 1; display: flex; overflow: hidden; opacity: 0.4;  }\n    .sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; display: flex; flex-direction: column; padding: 18px 0 16px; flex-shrink: 0; }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #333; }\n    .vault-row.all-vaults { font-size: 12px; color: #888; padding: 5px 16px 10px; border-bottom: 1px solid #ebebeb; margin-bottom: 6px; }\n    .vault-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .unread-pip { margin-left: auto; width: 6px; height: 6px; background: #C8922A; border-radius: 50%; }\n    .ledger-area { flex: 1; display: flex; flex-direction: column; overflow: hidden; }\n    .ledger-topbar { display: flex; align-items: center; justify-content: space-between; padding: 18px 24px 14px; flex-shrink: 0; }\n    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; }\n    .ledger-controls { display: flex; gap: 8px; }\n    .ctrl-btn { display: flex; align-items: center; gap: 5px; padding: 5px 12px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; font-size: 12px; color: #555; }\n    .table-wrap { flex: 1; overflow-y: hidden; padding: 0 24px; }\n    table { width: 100%; border-collapse: separate; border-spacing: 0; }\n    thead th { padding: 0 10px 10px; text-align: left; font-size: 10px; font-weight: 600; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; border-bottom: 1px solid #ebebeb; }\n    td { padding: 11px 10px; vertical-align: middle; border-bottom: 1px solid #f2f2f2; font-size: 13px; }\n    .row-selected td { background: #FBF0D8 !important; }\n    .sdot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }\n    .s-ok { background: #22C55E; } .s-pend { background: #F59E0B; }\n    .th-line1 { font-size: 13px; font-weight: 500; }\n    .th-type { font-size: 9px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-right: 6px; }\n    .th-line2 { font-size: 11px; color: #aaa; margin-top: 2px; }\n    .flow { display: flex; align-items: center; gap: 4px; }\n    .vb { width: 22px; height: 22px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 8px; font-weight: 700; color: #fff; flex-shrink: 0; }\n    .fa { color: #ccc; font-size: 10px; }\n    .pipeline { display: flex; gap: 3px; }\n    .pseg { height: 4px; width: 20px; border-radius: 2px; }\n    .pseg-done { background: #C8922A; } .pseg-wait { background: #e8e8e8; }\n    .ts { font-size: 12px; color: #bbb; }\n    .stellar-val { font-family: 'SF Mono', monospace; font-size: 10px; color: #ccc; }\n    .c3b { display: inline-flex; align-items: center; padding: 3px 8px; border-radius: 5px; font-size: 11px; font-weight: 600; }\n    .c3-ext { background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; }\n    .c3-linked { background: #FBF0D8; border: 1px solid #E8C87A; color: #C8922A; }\n    .ledger-footer { padding: 12px 24px; border-top: 1px solid #ebebeb; font-size: 12px; color: #bbb; flex-shrink: 0; }\n\n    /* ---- DETAIL PANEL ---- */\n    .detail-panel {  height: 100%; background: #fff; border-left: 1px solid #e0e0e0; display: flex; flex-direction: column; flex-shrink: 0; box-shadow: -6px 0 24px rgba(0,0,0,0.08); overflow: hidden; z-index: 10; }\n\n    .dp-header { padding: 20px 22px 16px; border-bottom: 1px solid #ebebeb; flex-shrink: 0; }\n    .dp-close-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }\n    .dp-badge { font-size: 10px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #C8922A; background: #FBF0D8; padding: 3px 8px; border-radius: 4px; }\n    .dp-close { width: 24px; height: 24px; border-radius: 6px; background: #f5f5f5; border: 1px solid #e5e5e5; display: flex; align-items: center; justify-content: center; font-size: 12px; color: #888; cursor: pointer; }\n    .dp-title { font-size: 15px; font-weight: 700; color: #1a1a1a; margin-bottom: 6px; line-height: 1.3; }\n    .dp-meta { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 8px; }\n    .dp-meta-chip { font-size: 11px; color: #999; background: #f8f8f8; border: 1px solid #ebebeb; padding: 2px 8px; border-radius: 4px; font-family: 'SF Mono', monospace; }\n    .dp-custom { margin-top: 8px; font-size: 11px; color: #999; }\n\n    .dp-body { flex: 1; overflow-y: auto; padding: 0 22px 20px; }\n    .dp-section { padding: 16px 0 12px; border-bottom: 1px solid #f0f0f0; }\n    .dp-section:last-child { border-bottom: none; }\n    .dp-section-title { font-size: 10px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-bottom: 12px; }\n\n    /* Pipeline steps */\n    .pipeline-step { display: flex; align-items: flex-start; gap: 10px; margin-bottom: 10px; }\n    .ps-icon { width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 9px; font-weight: 700; color: #fff; flex-shrink: 0; margin-top: 1px; }\n    .ps-icon-wait { background: #f0f0f0 !important; border: 2px dashed #ddd; color: #ccc !important; }\n    .ps-content { flex: 1; }\n    .ps-label { font-size: 12px; font-weight: 500; color: #333; }\n    .ps-label.wait { color: #ccc; font-style: italic; }\n    .ps-sublabel { font-size: 11px; color: #bbb; margin-top: 2px; }\n    .ps-ts { font-size: 10px; color: #ccc; margin-top: 2px; }\n    .ps-check { font-size: 13px; flex-shrink: 0; }\n    .ps-check.done { color: #22C55E; }\n    .ps-check.wait { color: #e0e0e0; }\n    .ps-connector { width: 1px; height: 10px; background: #e8e8e8; margin-left: 12px; }\n\n    /* C3 divider in pipeline */\n    .c3-pipe-divider { display: flex; align-items: center; gap: 10px; margin: 8px 0 10px; }\n    .c3-pipe-line { flex: 1; height: 1px; background: #E8C87A; }\n    .c3-pipe-label { font-size: 10px; font-weight: 600; color: #C8922A; letter-spacing: 0.08em; text-transform: uppercase; white-space: nowrap; }\n\n    /* External vault step - gold */\n    .ps-icon-external { background: #FBF0D8 !important; border: 2px solid #C8922A !important; color: #C8922A !important; font-size: 8px !important; }\n    .ps-label-external { color: #C8922A !important; font-weight: 600 !important; }\n    .ps-sublabel-external { color: #C8922A !important; opacity: 0.7; }\n    .ext-role-badge { display: inline-block; font-size: 9px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: #C8922A; background: #FBF0D8; border: 1px solid #E8C87A; padding: 1px 5px; border-radius: 3px; margin-left: 6px; }\n\n    /* Commits */\n    .commit-row { display: flex; align-items: flex-start; gap: 10px; padding: 8px 0; border-bottom: 1px solid #f5f5f5; }\n    .commit-row:last-child { border-bottom: none; }\n    .commit-row.external { background: #FFFDF8; margin: 0 -4px; padding: 8px 4px; border-radius: 5px; }\n    .commit-ts { font-size: 10px; color: #bbb; white-space: nowrap; padding-top: 1px; min-width: 70px; }\n    .commit-body { flex: 1; }\n    .commit-actor { font-size: 11px; font-weight: 600; color: #555; }\n    .commit-actor.ext { color: #C8922A; }\n    .commit-action { font-size: 11px; color: #333; margin-left: 4px; }\n    .commit-cid { font-size: 10px; color: #bbb; font-family: 'SF Mono', monospace; margin-top: 2px; }\n    .commit-verify { font-size: 10px; color: #C8922A; cursor: pointer; white-space: nowrap; padding-top: 1px; }\n\n    /* C3 extended section */\n    .c3-extended-box { padding: 14px 16px; background: #FBF0D8; border: 1px solid #E8C87A; border-radius: 8px; margin-top: 4px; }\n    .c3-ext-header { display: flex; align-items: center; gap: 7px; margin-bottom: 8px; }\n    .c3-ext-icon { font-size: 14px; }\n    .c3-ext-title { font-size: 12px; font-weight: 600; color: #C8922A; }\n    .c3-ext-active-badge { font-size: 9px; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; background: #C8922A; color: #fff; padding: 2px 6px; border-radius: 3px; margin-left: auto; }\n    .c3-ext-detail { font-size: 12px; color: #7A5A1A; margin-bottom: 10px; line-height: 1.5; }\n    .c3-ext-detail span { font-weight: 600; }\n    .c3-view-btn { display: inline-flex; align-items: center; gap: 6px; padding: 7px 14px; border: 1px solid #C8922A; border-radius: 6px; font-size: 12px; font-weight: 500; color: #C8922A; background: #fff; cursor: pointer; }\n    .c3-view-btn:hover { background: #FBF0D8; }\n\n    /* Stellar ref */\n    .stellar-ref-row { display: flex; align-items: center; gap: 8px; padding: 10px 12px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 7px; }\n    .stellar-hash-full { font-family: 'SF Mono', monospace; font-size: 10px; color: #888; flex: 1; overflow: hidden; text-overflow: ellipsis; }\n    .copy-btn { font-size: 10px; color: #C8922A; cursor: pointer; white-space: nowrap; }\n\n    /* Actions */\n    .dp-actions { padding: 14px 22px; border-top: 1px solid #ebebeb; display: flex; gap: 8px; flex-shrink: 0; }\n    .action-btn { flex: 1; padding: 8px 10px; border: 1px solid #e5e5e5; border-radius: 7px; background: #fff; font-size: 11px; font-weight: 500; color: #444; cursor: pointer; text-align: center; }\n    .action-btn:hover { border-color: #C8922A; color: #C8922A; }\n    .action-btn.danger:hover { border-color: #EF4444; color: #EF4444; }\n  "
    }}
  />
   <style
    dangerouslySetInnerHTML={{
      __html:
        "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; position: relative; }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n    .topbar { height: 52px; background: #fff; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; flex-shrink: 0; opacity: 0.45;  }\n    .tb-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .tb-sp { flex: 1; }\n    .btn-g { background: #f5f5f5; border: 1px solid #e5e5e5; color: #444; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n    .btn-p { background: #C8922A; color: #fff; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n    .layout { flex: 1; display: flex; overflow: hidden; }\n\n    /* Dimmed ledger */\n    .ledger-wrapper { flex: 1; display: flex; overflow: hidden; opacity: 0.38;  }\n    .sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; display: flex; flex-direction: column; padding: 18px 0; flex-shrink: 0; }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #333; }\n    .vault-row.all { font-size: 12px; color: #888; padding: 5px 16px 10px; border-bottom: 1px solid #ebebeb; margin-bottom: 6px; }\n    .vdot { width: 8px; height: 8px; border-radius: 50%; }\n    .pip { margin-left: auto; width: 6px; height: 6px; background: #EF4444; border-radius: 50%; }\n    .ledger-area { flex: 1; display: flex; flex-direction: column; overflow: hidden; }\n    .ledger-topbar { display: flex; align-items: center; justify-content: space-between; padding: 18px 24px 14px; }\n    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; }\n    .table-wrap { flex: 1; overflow: hidden; padding: 0 24px; }\n    table { width: 100%; border-collapse: separate; border-spacing: 0; }\n    thead th { padding: 0 10px 10px; text-align: left; font-size: 10px; font-weight: 600; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; border-bottom: 1px solid #ebebeb; }\n    td { padding: 11px 10px; vertical-align: middle; border-bottom: 1px solid #f2f2f2; font-size: 13px; }\n    .row-sel td { background: #FEF2F2 !important; }\n    .sdot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }\n    .s-ok { background: #22C55E; } .s-pend { background: #F59E0B; } .s-dispute { background: #EF4444; }\n    .th-type { font-size: 9px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-right: 6px; }\n    .flow { display: flex; align-items: center; gap: 4px; }\n    .vb { width: 22px; height: 22px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 8px; font-weight: 700; color: #fff; }\n    .fa { color: #ccc; font-size: 10px; }\n    .pipeline { display: flex; gap: 3px; }\n    .pseg { height: 4px; width: 20px; border-radius: 2px; }\n    .pseg-done { background: #C8922A; } .pseg-wait { background: #e8e8e8; } .pseg-reject { background: #EF4444; }\n    .ts { font-size: 12px; color: #bbb; }\n    .stellar-val { font-family: 'SF Mono', monospace; font-size: 10px; color: #ccc; }\n    .c3b { display: inline-flex; align-items: center; padding: 3px 8px; border-radius: 5px; font-size: 11px; font-weight: 600; background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; }\n    .dispute-badge-row { display: inline-flex; align-items: center; gap: 4px; background: #FEF2F2; border: 1px solid #FECACA; color: #DC2626; padding: 2px 7px; border-radius: 4px; font-size: 10px; font-weight: 700; letter-spacing: 0.05em; }\n    .ledger-footer { padding: 12px 24px; border-top: 1px solid #ebebeb; font-size: 12px; color: #bbb; flex-shrink: 0; }\n\n    /* ---- DETAIL PANEL ---- */\n    .detail-panel { height: 100%; background: #fff; border-left: 1px solid #e0e0e0; display: flex; flex-direction: column; flex-shrink: 0; box-shadow: -6px 0 24px rgba(0,0,0,0.08); overflow: hidden; z-index: 10; }\n\n    /* Dispute alert banner */\n    .dispute-banner { background: #FEF2F2; border-bottom: 1px solid #FECACA; padding: 10px 20px; display: flex; align-items: center; gap: 10px; flex-shrink: 0; }\n    .db-icon { font-size: 14px; }\n    .db-text { font-size: 12px; color: #991B1B; flex: 1; font-weight: 500; }\n    .db-text span { font-weight: 700; }\n    .db-ts { font-size: 11px; color: #EF4444; white-space: nowrap; }\n\n    .dp-header { padding: 16px 22px 14px; border-bottom: 1px solid #ebebeb; flex-shrink: 0; }\n    .dp-close-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }\n    .dp-badge { font-size: 10px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #DC2626; background: #FEF2F2; padding: 3px 8px; border-radius: 4px; border: 1px solid #FECACA; }\n    .dp-close { width: 24px; height: 24px; border-radius: 6px; background: #f5f5f5; border: 1px solid #e5e5e5; display: flex; align-items: center; justify-content: center; font-size: 12px; color: #888; cursor: pointer; }\n    .dp-title { font-size: 15px; font-weight: 700; color: #1a1a1a; margin-bottom: 6px; }\n    .dp-meta { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 6px; }\n    .dp-meta-chip { font-size: 11px; color: #999; background: #f8f8f8; border: 1px solid #ebebeb; padding: 2px 8px; border-radius: 4px; font-family: 'SF Mono', monospace; }\n    .dp-custom { margin-top: 6px; font-size: 11px; color: #999; }\n\n    .dp-body { flex: 1; overflow-y: auto; padding: 0 22px 20px; }\n    .dp-section { padding: 14px 0 12px; border-bottom: 1px solid #f0f0f0; }\n    .dp-section:last-child { border-bottom: none; }\n    .dp-section-title { font-size: 10px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-bottom: 12px; }\n\n    /* Pipeline steps */\n    .pipeline-step { display: flex; align-items: flex-start; gap: 10px; margin-bottom: 10px; }\n    .ps-icon { width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 9px; font-weight: 700; color: #fff; flex-shrink: 0; margin-top: 1px; }\n    .ps-icon-wait { background: #f0f0f0 !important; border: 2px dashed #ddd; color: #ccc !important; }\n    .ps-icon-rejected { background: #FEF2F2 !important; border: 2px solid #EF4444 !important; color: #EF4444 !important; }\n    .ps-content { flex: 1; }\n    .ps-label { font-size: 12px; font-weight: 500; color: #333; }\n    .ps-label.rejected { color: #DC2626; font-weight: 600; }\n    .ps-label.wait { color: #ccc; font-style: italic; }\n    .ps-sublabel { font-size: 11px; color: #bbb; margin-top: 2px; }\n    .ps-sublabel.rejected { color: #DC2626; opacity: 0.7; }\n    .ps-ts { font-size: 10px; color: #ccc; margin-top: 2px; }\n    .ps-check { font-size: 13px; flex-shrink: 0; }\n    .ps-check.done { color: #22C55E; }\n    .ps-check.reject { color: #EF4444; }\n    .ps-check.wait { color: #e0e0e0; }\n    .ps-connector { width: 1px; height: 8px; background: #e8e8e8; margin-left: 12px; }\n    .ps-connector.red { background: #FECACA; }\n    .reject-reason-inline { margin-top: 6px; padding: 8px 10px; background: #FEF2F2; border: 1px solid #FECACA; border-radius: 6px; font-size: 11px; color: #991B1B; line-height: 1.5; }\n    .rri-label { font-size: 9px; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; color: #EF4444; margin-bottom: 4px; }\n\n    /* Commits */\n    .commit-row { display: flex; align-items: flex-start; gap: 10px; padding: 8px 0; border-bottom: 1px solid #f5f5f5; }\n    .commit-row:last-child { border-bottom: none; }\n    .commit-row.rejected-row { background: #FEF2F2; margin: 0 -4px; padding: 8px 4px; border-radius: 5px; border-bottom: none; }\n    .commit-ts { font-size: 10px; color: #bbb; white-space: nowrap; padding-top: 1px; min-width: 70px; }\n    .commit-ts.red { color: #DC2626; opacity: 0.8; }\n    .commit-body { flex: 1; }\n    .commit-actor { font-size: 11px; font-weight: 600; color: #555; }\n    .commit-actor.red { color: #DC2626; }\n    .commit-action { font-size: 11px; color: #333; margin-left: 4px; }\n    .commit-action.red { color: #DC2626; }\n    .commit-cid { font-size: 10px; color: #bbb; font-family: 'SF Mono', monospace; margin-top: 2px; }\n    .commit-reason { font-size: 11px; color: #991B1B; margin-top: 5px; line-height: 1.4; font-style: italic; }\n    .commit-reason::before { content: '\"'; } .commit-reason::after { content: '\"'; }\n    .commit-verify { font-size: 10px; color: #C8922A; cursor: pointer; white-space: nowrap; padding-top: 1px; }\n    .commit-verify.red { color: #EF4444; }\n    .reject-badge { display: inline-flex; align-items: center; gap: 3px; font-size: 9px; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase; background: #FEF2F2; border: 1px solid #FECACA; color: #DC2626; padding: 1px 5px; border-radius: 3px; margin-left: 5px; }\n\n    /* Stellar ref */\n    .stellar-ref-row { display: flex; align-items: center; gap: 8px; padding: 10px 12px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 7px; }\n    .stellar-hash-full { font-family: 'SF Mono', monospace; font-size: 10px; color: #888; flex: 1; overflow: hidden; text-overflow: ellipsis; }\n    .copy-btn { font-size: 10px; color: #C8922A; cursor: pointer; white-space: nowrap; }\n\n    /* Resolution callout */\n    .resolution-box { padding: 14px 16px; background: #FFFBEB; border: 1px solid #FCD34D; border-radius: 8px; margin-top: 4px; }\n    .rb-header { display: flex; align-items: center; gap: 7px; margin-bottom: 8px; }\n    .rb-icon { font-size: 13px; }\n    .rb-title { font-size: 12px; font-weight: 600; color: #92400E; }\n    .rb-steps { display: flex; flex-direction: column; gap: 6px; }\n    .rb-step { display: flex; align-items: flex-start; gap: 7px; font-size: 11px; color: #78350F; line-height: 1.4; }\n    .rb-step-num { font-size: 10px; font-weight: 700; background: #FCD34D; color: #78350F; width: 16px; height: 16px; border-radius: 50%; display: flex; align-items: center; justify-content: center; flex-shrink: 0; margin-top: 0px; }\n\n    /* Actions */\n    .dp-actions { padding: 14px 22px; border-top: 1px solid #ebebeb; display: flex; gap: 8px; flex-shrink: 0; }\n    .action-btn { flex: 1; padding: 8px 10px; border: 1px solid #e5e5e5; border-radius: 7px; background: #fff; font-size: 11px; font-weight: 500; color: #444; cursor: pointer; text-align: center; }\n    .action-btn:hover { border-color: #C8922A; color: #C8922A; }\n    .action-btn.primary { background: #C8922A; color: #fff; border-color: #C8922A; }\n    .action-btn.primary:hover { background: #b8821a; }\n  "
    }}
  />

  <style
    dangerouslySetInnerHTML={{
      __html:
        "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window {\n      width: 1320px; height: 820px; background: #fff; border-radius: 12px;\n      box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12);\n      overflow: hidden; display: flex; flex-direction: column; position: relative;\n    }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    /* Dimmed bg */\n    .bg-app { position: absolute; inset: 40px 0 0 0; opacity: 0.18;  display: flex; flex-direction: column; }\n    .bg-topbar { height: 52px; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; }\n    .bg-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .bg-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .bg-sp { flex: 1; }\n    .bg-body { flex: 1; display: flex; }\n    .bg-sidebar { ; background: #fafafa; border-right: 1px solid #ebebeb; padding: 18px 0; }\n    .bg-sidebar-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .bg-vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; }\n    .bg-dot { width: 8px; height: 8px; border-radius: 50%; }\n    .bg-main { flex: 1; padding: 18px 24px; }\n    .bg-ledger-title { font-size: 17px; font-weight: 700; margin-bottom: 20px; }\n    .bg-row { display: flex; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f2f2f2; }\n    .bg-sdot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .bg-tname { font-size: 13px; font-weight: 500; flex: 1; }\n    .bg-ts { font-size: 11px; color: #aaa; }\n\n    .scrim { position: absolute; inset: 40px 0 0 0; background: rgba(0,0,0,0.28); z-index: 2; }\n    .modal-wrap { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; z-index: 3; }\n    .modal { width: 660px; background: #fff; border-radius: 12px; box-shadow: 0 24px 80px rgba(0,0,0,0.28), 0 8px 24px rgba(0,0,0,0.14); overflow: hidden; display: flex; flex-direction: column; }\n\n    .modal-header { padding: 22px 26px 18px; border-bottom: 1px solid #ebebeb; display: flex; align-items: flex-start; justify-content: space-between; }\n    .modal-header-left {}\n    .modal-title { font-size: 13px; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; color: #1a1a1a; }\n    .modal-subtitle { font-size: 12px; color: #aaa; margin-top: 3px; }\n    .step-indicator { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }\n    .step-label { font-size: 11px; color: #bbb; }\n    .step-dots { display: flex; gap: 5px; }\n    .sdot-i { width: 7px; height: 7px; border-radius: 50%; background: #e5e5e5; }\n    .sdot-i.active { background: #C8922A; }\n    .sdot-i.done { background: #C8922A; opacity: 0.4; }\n\n    .modal-body { padding: 22px 26px; overflow-y: auto; }\n\n    /* Property slots table */\n    .section-label { font-size: 10px; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; }\n\n    .slots-table { width: 100%; border-collapse: separate; border-spacing: 0; margin-bottom: 4px; }\n    .slots-table thead th {\n      font-size: 10px; font-weight: 600; letter-spacing: 0.08em; text-transform: uppercase;\n      color: #bbb; padding: 0 10px 8px; text-align: left; border-bottom: 1px solid #ebebeb;\n    }\n    .slots-table thead th:nth-child(3), .slots-table thead th:nth-child(4) { text-align: center; }\n    .slots-table tbody tr:hover { background: #fafafa; }\n    .slots-table td { padding: 10px; border-bottom: 1px solid #f5f5f5; vertical-align: middle; }\n\n    .slot-name-input {\n      font-size: 13px; font-weight: 500; color: #1a1a1a;\n      border: none; outline: none;\n      background: transparent; width: 100%;\n      font-family: 'Inter', sans-serif;\n      padding: 4px 6px; border-radius: 4px;\n    }\n    .slot-name-input:hover { background: #f5f5f5; }\n    .slot-name-input:focus { background: #FBF0D8; }\n\n    .slot-hint { font-size: 11px; color: #bbb; }\n\n    /* Gate toggle */\n    .gate-toggle-wrap { display: flex; justify-content: center; }\n    .gate-toggle {\n      width: 34px; height: 18px;\n      border-radius: 9px;\n      background: #C8922A;\n      position: relative;\n      cursor: pointer;\n      flex-shrink: 0;\n    }\n    .gate-toggle.off { background: #e0e0e0; }\n    .gate-toggle.off::after { right: auto; left: 2px; }\n\n    /* Custom prop cell */\n    .custom-kv { display: flex; gap: 4px; }\n    .kv-input {\n      font-size: 11px; color: #999;\n      border: 1px solid #ebebeb; border-radius: 4px;\n      padding: 3px 7px; outline: none;\n      background: #fafafa; width: 70px;\n      font-family: 'Inter', sans-serif;\n    }\n    .kv-sep { font-size: 10px; color: #ccc; align-self: center; }\n\n    /* Gate note */\n    .gate-note { font-size: 11px; color: #bbb; margin: 6px 0 20px; line-height: 1.5; }\n\n    /* Channel custom property */\n    .custom-prop-row { display: flex; gap: 10px; }\n    .cp-field { flex: 1; }\n    .cp-label { font-size: 10px; font-weight: 600; color: #bbb; letter-spacing: 0.06em; text-transform: uppercase; margin-bottom: 5px; }\n    .cp-input {\n      width: 100%; padding: 9px 12px;\n      border: 1px solid #e5e5e5; border-radius: 7px;\n      font-size: 13px; color: #1a1a1a;\n      outline: none; font-family: 'Inter', sans-serif;\n    }\n    .cp-input:focus { border-color: #C8922A; }\n    .cp-hint { font-size: 11px; color: #bbb; margin-top: 6px; }\n\n    .modal-footer { padding: 16px 26px; border-top: 1px solid #ebebeb; display: flex; align-items: center; justify-content: space-between; gap: 10px; }\n    .btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 18px; border-radius: 7px; font-size: 13px; font-weight: 500; cursor: pointer; border: none; }\n    .btn-ghost { background: #f5f5f5; border: 1px solid #e5e5e5; color: #555; }\n    .btn-primary { background: #C8922A; color: #fff; }\n  "
    }}
  />
    </>
  )
}

export const c3Theme = {
  gold: '#C8922A',
  goldSoft: '#FBF0D8',
  goldBorder: '#E8C87A',
  text: '#1A1A1A',
  muted: '#888888',
  soft: '#F5F5F5',
  border: '#E5E5E5',
  borderSoft: '#EBEBEB',
  bg: '#FFFFFF',
};

export const C3BaseStyles = css`
  :root {
    --c3-gold: ${c3Theme.gold};
    --c3-gold-soft: ${c3Theme.goldSoft};
    --c3-gold-border: ${c3Theme.goldBorder};
    --c3-text: ${c3Theme.text};
    --c3-muted: ${c3Theme.muted};
    --c3-soft: ${c3Theme.soft};
    --c3-border: ${c3Theme.border};
    --c3-border-soft: ${c3Theme.borderSoft};
    --c3-bg: ${c3Theme.bg};
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
    position: relative;
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

  .tl-red {
    background: #ff5f57;
  }

  .tl-yellow {
    background: #ffbd2e;
  }

  .tl-green {
    background: #28ca41;
  }

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
  .view-btn:hover {
    transform: translateY(-1px);
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

  .layout {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

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

  .ledger-area,
  .inbox-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-width: 0;
  }

  .ledger-topbar,
  .inbox-topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 24px 14px;
    flex-shrink: 0;
    background: rgba(255, 255, 255, 0.68);
  }

  .ledger-title,
  .inbox-title {
    font-size: 17px;
    font-weight: 700;
    letter-spacing: 0.06em;
    color: var(--c3-text);
  }

  .ledger-controls,
  .inbox-controls {
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

  .table-wrap::-webkit-scrollbar,
  .ledger-area::-webkit-scrollbar,
  .sidebar::-webkit-scrollbar,
  .empty-area::-webkit-scrollbar,
  .dp-body::-webkit-scrollbar,
  .inbox-list::-webkit-scrollbar {
    width: 10px;
    height: 10px;
  }

  .table-wrap::-webkit-scrollbar-thumb,
  .ledger-area::-webkit-scrollbar-thumb,
  .sidebar::-webkit-scrollbar-thumb,
  .empty-area::-webkit-scrollbar-thumb,
  .dp-body::-webkit-scrollbar-thumb,
  .inbox-list::-webkit-scrollbar-thumb {
    background: rgba(200, 146, 42, 0.24);
    border: 2px solid transparent;
    background-clip: content-box;
    border-radius: 999px;
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

  .sdot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    display: inline-block;
  }

  .s-ok {
    background: #22c55e;
  }

  .s-pend {
    background: #f59e0b;
  }

  .s-flag {
    background: #ef4444;
  }

  .s-dispute {
    background: #dc2626;
  }

  .th-line1 {
    font-size: 13px;
    font-weight: 500;
    color: var(--c3-text);
    white-space: nowrap;
  }

  .th-line1.new-thread {
    color: var(--c3-gold);
  }

  .th-type {
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #bbb;
    margin-right: 6px;
  }

  .th-line2 {
    font-size: 11px;
    color: #aaa;
    margin-top: 2px;
  }

  .ts {
    font-size: 12px;
    color: #bbb;
    white-space: nowrap;
  }

  .stellar-val {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 10px;
    color: #bbb;
    white-space: nowrap;
  }

  .flow {
    display: flex;
    align-items: center;
    gap: 4px;
  }

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

  .vb-sm {
    width: 18px;
    height: 18px;
    font-size: 7px;
  }

  .fa {
    color: #ccc;
    font-size: 10px;
    font-weight: 400;
  }

  .pipeline {
    display: flex;
    gap: 3px;
    align-items: center;
  }

  .pseg {
    height: 4px;
    width: 20px;
    border-radius: 999px;
  }

  .pseg-done {
    background: var(--c3-gold);
  }

  .pseg-wait {
    background: #e8e8e8;
  }

  .pseg-reject {
    background: #ef4444;
  }

  .row-sel td,
  .row-new td {
    background: rgba(251, 240, 216, 0.5);
  }

  .row-new-left-border td:first-child {
    border-left: 3px solid var(--c3-gold);
  }

  .new-label {
    display: inline-block;
    margin-left: 4px;
    padding: 1px 6px;
    border-radius: 3px;
    background: var(--c3-gold);
    color: #fff;
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    vertical-align: middle;
  }

  .row-hover-actions {
    display: none;
    gap: 5px;
  }

  .rha-btn,
  .view-btn,
  .action-btn {
    padding: 3px 9px;
    background: #fff;
    border-color: var(--c3-border);
    font-size: 11px;
    color: #555;
    white-space: nowrap;
  }

  .rha-btn:hover,
  .view-btn:hover,
  .action-btn:hover {
    border-color: var(--c3-gold);
    color: var(--c3-gold);
  }

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

  .ledger-footer-hint {
    text-align: center;
  }

  .empty-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    padding: 32px 24px;
    overflow: auto;
  }

  .empty-headline {
    font-size: 15px;
    font-weight: 600;
    color: var(--c3-text);
    text-align: center;
  }

  .empty-subtext {
    max-width: 500px;
    margin-bottom: 36px;
    color: #aaa;
    font-size: 13px;
    line-height: 1.6;
    text-align: center;
  }

  .template-grid {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 16px;
    width: 100%;
    max-width: 880px;
  }

  .template-card {
    width: 240px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 22px 20px 18px;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.82);
    border: 1px solid var(--c3-border-soft);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  }

  .template-card:hover {
    border-color: var(--c3-gold);
    box-shadow: 0 14px 34px rgba(200, 146, 42, 0.12);
  }

  .tc-icon {
    width: 38px;
    height: 38px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    margin-bottom: 4px;
  }

  .tc-type-label {
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #bbb;
  }

  .tc-name {
    font-size: 14px;
    font-weight: 600;
    color: var(--c3-text);
    line-height: 1.3;
  }

  .tc-desc {
    font-size: 11px;
    color: #aaa;
    line-height: 1.5;
    flex: 1;
  }

  .tc-flow {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 6px 0 2px;
    border-top: 1px solid #f0f0f0;
  }

  .tc-vault {
    font-size: 10px;
    color: #888;
    background: var(--c3-soft);
    padding: 2px 6px;
    border-radius: 4px;
  }

  .tc-arrow {
    font-size: 10px;
    color: #ddd;
  }

  .tc-start {
    width: 100%;
    padding: 8px;
    margin-top: 4px;
    border-color: var(--c3-border);
    background: rgba(250, 250, 250, 0.95);
    color: #888;
    font-size: 12px;
    font-weight: 500;
  }

  .template-card:hover .tc-start {
    background: var(--c3-gold);
    border-color: var(--c3-gold);
    color: #fff;
  }

  .toast-bar {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 20px;
    flex-shrink: 0;
    background: rgba(26, 26, 26, 0.96);
    color: #fff;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  }

  .toast-icon {
    font-size: 14px;
  }

  .toast-text {
    flex: 1;
    font-size: 12px;
  }

  .toast-text strong {
    color: var(--c3-gold);
  }

  .toast-action {
    color: var(--c3-gold);
    font-weight: 600;
    white-space: nowrap;
  }

  .toast-close {
    color: #888;
    cursor: pointer;
    font-size: 14px;
  }

  .modal-wrap {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3;
  }

  .modal {
    width: 560px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border-radius: 18px;
    background: rgba(255, 255, 255, 0.92);
    box-shadow: 0 24px 80px rgba(0, 0, 0, 0.28), 0 8px 24px rgba(0, 0, 0, 0.16);
  }

  .dp-header {
    padding: 20px 22px 16px;
    flex-shrink: 0;
  }

  .dp-body {
    flex: 1;
    overflow: auto;
    padding: 0 22px 20px;
  }

  .dp-actions {
    padding: 14px 22px;
    border-top: 1px solid var(--c3-border-soft);
    display: flex;
    gap: 8px;
    flex-shrink: 0;
  }

  .c3-section {
    padding: 14px 16px;
    background: rgba(250, 250, 250, 0.88);
    border: 1px solid var(--c3-border-soft);
    border-radius: 12px;
    margin-top: 4px;
  }

  .screen-label {
    position: absolute;
    top: 8px;
    right: 16px;
    font-size: 10px;
    color: #ccc;
    letter-spacing: 0.06em;
    font-family: 'SF Mono', monospace;
  }
`;

export const PanelStyles = () => {
  return (
    <>
          <style
        dangerouslySetInnerHTML={{
          __html:
            "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; position: relative; }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n\n    /* Topbar dimmed */\n    .topbar { height: 52px; background: #fff; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; flex-shrink: 0; opacity: 0.45;  }\n    .tb-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .tb-wp { background: #f5f5f5; border: 1px solid #e8e8e8; border-radius: 6px; padding: 4px 10px; font-size: 12px; color: #555; }\n    .tb-spacer { flex: 1; }\n    .btn-g { background: #f5f5f5; border: 1px solid #e5e5e5; color: #444; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n    .btn-p { background: #C8922A; color: #fff; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n\n    .layout { flex: 1; display: flex; overflow: hidden; }\n\n    /* Dimmed ledger wrapper */\n    .ledger-wrapper { flex: 1; display: flex; overflow: hidden; opacity: 0.35;  }\n    .sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; display: flex; flex-direction: column; padding: 18px 0 16px; flex-shrink: 0; }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #333; }\n    .vault-row.all-vaults { font-size: 12px; color: #888; padding: 5px 16px 10px; border-bottom: 1px solid #ebebeb; margin-bottom: 6px; }\n    .vault-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }\n    .unread-pip { margin-left: auto; width: 6px; height: 6px; background: #C8922A; border-radius: 50%; }\n    .ledger-area { flex: 1; display: flex; flex-direction: column; overflow: hidden; }\n    .ledger-topbar { display: flex; align-items: center; justify-content: space-between; padding: 18px 24px 14px; flex-shrink: 0; }\n    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; }\n    .table-wrap { flex: 1; overflow: hidden; padding: 0 24px; }\n    table { width: 100%; border-collapse: separate; border-spacing: 0; }\n    thead th { padding: 0 10px 10px; text-align: left; font-size: 10px; font-weight: 600; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; border-bottom: 1px solid #ebebeb; }\n    td { padding: 11px 10px; vertical-align: middle; border-bottom: 1px solid #f2f2f2; font-size: 13px; }\n    .sdot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }\n    .s-ok { background: #22C55E; } .s-pend { background: #F59E0B; }\n    .th-type { font-size: 9px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-right: 6px; }\n    .flow { display: flex; align-items: center; gap: 4px; }\n    .vb { width: 22px; height: 22px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 8px; font-weight: 700; color: #fff; }\n    .fa { color: #ccc; font-size: 10px; }\n    .pipeline { display: flex; gap: 3px; }\n    .pseg { height: 4px; width: 20px; border-radius: 2px; }\n    .pseg-done { background: #C8922A; } .pseg-wait { background: #e8e8e8; }\n    .ts { font-size: 12px; color: #bbb; }\n    .stellar-val { font-family: 'SF Mono', monospace; font-size: 10px; color: #ccc; }\n    .c3b { display: inline-flex; align-items: center; padding: 3px 8px; border-radius: 5px; font-size: 11px; font-weight: 600; }\n    .c3-ext { background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; }\n    .c3-linked { background: #FBF0D8; border: 1px solid #E8C87A; color: #C8922A; }\n    .ledger-footer { padding: 12px 24px; border-top: 1px solid #ebebeb; font-size: 12px; color: #bbb; flex-shrink: 0; }\n\n    /* ===== SLIDE PANEL ===== */\n    .slide-panel { margin-top: -1rem ; height: 100%; background: #fff; border-left: 1px solid #e0e0e0; display: flex; flex-direction: column; flex-shrink: 0; box-shadow: -8px 0 32px rgba(0,0,0,0.10); z-index: 10; overflow: hidden; }\n\n    /* Panel header */\n    .sp-header { padding: 20px 24px 16px; border-bottom: 1px solid #ebebeb; flex-shrink: 0; }\n    .sp-header-row { display: flex; align-items: center; justify-content: space-between; }\n    .sp-title { font-size: 15px; font-weight: 700; color: #1a1a1a; letter-spacing: 0.01em; }\n    .sp-subtitle { font-size: 12px; color: #bbb; margin-top: 3px; }\n    .sp-close { width: 26px; height: 26px; border-radius: 7px; background: #f5f5f5; border: 1px solid #e5e5e5; display: flex; align-items: center; justify-content: center; font-size: 13px; color: #888; cursor: pointer; flex-shrink: 0; }\n\n    /* Panel body */\n    .sp-body { flex: 1; overflow-y: auto; padding: 20px 24px; display: flex; flex-direction: column; gap: 20px; }\n\n    /* Field label style */\n    .fl { font-size: 10px; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; margin-bottom: 6px; }\n    .fl-hint { font-size: 11px; color: #bbb; font-weight: 400; margin-left: 6px; text-transform: none; letter-spacing: 0; }\n\n    /* Channel selector */\n    .channel-selected { display: flex; align-items: center; gap: 10px; padding: 12px 14px; border: 1.5px solid #C8922A; border-radius: 9px; background: #FFFDF8; cursor: pointer; }\n    .cs-icon-wrap { width: 34px; height: 34px; border-radius: 8px; background: #FBF0D8; display: flex; align-items: center; justify-content: center; font-size: 15px; flex-shrink: 0; }\n    .cs-content { flex: 1; }\n    .cs-name { font-size: 13px; font-weight: 600; color: #1a1a1a; }\n    .cs-desc { font-size: 11px; color: #bbb; margin-top: 2px; }\n    .cs-arrow { font-size: 11px; color: #C8922A; }\n\n    /* Channel flow preview */\n    .channel-flow-box { background: #fafafa; border: 1px solid #ebebeb; border-radius: 8px; padding: 14px 16px; }\n    .cfb-row { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }\n    .cfb-vault { display: flex; align-items: center; gap: 5px; padding: 5px 9px; border: 1px solid #e5e5e5; border-radius: 6px; background: #fff; font-size: 12px; font-weight: 500; color: #333; }\n    .cfb-dot { width: 7px; height: 7px; border-radius: 50%; }\n    .cfb-arrow { color: #ccc; font-size: 12px; }\n    .cfb-meta { display: flex; gap: 12px; margin-top: 10px; padding-top: 10px; border-top: 1px solid #f0f0f0; }\n    .cfb-metaitem { display: flex; align-items: center; gap: 5px; font-size: 11px; color: #bbb; }\n    .cfb-metaitem strong { color: #555; }\n\n    /* Thread name input */\n    .thread-name-wrap { display: flex; align-items: center; border: 1.5px solid #e5e5e5; border-radius: 8px; overflow: hidden; }\n    .thread-name-wrap:focus-within { border-color: #C8922A; background: #FFFDF8; }\n    .thread-name-prefix { padding: 10px 12px; background: #fafafa; border-right: 1px solid #e5e5e5; font-size: 12px; color: #bbb; white-space: nowrap; font-family: 'SF Mono', monospace; }\n    .thread-name-input { flex: 1; padding: 10px 12px; border: none; outline: none; font-size: 13px; color: #1a1a1a; font-family: 'Inter', sans-serif; background: transparent; }\n\n    /* Property fields */\n    .props-grid { display: flex; flex-direction: column; gap: 10px; }\n    .prop-row { display: flex; gap: 8px; align-items: flex-end; }\n    .prop-key-wrap { flex: 1; }\n    .prop-val-wrap { flex: 2; }\n    .prop-label { font-size: 10px; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase; color: #bbb; margin-bottom: 4px; }\n    .prop-input { width: 100%; padding: 8px 10px; border: 1px solid #e5e5e5; border-radius: 7px; font-size: 12px; color: #1a1a1a; outline: none; font-family: 'Inter', sans-serif; background: #fff; }\n    .prop-input:focus { border-color: #C8922A; background: #FFFDF8; }\n    .prop-input.prefilled { background: #fafafa; color: #888; }\n    .prop-remove { width: 26px; height: 32px; display: flex; align-items: center; justify-content: center; font-size: 14px; color: #ccc; cursor: pointer; flex-shrink: 0; border-radius: 5px; }\n    .prop-remove:hover { color: #EF4444; background: #FEE2E2; }\n    .add-prop-btn { display: inline-flex; align-items: center; gap: 5px; font-size: 11px; color: #C8922A; cursor: pointer; padding: 4px 0; }\n    .add-prop-btn:hover { opacity: 0.8; }\n\n    /* Vault overrides */\n    .vault-overrides { display: flex; flex-direction: column; gap: 8px; }\n    .vault-override-row { display: flex; align-items: center; gap: 10px; }\n    .vor-role { font-size: 12px; color: #888; width: 68px; flex-shrink: 0; }\n    .vor-select { flex: 1; display: flex; align-items: center; gap: 7px; padding: 7px 10px; border: 1px solid #e5e5e5; border-radius: 7px; background: #fafafa; cursor: pointer; }\n    .vor-select:hover { border-color: #C8922A; }\n    .vor-dot { width: 7px; height: 7px; border-radius: 50%; }\n    .vor-name { font-size: 12px; font-weight: 500; color: #333; flex: 1; }\n    .vor-arrow { font-size: 10px; color: #bbb; }\n\n    /* Stellar info bar */\n    .stellar-info { display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 7px; }\n    .si-icon { font-size: 12px; }\n    .si-text { font-size: 11px; color: #888; flex: 1; line-height: 1.4; }\n    .si-text strong { color: #1a1a1a; }\n    .si-status { display: flex; align-items: center; gap: 5px; }\n    .si-dot { width: 6px; height: 6px; border-radius: 50%; background: #22C55E; }\n    .si-label { font-size: 10px; color: #22C55E; font-weight: 600; }\n\n    /* Footer actions */\n    .sp-footer { padding: 16px 24px; border-top: 1px solid #ebebeb; flex-shrink: 0; display: flex; flex-direction: column; gap: 10px; }\n    .start-btn { width: 100%; padding: 12px; background: #C8922A; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 700; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 8px; letter-spacing: 0.02em; }\n    .start-btn:hover { background: #b8821a; }\n    .footer-note { font-size: 11px; color: #bbb; text-align: center; line-height: 1.4; }\n    .footer-note strong { color: #888; }\n  "
        }}
      />
      <style
    dangerouslySetInnerHTML={{
      __html:
        "\n    * { margin: 0; padding: 0; box-sizing: border-box; }\n    .window { width: 1320px; height: 820px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); overflow: hidden; display: flex; flex-direction: column; position: relative; }\n    .titlebar { height: 40px; background: #f5f5f5; border-bottom: 1px solid #e0e0e0; display: flex; align-items: center; padding: 0 16px; flex-shrink: 0; z-index: 1; }\n    .tls { display: flex; gap: 8px; }\n    .tl { width: 12px; height: 12px; border-radius: 50%; }\n    .tl-r { background: #FF5F57; } .tl-y { background: #FFBD2E; } .tl-g { background: #28CA41; }\n    .tb-name { flex: 1; text-align: center; font-size: 12px; font-weight: 500; color: #666; letter-spacing: 0.04em; }\n    .topbar { height: 52px; background: #fff; border-bottom: 1px solid #ebebeb; display: flex; align-items: center; padding: 0 20px; gap: 16px; flex-shrink: 0; opacity: 0.45;  }\n    .tb-logo { font-weight: 700; font-size: 14px; letter-spacing: 0.12em; }\n    .tb-sp { flex: 1; }\n    .btn-g { background: #f5f5f5; border: 1px solid #e5e5e5; color: #444; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n    .btn-p { background: #C8922A; color: #fff; padding: 6px 14px; border-radius: 7px; font-size: 12px; }\n    .layout { flex: 1; display: flex; overflow: hidden; }\n\n    /* Dimmed ledger */\n    .ledger-wrapper { flex: 1; display: flex; overflow: hidden; opacity: 0.38;  }\n    .sidebar {  background: #fafafa; border-right: 1px solid #ebebeb; display: flex; flex-direction: column; padding: 18px 0; flex-shrink: 0; }\n    .sidebar-section-label { padding: 0 16px 8px; font-size: 10px; font-weight: 600; letter-spacing: 0.1em; color: #bbb; text-transform: uppercase; }\n    .vault-row { display: flex; align-items: center; gap: 9px; padding: 7px 16px; font-size: 13px; color: #333; }\n    .vault-row.all { font-size: 12px; color: #888; padding: 5px 16px 10px; border-bottom: 1px solid #ebebeb; margin-bottom: 6px; }\n    .vdot { width: 8px; height: 8px; border-radius: 50%; }\n    .pip { margin-left: auto; width: 6px; height: 6px; background: #EF4444; border-radius: 50%; }\n    .ledger-area { flex: 1; display: flex; flex-direction: column; overflow: hidden; }\n    .ledger-topbar { display: flex; align-items: center; justify-content: space-between; padding: 18px 24px 14px; }\n    .ledger-title { font-size: 17px; font-weight: 700; letter-spacing: 0.06em; }\n    .table-wrap { flex: 1; overflow: hidden; padding: 0 24px; }\n    table { width: 100%; border-collapse: separate; border-spacing: 0; }\n    thead th { padding: 0 10px 10px; text-align: left; font-size: 10px; font-weight: 600; letter-spacing: 0.09em; text-transform: uppercase; color: #bbb; border-bottom: 1px solid #ebebeb; }\n    td { padding: 11px 10px; vertical-align: middle; border-bottom: 1px solid #f2f2f2; font-size: 13px; }\n    .row-sel td { background: #FEF2F2 !important; }\n    .sdot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }\n    .s-ok { background: #22C55E; } .s-pend { background: #F59E0B; } .s-dispute { background: #EF4444; }\n    .th-type { font-size: 9px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-right: 6px; }\n    .flow { display: flex; align-items: center; gap: 4px; }\n    .vb { width: 22px; height: 22px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 8px; font-weight: 700; color: #fff; }\n    .fa { color: #ccc; font-size: 10px; }\n    .pipeline { display: flex; gap: 3px; }\n    .pseg { height: 4px; width: 20px; border-radius: 2px; }\n    .pseg-done { background: #C8922A; } .pseg-wait { background: #e8e8e8; } .pseg-reject { background: #EF4444; }\n    .ts { font-size: 12px; color: #bbb; }\n    .stellar-val { font-family: 'SF Mono', monospace; font-size: 10px; color: #ccc; }\n    .c3b { display: inline-flex; align-items: center; padding: 3px 8px; border-radius: 5px; font-size: 11px; font-weight: 600; background: #f5f5f5; border: 1px solid #e5e5e5; color: #999; }\n    .dispute-badge-row { display: inline-flex; align-items: center; gap: 4px; background: #FEF2F2; border: 1px solid #FECACA; color: #DC2626; padding: 2px 7px; border-radius: 4px; font-size: 10px; font-weight: 700; letter-spacing: 0.05em; }\n    .ledger-footer { padding: 12px 24px; border-top: 1px solid #ebebeb; font-size: 12px; color: #bbb; flex-shrink: 0; }\n\n    /* ---- DETAIL PANEL ---- */\n    .detail-panel { margin-top: -1rem ; height: 100%; background: #fff; border-left: 1px solid #e0e0e0; display: flex; flex-direction: column; flex-shrink: 0; box-shadow: -6px 0 24px rgba(0,0,0,0.08); overflow: hidden; z-index: 10; }\n\n    /* Dispute alert banner */\n    .dispute-banner { background: #FEF2F2; border-bottom: 1px solid #FECACA; padding: 10px 20px; display: flex; align-items: center; gap: 10px; flex-shrink: 0; }\n    .db-icon { font-size: 14px; }\n    .db-text { font-size: 12px; color: #991B1B; flex: 1; font-weight: 500; }\n    .db-text span { font-weight: 700; }\n    .db-ts { font-size: 11px; color: #EF4444; white-space: nowrap; }\n\n    .dp-header { padding: 16px 22px 14px; border-bottom: 1px solid #ebebeb; flex-shrink: 0; }\n    .dp-close-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }\n    .dp-badge { font-size: 10px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #DC2626; background: #FEF2F2; padding: 3px 8px; border-radius: 4px; border: 1px solid #FECACA; }\n    .dp-close { width: 24px; height: 24px; border-radius: 6px; background: #f5f5f5; border: 1px solid #e5e5e5; display: flex; align-items: center; justify-content: center; font-size: 12px; color: #888; cursor: pointer; }\n    .dp-title { font-size: 15px; font-weight: 700; color: #1a1a1a; margin-bottom: 6px; }\n    .dp-meta { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 6px; }\n    .dp-meta-chip { font-size: 11px; color: #999; background: #f8f8f8; border: 1px solid #ebebeb; padding: 2px 8px; border-radius: 4px; font-family: 'SF Mono', monospace; }\n    .dp-custom { margin-top: 6px; font-size: 11px; color: #999; }\n\n    .dp-body { flex: 1; overflow-y: auto; padding: 0 22px 20px; }\n    .dp-section { padding: 14px 0 12px; border-bottom: 1px solid #f0f0f0; }\n    .dp-section:last-child { border-bottom: none; }\n    .dp-section-title { font-size: 10px; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #bbb; margin-bottom: 12px; }\n\n    /* Pipeline steps */\n    .pipeline-step { display: flex; align-items: flex-start; gap: 10px; margin-bottom: 10px; }\n    .ps-icon { width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 9px; font-weight: 700; color: #fff; flex-shrink: 0; margin-top: 1px; }\n    .ps-icon-wait { background: #f0f0f0 !important; border: 2px dashed #ddd; color: #ccc !important; }\n    .ps-icon-rejected { background: #FEF2F2 !important; border: 2px solid #EF4444 !important; color: #EF4444 !important; }\n    .ps-content { flex: 1; }\n    .ps-label { font-size: 12px; font-weight: 500; color: #333; }\n    .ps-label.rejected { color: #DC2626; font-weight: 600; }\n    .ps-label.wait { color: #ccc; font-style: italic; }\n    .ps-sublabel { font-size: 11px; color: #bbb; margin-top: 2px; }\n    .ps-sublabel.rejected { color: #DC2626; opacity: 0.7; }\n    .ps-ts { font-size: 10px; color: #ccc; margin-top: 2px; }\n    .ps-check { font-size: 13px; flex-shrink: 0; }\n    .ps-check.done { color: #22C55E; }\n    .ps-check.reject { color: #EF4444; }\n    .ps-check.wait { color: #e0e0e0; }\n    .ps-connector { width: 1px; height: 8px; background: #e8e8e8; margin-left: 12px; }\n    .ps-connector.red { background: #FECACA; }\n    .reject-reason-inline { margin-top: 6px; padding: 8px 10px; background: #FEF2F2; border: 1px solid #FECACA; border-radius: 6px; font-size: 11px; color: #991B1B; line-height: 1.5; }\n    .rri-label { font-size: 9px; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; color: #EF4444; margin-bottom: 4px; }\n\n    /* Commits */\n    .commit-row { display: flex; align-items: flex-start; gap: 10px; padding: 8px 0; border-bottom: 1px solid #f5f5f5; }\n    .commit-row:last-child { border-bottom: none; }\n    .commit-row.rejected-row { background: #FEF2F2; margin: 0 -4px; padding: 8px 4px; border-radius: 5px; border-bottom: none; }\n    .commit-ts { font-size: 10px; color: #bbb; white-space: nowrap; padding-top: 1px; min-width: 70px; }\n    .commit-ts.red { color: #DC2626; opacity: 0.8; }\n    .commit-body { flex: 1; }\n    .commit-actor { font-size: 11px; font-weight: 600; color: #555; }\n    .commit-actor.red { color: #DC2626; }\n    .commit-action { font-size: 11px; color: #333; margin-left: 4px; }\n    .commit-action.red { color: #DC2626; }\n    .commit-cid { font-size: 10px; color: #bbb; font-family: 'SF Mono', monospace; margin-top: 2px; }\n    .commit-reason { font-size: 11px; color: #991B1B; margin-top: 5px; line-height: 1.4; font-style: italic; }\n    .commit-reason::before { content: '\"'; } .commit-reason::after { content: '\"'; }\n    .commit-verify { font-size: 10px; color: #C8922A; cursor: pointer; white-space: nowrap; padding-top: 1px; }\n    .commit-verify.red { color: #EF4444; }\n    .reject-badge { display: inline-flex; align-items: center; gap: 3px; font-size: 9px; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase; background: #FEF2F2; border: 1px solid #FECACA; color: #DC2626; padding: 1px 5px; border-radius: 3px; margin-left: 5px; }\n\n    /* Stellar ref */\n    .stellar-ref-row { display: flex; align-items: center; gap: 8px; padding: 10px 12px; background: #fafafa; border: 1px solid #ebebeb; border-radius: 7px; }\n    .stellar-hash-full { font-family: 'SF Mono', monospace; font-size: 10px; color: #888; flex: 1; overflow: hidden; text-overflow: ellipsis; }\n    .copy-btn { font-size: 10px; color: #C8922A; cursor: pointer; white-space: nowrap; }\n\n    /* Resolution callout */\n    .resolution-box { padding: 14px 16px; background: #FFFBEB; border: 1px solid #FCD34D; border-radius: 8px; margin-top: 4px; }\n    .rb-header { display: flex; align-items: center; gap: 7px; margin-bottom: 8px; }\n    .rb-icon { font-size: 13px; }\n    .rb-title { font-size: 12px; font-weight: 600; color: #92400E; }\n    .rb-steps { display: flex; flex-direction: column; gap: 6px; }\n    .rb-step { display: flex; align-items: flex-start; gap: 7px; font-size: 11px; color: #78350F; line-height: 1.4; }\n    .rb-step-num { font-size: 10px; font-weight: 700; background: #FCD34D; color: #78350F; width: 16px; height: 16px; border-radius: 50%; display: flex; align-items: center; justify-content: center; flex-shrink: 0; margin-top: 0px; }\n\n    /* Actions */\n    .dp-actions { padding: 14px 22px; border-top: 1px solid #ebebeb; display: flex; gap: 8px; flex-shrink: 0; }\n    .action-btn { flex: 1; padding: 8px 10px; border: 1px solid #e5e5e5; border-radius: 7px; background: #fff; font-size: 11px; font-weight: 500; color: #444; cursor: pointer; text-align: center; }\n    .action-btn:hover { border-color: #C8922A; color: #C8922A; }\n    .action-btn.primary { background: #C8922A; color: #fff; border-color: #C8922A; }\n    .action-btn.primary:hover { background: #b8821a; }\n  "
    }}
  />

  

</>
      
  )
}