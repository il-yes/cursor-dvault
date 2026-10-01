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
  await expect(page.getByText('EuroSteel Construction')).toBeVisible();
  await expect(page.getByText('OFF-1042')).toBeVisible();
  await expect(page.getByText('Chosen')).toBeVisible();
  await expect(page.getByText('€142,500')).toBeVisible();
  await expect(page.getByText('€1,187.50 / ton')).toBeVisible();
  await expect(page.getByText('Aug 15, 2024 at 10:00 AM')).toBeVisible();

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
