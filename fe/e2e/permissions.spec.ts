import { test, expect } from '@playwright/test'
import { authenticate, SUPER_ADMIN } from './helpers/auth'
import { dismissSwal } from './helpers/swal'
import { goToLastPage } from './helpers/pagination'

test.beforeEach(async ({ page, request }) => {
  await authenticate(page, request, SUPER_ADMIN)
  await page.goto('/uam/permissions')
})

test('creates, edits, and deletes a permission', async ({ page }) => {
  const name = `e2e.perm.${Date.now()}`
  const renamed = `${name}.renamed`

  await page.getByRole('button', { name: 'Tambah Permission' }).click()
  await page.getByRole('textbox', { name: 'Nama Permission' }).fill(name)
  await page.getByRole('textbox', { name: 'Deskripsi (opsional)' }).fill('Created by e2e test')
  await page.getByRole('button', { name: 'Simpan' }).click()
  await dismissSwal(page)
  await goToLastPage(page)

  const card = page.locator('.group', { hasText: name }).first()
  await expect(card).toBeVisible()

  await card.getByTitle('Edit').click()
  await page.getByRole('textbox', { name: 'Nama Permission' }).fill(renamed)
  await page.getByRole('button', { name: 'Perbarui' }).click()
  await dismissSwal(page)

  const renamedCard = page.locator('.group', { hasText: renamed }).first()
  await expect(renamedCard).toBeVisible()

  await renamedCard.getByTitle('Hapus').click()
  await page.getByRole('button', { name: 'Yes' }).click()
  await dismissSwal(page)

  await expect(page.locator('.group', { hasText: renamed })).not.toBeVisible()
})
