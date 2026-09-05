import { useQuery } from '@tanstack/react-query'

import { alertService } from './alert.service'
import { vehicleService } from '../vehicles/vehicle.service'

export function useDashboardVehicles() {
  return useQuery({
    queryKey: ['dashboard', 'vehicles'],
    queryFn: () => vehicleService.getAll(),
    refetchInterval: 30_000,
  })
}

export function useDashboardAlerts() {
  return useQuery({
    queryKey: ['dashboard', 'alerts'],
    queryFn: () => alertService.getAll(),
    refetchInterval: 15_000,
  })
}
