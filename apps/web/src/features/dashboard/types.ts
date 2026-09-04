export interface Alert {
  id: string
  vehicle_id: string
  type: string
  status: string
  started_at: string
  resolved_at?: string | null
  created_at: string
  updated_at: string
}

export interface AlertListResponse {
  alerts: Alert[]
}

export interface VehiclePosition {
  vehicleId: string
  latitude: number
  longitude: number
  timestamp: string
}

export type VehicleStatus = 'MOVING' | 'STOPPED' | 'ALERT' | 'UNKNOWN'

export interface VehicleState {
  vehicleId: string
  status: VehicleStatus
  lastLatitude?: number
  lastLongitude?: number
  lastUpdate?: string
}

export interface SSEEvent {
  type: string
  payload: Record<string, unknown>
}
