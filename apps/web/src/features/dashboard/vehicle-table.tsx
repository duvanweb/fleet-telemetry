import { Link } from 'react-router-dom'

import type { Vehicle } from '../vehicles/types'
import type { VehicleState, VehicleStatus } from './types'

interface Props {
  vehicles: Vehicle[]
  vehicleStates: Map<string, VehicleState>
}

const STATUS_BADGE: Record<VehicleStatus, { label: string; className: string }> = {
  MOVING: { label: 'Moving', className: 'bg-green-100 text-green-700' },
  STOPPED: { label: 'Stopped', className: 'bg-yellow-100 text-yellow-700' },
  ALERT: { label: 'Alert', className: 'bg-red-100 text-red-700' },
  UNKNOWN: { label: 'Unknown', className: 'bg-gray-100 text-gray-500' },
}

function StatusBadge({ status }: { status: VehicleStatus }) {
  const { label, className } = STATUS_BADGE[status]
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${className}`}>
      {label}
    </span>
  )
}

export function VehicleTable({ vehicles, vehicleStates }: Props) {
  const active = vehicles.filter((v) => !v.deleted_at)

  if (active.length === 0) {
    return (
      <div className='bg-white rounded-xl border border-gray-100 shadow-sm p-8 text-center text-gray-400'>
        No vehicles registered.
      </div>
    )
  }

  return (
    <div className='bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden'>
      <div className='overflow-x-auto'>
        <table className='w-full text-sm'>
          <thead className='bg-gray-50 border-b border-gray-100'>
            <tr>
              <th className='px-4 py-3 text-left font-medium text-gray-500'>Name</th>
              <th className='px-4 py-3 text-left font-medium text-gray-500'>Plate</th>
              <th className='px-4 py-3 text-left font-medium text-gray-500'>Status</th>
              <th className='px-4 py-3 text-left font-medium text-gray-500'>Last update</th>
            </tr>
          </thead>
          <tbody className='divide-y divide-gray-50'>
            {active.map((v) => {
              const state = vehicleStates.get(v.id)
              const status: VehicleStatus = state?.status ?? 'UNKNOWN'
              const lastUpdate = state?.lastUpdate
                ? new Date(state.lastUpdate).toLocaleTimeString()
                : '—'
              return (
                <tr key={v.id} className='hover:bg-gray-50 transition-colors'>
                  <td className='px-4 py-3 font-medium'>
                    <Link to={`/vehicles/${v.id}`} className='text-blue-600 hover:underline'>
                      {v.name}
                    </Link>
                  </td>
                  <td className='px-4 py-3 text-gray-600 font-mono'>{v.plate}</td>
                  <td className='px-4 py-3'>
                    <StatusBadge status={status} />
                  </td>
                  <td className='px-4 py-3 text-gray-400'>{lastUpdate}</td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
