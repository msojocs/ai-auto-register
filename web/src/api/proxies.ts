import api from './axios'

export interface Proxy {
  id: number
  host: string
  port: string
  username: string
  password: string
  protocol: string
  status: string
}

export interface CreateProxyPayload {
  host: string
  port: string
  username?: string
  password?: string
  protocol?: string
}

export function getProxies() {
  return api.get<{ proxies: Proxy[]; total: number }>('/proxies')
}

export function createProxy(payload: CreateProxyPayload) {
  return api.post<{ proxy: Proxy }>('/proxies', payload)
}

export function updateProxy(id: number, payload: CreateProxyPayload) {
  return api.put<{ proxy: Proxy }>(`/proxies/${id}`, payload)
}

export function deleteProxy(id: number) {
  return api.delete(`/proxies/${id}`)
}

export function testProxy(id: number) {
  return api.post<{ ok: boolean }>(`/proxies/${id}/test`)
}
