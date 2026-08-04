import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'
import { post } from '@/plugins/axios'
import storage from '@/helpers/storage'
import type { IApiResponse } from '@/plugins/axios'
import { required, email } from '@vuelidate/validators'

export interface ILoginPayload {
  email: string
  password: string
}

export interface ILoginResponse {
  twofa_required: boolean
  token?: string
  refresh_token?: string
  user?: {
    id: number
    name: string
    email: string
    avatar: string | null
    created_at: string
    roles: { id: number; name: string; description: string }[]
    permissions: { id: number; name: string; description: string }[]
  }
}

export interface IRefreshTokenPayload {
  refresh_token: string
}

export interface IRefreshTokenResponse {
  token: string
  refresh_token: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(storage.getItem<string>('token'))
  const refreshToken = ref<string | null>(storage.getItem<string>('refresh_token'))
  const user = ref<ILoginResponse['user'] | null>(storage.getItem<ILoginResponse['user']>('user'))
  const isLoading = ref(false)
  const isAuthenticated = computed((): boolean => !!token.value)
  // Set by login() when the account has 2FA enabled — the email a pending /2fa/verify belongs to.
  const pendingEmail = ref<string | null>(null)

  const form = reactive<ILoginPayload>({
    email: '',
    password: '',
  })

  const formRules = {
    email: {
      required,
      email,
    },
    password: {
      required,
    },
  }

  // applySession stores a completed login's tokens — shared by login() and verifyTwoFA().
  function applySession(result: ILoginResponse) {
    if (!result.token || !result.refresh_token || !result.user) return

    token.value = result.token
    refreshToken.value = result.refresh_token
    user.value = result.user

    storage.setItem('token', result.token)
    storage.setItem('refresh_token', result.refresh_token)
    storage.setItem('user', result.user)
  }

  // login returns true when the account has 2FA enabled — the caller should route to /2fa
  // instead of treating this as a completed login.
  async function login(): Promise<boolean> {
    isLoading.value = true
    try {
      const response = await post<IApiResponse<ILoginResponse>>(
        '/api/auth/login',
        { email: form.email, password: form.password },
        { hideError: true },
      )
      const result = response.data.data

      if (result.twofa_required) {
        pendingEmail.value = form.email
        return true
      }

      applySession(result)
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function verifyTwoFA(code: string): Promise<void> {
    if (!pendingEmail.value) throw new Error('No pending 2FA login')

    isLoading.value = true
    try {
      const response = await post<IApiResponse<ILoginResponse>>(
        '/api/auth/2fa/verify',
        { email: pendingEmail.value, code },
        { hideError: true },
      )
      applySession(response.data.data)
      pendingEmail.value = null
    } finally {
      isLoading.value = false
    }
  }

  async function refreshTokenFn(): Promise<string> {
    if (!refreshToken.value) throw new Error('No refresh token available')

    const response = await post<IApiResponse<IRefreshTokenResponse>>(
      '/api/auth/refresh',
      { refresh_token: refreshToken.value },
      { hideError401: true },
    )
    const { token: newToken, refresh_token: newRefreshToken } = response.data.data

    token.value = newToken
    refreshToken.value = newRefreshToken

    storage.setItem('token', newToken)
    storage.setItem('refresh_token', newRefreshToken)

    return newToken
  }

  async function logout(): Promise<void> {
    try {
      if (token.value) {
        await post('/api/auth/logout')
      }
    } catch (error) {
      console.error('Logout error', error)
    } finally {
      token.value = null
      refreshToken.value = null
      user.value = null
      storage.clearAll()
    }
  }

  async function forgotPassword(email: string): Promise<void> {
    await post('/api/auth/forgot-password', { email })
  }

  async function resetPassword(token: string, new_password: string): Promise<void> {
    await post('/api/auth/reset-password', { token, new_password })
  }

  return {
    token,
    refreshToken,
    user,
    isLoading,
    pendingEmail,
    form,
    formRules,
    login,
    verifyTwoFA,
    refreshTokenFn,
    logout,
    forgotPassword,
    resetPassword,
    isAuthenticated,
  }
})
