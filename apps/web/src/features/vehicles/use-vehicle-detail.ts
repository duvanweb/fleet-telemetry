import { useQuery } from '@tanstack/react-query'

import { vehicleService } from './vehicle.service'
import { telemetryService } from './telemetry.service'

export function useVehicleDetail(id: string) {
  return useQuery({
    queryKey: ['vehicles', id],
    queryFn: () => vehicleService.getById(id),
    enabled: !!id,
  })
}

export function useVehicleTelemetry(vehicleId: string, from?: string, to?: string) {
  return useQuery({
    queryKey: ['telemetry', vehicleId, from, to],
    queryFn: () => telemetryService.getByVehicleID(vehicleId, from, to),
    enabled: !!vehicleId,
    staleTime: 30_000,
  })
}
