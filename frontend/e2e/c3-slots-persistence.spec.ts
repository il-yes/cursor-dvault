import { test, expect } from '@playwright/test';

test.describe('C3 Channel Slots Real Integration Test (Unmocked)', () => {
  const testSlotName = `Slot_4_QA_Review_${Date.now()}`;

  test.beforeEach(async ({ page }) => {
    // Only set session auth in localStorage so DashboardLayout doesn't redirect to login screen;
    // DO NOT mock ListChannels, GetChannel, UpdateChannel, or store slots data!
    await page.addInitScript(() => {
      localStorage.setItem('auth-storage', JSON.stringify({
        state: {
          user: { id: "usr_dev", email: "ousman11@mail.com" },
          jwtToken: "mock_jwt",
          isAuthenticated: true,
        },
        version: 0,
      }));

      localStorage.setItem('app-storage', JSON.stringify({
        state: {
          user: { id: "usr_dev", email: "ousman11@mail.com" },
          loggedIn: true,
          session: { vault_runtime_context: { SessionSecrets: {} } },
        },
        version: 0,
      }));
    });
  });

  test('Real Backend Flow: Add Slot -> Save -> Backend Persist -> Page Reload -> Assert Persistence', async ({ page }) => {
    console.log("STEP 1: Navigating to C3 Configuration page...");
    await page.goto('/dashboard/c3/config');
    await page.waitForLoadState('domcontentloaded');

    const slotsTab = page.locator('.tab:has-text("Slots")').or(page.locator('text=Slots')).first();
    await expect(slotsTab).toBeVisible({ timeout: 15000 });
    await slotsTab.click();

    const slotRows = page.locator('.slot-area .slot-row');
    const initialCount = await slotRows.count();
    console.log(`[REAL_INTEGRATION] Real backend initial slots count: ${initialCount}`);

    // 2. Click "+ Add another slot"
    console.log("STEP 2: Clicking + Add another slot...");
    const addTrigger = page.locator('.add-trigger').or(page.locator('text=+ Add another slot'));
    await expect(addTrigger).toBeVisible();
    await addTrigger.click();

    // Fill form for Slot 4
    const nameInput = page.locator('input.slot-input').or(page.locator('input[placeholder="new_slot_name"]'));
    await expect(nameInput).toBeVisible();
    await nameInput.fill(testSlotName);

    const vaultSelect = page.locator('select.slot-select').first();
    await vaultSelect.selectOption('vault_finance');

    const addConfirmBtn = page.locator('button.btn-confirm').or(page.locator('text=+ Add'));
    await addConfirmBtn.click();

    // Expected: initialCount + 1 slots rendered in UI draft
    await expect(slotRows).toHaveCount(initialCount + 1);
    await expect(page.locator(`text=${testSlotName}`)).toBeVisible();
    console.log(`✓ Added slot verified: ${initialCount + 1} slots rendered in UI draft.`);

    // 3. Save Changes
    console.log("STEP 3: Saving Changes to backend...");
    const saveBtn = page.locator('button.btn-save');
    await expect(saveBtn).toBeVisible();
    await saveBtn.click();

    const saveFeedback = page.locator('.dirty-bar');
    await expect(saveFeedback).toContainText(/saved successfully/i, { timeout: 10000 });

    await expect(slotRows).toHaveCount(initialCount + 1);
    await expect(page.locator(`text=${testSlotName}`)).toBeVisible();
    console.log(`✓ Save complete: ${initialCount + 1} slots rendered post-save.`);

    // 4. Reload application page
    console.log("STEP 4: Reloading application page to test real backend reconstruction...");
    await page.reload();
    await page.waitForLoadState('domcontentloaded');

    if (await slotsTab.isVisible({ timeout: 5000 }).catch(() => false)) {
      await slotsTab.click();
    }

    // Assert that the real backend response after reload includes the newly saved slot
    await expect(slotRows).toHaveCount(initialCount + 1, { timeout: 10000 });
    await expect(page.locator(`text=${testSlotName}`)).toBeVisible();
    console.log(`✓ Real Backend Reload Verified: ${initialCount + 1} slots loaded directly from backend after reload.`);

    // 5. Delete slot & Save
    console.log("STEP 5: Cleaning up test slot...");
    const addedSlotRow = page.locator('.sg.slot-row', { hasText: testSlotName });
    const deleteBtn = addedSlotRow.locator('.icon-btn.del');
    await deleteBtn.click();

    await expect(slotRows).toHaveCount(initialCount);

    console.log("Saving deletion to backend...");
    await saveBtn.click();
    await expect(saveFeedback).toContainText(/saved successfully/i, { timeout: 10000 });

    // 6. Reload & final assert
    console.log("STEP 6: Reloading application after deletion...");
    await page.reload();
    await page.waitForLoadState('domcontentloaded');

    if (await slotsTab.isVisible({ timeout: 5000 }).catch(() => false)) {
      await slotsTab.click();
    }

    await expect(slotRows).toHaveCount(initialCount, { timeout: 10000 });
    await expect(page.locator(`text=${testSlotName}`)).toBeHidden();
    console.log("✓ Final reload verified: returned to original slot count directly from backend.");
  });
});
