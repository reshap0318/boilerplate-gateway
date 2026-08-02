import { test, expect } from '@playwright/test'
import { authenticate, ADMIN_GATEWAY, SUPER_ADMIN } from './helpers/auth'
import { dismissSwal } from './helpers/swal'

async function loginViaUi(page: import('@playwright/test').Page) {
  await page.goto('/login')
  await page.getByRole('textbox', { name: 'Email' }).fill(SUPER_ADMIN.email)
  await page.getByRole('textbox', { name: 'Password' }).fill(SUPER_ADMIN.password)
  await page.getByRole('button', { name: 'Login' }).click()
  await expect(page).toHaveURL('/')
  await dismissSwal(page)
}

test.describe('login', () => {
  test('correct credentials render the dashboard with no console errors', async ({ page }) => {
    const errors: string[] = []
    page.on('console', (msg) => {
      if (msg.type() === 'error') errors.push(msg.text())
    })

    await loginViaUi(page)
    expect(errors).toEqual([])
  })

  test('wrong password shows an error and stays on the login page', async ({ page }) => {
    await page.goto('/login')
    await page.getByRole('textbox', { name: 'Email' }).fill(SUPER_ADMIN.email)
    await page.getByRole('textbox', { name: 'Password' }).fill('wrong-password')
    await page.getByRole('button', { name: 'Login' }).click()

    await expect(page.getByText('Login Gagal')).toBeVisible()
    await expect(page).toHaveURL('/login')
  })
})

test.describe('route guards', () => {
  test('protected route redirects to /login when unauthenticated', async ({ page }) => {
    await page.goto('/users')
    await expect(page).toHaveURL('/login')
  })

  test('guest-only route redirects to / when already authenticated', async ({ page, request }) => {
    await authenticate(page, request, SUPER_ADMIN)
    await page.goto('/login')
    await expect(page).toHaveURL('/')
  })

  test('route requiring a permission the user lacks redirects home', async ({ page, request }) => {
    // Admin Gateway role has no user.index permission (see role_permission_seeder.go).
    await authenticate(page, request, ADMIN_GATEWAY)
    await page.goto('/users')
    // HomeView shows the warning then immediately strips ?accessDenied=true (see HomeView.vue) —
    // assert the toast instead of racing that transient query string.
    await expect(page.getByText('Access Denied')).toBeVisible()
    await expect(page).toHaveURL('/')
  })
})

test.describe('forgot password', () => {
  test('submits and redirects back to login', async ({ page }) => {
    await page.goto('/forgot-password')
    await page.getByRole('textbox', { name: 'Email' }).fill(SUPER_ADMIN.email)
    await page.getByRole('button', { name: 'Kirim Tautan Reset' }).click()

    await expect(page.getByText('Berhasil')).toBeVisible()
    await expect(page).toHaveURL('/login')
  })
})

test.describe('logout', () => {
  test('clears the session and returns to login', async ({ page }) => {
    // Real UI login here (not the authenticate() shortcut) — that helper's addInitScript would
    // re-inject the still-valid session on every subsequent page.goto(), masking a logout bug.
    await loginViaUi(page)

    await page.locator('button.group').first().click() // avatar dropdown toggle
    await page.getByRole('button', { name: 'Log out' }).click()

    await expect(page).toHaveURL('/login')
    // A protected route should now bounce back to login too.
    await page.goto('/users')
    await expect(page).toHaveURL('/login')
  })
})
