import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'

import { RouteHistory } from './route-history'
import { useVehicleDetail, useVehicleTelemetry } from './use-vehicle-detail'

function formatDate(iso: string) {
  return new Date(iso).toLocaleString()
}

export default function VehicleDetail() {
  const { id } = useParams<{ id: string }>()
  const vehicleId = id ?? ''

  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')

  const { data: vehicle, isLoading: loadingVehicle, isError: vehicleError } = useVehicleDetail(vehicleId)
  const { data: telemetry, isLoading: loadingTelemetry } = useVehicleTelemetry(vehicleId, from || undefined, to || undefined)

  const points = telemetry?.points ?? []

  if (loadingVehicle) {
    return (
      <div className='min-h-screen bg-gray-50 flex items-center justify-center'>
        <div className='w-8 h-8 border-2 border-blue-500 border-t-transparent rounded-full animate-spin' />
      </div>
    )
  }

  if (vehicleError || !vehicle) {
    return (
      <div className='min-h-screen bg-gray-50 flex items-center justify-center'>
        <p className='text-gray-500'>Vehicle not found.</p>
      </div>
    )
  }

  return (
    <div className='min-h-screen bg-gray-50'>
      <header className='bg-white border-b border-gray-100 shadow-sm'>
        <div className='max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 h-14 flex items-center gap-3'>
          <Link to='/dashboard' className='text-sm text-blue-500 hover:text-blue-700'>
            ← Dashboard
          </Link>
          <span className='text-gray-300'>/</span>
          <h1 className='text-lg font-semibold text-gray-800'>{vehicle.name}</h1>
        </div>
      </header>

      <main className='max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-6'>
        <div className='bg-white rounded-xl border border-gray-100 shadow-sm p-5 grid grid-cols-2 gap-4 sm:grid-cols-4'>
          <div>
            <p className='text-xs text-gray-400'>Plate</p>
            <p className='font-semibold text-gray-800 font-mono'>{vehicle.plate}</p>
          </div>
          <div>
            <p className='text-xs text-gray-400'>External ID</p>
            <p className='font-semibold text-gray-800 font-mono text-sm'>{vehicle.external_id}</p>
          </div>
          <div>
            <p className='text-xs text-gray-400'>Created</p>
            <p className='text-sm text-gray-600'>{formatDate(vehicle.created_at)}</p>
          </div>
          <div>
            <p className='text-xs text-gray-400'>Status</p>
            <p className='text-sm font-medium'>
              {vehicle.deleted_at ? (
                <span className='text-red-500'>Deleted</span>
              ) : (
                <span className='text-green-600'>Active</span>
              )}
            </p>
          </div>
        </div>

        <div className='bg-white rounded-xl border border-gray-100 shadow-sm p-5 space-y-4'>
          <h2 className='font-semibold text-gray-800'>Route history</h2>
          <div className='flex flex-wrap gap-3'>
            <div>
              <label className='block text-xs text-gray-400 mb-1'>From</label>
              <input
                type='datetime-local'
                value={from}
                onChange={(e) => setFrom(e.target.value)}
                className='border border-gray-200 rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-300'
              />
            </div>
            <div>
              <label className='block text-xs text-gray-400 mb-1'>To</label>
              <input
                type='datetime-local'
                value={to}
                onChange={(e) => setTo(e.target.value)}
                className='border border-gray-200 rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-300'
              />
            </div>
          </div>

          {loadingTelemetry ? (
            <div className='flex justify-center py-10'>
              <div className='w-6 h-6 border-2 border-blue-500 border-t-transparent rounded-full animate-spin' />
            </div>
          ) : (
            <>
              <p className='text-xs text-gray-400'>{points.length} points</p>
              <RouteHistory points={points} />
            </>
          )}
        </div>
      </main>
    </div>
  )
}
