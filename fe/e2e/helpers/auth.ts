import type { APIRequestContext, Page } from '@playwright/test'

// Seeded accounts (serv-uam/cmd/migration/seeders/user_seeder.go, role_permission_seeder.go).
export const SUPER_ADMIN = { email: 'suAdmin@app.com', password: '@dmin#123' } // all permissions
export const ADMIN_GATEWAY = { email: 'admin@app.com', password: 'Admin#123' } // gateway-only permissions
export const VIEWER = { email: 'viewer@app.com', password: 'Viewer#123' } // read-only

interface Session {
  token: string
  refresh_token: string
  user: unknown
}

async function loginViaApi(request: APIRequestContext, creds: { email: string; password: string }) {
  const res = await request.post('/api/auth/login', { data: creds })
  if (!res.ok()) {
    throw new Error(`login failed for ${creds.email}: ${res.status()} ${await res.text()}`)
  }
  return ((await res.json()).data as Session) ?? null
}

// Logs in over the API (fast, avoids re-driving the login form in every spec — that flow is
// already covered by auth.spec.ts) and seeds localStorage exactly like stores/auth.ts does, so
// the app boots already authenticated.
export async function authenticate(
  page: Page,
  request: APIRequestContext,
  creds: { email: string; password: string } = SUPER_ADMIN,
) {
  const session = await loginViaApi(request, creds)
  await page.addInitScript((s) => {
    localStorage.setItem('token', JSON.stringify(s.token))
    localStorage.setItem('refresh_token', JSON.stringify(s.refresh_token))
    localStorage.setItem('user', JSON.stringify(s.user))
  }, session)
  return session
}
