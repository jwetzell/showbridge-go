import { expect, test } from '@playwright/test';

test('add route', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('button', { name: 'Add Route' }).click();
  await page.getByRole('menuitem', { name: 'HTTP Server' }).click();
  const route = page.locator('#routes\\.0');
  await expect(route).toBeVisible();
  await expect(route).toContainText('Route 0');
});

test('remove route', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('button', { name: 'Add Route' }).click();
  const route = page.locator('#routes\\.0');
  await expect(route).toBeVisible();
  await route.getByRole('button', { name: 'Delete Route' }).click();
  await expect(route).not.toBeVisible();
});
