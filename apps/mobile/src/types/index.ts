export interface TelemetryPayload {
  vehicle_id: string
  latitude: number
  longitude: number
  device_timestamp: string
}

export interface QueuedTelemetry extends TelemetryPayload {
  id: string
  synced: boolean
  queued_at: string
}

export interface AlertRecord {
  id: string
  vehicle_id: string
  type: string
  status: string
  started_at: string
  received_at: string
}

export interface Coordinates {
  latitude: number
  longitude: number
  speed?: number | null
  accuracy?: number | null
}
