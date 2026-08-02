import type { Page } from '@playwright/test'

// Seed data alone can span multiple pages (e.g. uam ships ~20+ permissions), and newly
// created rows are appended at the end (ascending id) — so a just-created card can land past
// page 1. Click "Next" (last button in the nav) until it's disabled.
export async function goToLastPage(page: Page) {
  const nav = page.locator('nav[aria-label="Pagination"]')
  if (!(await nav.isVisible().catch(() => false))) return

  const next = nav.getByRole('button').last()
  for (let i = 0; i < 20 && (await next.isEnabled()); i++) {
    await next.click()
    await page.waitForTimeout(300)
  }
}
