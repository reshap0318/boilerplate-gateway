import { test, expect } from '@playwright/test'
import { authenticate, SUPER_ADMIN } from './helpers/auth'
import { dismissSwal } from './helpers/swal'

test.beforeEach(async ({ page, request }) => {
  await authenticate(page, request, SUPER_ADMIN)
  await page.goto('/profile')
})

test('renders the current user profile', async ({ page }) => {
  await expect(page.getByRole('heading', { name: 'Super Admin' })).toBeVisible()
  await expect(page.getByText(SUPER_ADMIN.email)).toBeVisible()
})

test('updates the profile name', async ({ page }) => {
  await page.getByRole('button', { name: 'Edit Profile' }).click()
  await page.getByRole('textbox', { name: 'Nama' }).fill('Super Admin')
  await page.getByRole('button', { name: 'Simpan' }).click()
  await dismissSwal(page)

  await expect(page.getByRole('heading', { name: 'Edit Profile' })).not.toBeVisible()
  await expect(page.getByRole('heading', { name: 'Super Admin' })).toBeVisible()
})
