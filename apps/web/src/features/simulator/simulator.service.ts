import type { Scenario, SimulationStatus, StartSimulationRequest } from './types'

const BASE_URL = 'http://localhost:8090/api'

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

export const simulatorService = {
  getScenarios(): Promise<{ scenarios: Scenario[] }> {
    return request('/simulator/scenarios')
  },

  getStatus(): Promise<SimulationStatus> {
    return request('/simulator/status')
  },

  start(data: StartSimulationRequest): Promise<void> {
    return request('/simulator/start', { method: 'POST', body: JSON.stringify(data) })
  },

  startScenario(scenario: string, data: Omit<StartSimulationRequest, 'scenario'>): Promise<void> {
    return request(`/simulator/scenarios/${scenario}/start`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  stop(): Promise<void> {
    return request('/simulator/stop', { method: 'POST', body: '{}' })
  },
}
