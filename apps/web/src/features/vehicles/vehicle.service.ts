import type { CreateVehicleRequest, UpdateVehicleRequest, Vehicle } from './types'

const BASE_URL = 'http://localhost:8081/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const text = await res.text().catch(() => res.statusText)
    throw new Error(text || `HTTP ${res.status}`)
  }
  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}

export const vehicleService = {
  getAll(): Promise<Vehicle[]> {
    return request<Vehicle[]>('/vehicles')
  },

  getById(id: string): Promise<Vehicle> {
    return request<Vehicle>(`/vehicles/${id}`)
  },

  create(data: CreateVehicleRequest): Promise<Vehicle> {
    return request<Vehicle>('/vehicles', { method: 'POST', body: JSON.stringify(data) })
  },

  update(id: string, data: UpdateVehicleRequest): Promise<Vehicle> {
    return request<Vehicle>(`/vehicles/${id}`, { method: 'PATCH', body: JSON.stringify(data) })
  },

  delete(id: string): Promise<void> {
    return request<void>(`/vehicles/${id}`, { method: 'DELETE' })
  },
}
