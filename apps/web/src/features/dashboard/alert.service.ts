import type { Alert, AlertListResponse } from './types'

const BASE_URL = 'http://localhost:8083/api'

async function request<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`)
  if (!res.ok) {
    const text = await res.text().catch(() => res.statusText)
    throw new Error(text || `HTTP ${res.status}`)
  }
  return res.json() as Promise<T>
}

export const alertService = {
  getAll(): Promise<AlertListResponse> {
    return request<AlertListResponse>('/alerts')
  },

  getOpenByVehicle(vehicleId: string): Promise<Alert> {
    return request<Alert>(`/vehicles/${vehicleId}/alerts`)
  },
}
