# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: construction-canonical-navigation.spec.ts >> material requirement renders exact Stitch layout, canonical material, supplier sourcing, and offer details
- Location: e2e/construction-canonical-navigation.spec.ts:129:1

# Error details

```
Test timeout of 45000ms exceeded.
```

```
Error: expect(locator).toBeVisible() failed

Locator: getByText('€142,500')
Expected: visible
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 10000ms
  - waiting for getByText('€142,500')

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
  - text: domain BuildFlow
  - button "Notifications": notifications
  - text: MV
  - main:
    - text: apartment PROJ-RT-104 • Metro Line 4 Expansion Active
    - heading "Material Requirement" [level=1]
    - text: OFFER SELECTED
    - paragraph: REQ-STRUCT-001 • Structural beams
    - paragraph: Structural beams required for viaduct section 4.
    - text: view_in_ar Material ID
    - paragraph: MAT-STRUCT-001
    - text: warning Critical Path precision_manufacturing Specification
    - paragraph: Grade S355JR
    - paragraph: Hot-Rolled (EN 10025-2)
    - text: scale Quantity
    - paragraph: 84 Beams
    - paragraph: 120 Metric Tons
    - text: calendar_today Required Date
    - paragraph: Aug 15, 2024
    - paragraph: "Slot: 08:00 - 12:00"
    - text: location_on Site Destination
    - paragraph: Site-001 (Paris)
    - paragraph: Riverside Tower
    - text: account_tree Construction Phase Foundation / Viaduct Section 4 check_circle hub
    - heading "Supplier Sourcing" [level=2]
    - text: "100% Filled 3 Invited 2 Offers Rec'd 1 Selected ES EuroSteel Construction verified Ref: OFF-1042 Chosen Total Contract Value €54,000.00 €1,187.50 / ton local_shipping Proposed Delivery:"
    - strong: Aug 15, 2024 at 10:00 AM
    - button "visibility View Offers (2)"
    - button "person_add Invite Supplier"
  - navigation:
    - button "home Home"
    - button "architecture Projects"
    - button "pending_actions Activity"
    - button "description Docs"
    - button "person_outline Profile"
- navigation:
  - button "home Home"
  - button "architecture Projects"
  - button "pending_actions Activity"
  - button "description Docs"
  - button "person_outline Profile"
```

# Test source

```ts
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
  129 | test('material requirement renders exact Stitch layout, canonical material, supplier sourcing, and offer details', async ({ page }) => {
  130 |   await visit(page, '/requirements');
  131 | 
  132 |   // Header & Context
  133 |   await expect(page.getByText('BuildFlow').first()).toBeVisible();
  134 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText('Material Requirement');
  135 |   await expect(page.getByText('PROJ-RT-104 • Metro Line 4 Expansion')).toBeVisible();
  136 |   await expect(page.getByText('OFFER SELECTED')).toBeVisible();
  137 | 
  138 |   // Requirement & Material Card
  139 |   await expect(page.getByText('REQ-STRUCT-001')).toBeVisible();
  140 |   await expect(page.getByText('MAT-STRUCT-001')).toBeVisible();
  141 |   await expect(page.getByText('Critical Path')).toBeVisible();
  142 |   await expect(page.getByText('Grade S355JR')).toBeVisible();
  143 |   await expect(page.getByText('84 Beams')).toBeVisible();
  144 |   await expect(page.getByText('120 Metric Tons')).toBeVisible();
  145 |   await expect(page.getByText('Aug 15, 2024').first()).toBeVisible();
  146 |   await expect(page.getByText('Slot: 08:00 - 12:00')).toBeVisible();
  147 |   await expect(page.getByText('Site-001 (Paris)')).toBeVisible();
  148 |   await expect(page.getByText('Foundation / Viaduct Section 4')).toBeVisible();
  149 | 
  150 |   // Supplier Sourcing Card
  151 |   await expect(page.getByRole('heading', { level: 2, name: 'Supplier Sourcing' })).toBeVisible();
  152 |   await expect(page.getByText('100% Filled')).toBeVisible();
  153 |   await expect(page.getByText('Invited')).toBeVisible();
  154 |   await expect(page.getByText("Offers Rec'd")).toBeVisible();
  155 |   await expect(page.getByText('Selected').first()).toBeVisible();
  156 | 
  157 |   // Selected Supplier Card
  158 |   await expect(page.getByText('EuroSteel Construction')).toBeVisible();
  159 |   await expect(page.getByText('OFF-1042')).toBeVisible();
  160 |   await expect(page.getByText('Chosen')).toBeVisible();
> 161 |   await expect(page.getByText('€142,500')).toBeVisible();
      |                                            ^ Error: expect(locator).toBeVisible() failed
  162 |   await expect(page.getByText('€1,187.50 / ton')).toBeVisible();
  163 |   await expect(page.getByText('Aug 15, 2024 at 10:00 AM')).toBeVisible();
  164 | 
  165 |   // Action Buttons
  166 |   const viewOffersBtn = page.getByRole('button', { name: 'View Offers (2)' });
  167 |   await expect(viewOffersBtn).toBeVisible();
  168 |   await expect(page.getByRole('button', { name: 'Invite Supplier' })).toBeVisible();
  169 | 
  170 |   // Navigation to Supplier Offers page
  171 |   await viewOffersBtn.click();
  172 |   await expect(page).toHaveURL(/.*\/offers/);
  173 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText('Supplier Offers');
  174 | });
  175 | 
  176 | test('supplier offers render exact Stitch layout, offer details, and suppliers', async ({ page }) => {
  177 |   await visit(page, '/offers');
  178 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText('Supplier Offers');
  179 |   await expect(page.getByText('120t Structural Beams', { exact: false }).first()).toBeVisible();
  180 |   await expect(page.getByText('EuroSteel Construction').first()).toBeVisible();
  181 |   await expect(page.getByText('Arcelor Infrastructure').first()).toBeVisible();
  182 |   await expect(page.getByText('ACCEPTED • DELIVERY PLANNED').first()).toBeVisible();
  183 |   await expect(page.getByText('CE Marked').first()).toBeVisible();
  184 |   await expect(page.getByText('EN 10204 3.1 Mill Cert').first()).toBeVisible();
  185 |   await expect(page.getByText(/OFF-REF-1042/).first()).toBeVisible();
  186 |   await expect(page.getByText(/OFF-1039/).first()).toBeVisible();
  187 | });
  188 | 
  189 | test('delivery detail renders exact Stitch layout and operational specifications', async ({ page }) => {
  190 |   await visit(page, '/deliveries');
  191 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText('Delivery #DEL-1042');
  192 |   await expect(page.getByText('DELAYED').first()).toBeVisible();
  193 |   await expect(page.getByText('Road Restriction on M1').first()).toBeVisible();
  194 |   await expect(page.getByText('Schedule Comparison').first()).toBeVisible();
  195 |   await expect(page.getByText('+21.5h Delay').first()).toBeVisible();
  196 |   await expect(page.getByText('Operational Specifications').first()).toBeVisible();
  197 |   await expect(page.getByText('REQ-REF-1042').first()).toBeVisible();
  198 |   await expect(page.getByText('MAT-STRUCT-001').first()).toBeVisible();
  199 |   await expect(page.getByText('EuroSteel').first()).toBeVisible();
  200 |   await expect(page.getByText('OFF-1042').first()).toBeVisible();
  201 |   await expect(page.getByText('FastBuild Logistics').first()).toBeVisible();
  202 |   await expect(page.getByText('ISS-1042 (Critical Delay)').first()).toBeVisible();
  203 |   await expect(page.getByText('Last Known Telemetry').first()).toBeVisible();
  204 |   await expect(page.getByText('GPS Active').first()).toBeVisible();
  205 |   await expect(page.getByText('Progress Timeline').first()).toBeVisible();
  206 |   await expect(page.getByText('Stage 4 of 8').first()).toBeVisible();
  207 |   await expect(page.getByText('Delayed on Route M1').first()).toBeVisible();
  208 |   await expect(page.getByText('Rerouted').first()).toBeVisible();
  209 |   await expect(page.getByText('Delivered').first()).toBeVisible();
  210 |   await expect(page.getByText('Inspected').first()).toBeVisible();
  211 |   await expect(page.getByText('Accepted').first()).toBeVisible();
  212 |   await expect(page.getByText('View Provenance & Root Cause Analysis').first()).toBeVisible();
  213 |   await expect(page.getByText('Print Consignment Manifest & Waybill').first()).toBeVisible();
  214 | });
  215 | 
  216 | test('transport renders exact Stitch Transport & Delay layout and details', async ({ page }) => {
  217 |   await visit(page, '/transport');
  218 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText('Transport & Delay');
  219 |   await expect(page.getByText('DEL-REF-1042').first()).toBeVisible();
  220 |   await expect(page.getByText('ACTIVE ALERT').first()).toBeVisible();
  221 |   await expect(page.getByText('Transport Delay Detected').first()).toBeVisible();
  222 |   await expect(page.getByText('Route M1 Blocked').first()).toBeVisible();
  223 |   await expect(page.getByText('+21h 30m').first()).toBeVisible();
  224 |   await expect(page.getByText('Live Corridor Geometry').first()).toBeVisible();
  225 |   await expect(page.getByText('Constraint Details').first()).toBeVisible();
  226 |   await expect(page.getByText('Detour Evaluation').first()).toBeVisible();
  227 |   await expect(page.getByText('ROUTE B APPROVED').first()).toBeVisible();
  228 |   await expect(page.getByRole('button', { name: /View Decision DEC-REF-1042/ })).toBeVisible();
  229 | });
  230 | 
  231 | test('issue renders exact Stitch layout and specifications', async ({ page }) => {
  232 |   await visit(page, '/issues');
  233 | 
  234 |   // Header & Status strip
  235 |   await expect(page.getByRole('heading', { level: 1 })).toHaveText('Construction Issue');
  236 |   await expect(page.getByText('SEVERITY: CRITICAL').first()).toBeVisible();
  237 |   await expect(page.getByText('RESOLVED').first()).toBeVisible();
  238 | 
  239 |   // Issue Title / Resolution Card
  240 |   await expect(page.getByText('LOGISTICS BOTTLENECK').first()).toBeVisible();
  241 |   await expect(page.getByText(/ISS-REF-1042/).first()).toBeVisible();
  242 |   await expect(page.getByText('Critical Beam Delivery Delayed by Route M1 Blockage').first()).toBeVisible();
  243 | 
  244 |   // Active Mitigation
  245 |   await expect(page.getByText('Active Mitigation').first()).toBeVisible();
  246 |   await expect(page.getByText('Rerouted via Decision DEC-1042').first()).toBeVisible();
  247 | 
  248 |   // Specifications
  249 |   await expect(page.getByText('Specifications').first()).toBeVisible();
  250 |   await expect(page.getByText('DEL-1042').first()).toBeVisible();
  251 |   await expect(page.getByText('Structural Beams (44t Prefabricated)').first()).toBeVisible();
  252 |   await expect(page.getByText('Foundation').first()).toBeVisible();
  253 |   await expect(page.getByText('Route M1 Alert').first()).toBeVisible();
  254 |   await expect(page.getByText('David K. (Logistics Mgr)').first()).toBeVisible();
  255 | 
  256 |   // Stakeholders
  257 |   await expect(page.getByText('Assigned Stakeholders').first()).toBeVisible();
  258 |   await expect(page.getByText('BuildCorp (Main)').first()).toBeVisible();
  259 |   await expect(page.getByText('FastBuild Logistics').first()).toBeVisible();
  260 |   await expect(page.getByText('Engineering Partners').first()).toBeVisible();
  261 | 
```