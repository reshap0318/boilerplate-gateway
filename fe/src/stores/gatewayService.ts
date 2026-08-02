import { defineStore } from 'pinia'
import { watch } from 'vue'
import { required, minLength, minValue, helpers } from '@vuelidate/validators'
import { get, post, type IApiResponse } from '@/plugins/axios'
import { useCrud } from '@/composables'

export interface IGatewayServiceHealth {
  health_status: 'unknown' | 'up' | 'down'
  health_checked_at: string | null
}

// @vuelidate/validators' `url` rule requires a TLD-style host and rejects things like
// "http://localhost:9000" — upstream services are frequently internal hosts/IPs, so a
// lenient scheme+host check (via the URL constructor) is used instead.
function isValidBaseUrl(value: string): boolean {
  if (!value) return true
  try {
    const parsed = new URL(value)
    return (parsed.protocol === 'http:' || parsed.protocol === 'https:') && !!parsed.hostname
  } catch {
    return false
  }
}

// Mirrors backend helpers.ValidateBasePath: must start with "/", must not end with "/",
// no empty segments ("//"), no "*"/":" segments, and confined to the reserved namespace
// for the selected protocol ("/api/svc/" for http, "/ws/svc/" for websocket).
function requiredBasePathPrefix(protocol: string): string {
  return protocol === 'websocket' ? '/ws/svc/' : '/api/svc/'
}

function isValidBasePath(value: string, protocol: string): boolean {
  if (!value) return true
  if (value === '/' || !value.startsWith('/') || value.endsWith('/')) return false
  if (!value.startsWith(requiredBasePathPrefix(protocol))) return false
  return value
    .slice(1)
    .split('/')
    .every((segment) => segment !== '' && !segment.includes('*') && !segment.startsWith(':'))
}

export interface IGatewayService {
  id: number
  name: string
  base_url: string
  base_path: string
  protocol: 'http' | 'websocket'
  is_active: boolean
  rate_limit_per_minute: number | null
  health_status: 'unknown' | 'up' | 'down'
  health_checked_at: string | null
  route_count: number
  created_at: string
  updated_at: string
}

export interface IGatewayServicePayload {
  id?: number
  name: string
  base_url: string
  base_path: string
  protocol: 'http' | 'websocket'
  rate_limit_per_minute: number | null
  is_active: boolean
}

export const useGatewayServiceStore = defineStore('gatewayService', () => {
  const crud = useCrud<IGatewayService, IGatewayServicePayload>({
    endpoint: '/api/services',
    entityName: 'service',
    initialForm: {
      name: '',
      base_url: '',
      base_path: '',
      protocol: 'http',
      rate_limit_per_minute: null,
      is_active: true,
    },
    formRules: {
      name: { required, minLength: minLength(3) },
      base_url: {
        required,
        validUrl: helpers.withMessage('Base URL tidak valid', isValidBaseUrl),
      },
      base_path: {
        required,
        validBasePath: helpers.withMessage(
          'Base path harus diawali /, tidak boleh diakhiri /, tidak boleh /api, dan tidak boleh mengandung * atau :',
          isValidBasePath,
        ),
      },
      protocol: { required },
      rate_limit_per_minute: {
        minValue: helpers.withMessage('Rate limit harus lebih dari 0', minValue(1)),
      },
    },
  })

  // Auto-swap the base_path prefix when protocol changes, so the field always points at the
  // namespace reserved for the currently-selected protocol instead of relying on the admin
  // to retype it correctly.
  watch(
    () => crud.form.protocol,
    (newProtocol, oldProtocol) => {
      if (newProtocol === oldProtocol) return
      const newPrefix = requiredBasePathPrefix(newProtocol)
      const oldPrefix = requiredBasePathPrefix(oldProtocol)
      if (!crud.form.base_path) {
        crud.form.base_path = newPrefix
      } else if (crud.form.base_path.startsWith(oldPrefix)) {
        crud.form.base_path = newPrefix + crud.form.base_path.slice(oldPrefix.length)
      }
    },
  )

  // Extended base_path rule — protocol-aware, so it must be defined after crud (references
  // crud.form.protocol) and overrides the static one passed into useCrud above.
  const formRules = {
    ...crud.formRules,
    base_path: {
      required,
      validBasePath: helpers.withMessage(
        'Base path harus diawali /api/svc/ (protocol http) atau /ws/svc/ (protocol websocket), sesuai protocol yang dipilih',
        (value: string) => isValidBasePath(value, crud.form.protocol),
      ),
    },
  }

  async function fetchAllServices(): Promise<IGatewayService[]> {
    try {
      const { data } = await get<IApiResponse<IGatewayService[]>>('/api/services')
      return data.data || []
    } catch (error: any) {
      console.error('Failed to fetch all services', error)
      return []
    }
  }

  async function healthCheck(id: number): Promise<void> {
    try {
      await post<IApiResponse<IGatewayServiceHealth>>(`/api/services/${id}/health-check`, {})
      await crud.fetchAll()
    } catch (error: any) {
      console.error('Failed to run health check', error)
      throw error
    }
  }

  return {
    ...crud,
    formRules,
    fetchAllServices,
    healthCheck,
  }
})
