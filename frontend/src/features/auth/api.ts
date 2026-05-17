import { apiClient } from '@/lib/api/client'

export interface LoginRequest {
  login: string
  password: string
}

export interface SystemUser {
  id: number
  login: string
  name: string
  email: string
  system_unit_id: number | null
  active: boolean
}

export interface LoginResponse {
  token: string
  expires_at: string
  user: SystemUser
}

export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  return apiClient.post<LoginResponse>('/api/v1/auth/login', credentials)
}
