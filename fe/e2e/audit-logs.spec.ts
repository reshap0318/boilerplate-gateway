import { test, expect } from '@playwright/test'
import { authenticate, SUPER_ADMIN } from './helpers/auth'
import { dismissSwal } from './helpers/swal'

test.beforeEach(async ({ page, request }) => {
  await authenticate(page, request, SUPER_ADMIN)
})

test('renders a create-service action after creating a service', async ({ page }) => {
  const name = `e2e-audit-svc-${Date.now()}`

  await page.goto('/gateway/services')
  await page.getByRole('button', { name: 'Tambah Service' }).click()
  await page.getByRole('textbox', { name: 'Nama Service', exact: true }).fill(name)
  await page.getByRole('textbox', { name: 'Base URL' }).fill('http://localhost:9999')
  await page.getByRole('textbox', { name: 'Base Path' }).fill(`/api/svc/e2e-${Date.now()}`)
  await page.getByRole('button', { name: 'Simpan' }).click()
  await dismissSwal(page)
  await expect(page.locator('.group', { hasText: name })).toBeVisible()

  await page.goto('/audit-logs')
  // Other specs (gateway-routes.spec.ts) run concurrently and also log CREATE rows — filter to
  // SERVICE specifically rather than assuming ours is the most recent one.
  const row = page.locator('tr', { hasText: 'CREATE' }).filter({ hasText: 'SERVICE' }).first()
  await expect(row).toBeVisible()

  // Row expands to show the JSON diff on click.
  await row.click()
  await expect(page.locator('pre')).toContainText('"after"')
})
