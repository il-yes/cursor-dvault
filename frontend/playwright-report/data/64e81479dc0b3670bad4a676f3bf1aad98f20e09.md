# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: construction-canonical-navigation.spec.ts >> material requirement renders canonical material
- Location: e2e/construction-canonical-navigation.spec.ts:129:1

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText('Precast Concrete Beam Heavy Grade').first()
Expected: visible
Timeout: 10000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 10000ms
  - waiting for getByText('Precast Concrete Beam Heavy Grade').first()

```

```yaml
- region "Notifications (F8)":
  - list
- region "Notifications alt+T"
- banner:
  - text: BuildFlow Sovereign Workspace manage_accounts Active Persona
  - combobox "Frontend Persona Switcher (controls UI visibility & available actions)":
    - option "Project Manager (PM)" [selected]
    - option "Architect (ARCHITECT)"
    - option "Engineer (ENGINEER)"
    - option "Supplier (SUPPLIER)"
    - option "Logistics Operator (LOGISTICS)"
    - option "Site Superintendent (SITE_SUPERINTENDENT)"
    - option "QA / Inspector (QA)"
  - button "Notifications": notifications
  - text: person
- main:
  - button "PRJ-METRO-001"
  - text: / Requirements / REQ-REF-1042 REQ-REF-1042 • Phase 2 Structural open
  - heading "120t Structural Beams" [level=1]
  - paragraph: C50/60 Concrete, 12m length
  - text: Required Quantity 12 units Target Date 2026-08-15 Material ID MAT-STRUCT-001 Accepted Supplier EuroSteel Construction
  - heading "Supplier Offer Secured" [level=3]
  - paragraph: Offer OFF-REF-1042 accepted for €54,000.00. View delivery tracking.
  - button "View Offer (OFF-REF-1042)"
  - button "Track Delivery →"
- navigation:
  - button "home Home"
  - button "architecture Projects"
  - button "pending_actions Activity"
  - button "description Docs"
  - button "person_outline Profile"
```

# Test source

```ts
  31  |   window.go = {
  32  |     main: { App },
  33  |     Stellar: {
  34  |       CheckWalletStatus: async () => 'connected',
  35  |     },
  36  |     IPFS: { CheckNodeStatus: async () => 'ready' },
  37  |   };
  38  |   window.runtime = window.runtime || {
  39  |     EventsOn: () => {},
  40  |     EventsOff: () => {},
  41  |     EventsEmit: () => {},
  42  |     WindowReloadApp: () => {},
  43  |   };
  44  | })();
  45  | `;
  46  | 
  47  | const BASE = '/dashboard/construction';
  48  | 
  49  | /** Console/page errors that are genuine defects in the construction rebase. */
  50  | function isRealConstructionError(text: string): boolean {
  51  |   if (/SCENARIO_DATA/i.test(text)) return true;
  52  |   if (/scenarioAccessors|scenarioMappers|constructionScenarioAdapter/.test(text) && /error|cannot read|undefined/i.test(text)) return true;
  53  |   if (/Cannot read properties of undefined/i.test(text)) return true;
  54  |   return false;
  55  | }
  56  | 
  57  | test.beforeEach(async ({ page }) => {
  58  |   await page.addInitScript(BOOTSTRAP);
  59  | });
  60  | 
  61  | /**
  62  |  * Navigate to a construction route and assert:
  63  |  *  1. the route mounts (header present)
  64  |  *  2. no uncaught page error
  65  |  *  3. no construction-related console error
  66  |  */
  67  | async function visit(page: import('@playwright/test').Page, path: string) {
  68  |   const pageErrors: string[] = [];
  69  |   const consoleErrors: string[] = [];
  70  | 
  71  |   const onPageError = (e: Error) => pageErrors.push(e.message);
  72  |   const onConsole = (m: ConsoleMessage) => {
  73  |     if (m.type() === 'error' && isRealConstructionError(m.text())) consoleErrors.push(m.text());
  74  |   };
  75  |   page.on('pageerror', onPageError);
  76  |   page.on('console', onConsole);
  77  | 
  78  |   await page.goto(`${BASE}${path}`, { waitUntil: 'networkidle' });
  79  |   await expect(page.getByText('BuildFlow', { exact: false }).first()).toBeVisible();
  80  | 
  81  |   page.off('pageerror', onPageError);
  82  |   page.off('console', onConsole);
  83  | 
  84  |   expect(pageErrors, `page errors on ${path}`).toEqual([]);
  85  |   expect(consoleErrors, `console errors on ${path}`).toEqual([]);
  86  | }
  87  | 
  88  | test('dashboard renders canonical project, no crash', async ({ page }) => {
  89  |   await visit(page, '/');
  90  |   // Canonical metrics, not hardcoded ones: the project is active under its
  91  |   // canonical reference, the issue is RESOLVED (so 0 open) and the decision is
  92  |   // APPROVED (so 0 pending). Cursor-ordered, not "Today".
  93  |   await expect(page.getByText(/1 Active \(PRJ-METRO-001\)/)).toBeVisible();
  94  |   await expect(page.getByText(/No open issues \(ISS-REF-1042 resolved\)/)).toBeVisible();
  95  |   await expect(page.getByText(/No approvals pending \(DEC-REF-1042 approved\)/)).toBeVisible();
  96  |   await expect(page.getByText('Cursor-Ordered')).toBeVisible();
  97  |   await expect(page.getByText('Today')).toHaveCount(0);
  98  |   // 23 canonical thread events.
  99  |   await expect(page.getByText('Recent Events')).toBeVisible();
  100 |   await expect(page.getByText('23').first()).toBeVisible();
  101 | });
  102 | 
  103 | test('project listing renders canonical project reference and 3 total projects', async ({ page }) => {
  104 |   await visit(page, '/projects');
  105 |   // code = canonical projectReference (PRJ-METRO-001), not the projectId.
  106 |   await expect(page.getByText('PRJ-METRO-001').first()).toBeVisible();
  107 |   await expect(page.getByText('Metro Line 4 Expansion').first()).toBeVisible();
  108 |   // PRJ-002 and PRJ-003 restored demo projects
  109 |   await expect(page.getByText('PRJ-002').first()).toBeVisible();
  110 |   await expect(page.getByText('Commercial Plaza North').first()).toBeVisible();
  111 |   await expect(page.getByText('PRJ-003').first()).toBeVisible();
  112 |   await expect(page.getByText('Riverside Logistics Hub').first()).toBeVisible();
  113 | });
  114 | 
  115 | test('project detail renders canonical project', async ({ page }) => {
  116 |   await visit(page, '/projects/PRJ-001');
  117 |   await expect(page.getByText('Metro Line 4 Expansion').first()).toBeVisible();
  118 | });
  119 | 
  120 | test('stakeholders renders the three canonical vault identities', async ({ page }) => {
  121 |   await visit(page, '/stakeholders');
  122 |   // Canonical participants are vault identities, not people. There is no
  123 |   // ConstructionStakeholder aggregate, so names are absent by design.
  124 |   await expect(page.getByText('vault_001-oem', { exact: false }).first()).toBeVisible();
  125 |   await expect(page.getByText('vault_002-michelin', { exact: false }).first()).toBeVisible();
  126 |   await expect(page.getByText('vault_003-faa', { exact: false }).first()).toBeVisible();
  127 | });
  128 | 
  129 | test('material requirement renders canonical material', async ({ page }) => {
  130 |   await visit(page, '/requirements');
> 131 |   await expect(page.getByText('Precast Concrete Beam Heavy Grade').first()).toBeVisible();
      |                                                                             ^ Error: expect(locator).toBeVisible() failed
  132 |   // canonical requirementReference, quantity 12 units, canonical specification.
  133 |   await expect(page.getByText('REQ-REF-1042').first()).toBeVisible();
  134 |   await expect(page.getByText(/12 units/).first()).toBeVisible();
  135 |   await expect(page.getByText('C50/60 Concrete, 12m length').first()).toBeVisible();
  136 | });
  137 | 
  138 | test('supplier offers render exact Stitch layout, offer details, and suppliers', async ({ page }) => {
  139 |   await visit(page, '/offers');
  140 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText('Supplier Offers');
  141 |   await expect(page.getByText('REQ-STRUCT-001 • 120t Structural Beams')).toBeVisible();
  142 |   await expect(page.getByText('EuroSteel Construction').first()).toBeVisible();
  143 |   await expect(page.getByText('Arcelor Infrastructure').first()).toBeVisible();
  144 |   await expect(page.getByText('ACCEPTED • DELIVERY PLANNED').first()).toBeVisible();
  145 |   await expect(page.getByText('CE Marked').first()).toBeVisible();
  146 |   await expect(page.getByText('EN 10204 3.1 Mill Cert').first()).toBeVisible();
  147 |   await expect(page.getByText('OFF-1042').first()).toBeVisible();
  148 |   await expect(page.getByText('OFF-1039').first()).toBeVisible();
  149 | });
  150 | 
  151 | test('delivery detail renders canonical site and route', async ({ page }) => {
  152 |   await visit(page, '/deliveries');
  153 |   await expect(page.getByText('North Hub Station Site', { exact: false }).first()).toBeVisible();
  154 | });
  155 | 
  156 | test('transport renders canonical route reference', async ({ page }) => {
  157 |   await visit(page, '/transport');
  158 |   await expect(page.getByText('Route-B', { exact: false }).first()).toBeVisible();
  159 | });
  160 | 
  161 | test('issue renders canonical issue reference', async ({ page }) => {
  162 |   await visit(page, '/issues');
  163 |   await expect(page.getByText('ISS-REF-1042', { exact: false }).first()).toBeVisible();
  164 | });
  165 | 
  166 | test('decision renders the canonical decision reference', async ({ page }) => {
  167 |   await visit(page, '/decisions');
  168 |   await expect(page.getByText('DEC-REF-1042', { exact: false }).first()).toBeVisible();
  169 | });
  170 | 
  171 | /**
  172 |  * Stitch DEC-1042 composition, projected from the canonical aggregate.
  173 |  *
  174 |  * This asserts two things at once: that each composition section is present, and
  175 |  * that the Stitch copy the scenario cannot support was replaced rather than
  176 |  * rendered. The negative assertions are the important half.
  177 |  */
  178 | test('decision renders Stitch composition from canonical data only', async ({ page }) => {
  179 |   await visit(page, '/decisions');
  180 | 
  181 |   // Status ribbon — canonical status, canonical id, canonical decision date.
  182 |   await expect(page.getByText('APPROVED').first()).toBeVisible();
  183 |   await expect(page.getByText('DEC-1042').first()).toBeVisible();
  184 |   await expect(page.getByText(/Finalized Aug 15/)).toBeVisible();
  185 | 
  186 |   // Header card — canonical subject, canonical context, derived related ids.
  187 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText(
  188 |     'Structural material delivery delay alternative route',
  189 |   );
  190 |   await expect(page.getByText('Route M1 blocked, ETA delayed to 2026-08-16T07:30:00Z')).toBeVisible();
  191 |   await expect(page.getByText('ISS-1042').first()).toBeVisible();
  192 |   await expect(page.getByText('DEL-1042').first()).toBeVisible();
  193 | 
  194 |   // Options — both canonical options, selection recovered from the decision text.
  195 |   await expect(page.getByRole('heading', { name: 'Options Evaluated' })).toBeVisible();
  196 |   await expect(page.getByText('2 Alternatives Assessed')).toBeVisible();
  197 |   await expect(page.getByRole('heading', { name: 'Wait for M1 clearance' })).toBeVisible();
  198 |   await expect(page.getByRole('heading', { name: 'Use Route B detour' })).toBeVisible();
  199 |   await expect(page.getByText('Selected & Authorized')).toBeVisible();
  200 |   await expect(page.getByText('Rejected').first()).toBeVisible();
  201 | 
  202 |   // The only canonical arrival: derived from transport, labelled UTC.
  203 |   await expect(page.getByText('Aug 16 @ 07:30 UTC').first()).toBeVisible();
  204 |   // The 21.5h slip is derivable from planned vs actual arrival.
  205 |   await expect(page.getByText(/Route-B detour · \+21\.5 h vs plan/)).toBeVisible();
  206 | 
  207 |   // Governance — canonical roles, no fabricated organisations.
  208 |   await expect(page.getByRole('heading', { name: 'Governance Chain' })).toBeVisible();
  209 |   for (const role of ['Contractor', 'Supplier', 'Logistics', 'Project Manager']) {
  210 |     await expect(page.getByText(role, { exact: true }).first()).toBeVisible();
  211 |   }
  212 |   await expect(page.getByText('Decided by Project Manager')).toBeVisible();
  213 |   await expect(page.getByText('Requested by Logistics Manager · 4 participants consulted')).toBeVisible();
  214 | 
  215 |   // Outcome — canonical decision and consequence, canonical transport reference.
  216 |   await expect(page.getByRole('heading', { name: 'Approve Route B alternative transport' })).toBeVisible();
  217 |   await expect(page.getByText(/Transport re-routed via Route B/)).toBeVisible();
  218 |   await expect(page.getByText('TR-1042').first()).toBeVisible();
  219 | 
  220 |   // Action bar routes on the real route constants.
  221 |   await expect(page.getByRole('button', { name: /View Updated Delivery \(DEL-1042\)/ })).toBeVisible();
  222 |   await expect(page.getByRole('button', { name: /View Audit Trail & Provenance/ })).toBeVisible();
  223 | 
  224 |   // ---- Negative assertions: Stitch copy the scenario cannot support ----
  225 |   // Clock time on a date-only decisionDate.
  226 |   await expect(page.getByText(/12:35/)).toHaveCount(0);
  227 |   await expect(page.getByText(/no time recorded/)).toBeVisible();
  228 |   // Commercial facts with no canonical field.
  229 |   await expect(page.getByText(/480/)).toHaveCount(0);
  230 |   await expect(page.getByText(/Logistics SLA/)).toHaveCount(0);
  231 |   // Distance/duration deltas, and the "48-72 hour" option impact.
```