import { expect, test } from '@playwright/test';

test('add module', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('button', { name: 'Add Module' }).click();
  await page.getByRole('menuitem', { name: 'HTTP Server' }).click();
  const module = page.locator('#modules\\.0');
  await expect(module).toBeVisible();
  await expect(module).toContainText('HTTP Server');
});

test('remove module', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('button', { name: 'Add Module' }).click();
  await page.getByRole('menuitem', { name: 'HTTP Server' }).click();
  const module = page.locator('#modules\\.0');
  await expect(module).toBeVisible();
  await module.getByRole('button', { name: 'Delete Module' }).click();
  await expect(module).not.toBeVisible();
});
