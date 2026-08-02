import type { Page } from '@playwright/test'

// swal.success()/swal.error() (src/plugins/swal.ts) render a SweetAlert2 popup with no
// auto-close timer — it sits on top of the page and blocks further clicks until dismissed.
export async function dismissSwal(page: Page) {
  await page.getByRole('button', { name: 'Confirm' }).click()
}
