import { test, expect } from '@playwright/test'
import { authenticate, SUPER_ADMIN } from './helpers/auth'
import { pickSingleSelect } from './helpers/select'
import { dismissSwal } from './helpers/swal'
import { goToLastPage } from './helpers/pagination'

test.beforeEach(async ({ page, request }) => {
  await authenticate(page, request, SUPER_ADMIN)
  await page.goto('/gateway/routes')
})

test('creates, edits, and deletes a route', async ({ page }) => {
  const pathPattern = `/e2e-route-${Date.now()}`
  const renamedPathPattern = `${pathPattern}-renamed`

  await page.getByRole('button', { name: 'Tambah Route' }).click()
  await pickSingleSelect(page, 'Service', 'serv-uam')
  // Method already defaults to "GET" (initialForm) — re-clicking it would toggle it OFF
  // (single-mode Multiselect deselects on a repeat click of the same option).
  await page.getByRole('textbox', { name: 'Path Pattern' }).fill(pathPattern)
  await page.getByRole('button', { name: 'Simpan' }).click()
  await dismissSwal(page)
  await goToLastPage(page)

  const row = page.locator('tr', { hasText: pathPattern }).first()
  await expect(row).toBeVisible()

  await row.getByTitle('Edit').click()
  await page.getByRole('textbox', { name: 'Path Pattern' }).fill(renamedPathPattern)
  await page.getByRole('button', { name: 'Perbarui' }).click()
  await dismissSwal(page)

  const renamedRow = page.locator('tr', { hasText: renamedPathPattern }).first()
  await expect(renamedRow).toBeVisible()

  await renamedRow.getByTitle('Hapus').click()
  await page.getByRole('button', { name: 'Yes' }).click()
  await dismissSwal(page)

  await expect(page.locator('tr', { hasText: renamedPathPattern })).not.toBeVisible()
})

test('rejects a path_pattern not starting with /', async ({ page }) => {
  await page.getByRole('button', { name: 'Tambah Route' }).click()
  await pickSingleSelect(page, 'Service', 'serv-uam')
  await page.getByRole('textbox', { name: 'Path Pattern' }).fill('missing-slash')
  await page.getByRole('button', { name: 'Simpan' }).click()

  await expect(page.getByText('Path harus diawali dengan /')).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Tambah Route' })).toBeVisible()
})

test('refresh cache button runs without error', async ({ page }) => {
  await page.getByRole('button', { name: 'Refresh Cache Now' }).click()
  await expect(page.getByText('Cache route berhasil di-refresh.')).toBeVisible()
})
