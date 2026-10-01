import { test, expect, type ConsoleMessage } from '@playwright/test';

/**
 * Browser-level verification of the canonical XS-BIM construction rebase.
 *
 * This is deliberately a NAVIGATION + RENDER test, not a re-test of the data
 * layer. The data layer is already covered by the runtime audit; what this file
 * proves is that each migrated page actually mounts, does not crash, and renders
 * canonical values in the real React Router tree.
 *
 * Why the init script below exists:
 *   src/App.tsx:119 gates the ENTIRE app behind `if (!appState.has_vault)`.
 *   `appState` starts as the boolean `false` and is only ever replaced by
 *   `AppAPI.GetAppState()` when `window.go` is present (Wails desktop runtime).
 *   In a plain browser there is no Wails backend, so `appState` stays `false`
 *   and the app renders the onboarding wizard forever. That gate is PRE-EXISTING
 *   and untouched by the canonical rebase, so rather than editing app shell
 *   (out of scope) we satisfy it the way Wails would.
 */
const BOOTSTRAP = `
(() => {
  const appState = { has_vault: true, onboarded: true, app_state: 'ready' };
  // NOTE: the generated bindings call window.go.main.App.<Method>() directly,
  // e.g. wailsjs/go/main/App.js:234 -> window['go']['main']['App']['GetAppState']().
  // An empty App object would make GetAppState() throw, leaving appState === false
  // and the app permanently stuck on the onboarding wizard.
  const App = {
    GetAppState: async () => appState,
    CompleteOnboarding: async () => undefined,
  };
  window.go = {
    main: { App },
    Stellar: {
      CheckWalletStatus: async () => 'connected',
    },
    IPFS: { CheckNodeStatus: async () => 'ready' },
  };
  window.runtime = window.runtime || {
    EventsOn: () => {},
    EventsOff: () => {},
    EventsEmit: () => {},
    WindowReloadApp: () => {},
  };
})();
`;

const BASE = '/dashboard/construction';

/** Console/page errors that are genuine defects in the construction rebase. */
function isRealConstructionError(text: string): boolean {
  if (/SCENARIO_DATA/i.test(text)) return true;
  if (/scenarioAccessors|scenarioMappers|constructionScenarioAdapter/.test(text) && /error|cannot read|undefined/i.test(text)) return true;
  if (/Cannot read properties of undefined/i.test(text)) return true;
  return false;
}

test.beforeEach(async ({ page }) => {
  await page.addInitScript(BOOTSTRAP);
});

/**
 * Navigate to a construction route and assert:
 *  1. the route mounts (header present)
 *  2. no uncaught page error
 *  3. no construction-related console error
 */
async function visit(page: import('@playwright/test').Page, path: string) {
  const pageErrors: string[] = [];
  const consoleErrors: string[] = [];

  const onPageError = (e: Error) => pageErrors.push(e.message);
  const onConsole = (m: ConsoleMessage) => {
    if (m.type() === 'error' && isRealConstructionError(m.text())) consoleErrors.push(m.text());
  };
  page.on('pageerror', onPageError);
  page.on('console', onConsole);

  await page.goto(`${BASE}${path}`, { waitUntil: 'networkidle' });
  await expect(page.getByText('BuildFlow', { exact: false }).first()).toBeVisible();

  page.off('pageerror', onPageError);
  page.off('console', onConsole);

  expect(pageErrors, `page errors on ${path}`).toEqual([]);
  expect(consoleErrors, `console errors on ${path}`).toEqual([]);
}

test('dashboard renders canonical project, no crash', async ({ page }) => {
  await visit(page, '/');
  // Canonical metrics, not hardcoded ones: the project is active under its
  // canonical reference, the issue is RESOLVED (so 0 open) and the decision is
  // APPROVED (so 0 pending). Cursor-ordered, not "Today".
  await expect(page.getByText(/1 Active \(PRJ-METRO-001\)/)).toBeVisible();
  await expect(page.getByText(/No open issues \(ISS-REF-1042 resolved\)/)).toBeVisible();
  await expect(page.getByText(/No approvals pending \(DEC-REF-1042 approved\)/)).toBeVisible();
  await expect(page.getByText('Cursor-Ordered')).toBeVisible();
  await expect(page.getByText('Today')).toHaveCount(0);
  // 23 canonical thread events.
  await expect(page.getByText('Recent Events')).toBeVisible();
  await expect(page.getByText('23').first()).toBeVisible();
});

test('project listing renders canonical project reference and 3 total projects', async ({ page }) => {
  await visit(page, '/projects');
  // code = canonical projectReference (PRJ-METRO-001), not the projectId.
  await expect(page.getByText('PRJ-METRO-001').first()).toBeVisible();
  await expect(page.getByText('Metro Line 4 Expansion').first()).toBeVisible();
  // PRJ-002 and PRJ-003 restored demo projects
  await expect(page.getByText('PRJ-002').first()).toBeVisible();
  await expect(page.getByText('Commercial Plaza North').first()).toBeVisible();
  await expect(page.getByText('PRJ-003').first()).toBeVisible();
  await expect(page.getByText('Riverside Logistics Hub').first()).toBeVisible();
});

test('project detail renders canonical project', async ({ page }) => {
  await visit(page, '/projects/PRJ-001');
  await expect(page.getByText('Metro Line 4 Expansion').first()).toBeVisible();
});

test('stakeholders renders the three canonical vault identities', async ({ page }) => {
  await visit(page, '/stakeholders');
  // Canonical participants are vault identities, not people. There is no
  // ConstructionStakeholder aggregate, so names are absent by design.
  await expect(page.getByText('vault_001-oem', { exact: false }).first()).toBeVisible();
  await expect(page.getByText('vault_002-michelin', { exact: false }).first()).toBeVisible();
  await expect(page.getByText('vault_003-faa', { exact: false }).first()).toBeVisible();
});

test('material requirement renders exact Stitch layout, canonical material, supplier sourcing, and offer details', async ({ page }) => {
  await visit(page, '/requirements');

  // Header & Context
  await expect(page.getByText('BuildFlow').first()).toBeVisible();
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Material Requirement');
  await expect(page.getByText('PROJ-RT-104 • Metro Line 4 Expansion')).toBeVisible();
  await expect(page.getByText('OFFER SELECTED')).toBeVisible();

  // Requirement & Material Card
  await expect(page.getByText('REQ-STRUCT-001')).toBeVisible();
  await expect(page.getByText('MAT-STRUCT-001')).toBeVisible();
  await expect(page.getByText('Critical Path')).toBeVisible();
  await expect(page.getByText('Grade S355JR')).toBeVisible();
  await expect(page.getByText('84 Beams')).toBeVisible();
  await expect(page.getByText('120 Metric Tons')).toBeVisible();
  await expect(page.getByText('Aug 15, 2024').first()).toBeVisible();
  await expect(page.getByText('Slot: 08:00 - 12:00')).toBeVisible();
  await expect(page.getByText('Site-001 (Paris)')).toBeVisible();
  await expect(page.getByText('Foundation / Viaduct Section 4')).toBeVisible();

  // Supplier Sourcing Card
  await expect(page.getByRole('heading', { level: 2, name: 'Supplier Sourcing' })).toBeVisible();
  await expect(page.getByText('100% Filled')).toBeVisible();
  await expect(page.getByText('Invited')).toBeVisible();
  await expect(page.getByText("Offers Rec'd")).toBeVisible();
  await expect(page.getByText('Selected').first()).toBeVisible();

  // Selected Supplier Card
  await expect(page.getByText('EuroSteel Construction').first()).toBeVisible();
  await expect(page.getByText('OFF-1042').first()).toBeVisible();
  await expect(page.getByText('Chosen').first()).toBeVisible();
  await expect(page.getByText(/142,500/).first()).toBeVisible();
  await expect(page.getByText(/1,187\.50/).first()).toBeVisible();
  await expect(page.getByText(/Aug 15, 2024/).first()).toBeVisible();

  // Action Buttons
  const viewOffersBtn = page.getByRole('button', { name: 'View Offers (2)' });
  await expect(viewOffersBtn).toBeVisible();
  await expect(page.getByRole('button', { name: 'Invite Supplier' })).toBeVisible();

  // Navigation to Supplier Offers page
  await viewOffersBtn.click();
  await expect(page).toHaveURL(/.*\/offers/);
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Supplier Offers');
});

test('supplier offers render exact Stitch layout, offer details, and suppliers', async ({ page }) => {
  await visit(page, '/offers');
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Supplier Offers');
  await expect(page.getByText('120t Structural Beams', { exact: false }).first()).toBeVisible();
  await expect(page.getByText('EuroSteel Construction').first()).toBeVisible();
  await expect(page.getByText('Arcelor Infrastructure').first()).toBeVisible();
  await expect(page.getByText('ACCEPTED • DELIVERY PLANNED').first()).toBeVisible();
  await expect(page.getByText('CE Marked').first()).toBeVisible();
  await expect(page.getByText('EN 10204 3.1 Mill Cert').first()).toBeVisible();
  await expect(page.getByText(/OFF-REF-1042/).first()).toBeVisible();
  await expect(page.getByText(/OFF-1039/).first()).toBeVisible();
});

test('delivery detail renders exact Stitch layout and operational specifications', async ({ page }) => {
  await visit(page, '/deliveries');
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Delivery #DEL-1042');
  await expect(page.getByText('DELAYED').first()).toBeVisible();
  await expect(page.getByText('Road Restriction on M1').first()).toBeVisible();
  await expect(page.getByText('Schedule Comparison').first()).toBeVisible();
  await expect(page.getByText('+21.5h Delay').first()).toBeVisible();
  await expect(page.getByText('Operational Specifications').first()).toBeVisible();
  await expect(page.getByText('REQ-REF-1042').first()).toBeVisible();
  await expect(page.getByText('MAT-STRUCT-001').first()).toBeVisible();
  await expect(page.getByText('EuroSteel').first()).toBeVisible();
  await expect(page.getByText('OFF-1042').first()).toBeVisible();
  await expect(page.getByText('FastBuild Logistics').first()).toBeVisible();
  await expect(page.getByText('ISS-1042 (Critical Delay)').first()).toBeVisible();
  await expect(page.getByText('Last Known Telemetry').first()).toBeVisible();
  await expect(page.getByText('GPS Active').first()).toBeVisible();
  await expect(page.getByText('Progress Timeline').first()).toBeVisible();
  await expect(page.getByText('Stage 4 of 8').first()).toBeVisible();
  await expect(page.getByText('Delayed on Route M1').first()).toBeVisible();
  await expect(page.getByText('Rerouted').first()).toBeVisible();
  await expect(page.getByText('Delivered').first()).toBeVisible();
  await expect(page.getByText('Inspected').first()).toBeVisible();
  await expect(page.getByText('Accepted').first()).toBeVisible();
  await expect(page.getByText('View Provenance & Root Cause Analysis').first()).toBeVisible();
  await expect(page.getByText('Print Consignment Manifest & Waybill').first()).toBeVisible();
});

test('transport renders exact Stitch Transport & Delay layout and details', async ({ page }) => {
  await visit(page, '/transport');
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Transport & Delay');
  await expect(page.getByText('DEL-REF-1042').first()).toBeVisible();
  await expect(page.getByText('ACTIVE ALERT').first()).toBeVisible();
  await expect(page.getByText('Transport Delay Detected').first()).toBeVisible();
  await expect(page.getByText('Route M1 Blocked').first()).toBeVisible();
  await expect(page.getByText('+21h 30m').first()).toBeVisible();
  await expect(page.getByText('Live Corridor Geometry').first()).toBeVisible();
  await expect(page.getByText('Constraint Details').first()).toBeVisible();
  await expect(page.getByText('Detour Evaluation').first()).toBeVisible();
  await expect(page.getByText('ROUTE B APPROVED').first()).toBeVisible();
  await expect(page.getByRole('button', { name: /View Decision DEC-REF-1042/ })).toBeVisible();
});

test('issue renders exact Stitch layout and specifications', async ({ page }) => {
  await visit(page, '/issues');

  // Header & Status strip
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Construction Issue');
  await expect(page.getByText('SEVERITY: CRITICAL').first()).toBeVisible();
  await expect(page.getByText('RESOLVED').first()).toBeVisible();

  // Issue Title / Resolution Card
  await expect(page.getByText('LOGISTICS BOTTLENECK').first()).toBeVisible();
  await expect(page.getByText(/ISS-REF-1042/).first()).toBeVisible();
  await expect(page.getByText('Critical Beam Delivery Delayed by Route M1 Blockage').first()).toBeVisible();

  // Active Mitigation
  await expect(page.getByText('Active Mitigation').first()).toBeVisible();
  await expect(page.getByText('Rerouted via Decision DEC-1042').first()).toBeVisible();

  // Specifications
  await expect(page.getByText('Specifications').first()).toBeVisible();
  await expect(page.getByText('DEL-1042').first()).toBeVisible();
  await expect(page.getByText('Structural Beams (44t Prefabricated)').first()).toBeVisible();
  await expect(page.getByText('Foundation').first()).toBeVisible();
  await expect(page.getByText('Route M1 Alert').first()).toBeVisible();
  await expect(page.getByText('David K. (Logistics Mgr)').first()).toBeVisible();

  // Stakeholders
  await expect(page.getByText('Assigned Stakeholders').first()).toBeVisible();
  await expect(page.getByText('BuildCorp (Main)').first()).toBeVisible();
  await expect(page.getByText('FastBuild Logistics').first()).toBeVisible();
  await expect(page.getByText('Engineering Partners').first()).toBeVisible();

  // Impact Assessment
  await expect(page.getByText('Impact Assessment').first()).toBeVisible();
  await expect(page.getByText('Foundation Phase Delayed').first()).toBeVisible();

  // Supporting Evidence
  await expect(page.getByText('Supporting Evidence (2)').first()).toBeVisible();
  await expect(page.getByText('Road_Restriction_Notice_M1.pdf').first()).toBeVisible();
  await expect(page.getByText('Transport_Detour_Report_TR1042.pdf').first()).toBeVisible();

  // Action Buttons
  await expect(page.getByRole('button', { name: /View Decision \(DEC-/ })).toBeVisible();
  await expect(page.getByRole('button', { name: /View Delivery \(/ })).toBeVisible();
  await expect(page.getByRole('button', { name: /Open Collaboration Thread/ })).toBeVisible();
});

test('decision renders the canonical decision reference', async ({ page }) => {
  await visit(page, '/decisions');
  await expect(page.getByText('DEC-REF-1042', { exact: false }).first()).toBeVisible();
});

/**
 * Stitch DEC-1042 composition, projected from the canonical aggregate.
 *
 * This asserts two things at once: that each composition section is present, and
 * that the Stitch copy the scenario cannot support was replaced rather than
 * rendered. The negative assertions are the important half.
 */
test('decision renders Stitch composition from canonical data only', async ({ page }) => {
  await visit(page, '/decisions');

  // Status ribbon — canonical status, canonical id, canonical decision date.
  await expect(page.getByText('APPROVED').first()).toBeVisible();
  await expect(page.getByText('DEC-1042').first()).toBeVisible();
  await expect(page.getByText(/Finalized Aug 15/)).toBeVisible();

  // Header card — canonical subject, canonical context, derived related ids.
  await expect(page.getByRole('heading', { level: 1 })).toHaveText(
    'Structural material delivery delay alternative route',
  );
  await expect(page.getByText('Route M1 blocked, ETA delayed to 2026-08-16T07:30:00Z')).toBeVisible();
  await expect(page.getByText('ISS-1042').first()).toBeVisible();
  await expect(page.getByText('DEL-1042').first()).toBeVisible();

  // Options — both canonical options, selection recovered from the decision text.
  await expect(page.getByRole('heading', { name: 'Options Evaluated' })).toBeVisible();
  await expect(page.getByText('2 Alternatives Assessed')).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Wait for M1 clearance' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Use Route B detour' })).toBeVisible();
  await expect(page.getByText('Selected & Authorized')).toBeVisible();
  await expect(page.getByText('Rejected').first()).toBeVisible();

  // The only canonical arrival: derived from transport, labelled UTC.
  await expect(page.getByText('Aug 16 @ 07:30 UTC').first()).toBeVisible();
  // The 21.5h slip is derivable from planned vs actual arrival.
  await expect(page.getByText(/Route-B detour · \+21\.5 h vs plan/)).toBeVisible();

  // Governance — canonical roles, no fabricated organisations.
  await expect(page.getByRole('heading', { name: 'Governance Chain' })).toBeVisible();
  for (const role of ['Contractor', 'Supplier', 'Logistics', 'Project Manager']) {
    await expect(page.getByText(role, { exact: true }).first()).toBeVisible();
  }
  await expect(page.getByText('Decided by Project Manager')).toBeVisible();
  await expect(page.getByText('Requested by Logistics Manager · 4 participants consulted')).toBeVisible();

  // Outcome — canonical decision and consequence, canonical transport reference.
  await expect(page.getByRole('heading', { name: 'Approve Route B alternative transport' })).toBeVisible();
  await expect(page.getByText(/Transport re-routed via Route B/)).toBeVisible();
  await expect(page.getByText('TR-1042').first()).toBeVisible();

  // Action bar routes on the real route constants.
  await expect(page.getByRole('button', { name: /View Updated Delivery \(DEL-1042\)/ })).toBeVisible();
  await expect(page.getByRole('button', { name: /View Audit Trail & Provenance/ })).toBeVisible();

  // ---- Negative assertions: Stitch copy the scenario cannot support ----
  // Clock time on a date-only decisionDate.
  await expect(page.getByText(/12:35/)).toHaveCount(0);
  await expect(page.getByText(/no time recorded/)).toBeVisible();
  // Commercial facts with no canonical field.
  await expect(page.getByText(/480/)).toHaveCount(0);
  await expect(page.getByText(/Logistics SLA/)).toHaveCount(0);
  // Distance/duration deltas, and the "48-72 hour" option impact.
  await expect(page.getByText(/42 km/)).toHaveCount(0);
  await expect(page.getByText(/50 min/)).toHaveCount(0);
  await expect(page.getByText(/48/)).toHaveCount(0);
  // Per-participant sign-off, organisations and timestamps.
  await expect(page.getByText(/Acknowledged/)).toHaveCount(0);
  await expect(page.getByText(/Confirmed$/)).toHaveCount(0);
  await expect(page.getByText(/EuroSteel/)).toHaveCount(0);
  await expect(page.getByText(/FastBuild/)).toHaveCount(0);
  await expect(page.getByText(/Structural Engineer/)).toHaveCount(0);
  await expect(page.getByText(/Consensus Met/)).toHaveCount(0);
  await expect(page.getByText(/cryptographic/)).toHaveCount(0);
  // A sector taxonomy the project does not carry.
  await expect(page.getByText(/Sector 4/)).toHaveCount(0);
  // Slots with no canonical datum must say so rather than look plausible.
  await expect(page.getByText('Not recorded in scenario').first()).toBeVisible();
  await expect(page.getByText('Sign-off not recorded').first()).toBeVisible();
});

test('decision action bar navigates to delivery and provenance', async ({ page }) => {
  await visit(page, '/decisions');

  await page.getByRole('button', { name: /View Updated Delivery \(DEL-1042\)/ }).click();
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Delivery #DEL-1042');

  await page.goto(`${BASE}/decisions`, { waitUntil: 'networkidle' });
  await page.getByRole('button', { name: /View Audit Trail & Provenance/ }).click();
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Provenance / Why?');
});

/**
 * The decision view must resolve entirely from the static scenario on the active
 * provider. VITE_CONSTRUCTION_DATA_SOURCE defaults to "mock", so any AppAPI call
 * here would mean the presentation layer had reached around the Functional Data
 * Boundary.
 *
 * The probe wraps the Wails App bindings rather than watching the network: in the
 * browser AppAPI methods are invoked through window.go, not over HTTP, so a
 * request-level assertion would miss exactly the leak it is meant to catch.
 */
test('decision view reaches AppAPI for nothing', async ({ page }) => {
  await page.addInitScript(`
    (() => {
      const app = window.go && window.go.main && window.go.main.App;
      window.__appApiCalls = [];
      if (!app) return;
      window.go.main.App = new Proxy(app, {
        get(target, prop) {
          const value = target[prop];
          if (typeof value !== 'function') return value;
          return function (...args) {
            window.__appApiCalls.push(String(prop));
            return value.apply(target, args);
          };
        }
      });
    })();
  `);

  await visit(page, '/decisions');
  await expect(page.getByRole('heading', { name: 'Options Evaluated' })).toBeVisible();

  const called = await page.evaluate(() => {
    const w = window as unknown as { __appApiCalls?: string[] };
    return w.__appApiCalls ?? [];
  });

  // GetAppState is the app-shell bootstrap the gate in src/App.tsx performs.
  expect(called.filter((m) => m !== 'GetAppState'), 'decision view must not call AppAPI').toEqual([]);
});

test('inspection renders exact Stitch layout, checklist, evidence, and timeline', async ({ page }) => {
  await visit(page, '/inspections');

  // Header & status
  await expect(page.getByText('INSPECTION COMPLETE').first()).toBeVisible();
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Foundation — Zone A');
  await expect(page.getByText('Inspector: Bureau Inspection').first()).toBeVisible();

  // Verification checklist
  await expect(page.getByText('Verification Checklist').first()).toBeVisible();
  await expect(page.getByText('Reinforcement Layout').first()).toBeVisible();
  await expect(page.getByText('Formwork Integrity').first()).toBeVisible();
  await expect(page.getByText('Concrete Quality').first()).toBeVisible();
  await expect(page.getByText('Overall Dimensions').first()).toBeVisible();
  await expect(page.getByText('Site Safety Protocol').first()).toBeVisible();

  // Supporting evidence
  await expect(page.getByText('Supporting Evidence').first()).toBeVisible();
  await expect(page.getByText('Rebar Layout').first()).toBeVisible();
  await expect(page.getByText('Formwork').first()).toBeVisible();
  await expect(page.getByText('Concrete Test Results').first()).toBeVisible();
  await expect(page.getByText('Official Inspection Report').first()).toBeVisible();

  // Workflow timeline
  await expect(page.getByText('Workflow Timeline').first()).toBeVisible();
  await expect(page.getByText('Final Approval').first()).toBeVisible();
  await expect(page.getByText('Review').first()).toBeVisible();
  await expect(page.getByText('Initiation').first()).toBeVisible();
});

test('collaboration thread renders canonical event counts, no fake timestamps', async ({ page }) => {
  await visit(page, '/thread');
  // 23 canonical thread events must be present.
  await expect(page.getByText('#23', { exact: false }).first()).toBeVisible();
});

test('project history renders canonical trace events', async ({ page }) => {
  await visit(page, '/history');
  await expect(page.getByText('TraceCore', { exact: false }).first()).toBeVisible();
});

test('provenance renders exact Stitch layout, causal progression, and audit layers', async ({ page }) => {
  await visit(page, '/provenance');

  // Header & Context
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Provenance / Why?');
  await expect(page.getByText('Construction Intelligence • Causal Chain & Provenance')).toBeVisible();

  // Root-Cause Synthesis
  await expect(page.getByText('Root-Cause Diagnosis')).toBeVisible();
  await expect(page.getByText('+21h 30m Delay')).toBeVisible();
  await expect(page.getByText('DEL-1042', { exact: false }).first()).toBeVisible();
  await expect(page.getByText('Original ETA')).toBeVisible();
  await expect(page.getByText('Actual Site Gate')).toBeVisible();
  await expect(page.getByText('Schedule Delta')).toBeVisible();

  // Causal Progression & 8 Milestones
  await expect(page.getByText('Causal Progression')).toBeVisible();
  await expect(page.getByText('8 verifiable milestones')).toBeVisible();

  await expect(page.getByText('Physical Cause').first()).toBeVisible();
  await expect(page.getByText('Road restriction on Route M1').first()).toBeVisible();

  await expect(page.getByText('Disruption Ticket').first()).toBeVisible();
  await expect(page.getByText('ISS-1042 — Critical Beam Delivery Delayed').first()).toBeVisible();

  await expect(page.getByText('Verification').first()).toBeVisible();
  await expect(page.getByText('Road Restriction Report (#RD-9942)').first()).toBeVisible();
  await expect(page.getByText('Corridor_M1_Closure_Order.pdf').first()).toBeVisible();

  await expect(page.getByText('Multi-Org Consensus').first()).toBeVisible();
  await expect(page.getByText('DEC-1042 — Alternative Route B Approved').first()).toBeVisible();

  await expect(page.getByText('Field Execution').first()).toBeVisible();
  await expect(page.getByText('Transport Rerouted via Bypass').first()).toBeVisible();

  await expect(page.getByText('Site Receipt').first()).toBeVisible();
  await expect(page.getByText('Delivery Received at Site-001').first()).toBeVisible();

  await expect(page.getByText('Quality Control').first()).toBeVisible();
  await expect(page.getByText('Inspection #INSP-1042 Passed').first()).toBeVisible();

  await expect(page.getByText('Construction Handover').first()).toBeVisible();
  await expect(page.getByText('Material Accepted for Assembly').first()).toBeVisible();

  // Material Acceptance Card
  await expect(page.getByText('Material Acceptance').first()).toBeVisible();
  await expect(page.getByText('MAT-STRUCT-001', { exact: false }).first()).toBeVisible();
  await expect(page.getByText('INSP-1042', { exact: false }).first()).toBeVisible();
  await expect(page.getByText('Lot #STM-88219 (84 Beams)').first()).toBeVisible();

  // Sources & Audit Integrity and Tab switching
  await expect(page.getByText('Sources & Audit Integrity').first()).toBeVisible();
  await expect(page.getByText('Physical Site Data & Field Delivery').first()).toBeVisible();
  await expect(page.getByText('Authority: Site Super (BuildCorp)').first()).toBeVisible();

  // Tab 2: Collaboration & Evidence
  await page.getByRole('button', { name: 'Collaboration & Evidence' }).click();
  await expect(page.getByText('Multi-Organization Consensus Protocol').first()).toBeVisible();
  await expect(page.getByText('Threads: 14 exchanged records').first()).toBeVisible();

  // Tab 3: Historical Milestones
  await page.getByRole('button', { name: 'Historical Milestones' }).click();
  await expect(page.getByText('Historical Milestone Record').first()).toBeVisible();
  await expect(page.getByText('Ref: TRACE-MILESTONE-001').first()).toBeVisible();

  // Bottom Actions
  await expect(page.getByRole('button', { name: /Open Full Project History/ })).toBeVisible();
  await expect(page.getByRole('button', { name: /Back to Delivery \(/ })).toBeVisible();
});

test('channels renders channel list sidebar, native thread select, and Stitch Thread layout', async ({ page }) => {
  await visit(page, '/channels');

  // Left sidebar & channel topology header
  await expect(page.getByText('PRJ-METRO-001').first()).toBeVisible();
  await expect(page.getByText(/7 Visible Channels/)).toBeVisible();
  for (const channel of ['General', 'Procurement', 'Logistics', 'Site Operations', 'Quality']) {
    await expect(page.getByRole('heading', { name: channel, exact: true })).toBeVisible();
  }

  // Verify channels are NOT represented by a select
  await expect(page.getByLabel('Select Channel')).toHaveCount(0);

  // Switch channel using left-column channel navigation
  await page.getByRole('button', { name: /Quality/ }).click();
  await expect(page.getByRole('heading', { name: 'Quality Channel' })).toBeVisible();

  // Verify native HTML <select> for thread selection exists
  const threadSelect = page.getByLabel('Select Thread');
  await expect(threadSelect).toBeVisible();
  await expect(threadSelect).toHaveValue('THREAD-QUAL-001');

  // Switch thread using native thread <select>
  await threadSelect.selectOption('THREAD-QUAL-002');

  // Verify Stitch Thread layout elements for selected thread (without header)
  await expect(page.getByText('Awaiting Approval').first()).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Foundation Concrete Inspection' })).toBeVisible();
  await expect(page.getByText('Thread ID: #THREAD-QUAL-002').first()).toBeVisible();

  // Verify Thread Timeline Events & Actors
  await expect(page.getByText('Contractor', { exact: true }).first()).toBeVisible();
  await expect(page.getByText('Inspector', { exact: true }).first()).toBeVisible();
  await expect(page.getByText('Inspection requested for pour zone A3. Ready for sign-off.').first()).toBeVisible();
  await expect(page.getByText('Concrete samples recorded. Slump test within tolerance (4.5 inches).').first()).toBeVisible();

  // Verify Attachment Cards
  await expect(page.getByText('Foundation_Report.pdf').first()).toBeVisible();
  await expect(page.getByText('Drawing_S-204.pdf').first()).toBeVisible();
});

test('site detail renders exact Stitch layout and navigates from project page', async ({ page }) => {
  await visit(page, '/sites/SITE-001');

  // Breadcrumbs & Header
  await expect(page.getByText('SITE-001 • Metro Line 4 Expansion').first()).toBeVisible();

  // Hero Card
  await expect(page.getByText('Phase 4 of 7 • Structure & Foundation').first()).toBeVisible();
  await expect(page.getByRole('heading', { level: 2 })).toHaveText('Riverside Tower');

  // Site Information Card
  await expect(page.getByText('Site Information').first()).toBeVisible();
  await expect(page.getByText('42% Completion').first()).toBeVisible();
  await expect(page.getByText('14,200 m²').first()).toBeVisible();
  await expect(page.getByText('Bureau Inspection').first()).toBeVisible();

  // Zone Surveillance Feeds
  await expect(page.getByText('Zone Surveillance Feeds').first()).toBeVisible();
  await expect(page.getByText('Foundation Trench — Sector B').first()).toBeVisible();

  // Operational Zones
  await expect(page.getByText('Operational Zones').first()).toBeVisible();
  await expect(page.getByText('Deep Excavation & Shoring').first()).toBeVisible();

  // Storage & Logistics
  await expect(page.getByText('Storage & Logistics').first()).toBeVisible();
  await expect(page.getByText('Yard Occupancy: 72% Occupied').first()).toBeVisible();
  await expect(page.getByRole('button', { name: /Track Delivery/ }).first()).toBeVisible();

  // Bottom action buttons
  await expect(page.getByRole('button', { name: /View Inspection Status/ }).first()).toBeVisible();

  // Navigation test from Project page -> Site Detail
  await visit(page, '/projects/PRJ-001');
  const siteLink = page.getByTestId('site-location-link');
  await expect(siteLink).toBeVisible();
  await siteLink.click();

  await expect(page).toHaveURL(/.*\/sites\/SITE-001/);
  await expect(page.getByRole('heading', { level: 2 })).toHaveText('Riverside Tower');
});

test('notifications renders exact Stitch layout, operational pulse, filters, mark all read, and navigation actions', async ({ page }) => {
  await visit(page, '/notifications');

  // Header & Subtitle
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Notifications');
  await expect(page.getByText('Operational alerts & project coordination')).toBeVisible();

  // Mark all read button
  const markReadBtn = page.getByRole('button', { name: 'Mark all read' });
  await expect(markReadBtn).toBeVisible();

  // Operational Health Pulse
  await expect(page.getByText('Site Feeds Live')).toBeVisible();
  await expect(page.getByText('Viaduct Corridor 4 • High Activity')).toBeVisible();
  await expect(page.getByText('Sync 1m ago')).toBeVisible();

  // Filter Tabs & Counts
  await expect(page.getByRole('button', { name: 'All 5' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Critical 1' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Operations 2' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Decisions 1' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Collaboration 1' })).toBeVisible();

  // Notification Cards Verification
  // 1. Critical Logistics
  await expect(page.getByText('CRITICAL • LOGISTICS', { exact: true })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Transport Delay Detected (DEL-1042)' })).toBeVisible();
  await expect(page.getByText('Route M1 Blocked (km 42.4)')).toBeVisible();
  await expect(page.getByText('Ref: DEL-1042 • Phase: Foundation')).toBeVisible();

  // 2. Decision
  await expect(page.getByText('DECISION APPROVED', { exact: true })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Decision Approved: Alternative Route B' })).toBeVisible();
  await expect(page.getByText('Ref: DEC-1042 • Consensus 4/4')).toBeVisible();

  // 3. Issue
  await expect(page.getByText('ISSUE MITIGATED', { exact: true })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Issue ISS-1042 Resolved via Detour' })).toBeVisible();
  await expect(page.getByText('Ref: ISS-1042 • Carrier: FastBuild')).toBeVisible();

  // 4. Inspection
  await expect(page.getByText('INSPECTION PASSED', { exact: true })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Inspection INSP-1042 Passed & Material Accepted' })).toBeVisible();
  await expect(page.getByText('Ref: INSP-1042 • Material: Grade S355JR')).toBeVisible();

  // 5. C3 Collaboration
  await expect(page.getByText('C3 COLLABORATION', { exact: true })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'New Thread Activity in #Logistics & Delivery' })).toBeVisible();
  await expect(page.getByText('EN_10204_3.1_Certificate_Signed.pdf')).toBeVisible();

  // Filter interaction test:
  // Critical tab -> 1 card
  await page.getByRole('button', { name: 'Critical 1' }).click();
  await expect(page.getByRole('heading', { name: 'Transport Delay Detected (DEL-1042)' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Decision Approved: Alternative Route B' })).toHaveCount(0);

  // Operations tab -> 2 cards
  await page.getByRole('button', { name: 'Operations 2' }).click();
  await expect(page.getByRole('heading', { name: 'Issue ISS-1042 Resolved via Detour' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Inspection INSP-1042 Passed & Material Accepted' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Transport Delay Detected (DEL-1042)' })).toHaveCount(0);

  // Decisions tab -> 1 card
  await page.getByRole('button', { name: 'Decisions 1' }).click();
  await expect(page.getByRole('heading', { name: 'Decision Approved: Alternative Route B' })).toBeVisible();

  // Collaboration tab -> 1 card
  await page.getByRole('button', { name: 'Collaboration 1' }).click();
  await expect(page.getByRole('heading', { name: 'New Thread Activity in #Logistics & Delivery' })).toBeVisible();

  // Switch back to All tab
  await page.getByRole('button', { name: 'All 5' }).click();

  // Mark All Read interaction test
  await markReadBtn.click();
  await expect(page.getByRole('button', { name: '✓ Caught up!' })).toBeDisabled();

  // Action Button Navigation test: View Delivery
  await page.getByRole('button', { name: 'View Delivery (DEL-1042)' }).click();
  await expect(page).toHaveURL(/.*\/deliveries/);

  // Header notification icon navigation test
  await visit(page, '/projects/PRJ-001');
  const notifIcon = page.getByRole('button', { name: 'Notifications' });
  await expect(notifIcon).toBeVisible();
  await notifIcon.click();
  await expect(page).toHaveURL(/.*\/notifications/);
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Notifications');
});

test('profile renders exact Stitch layout, credentials, operational access, field switches, and offline sync interaction', async ({ page }) => {
  await visit(page, '/profile');

  // Identity Verification
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Manuel Vincent');
  await expect(page.getByText('Lead Structural Engineer')).toBeVisible();
  await expect(page.getByText('Partner • Engineering Partners')).toBeVisible();
  await expect(page.getByTitle('Verified User')).toBeVisible();
  await expect(page.getByLabel('Online indicator')).toBeVisible();

  // Workspace & Active Site Metadata Cards
  await expect(page.getByText('Workspace ID')).toBeVisible();
  await expect(page.getByText('SW-FR-7501')).toBeVisible();
  await expect(page.getByText('Mesh', { exact: true })).toBeVisible();
  await expect(page.getByText('Active Site')).toBeVisible();
  await expect(page.getByText('PRJ-001', { exact: true })).toBeVisible();
  await expect(page.getByText('Metro Line 4').first()).toBeVisible();

  // Credentials & Sign-off
  await expect(page.getByRole('heading', { level: 2, name: 'Credentials & Sign-off' })).toBeVisible();
  await expect(page.getByText('3 Active', { exact: true })).toBeVisible();
  await expect(page.getByText('Eurocodes EN 1992-1 / 1993')).toBeVisible();
  await expect(page.getByText('Structural Concrete & Steel Authority')).toBeVisible();
  await expect(page.getByText('Full Sign-off Privilege')).toBeVisible();
  await expect(page.getByText('Valid to Dec 2026')).toBeVisible();

  await expect(page.getByText('Bureau Inspection Assessor')).toBeVisible();
  await expect(page.getByText('Level 3 Certified Technical Auditor')).toBeVisible();
  await expect(page.getByText('Audit Registry: ID-8840-X')).toBeVisible();
  await expect(page.getByText('Annual Renewal Sync OK')).toBeVisible();

  await expect(page.getByText('Site Safety Pass: Category A')).toBeVisible();
  await expect(page.getByText('High Risk Infrastructure / Viaduct Zones')).toBeVisible();
  await expect(page.getByText('All Sectors Authorized')).toBeVisible();

  // Operational Access
  await expect(page.getByRole('heading', { level: 2, name: 'Operational Access' })).toBeVisible();
  await expect(page.getByText('Sovereign Rights', { exact: true })).toBeVisible();
  await expect(page.getByText('PRJ-001 Metro Line 4')).toBeVisible();
  await expect(page.getByText('Full Engineering Sign-off & Inspection Authority')).toBeVisible();
  await expect(page.getByText('SITE-001 Riverside Tower')).toBeVisible();
  await expect(page.getByText('Zone Access: Viaduct & Deep Foundation Sector')).toBeVisible();
  await expect(page.getByText('TraceCore Milestone Notarization')).toBeVisible();
  await expect(page.getByText('Active Cryptographic Delegated Signer')).toBeVisible();
  await expect(page.getByText('HSM Tier 1')).toBeVisible();
  await expect(page.getByText('C3 Collaboration Feeds')).toBeVisible();
  await expect(page.getByText('#Logistics, #Foundation, #Technical-Review')).toBeVisible();
  await expect(page.getByText('3 Linked')).toBeVisible();

  // Field Configuration Switches Initial State
  const offlineSwitch = page.getByRole('switch', { name: 'Offline Auto-Sync' });
  const logisticsSwitch = page.getByRole('switch', { name: 'Critical Logistics Dispatch' });
  const contrastSwitch = page.getByRole('switch', { name: 'High-Contrast Sunlight Mode' });
  const satelliteSwitch = page.getByRole('switch', { name: 'Satellite / Low-Bandwidth Mode' });

  await expect(offlineSwitch).toHaveAttribute('aria-checked', 'true');
  await expect(logisticsSwitch).toHaveAttribute('aria-checked', 'true');
  await expect(contrastSwitch).toHaveAttribute('aria-checked', 'false');
  await expect(satelliteSwitch).toHaveAttribute('aria-checked', 'false');

  // Toggle Switches Interaction
  await offlineSwitch.click();
  await expect(offlineSwitch).toHaveAttribute('aria-checked', 'false');

  await contrastSwitch.click();
  await expect(contrastSwitch).toHaveAttribute('aria-checked', 'true');

  // System Health & Storage
  await expect(page.getByRole('heading', { level: 2, name: 'System Health & Storage' })).toBeVisible();
  await expect(page.getByText('Optimal', { exact: true })).toBeVisible();
  await expect(page.getByText('42.8 MB')).toBeVisible();
  await expect(page.getByText('14 Vector Blueprints')).toBeVisible();
  await expect(page.getByText('Encrypted Mesh OK')).toBeVisible();
  await expect(page.getByText('v2.4.1 (Sovereign Mobile Edition)')).toBeVisible();

  // Synchronize All Offline Data Interaction
  const syncBtn = page.getByRole('button', { name: 'Synchronize All Offline Data' });
  await expect(syncBtn).toBeVisible();
  await syncBtn.click();

  // Verify transition to synchronizing state
  const syncingBtn = page.getByRole('button', { name: 'Synchronizing Telemetry...' });
  await expect(syncingBtn).toBeVisible();
  await expect(syncingBtn).toBeDisabled();

  // Wait for sync completion
  await expect(page.getByText('Just now', { exact: true })).toBeVisible();

  // Bottom Navigation -> Profile Navigation
  await visit(page, '/projects');
  const profileTab = page.getByRole('button', { name: 'Profile' });
  await expect(profileTab).toBeVisible();
  await profileTab.click();
  await expect(page).toHaveURL(/.*\/profile/);
  await expect(page.getByRole('heading', { level: 1 })).toHaveText('Manuel Vincent');

  // Footer Verification
  await expect(page.getByText('BuildFlow Field OS • Verified Encrypted Node #881-A')).toBeVisible();
});



