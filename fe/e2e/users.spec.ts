import { test, expect } from '@playwright/test'
import { authenticate, SUPER_ADMIN } from './helpers/auth'
import { pickTagSelect } from './helpers/select'
import { dismissSwal } from './helpers/swal'

test.beforeEach(async ({ page, request }) => {
  await authenticate(page, request, SUPER_ADMIN)
  await page.goto('/users')
})

test('creates a user, suspends, activates, edits, then deletes it', async ({ page }) => {
  const email = `e2e.user.${Date.now()}@example.com`
  const name = 'E2E Test User'
  const renamed = 'E2E Test User Renamed'

  await page.getByRole('button', { name: 'Tambah User' }).click()
  await page.getByRole('textbox', { name: 'Nama' }).fill(name)
  await page.getByRole('textbox', { name: 'Email' }).fill(email)
  await page.getByRole('textbox', { name: 'Password', exact: true }).fill('Password123!')
  await page.getByRole('textbox', { name: 'Konfirmasi Password' }).fill('Password123!')
  await pickTagSelect(page, 'Roles', 'Viewer')
  await page.getByRole('button', { name: 'Simpan' }).click()
  await dismissSwal(page)

  const card = page.locator('.group', { hasText: name }).first()
  await expect(card).toBeVisible()
  await expect(card.getByText(email)).toBeVisible()

  await card.getByRole('button', { name: 'Suspend' }).click()
  await expect(card.getByRole('button', { name: 'Aktifkan' })).toBeVisible()
  await card.getByRole('button', { name: 'Aktifkan' }).click()
  await expect(card.getByRole('button', { name: 'Suspend' })).toBeVisible()

  await card.getByRole('button', { name: 'Edit' }).click()
  await page.getByRole('textbox', { name: 'Nama' }).fill(renamed)
  await page.getByRole('button', { name: 'Perbarui' }).click()
  await dismissSwal(page)

  const renamedCard = page.locator('.group', { hasText: renamed }).first()
  await expect(renamedCard).toBeVisible()

  await renamedCard.getByRole('button', { name: 'Hapus' }).click()
  await page.getByRole('button', { name: 'Yes' }).click()
  await dismissSwal(page)

  await expect(page.locator('.group', { hasText: renamed })).not.toBeVisible()
})
