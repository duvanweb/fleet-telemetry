import { useEffect } from 'react'
import { useQueryClient } from '@tanstack/react-query'

import type { VehicleState } from './types'

const SSE_URL = 'http://localhost:8082/api/events'

export function useFleetSSE(
  onVehicleUpdate: (state: VehicleState) => void,
) {
  const queryClient = useQueryClient()

  useEffect(() => {
    const source = new EventSource(SSE_URL)

    source.addEventListener('telemetry.received', (e) => {
      try {
        const payload = JSON.parse(e.data)
        const data = payload.payload ?? payload
        onVehicleUpdate({
          vehicleId: data.vehicle_id,
          status: 'MOVING',
          lastLatitude: data.latitude,
          lastLongitude: data.longitude,
          lastUpdate: data.received_at ?? new Date().toISOString(),
        })
      } catch {
        // ignore malformed frames
      }
    })

    source.addEventListener('alert.created', () => {
      queryClient.invalidateQueries({ queryKey: ['dashboard', 'alerts'] })
    })

    source.addEventListener('alert.resolved', () => {
      queryClient.invalidateQueries({ queryKey: ['dashboard', 'alerts'] })
    })

    source.addEventListener('vehicle.updated', () => {
      queryClient.invalidateQueries({ queryKey: ['dashboard', 'vehicles'] })
    })

    source.onerror = () => {
      source.close()
    }

    return () => source.close()
  }, [queryClient, onVehicleUpdate])
}
