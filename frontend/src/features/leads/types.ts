export interface Lead {
  id: number
  fone1: string
  fone2?: string | null
  nome: string
  status_id: number
  unidd_id?: number | null
  email?: string | null
  contato_ok: string
  criado_em: string
  alterado_em?: string | null
}

export interface LeadsListResponse {
  leads: Lead[]
  limit: number
  offset: number
}

export interface CreateLeadPayload {
  fone1: string
  nome: string
  status_id: number
  contato_ok?: string
  email?: string
}

export interface UpdateLeadPayload {
  fone1?: string
  nome?: string
  status_id?: number
  contato_ok?: string
  email?: string
}
