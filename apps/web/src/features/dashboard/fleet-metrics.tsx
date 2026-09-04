import type { Vehicle } from '../vehicles/types'
import type { Alert, VehicleState } from './types'

interface Props {
  vehicles: Vehicle[]
  alerts: Alert[]
  vehicleStates: Map<string, VehicleState>
}

interface MetricTileProps {
  label: string
  value: number
  color: string
}

function MetricTile({ label, value, color }: MetricTileProps) {
  return (
    <div className='bg-white rounded-xl border border-gray-100 shadow-sm p-5'>
      <p className='text-sm text-gray-500'>{label}</p>
      <p className={`text-3xl font-bold mt-1 ${color}`}>{value}</p>
    </div>
  )
}

export function FleetMetrics({ vehicles, alerts, vehicleStates }: Props) {
  const active = vehicles.filter((v) => !v.deleted_at)
  const stopped = [...vehicleStates.values()].filter((s) => s.status === 'STOPPED').length
  const moving = [...vehicleStates.values()].filter((s) => s.status === 'MOVING').length
  const openAlerts = alerts.filter((a) => a.status === 'OPEN').length

  return (
    <div className='grid grid-cols-2 gap-4 sm:grid-cols-4'>
      <MetricTile label='Total vehicles' value={active.length} color='text-gray-800' />
      <MetricTile label='Moving' value={moving} color='text-green-600' />
      <MetricTile label='Stopped' value={stopped} color='text-yellow-600' />
      <MetricTile label='Open alerts' value={openAlerts} color='text-red-600' />
    </div>
  )
}
