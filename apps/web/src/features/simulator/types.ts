export interface Scenario {
  name: string
  description: string
  point_count: number
}

export interface SimulationStatus {
  running: boolean
  scenario?: string
  vehicle_count?: number
  interval_ms?: number
  duplicate_rate?: number
  invalid_rate?: number
  started_at?: string
}

export interface StartSimulationRequest {
  vehicle_ids: string[]
  scenario: string
  interval_ms: number
  duplicate_rate: number
  invalid_rate: number
}

export type LogEntryType = 'GPS_SENT' | 'DUPLICATE' | 'INVALID_PAYLOAD' | 'ERROR' | 'STARTED' | 'STOPPED'

export interface LogEntry {
  id: string
  timestamp: string
  type: LogEntryType
  message: string
}
