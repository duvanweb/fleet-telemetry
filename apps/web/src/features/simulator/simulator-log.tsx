import { useEffect, useRef } from 'react'

import type { LogEntry, LogEntryType } from './types'

interface Props {
  entries: LogEntry[]
}

const TYPE_STYLES: Record<LogEntryType, string> = {
  GPS_SENT: 'text-green-400',
  DUPLICATE: 'text-yellow-400',
  INVALID_PAYLOAD: 'text-red-400',
  ERROR: 'text-red-500',
  STARTED: 'text-blue-400',
  STOPPED: 'text-gray-400',
}

export function SimulatorLog({ entries }: Props) {
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [entries.length])

  return (
    <div className='bg-gray-900 rounded-xl font-mono text-xs overflow-y-auto' style={{ height: 280 }}>
      {entries.length === 0 ? (
        <p className='text-gray-600 p-4'>No events yet. Start a simulation to see activity.</p>
      ) : (
        <div className='p-3 space-y-0.5'>
          {entries.map((entry) => (
            <div key={entry.id} className='flex gap-3'>
              <span className='text-gray-600 shrink-0'>
                {new Date(entry.timestamp).toLocaleTimeString()}
              </span>
              <span className={`shrink-0 w-24 ${TYPE_STYLES[entry.type]}`}>{entry.type}</span>
              <span className='text-gray-300 truncate'>{entry.message}</span>
            </div>
          ))}
          <div ref={bottomRef} />
        </div>
      )}
    </div>
  )
}
