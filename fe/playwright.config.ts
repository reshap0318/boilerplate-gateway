import { defineConfig } from '@playwright/test'

// Regression e2e against the running docker-compose stack (fe's nginx on :3000, proxying to
// gateway/uam/message). Requires `docker compose up -d` beforehand — this config does not
// start the stack itself.
export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  // Sequential on purpose: specs share one gateway DB/route cache (CRUD tests mutate real
  // rows) and the gateway itself rate-limits /api/* — parallel workers both race each other's
  // data and trip that limiter with concurrent logins.
  workers: 1,
  // 'html' always writes the report to disk; open:'never' stops it from trying to launch a
  // browser after every run — view it on demand via `yarn test:e2e:report`.
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',
  },
})
