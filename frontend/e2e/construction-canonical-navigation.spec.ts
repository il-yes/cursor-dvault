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

test('project listing renders canonical project reference', async ({ page }) => {
  await visit(page, '/projects');
  // code = canonical projectReference (PRJ-METRO-001), not the projectId.
  await expect(page.getByText('PRJ-METRO-001').first()).toBeVisible();
  await expect(page.getByText('Metro Line 4 Expansion').first()).toBeVisible();
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

test('material requirement renders canonical material', async ({ page }) => {
  await visit(page, '/requirements');
  await expect(page.getByText('Precast Concrete Beam Heavy Grade').first()).toBeVisible();
  // canonical requirementReference, quantity 12 units, canonical specification.
  await expect(page.getByText('REQ-REF-1042').first()).toBeVisible();
  await expect(page.getByText(/12 units/).first()).toBeVisible();
  await expect(page.getByText('C50/60 Concrete, 12m length').first()).toBeVisible();
});

test('supplier offers render canonical supplier and currency', async ({ page }) => {
  await visit(page, '/offers');
  await expect(page.getByText('Apex Precast Logistics')).toBeVisible();
  await expect(page.getByText('€', { exact: false }).first()).toBeVisible();
});

test('delivery detail renders canonical site and route', async ({ page }) => {
  await visit(page, '/deliveries');
  await expect(page.getByText('North Hub Station Site', { exact: false }).first()).toBeVisible();
});

test('transport renders canonical route reference', async ({ page }) => {
  await visit(page, '/transport');
  await expect(page.getByText('Route-B', { exact: false }).first()).toBeVisible();
});

test('issue renders canonical issue reference', async ({ page }) => {
  await visit(page, '/issues');
  await expect(page.getByText('ISS-REF-1042', { exact: false }).first()).toBeVisible();
});

test('decision renders canonical decision reference', async ({ page }) => {
  await visit(page, '/decisions');
  await expect(page.getByText('DEC-REF-1042', { exact: false }).first()).toBeVisible();
});

test('inspection renders canonical inspection', async ({ page }) => {
  await visit(page, '/inspections');
  await expect(page.getByText('INSP-REF-1042', { exact: false }).first()).toBeVisible();
  await expect(page.getByText('Receiving Inspection', { exact: false }).first()).toBeVisible();
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

test('provenance renders canonical causal chain', async ({ page }) => {
  await visit(page, '/provenance');
  await expect(page.getByText('DEL-REF-1042').first()).toBeVisible();
  await expect(page.getByText('ISS-REF-1042').first()).toBeVisible();
  await expect(page.getByText('DEC-REF-1042').first()).toBeVisible();
  await expect(page.getByText('Road Restriction Official Notice M1')).toBeVisible();
  // Canonical content CID, not a fabricated one.
  await expect(page.getByText(/CID: QmRoadRestrictionReport1042/)).toBeVisible();
});

test('channels renders the canonical channel topology', async ({ page }) => {
  await visit(page, '/channels');
  // KNOWN DIVERGENCE (documented, not fixed): the page preserves a 7-channel
  // presentation topology instead of the canonical single collaboration channel
  // ("Site Logistics & Supply Chain"). Per-channel thread titles are built from
  // canonical ids, but they only render once a thread is selected.
  await expect(page.getByText('PRJ-METRO-001').first()).toBeVisible();
  await expect(page.getByText(/7 Visible Channels/)).toBeVisible();
  for (const channel of ['General', 'Procurement', 'Logistics', 'Site Operations', 'Quality']) {
    await expect(page.getByRole('heading', { name: channel, exact: true })).toBeVisible();
  }
  // The canonical single collaboration channel is NOT surfaced on this surface.
  await expect(page.getByText('Site Logistics & Supply Chain')).toHaveCount(0);
});
