import { apiClient } from '@/lib/api/client'
import type {
  CreateLeadPayload,
  Lead,
  LeadsListResponse,
  UpdateLeadPayload,
} from './types'

export async function fetchLeads(
  limit = 50,
  offset = 0
): Promise<LeadsListResponse> {
  return apiClient.get<LeadsListResponse>(
    `/api/v1/leads?limit=${limit}&offset=${offset}`
  )
}

export async function createLead(payload: CreateLeadPayload): Promise<Lead> {
  return apiClient.post<Lead>('/api/v1/leads', {
    contato_ok: 'N',
    ...payload,
  })
}

export async function updateLead(
  id: number,
  payload: UpdateLeadPayload
): Promise<Lead> {
  return apiClient.put<Lead>(`/api/v1/leads/${id}`, payload)
}
