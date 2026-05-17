import { apiClient } from '@/lib/api/client'
import { useAuthStore } from '@/stores/auth-store'

/** Restores API auth header from persisted token (cookie). */
export function restoreAuthSession(): void {
  const token = useAuthStore.getState().auth.accessToken
  if (token) {
    apiClient.setAuthToken(token)
  }
}
