# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: c3-slots-persistence.spec.ts >> C3 Channel Slots Real Integration Test (Unmocked) >> Real Backend Flow: Add Slot -> Save -> Backend Persist -> Page Reload -> Assert Persistence
- Location: e2e/c3-slots-persistence.spec.ts:30:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: locator('.tab:has-text("Slots")').or(locator('text=Slots')).first()
Expected: visible
Timeout: 15000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 15000ms
  - waiting for locator('.tab:has-text("Slots")').or(locator('text=Slots')).first()

```

```yaml
- region "Notifications (F8)":
  - list
- region "Notifications alt+T"
- heading "Welcome to Ankhora" [level=1]
- text: Step 1 of 6
- heading "How do you want to use Ankhora?" [level=2]
- paragraph: All plans include zero-knowledge encryption and Stellar verification.
- button "Personal / Professional Use Ankhora with an email, password, and traditional account."
- button "Fully Anonymous “Even your subscription is anonymous”. A Stellar keypair is generated for you."
- button "Continue with Current Account"
```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test';
  2   | 
  3   | test.describe('C3 Channel Slots Real Integration Test (Unmocked)', () => {
  4   |   const testSlotName = `Slot_4_QA_Review_${Date.now()}`;
  5   | 
  6   |   test.beforeEach(async ({ page }) => {
  7   |     // Only set session auth in localStorage so DashboardLayout doesn't redirect to login screen;
  8   |     // DO NOT mock ListChannels, GetChannel, UpdateChannel, or store slots data!
  9   |     await page.addInitScript(() => {
  10  |       localStorage.setItem('auth-storage', JSON.stringify({
  11  |         state: {
  12  |           user: { id: "usr_dev", email: "ousman11@mail.com" },
  13  |           jwtToken: "mock_jwt",
  14  |           isAuthenticated: true,
  15  |         },
  16  |         version: 0,
  17  |       }));
  18  | 
  19  |       localStorage.setItem('app-storage', JSON.stringify({
  20  |         state: {
  21  |           user: { id: "usr_dev", email: "ousman11@mail.com" },
  22  |           loggedIn: true,
  23  |           session: { vault_runtime_context: { SessionSecrets: {} } },
  24  |         },
  25  |         version: 0,
  26  |       }));
  27  |     });
  28  |   });
  29  | 
  30  |   test('Real Backend Flow: Add Slot -> Save -> Backend Persist -> Page Reload -> Assert Persistence', async ({ page }) => {
  31  |     console.log("STEP 1: Navigating to C3 Configuration page...");
  32  |     await page.goto('/dashboard/c3/config');
  33  |     await page.waitForLoadState('domcontentloaded');
  34  | 
  35  |     const slotsTab = page.locator('.tab:has-text("Slots")').or(page.locator('text=Slots')).first();
> 36  |     await expect(slotsTab).toBeVisible({ timeout: 15000 });
      |                            ^ Error: expect(locator).toBeVisible() failed
  37  |     await slotsTab.click();
  38  | 
  39  |     const slotRows = page.locator('.slot-area .slot-row');
  40  |     const initialCount = await slotRows.count();
  41  |     console.log(`[REAL_INTEGRATION] Real backend initial slots count: ${initialCount}`);
  42  | 
  43  |     // 2. Click "+ Add another slot"
  44  |     console.log("STEP 2: Clicking + Add another slot...");
  45  |     const addTrigger = page.locator('.add-trigger').or(page.locator('text=+ Add another slot'));
  46  |     await expect(addTrigger).toBeVisible();
  47  |     await addTrigger.click();
  48  | 
  49  |     // Fill form for Slot 4
  50  |     const nameInput = page.locator('input.slot-input').or(page.locator('input[placeholder="new_slot_name"]'));
  51  |     await expect(nameInput).toBeVisible();
  52  |     await nameInput.fill(testSlotName);
  53  | 
  54  |     const vaultSelect = page.locator('select.slot-select').first();
  55  |     await vaultSelect.selectOption('vault_finance');
  56  | 
  57  |     const addConfirmBtn = page.locator('button.btn-confirm').or(page.locator('text=+ Add'));
  58  |     await addConfirmBtn.click();
  59  | 
  60  |     // Expected: initialCount + 1 slots rendered in UI draft
  61  |     await expect(slotRows).toHaveCount(initialCount + 1);
  62  |     await expect(page.locator(`text=${testSlotName}`)).toBeVisible();
  63  |     console.log(`✓ Added slot verified: ${initialCount + 1} slots rendered in UI draft.`);
  64  | 
  65  |     // 3. Save Changes
  66  |     console.log("STEP 3: Saving Changes to backend...");
  67  |     const saveBtn = page.locator('button.btn-save');
  68  |     await expect(saveBtn).toBeVisible();
  69  |     await saveBtn.click();
  70  | 
  71  |     const saveFeedback = page.locator('.dirty-bar');
  72  |     await expect(saveFeedback).toContainText(/saved successfully/i, { timeout: 10000 });
  73  | 
  74  |     await expect(slotRows).toHaveCount(initialCount + 1);
  75  |     await expect(page.locator(`text=${testSlotName}`)).toBeVisible();
  76  |     console.log(`✓ Save complete: ${initialCount + 1} slots rendered post-save.`);
  77  | 
  78  |     // 4. Reload application page
  79  |     console.log("STEP 4: Reloading application page to test real backend reconstruction...");
  80  |     await page.reload();
  81  |     await page.waitForLoadState('domcontentloaded');
  82  | 
  83  |     if (await slotsTab.isVisible({ timeout: 5000 }).catch(() => false)) {
  84  |       await slotsTab.click();
  85  |     }
  86  | 
  87  |     // Assert that the real backend response after reload includes the newly saved slot
  88  |     await expect(slotRows).toHaveCount(initialCount + 1, { timeout: 10000 });
  89  |     await expect(page.locator(`text=${testSlotName}`)).toBeVisible();
  90  |     console.log(`✓ Real Backend Reload Verified: ${initialCount + 1} slots loaded directly from backend after reload.`);
  91  | 
  92  |     // 5. Delete slot & Save
  93  |     console.log("STEP 5: Cleaning up test slot...");
  94  |     const addedSlotRow = page.locator('.sg.slot-row', { hasText: testSlotName });
  95  |     const deleteBtn = addedSlotRow.locator('.icon-btn.del');
  96  |     await deleteBtn.click();
  97  | 
  98  |     await expect(slotRows).toHaveCount(initialCount);
  99  | 
  100 |     console.log("Saving deletion to backend...");
  101 |     await saveBtn.click();
  102 |     await expect(saveFeedback).toContainText(/saved successfully/i, { timeout: 10000 });
  103 | 
  104 |     // 6. Reload & final assert
  105 |     console.log("STEP 6: Reloading application after deletion...");
  106 |     await page.reload();
  107 |     await page.waitForLoadState('domcontentloaded');
  108 | 
  109 |     if (await slotsTab.isVisible({ timeout: 5000 }).catch(() => false)) {
  110 |       await slotsTab.click();
  111 |     }
  112 | 
  113 |     await expect(slotRows).toHaveCount(initialCount, { timeout: 10000 });
  114 |     await expect(page.locator(`text=${testSlotName}`)).toBeHidden();
  115 |     console.log("✓ Final reload verified: returned to original slot count directly from backend.");
  116 |   });
  117 | });
  118 | 
```