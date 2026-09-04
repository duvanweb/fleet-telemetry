import { useCallback, useState } from 'react'

import { AlertsPanel } from './alerts-panel'
import { FleetMap } from './fleet-map'
import { FleetMetrics } from './fleet-metrics'
import { VehicleTable } from './vehicle-table'
import { useDashboardAlerts, useDashboardVehicles } from './use-dashboard'
import { useFleetSSE } from './use-fleet-sse'
import type { VehicleState } from './types'

export default function Dashboard() {
  const { data: vehicles = [], isLoading: loadingVehicles } = useDashboardVehicles()
  const { data: alertData, isLoading: loadingAlerts } = useDashboardAlerts()
  const alerts = alertData?.alerts ?? []

  const [vehicleStates, setVehicleStates] = useState<Map<string, VehicleState>>(new Map())

  const handleVehicleUpdate = useCallback((state: VehicleState) => {
    setVehicleStates((prev) => {
      const next = new Map(prev)
      const existing = next.get(state.vehicleId)
      next.set(state.vehicleId, { ...existing, ...state })
      return next
    })
  }, [])

  useFleetSSE(handleVehicleUpdate)

  const isLoading = loadingVehicles || loadingAlerts

  return (
    <div className='min-h-screen bg-gray-50'>
      <header className='bg-white border-b border-gray-100 shadow-sm'>
        <div className='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-14 flex items-center justify-between'>
          <h1 className='text-lg font-semibold text-gray-800'>Fleet Dashboard</h1>
          <div className='flex items-center gap-2 text-xs text-gray-400'>
            <span className='w-2 h-2 rounded-full bg-green-400 animate-pulse' />
            Live
          </div>
        </div>
      </header>

      <main className='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-6'>
        {isLoading ? (
          <div className='flex justify-center py-16'>
            <div className='w-8 h-8 border-2 border-blue-500 border-t-transparent rounded-full animate-spin' />
          </div>
        ) : (
          <>
            <FleetMetrics vehicles={vehicles} alerts={alerts} vehicleStates={vehicleStates} />

            <FleetMap vehicles={vehicles} vehicleStates={vehicleStates} />

            <div className='grid grid-cols-1 lg:grid-cols-3 gap-6'>
              <div className='lg:col-span-2'>
                <h2 className='text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3'>
                  Vehicles
                </h2>
                <VehicleTable vehicles={vehicles} vehicleStates={vehicleStates} />
              </div>

              <div>
                <h2 className='text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3'>
                  Alerts
                </h2>
                <AlertsPanel alerts={alerts} />
              </div>
            </div>
          </>
        )}
      </main>
    </div>
  )
}
