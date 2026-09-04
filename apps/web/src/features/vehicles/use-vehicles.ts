import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { vehicleService } from './vehicle.service'
import type { CreateVehicleRequest, UpdateVehicleRequest } from './types'

const QUERY_KEY = ['vehicles']

export function useVehicles() {
  return useQuery({
    queryKey: QUERY_KEY,
    queryFn: () => vehicleService.getAll(),
  })
}

export function useCreateVehicle() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateVehicleRequest) => vehicleService.create(data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })
}

export function useUpdateVehicle() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateVehicleRequest }) =>
      vehicleService.update(id, data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })
}

export function useDeleteVehicle() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => vehicleService.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: QUERY_KEY }),
  })
}
