import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { simulatorService } from './simulator.service'
import type { StartSimulationRequest } from './types'

export function useSimulatorScenarios() {
  return useQuery({
    queryKey: ['simulator', 'scenarios'],
    queryFn: () => simulatorService.getScenarios(),
    staleTime: Infinity,
  })
}

export function useSimulatorStatus() {
  return useQuery({
    queryKey: ['simulator', 'status'],
    queryFn: () => simulatorService.getStatus(),
    refetchInterval: 2000,
  })
}

export function useStartSimulation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: StartSimulationRequest) => simulatorService.start(data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['simulator', 'status'] }),
  })
}

export function useStopSimulation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => simulatorService.stop(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['simulator', 'status'] }),
  })
}
