const BASE_URL = 'http://localhost:8082/api'

export interface TelemetryPoint {
  latitude: number
  longitude: number
  device_timestamp: string
}

export interface TelemetryHistoryResponse {
  vehicle_id: string
  points: TelemetryPoint[]
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const text = await res.text().catch(() => res.statusText)
    throw new Error(text || `HTTP ${res.status}`)
  }
  return res.json() as Promise<T>
}

export const telemetryService = {
  getByVehicleID(vehicleId: string, from?: string, to?: string): Promise<TelemetryHistoryResponse> {
    const params = new URLSearchParams()
    if (from) params.set('from', from)
    if (to) params.set('to', to)
    const qs = params.toString()
    return request<TelemetryHistoryResponse>(`/vehicles/${vehicleId}/telemetry${qs ? `?${qs}` : ''}`)
  },
}
