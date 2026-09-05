import type { Dispatch, SetStateAction } from 'react'

import type { Scenario, StartSimulationRequest } from './types'

interface Props {
  scenarios: Scenario[]
  form: StartSimulationRequest
  onChange: Dispatch<SetStateAction<StartSimulationRequest>>
  running: boolean
  onStart: () => void
  onStop: () => void
  vehicles: { id: string; name: string }[]
}

export function SimulatorControls({ scenarios, form, onChange, running, onStart, onStop, vehicles }: Props) {
  function toggleVehicle(id: string) {
    onChange((prev) => ({
      ...prev,
      vehicle_ids: prev.vehicle_ids.includes(id)
        ? prev.vehicle_ids.filter((v) => v !== id)
        : [...prev.vehicle_ids, id],
    }))
  }

  return (
    <div className='bg-white rounded-xl border border-gray-100 shadow-sm p-5 space-y-4'>
      <div>
        <label className='block text-xs text-gray-400 mb-1'>Scenario</label>
        <select
          value={form.scenario}
          onChange={(e) => onChange((p) => ({ ...p, scenario: e.target.value }))}
          disabled={running}
          className='border border-gray-200 rounded-lg px-3 py-2 text-sm w-full focus:outline-none focus:ring-2 focus:ring-blue-300 disabled:opacity-50'
        >
          {scenarios.map((s) => (
            <option key={s.name} value={s.name}>
              {s.name} — {s.description}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className='block text-xs text-gray-400 mb-1'>Vehicles</label>
        <div className='flex flex-wrap gap-2'>
          {vehicles.map((v) => (
            <button
              key={v.id}
              type='button'
              disabled={running}
              onClick={() => toggleVehicle(v.id)}
              className={`px-3 py-1 rounded-full text-xs border transition-colors disabled:opacity-50 ${
                form.vehicle_ids.includes(v.id)
                  ? 'bg-blue-500 text-white border-blue-500'
                  : 'bg-white text-gray-600 border-gray-200 hover:border-blue-300'
              }`}
            >
              {v.name}
            </button>
          ))}
          {vehicles.length === 0 && (
            <p className='text-xs text-gray-400'>No vehicles available. Create some first.</p>
          )}
        </div>
      </div>

      <div className='grid grid-cols-1 sm:grid-cols-3 gap-4'>
        <div>
          <label className='block text-xs text-gray-400 mb-1'>
            Interval: {form.interval_ms}ms
          </label>
          <input
            type='range'
            min={1000}
            max={30000}
            step={500}
            value={form.interval_ms}
            disabled={running}
            onChange={(e) => onChange((p) => ({ ...p, interval_ms: Number(e.target.value) }))}
            className='w-full accent-blue-500 disabled:opacity-50'
          />
        </div>

        <div>
          <label className='block text-xs text-gray-400 mb-1'>
            Duplicate rate: {(form.duplicate_rate * 100).toFixed(0)}%
          </label>
          <input
            type='range'
            min={0}
            max={1}
            step={0.05}
            value={form.duplicate_rate}
            disabled={running}
            onChange={(e) => onChange((p) => ({ ...p, duplicate_rate: Number(e.target.value) }))}
            className='w-full accent-yellow-400 disabled:opacity-50'
          />
        </div>

        <div>
          <label className='block text-xs text-gray-400 mb-1'>
            Invalid rate: {(form.invalid_rate * 100).toFixed(0)}%
          </label>
          <input
            type='range'
            min={0}
            max={1}
            step={0.05}
            value={form.invalid_rate}
            disabled={running}
            onChange={(e) => onChange((p) => ({ ...p, invalid_rate: Number(e.target.value) }))}
            className='w-full accent-red-400 disabled:opacity-50'
          />
        </div>
      </div>

      <div className='flex gap-3'>
        <button
          type='button'
          onClick={onStart}
          disabled={running || form.vehicle_ids.length === 0}
          className='flex-1 bg-green-500 hover:bg-green-600 text-white rounded-lg py-2 text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed'
        >
          Start
        </button>
        <button
          type='button'
          onClick={onStop}
          disabled={!running}
          className='flex-1 bg-red-500 hover:bg-red-600 text-white rounded-lg py-2 text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed'
        >
          Stop
        </button>
      </div>
    </div>
  )
}
