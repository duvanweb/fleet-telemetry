import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { SimulatorControls } from './simulator-controls'
import { SimulatorLog } from './simulator-log'
import {
  useSimulatorScenarios,
  useSimulatorStatus,
  useStartSimulation,
  useStopSimulation,
} from './use-simulator'
import type { LogEntry, StartSimulationRequest } from './types'
import { useDashboardVehicles } from '../dashboard/use-dashboard'

const MAX_LOG_ENTRIES = 200

const DEFAULT_FORM: StartSimulationRequest = {
  vehicle_ids: [],
  scenario: 'normal',
  interval_ms: 5000,
  duplicate_rate: 0,
  invalid_rate: 0,
}

function makeEntry(type: LogEntry['type'], message: string): LogEntry {
  return { id: `${Date.now()}-${Math.random()}`, timestamp: new Date().toISOString(), type, message }
}

export default function Simulator() {
  const { data: scenariosData, isLoading: loadingScenarios } = useSimulatorScenarios()
  const { data: status } = useSimulatorStatus()
  const { data: vehiclesData } = useDashboardVehicles()

  const startMutation = useStartSimulation()
  const stopMutation = useStopSimulation()

  const [form, setForm] = useState<StartSimulationRequest>(DEFAULT_FORM)
  const [logs, setLogs] = useState<LogEntry[]>([])

  const scenarios = scenariosData?.scenarios ?? []
  const vehicles = (vehiclesData ?? []).filter((v) => !v.deleted_at)
  const running = status?.running ?? false

  useEffect(() => {
    if (scenarios.length > 0 && !form.scenario) {
      setForm((p) => ({ ...p, scenario: scenarios[0].name }))
    }
  }, [scenarios, form.scenario])

  const addLog = useCallback((entry: LogEntry) => {
    setLogs((prev) => {
      const next = [...prev, entry]
      return next.length > MAX_LOG_ENTRIES ? next.slice(next.length - MAX_LOG_ENTRIES) : next
    })
  }, [])

  function handleStart() {
    startMutation.mutate(form, {
      onSuccess: () => addLog(makeEntry('STARTED', `Scenario "${form.scenario}" started with ${form.vehicle_ids.length} vehicle(s)`)),
      onError: (err) => addLog(makeEntry('ERROR', err.message)),
    })
  }

  function handleStop() {
    stopMutation.mutate(undefined, {
      onSuccess: () => addLog(makeEntry('STOPPED', 'Simulation stopped')),
      onError: (err) => addLog(makeEntry('ERROR', err.message)),
    })
  }

  return (
    <div className='min-h-screen bg-gray-50'>
      <header className='bg-white border-b border-gray-100 shadow-sm'>
        <div className='max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 h-14 flex items-center gap-3'>
          <Link to='/dashboard' className='text-sm text-blue-500 hover:text-blue-700'>
            ← Dashboard
          </Link>
          <span className='text-gray-300'>/</span>
          <h1 className='text-lg font-semibold text-gray-800'>Simulator</h1>
          {running && (
            <span className='ml-auto flex items-center gap-1.5 text-xs text-green-600 font-medium'>
              <span className='w-2 h-2 rounded-full bg-green-500 animate-pulse' />
              Running — {status?.scenario} ({status?.vehicle_count} vehicle{status?.vehicle_count !== 1 ? 's' : ''})
            </span>
          )}
        </div>
      </header>

      <main className='max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-4'>
        {loadingScenarios ? (
          <div className='flex justify-center py-16'>
            <div className='w-8 h-8 border-2 border-blue-500 border-t-transparent rounded-full animate-spin' />
          </div>
        ) : (
          <SimulatorControls
            scenarios={scenarios}
            form={form}
            onChange={setForm}
            running={running}
            onStart={handleStart}
            onStop={handleStop}
            vehicles={vehicles.map((v) => ({ id: v.id, name: v.name }))}
          />
        )}

        <div>
          <h2 className='text-sm font-semibold text-gray-500 uppercase tracking-wide mb-2'>
            Event log
          </h2>
          <SimulatorLog entries={logs} />
        </div>
      </main>
    </div>
  )
}
