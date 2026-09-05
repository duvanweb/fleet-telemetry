import type { Alert } from './types'

interface Props {
  alerts: Alert[]
}

function formatDuration(startedAt: string): string {
  const ms = Date.now() - new Date(startedAt).getTime()
  const minutes = Math.floor(ms / 60_000)
  if (minutes < 60) return `${minutes}m`
  const hours = Math.floor(minutes / 60)
  return `${hours}h ${minutes % 60}m`
}

export function AlertsPanel({ alerts }: Props) {
  const open = alerts.filter((a) => a.status === 'OPEN')

  return (
    <div className='bg-white rounded-xl border border-gray-100 shadow-sm'>
      <div className='px-4 py-3 border-b border-gray-100 flex items-center justify-between'>
        <h3 className='font-semibold text-gray-800'>Open alerts</h3>
        {open.length > 0 && (
          <span className='bg-red-100 text-red-600 text-xs font-medium px-2 py-0.5 rounded-full'>
            {open.length}
          </span>
        )}
      </div>

      {open.length === 0 ? (
        <div className='p-6 text-center text-gray-400 text-sm'>No open alerts.</div>
      ) : (
        <ul className='divide-y divide-gray-50'>
          {open.map((alert) => (
            <li key={alert.id} className='px-4 py-3 flex items-center justify-between'>
              <div>
                <p className='text-sm font-medium text-gray-800'>{alert.type.replace('_', ' ')}</p>
                <p className='text-xs text-gray-400 mt-0.5'>Vehicle {alert.vehicle_id.slice(0, 8)}</p>
              </div>
              <span className='text-xs text-red-500 font-medium'>
                {formatDuration(alert.started_at)}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
