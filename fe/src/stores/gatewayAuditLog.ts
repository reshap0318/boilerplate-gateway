import { defineStore } from 'pinia'
import { useCrud } from '@/composables'

export interface IGatewayAuditLog {
  id: number
  entity_type: string
  entity_id: number
  action: string
  actor_id: number | null
  actor_name: string
  actor_email: string
  // {"before": {...}, "after": {...}} — create sends after only, delete sends before only.
  // Omitted by the backend entirely when empty (omitempty).
  payloads?: Record<string, unknown>
  created_at: string
}

// Audit Log is read-only (no create/update/delete from the UI) — useCrud is still used for
// its fetchAll/pagination/loading plumbing, but the payload type is intentionally empty.
export type IGatewayAuditLogPayload = Record<string, never>

export const useGatewayAuditLogStore = defineStore('gatewayAuditLog', () => {
  const crud = useCrud<IGatewayAuditLog, IGatewayAuditLogPayload>({
    endpoint: '/api/svc/uam/audit-logs',
    entityName: 'audit log',
    initialForm: {} as IGatewayAuditLogPayload,
    formRules: {},
  })

  return { ...crud }
})
