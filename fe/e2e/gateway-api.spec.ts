import { test, expect } from '@playwright/test'

// API-level regression for the gateway's routing redesign:
//   - static Management API routes (gateway's own /api/*)
//   - dynamic DB-driven proxy routes (/api/svc/uam/*, /api/svc/message/*)
//   - forgot/reset-password bookkeeping rows (hidden from the route list, still functional)
//   - retired paths and nginx's /ws/ location reaching the gateway
// Uses the `request` fixture only — no browser needed, runs against the live stack.

let token: string

test.describe('system', () => {
  test('gateway health via nginx', async ({ request }) => {
    const res = await request.get('/health')
    expect(res.status()).toBe(200)
  })

  test('SPA fallback serves index.html for an unknown client route', async ({ request }) => {
    const res = await request.get('/dashboard')
    expect(res.status()).toBe(200)
    expect(await res.text()).toContain('<!doctype html')
  })
})

test.describe('auth (static Management API)', () => {
  test('login with correct credentials returns a token', async ({ request }) => {
    const res = await request.post('/api/auth/login', {
      data: { email: 'admin@app.com', password: 'Admin#123' },
    })
    expect(res.status()).toBe(200)
    const body = await res.json()
    expect(body.data.token).toBeTruthy()
    token = body.data.token
  })

  test('login with wrong password is rejected', async ({ request }) => {
    const res = await request.post('/api/auth/login', {
      data: { email: 'admin@app.com', password: 'wrong' },
    })
    expect(res.status()).toBe(422)
  })

  test('forgot-password (server-to-server forward) still works', async ({ request }) => {
    const res = await request.post('/api/auth/forgot-password', {
      data: { email: 'admin@app.com' },
    })
    expect(res.status()).toBe(200)
  })
})

test.describe('management API: services & routes', () => {
  test.beforeAll(async ({ request }) => {
    const res = await request.post('/api/auth/login', {
      data: { email: 'admin@app.com', password: 'Admin#123' },
    })
    token = (await res.json()).data.token
  })

  test('services list contains serv-uam and serv-message under /api/svc/*', async ({ request }) => {
    const res = await request.get('/api/services', {
      headers: { Authorization: `Bearer ${token}` },
    })
    expect(res.status()).toBe(200)
    const body = await res.json()
    const names = body.data.map((s: { name: string }) => s.name)
    expect(names).toContain('serv-uam')
    expect(names).toContain('serv-message')
    for (const svc of body.data) {
      expect(svc.base_path).toMatch(/^\/(api|ws)\/svc\//)
    }
  })

  test('route list hides the forgot/reset-password bookkeeping rows (IDs 1/2)', async ({
    request,
  }) => {
    const res = await request.get('/api/routes?page_size=-1', {
      headers: { Authorization: `Bearer ${token}` },
    })
    expect(res.status()).toBe(200)
    const body = await res.json()
    const patterns = body.data.map((r: { path_pattern: string }) => r.path_pattern)
    expect(patterns).not.toContain('/auth/forgot-password')
    expect(patterns).not.toContain('/auth/reset-password')
    expect(patterns).toContain('/me')
  })
})

test.describe('dynamic proxy', () => {
  test.beforeAll(async ({ request }) => {
    const res = await request.post('/api/auth/login', {
      data: { email: 'admin@app.com', password: 'Admin#123' },
    })
    token = (await res.json()).data.token
  })

  test('uam/me without token is rejected', async ({ request }) => {
    const res = await request.get('/api/svc/uam/me')
    expect(res.status()).toBe(401)
  })

  test('uam/me with token is proxied through to serv-uam', async ({ request }) => {
    const res = await request.get('/api/svc/uam/me', {
      headers: { Authorization: `Bearer ${token}` },
    })
    expect(res.status()).toBe(200)
  })

  test('message/notifications with token is proxied through to serv-message', async ({
    request,
  }) => {
    const res = await request.get('/api/svc/message/notifications', {
      headers: { Authorization: `Bearer ${token}` },
    })
    expect(res.status()).toBe(200)
  })
})

test.describe('retired paths & nginx routing', () => {
  test('old pre-migration /uam/me path is no longer proxied (falls back to the SPA)', async ({
    request,
  }) => {
    const res = await request.get('/uam/me')
    expect(res.status()).toBe(200)
    expect(await res.text()).toContain('<!doctype html')
  })

  test("/ws/ location reaches the gateway (JSON 404, not nginx's own 404 page)", async ({
    request,
  }) => {
    const res = await request.get('/ws/svc/nonexistent')
    const body = await res.json().catch(() => null)
    expect(body).not.toBeNull()
    expect(body.message).toBeTruthy()
  })
})
