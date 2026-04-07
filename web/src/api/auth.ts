import api from './axios'

export interface LoginResponse {
  token: string
  user: { id: number; username: string; avatar_url?: string }
}

export function login(username: string, password: string) {
  return api.post<LoginResponse>('/auth/login', { username, password })
}

export function changePassword(currentPassword: string, newPassword: string) {
  return api.post('/auth/change-password', {
    current_password: currentPassword,
    new_password: newPassword,
  })
}
