import { test, expect } from '@playwright/test'
import { authenticate, SUPER_ADMIN } from './helpers/auth'
import { pickSingleSelect } from './helpers/select'
import { dismissSwal } from './helpers/swal'

test.beforeEach(async ({ page, request }) => {
  await authenticate(page, request, SUPER_ADMIN)
  await page.goto('/gateway/services')
})

test('rejects a base_path outside the protocol-reserved namespace', async ({ page }) => {
  await page.getByRole('button', { name: 'Tambah Service' }).click()
  await page
    .getByRole('textbox', { name: 'Nama Service', exact: true })
    .fill(`e2e-svc-${Date.now()}`)
  await page.getByRole('textbox', { name: 'Base URL' }).fill('http://localhost:9999')
  await page.getByRole('textbox', { name: 'Base Path' }).fill('/not-under-svc')
  await page.getByRole('button', { name: 'Simpan' }).click()

  await expect(
    page.getByText('Base path harus diawali /api/svc/ (protocol http) atau /ws/svc/'),
  ).toBeVisible()
  // Modal stays open — the request never went out.
  await expect(page.getByRole('heading', { name: 'Tambah Service' })).toBeVisible()
})

test('auto-swaps the base_path prefix when protocol changes', async ({ page }) => {
  await page.getByRole('button', { name: 'Tambah Service' }).click()
  await page.getByRole('textbox', { name: 'Base Path' }).fill('/api/svc/order')

  await pickSingleSelect(page, 'Protocol', 'WebSocket')
  await expect(page.getByRole('textbox', { name: 'Base Path' })).toHaveValue('/ws/svc/order')

  await pickSingleSelect(page, 'Protocol', 'HTTP')
  await expect(page.getByRole('textbox', { name: 'Base Path' })).toHaveValue('/api/svc/order')
})

test('creates, health-checks, edits, and deletes a service', async ({ page }) => {
  const name = `e2e-svc-${Date.now()}`
  const renamed = `${name}-renamed`

  await page.getByRole('button', { name: 'Tambah Service' }).click()
  await page.getByRole('textbox', { name: 'Nama Service', exact: true }).fill(name)
  await page.getByRole('textbox', { name: 'Base URL' }).fill('http://localhost:9999')
  await page.getByRole('textbox', { name: 'Base Path' }).fill(`/api/svc/e2e-${Date.now()}`)
  await page.getByRole('button', { name: 'Simpan' }).click()
  await dismissSwal(page)

  const card = page.locator('.group', { hasText: name }).first()
  await expect(card).toBeVisible()

  await card.getByRole('button', { name: 'Cek' }).click()
  await expect(card.getByText('Belum pernah dicek')).not.toBeVisible({ timeout: 15000 })

  await card.getByRole('button', { name: 'Edit' }).click()
  await page.getByRole('textbox', { name: 'Nama Service', exact: true }).fill(renamed)
  await page.getByRole('button', { name: 'Perbarui' }).click()
  await dismissSwal(page)

  const renamedCard = page.locator('.group', { hasText: renamed }).first()
  await expect(renamedCard).toBeVisible()

  await renamedCard.getByRole('button', { name: 'Hapus' }).click()
  await page.getByRole('button', { name: 'Yes' }).click()
  await dismissSwal(page)

  await expect(page.locator('.group', { hasText: renamed })).not.toBeVisible()
})
