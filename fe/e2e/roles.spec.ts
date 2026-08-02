import { test, expect } from '@playwright/test'
import { authenticate, SUPER_ADMIN } from './helpers/auth'
import { dismissSwal } from './helpers/swal'

test.beforeEach(async ({ page, request }) => {
  await authenticate(page, request, SUPER_ADMIN)
  await page.goto('/uam/roles')
})

test('creates, edits, and deletes a role', async ({ page }) => {
  const name = `e2e-role-${Date.now()}`
  const renamed = `${name}-renamed`

  await page.getByRole('button', { name: 'Tambah Role' }).click()
  await page.getByRole('textbox', { name: 'Nama Role' }).fill(name)
  await page.getByRole('textbox', { name: 'Deskripsi (opsional)' }).fill('Created by e2e test')
  // Grant at least one permission so the role isn't empty.
  await page
    .locator('label', { hasText: 'user.index' })
    .first()
    .locator('input[type="checkbox"]')
    .check()
  await page.getByRole('button', { name: 'Simpan' }).click()
  await dismissSwal(page)

  const card = page.locator('.group', { hasText: name }).first()
  await expect(card).toBeVisible()
  await expect(card.getByText('user.index')).toBeVisible()

  await card.getByTitle('Edit').click()
  await page.getByRole('textbox', { name: 'Nama Role' }).fill(renamed)
  await page.getByRole('button', { name: 'Perbarui' }).click()
  await dismissSwal(page)

  const renamedCard = page.locator('.group', { hasText: renamed }).first()
  await expect(renamedCard).toBeVisible()

  await renamedCard.getByTitle('Hapus').click()
  await page.getByRole('button', { name: 'Yes' }).click()
  await dismissSwal(page)

  await expect(page.locator('.group', { hasText: renamed })).not.toBeVisible()
})
