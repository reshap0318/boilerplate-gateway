import type { Page } from '@playwright/test'

// FormSelect (src/components/utils/FormSelect.vue) always renders @vueform/multiselect with
// append-to-body:true, so a page with more than one Multiselect (an IndexView filter AND a
// modal field, both built from the same options) ends up with duplicate `.multiselect-option`
// elements sharing the same value/id in the DOM — only the currently-open dropdown's copy is
// actually visible, so option lookups are scoped to `.multiselect-dropdown:visible` rather than
// matched globally, which would otherwise resolve the closed duplicate.
async function openMultiselect(page: Page, labelText: string) {
  const control = page.locator(`label:text-is("${labelText}") + .multiselect`)
  // The previous dropdown's close transition (~150ms) can leave its container matching
  // `:visible` for a moment after a fresh click — settle before opening the next one so
  // `visibleDropdownOption` doesn't resolve against a still-fading, unrelated dropdown.
  await page.waitForTimeout(250)
  await control.click()
  await page.locator('.multiselect-dropdown:visible').waitFor()
  return control
}

function visibleDropdownOption(page: Page, optionText: string) {
  return page
    .locator('.multiselect-dropdown:visible .multiselect-option', { hasText: optionText })
    .first()
}

export async function pickSingleSelect(page: Page, labelText: string, optionText: string) {
  await openMultiselect(page, labelText)
  await visibleDropdownOption(page, optionText).click()
}

export async function pickTagSelect(page: Page, labelText: string, optionText: string) {
  const control = await openMultiselect(page, labelText)
  await visibleDropdownOption(page, optionText).click()
  await control.press('Escape')
}
