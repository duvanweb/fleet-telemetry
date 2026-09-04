import { useState } from 'react'

import type { Vehicle } from './types'
import { VehicleDeleteDialog } from './vehicle-delete-dialog'
import { VehicleForm } from './vehicle-form'
import { useCreateVehicle, useDeleteVehicle, useUpdateVehicle, useVehicles } from './use-vehicles'

type Modal =
  | { type: 'create' }
  | { type: 'edit'; vehicle: Vehicle }
  | { type: 'delete'; vehicle: Vehicle }
  | null

function TableSkeleton() {
  return (
    <tbody>
      {Array.from({ length: 4 }).map((_, i) => (
        <tr key={i} className='border-b border-gray-100'>
          {Array.from({ length: 4 }).map((_, j) => (
            <td key={j} className='px-4 py-3'>
              <div className='h-4 bg-gray-200 rounded animate-pulse' />
            </td>
          ))}
        </tr>
      ))}
    </tbody>
  )
}

export default function VehicleList() {
  const { data: vehicles, isLoading, isError, error } = useVehicles()
  const createVehicle = useCreateVehicle()
  const updateVehicle = useUpdateVehicle()
  const deleteVehicle = useDeleteVehicle()

  const [modal, setModal] = useState<Modal>(null)
  const [apiError, setApiError] = useState<string | null>(null)

  function handleClose() {
    setModal(null)
    setApiError(null)
  }

  function handleCreate(values: { external_id: string; plate: string; name: string }) {
    createVehicle.mutate(values, {
      onSuccess: handleClose,
      onError: (err) => setApiError(err.message),
    })
  }

  function handleUpdate(values: { plate: string; name: string; external_id: string }) {
    if (modal?.type !== 'edit') return
    updateVehicle.mutate(
      { id: modal.vehicle.id, data: { plate: values.plate, name: values.name } },
      { onSuccess: handleClose, onError: (err) => setApiError(err.message) },
    )
  }

  function handleDelete() {
    if (modal?.type !== 'delete') return
    deleteVehicle.mutate(modal.vehicle.id, {
      onSuccess: handleClose,
      onError: (err) => setApiError(err.message),
    })
  }

  const activeVehicles = vehicles?.filter((v) => !v.deleted_at) ?? []

  return (
    <div className='min-h-screen bg-gray-50 py-8'>
      <div className='max-w-5xl mx-auto px-4'>
        <div className='flex items-center justify-between mb-6'>
          <div>
            <h1 className='text-2xl font-bold text-gray-900'>Vehicles</h1>
            <p className='text-sm text-gray-500 mt-1'>Manage your fleet vehicles</p>
          </div>
          <button
            onClick={() => setModal({ type: 'create' })}
            className='px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors'
          >
            + Add vehicle
          </button>
        </div>

        {isError && (
          <div className='mb-4 p-4 bg-red-50 border border-red-200 rounded-md'>
            <p className='text-sm text-red-700'>
              Failed to load vehicles: {(error as Error).message}
            </p>
          </div>
        )}

        <div className='bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden'>
          <table className='w-full text-sm'>
            <thead>
              <tr className='bg-gray-50 border-b border-gray-200'>
                <th className='px-4 py-3 text-left font-medium text-gray-600'>Name</th>
                <th className='px-4 py-3 text-left font-medium text-gray-600'>Plate</th>
                <th className='px-4 py-3 text-left font-medium text-gray-600'>External ID</th>
                <th className='px-4 py-3 text-left font-medium text-gray-600'>Actions</th>
              </tr>
            </thead>

            {isLoading ? (
              <TableSkeleton />
            ) : activeVehicles.length === 0 ? (
              <tbody>
                <tr>
                  <td colSpan={4} className='px-4 py-12 text-center text-gray-400'>
                    No vehicles yet. Add your first vehicle to get started.
                  </td>
                </tr>
              </tbody>
            ) : (
              <tbody>
                {activeVehicles.map((vehicle) => (
                  <tr
                    key={vehicle.id}
                    className='border-b border-gray-100 last:border-0 hover:bg-gray-50 transition-colors'
                  >
                    <td className='px-4 py-3 font-medium text-gray-900'>{vehicle.name}</td>
                    <td className='px-4 py-3 text-gray-600'>
                      <span className='inline-block px-2 py-0.5 bg-gray-100 rounded text-xs font-mono'>
                        {vehicle.plate}
                      </span>
                    </td>
                    <td className='px-4 py-3 text-gray-500 font-mono text-xs'>
                      {vehicle.external_id}
                    </td>
                    <td className='px-4 py-3'>
                      <div className='flex gap-2'>
                        <button
                          onClick={() => setModal({ type: 'edit', vehicle })}
                          className='text-xs text-blue-600 hover:text-blue-800 font-medium transition-colors'
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => setModal({ type: 'delete', vehicle })}
                          className='text-xs text-red-500 hover:text-red-700 font-medium transition-colors'
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            )}
          </table>
        </div>
      </div>

      {(modal?.type === 'create' || modal?.type === 'edit') && (
        <div className='fixed inset-0 z-50 flex items-center justify-center'>
          <div className='absolute inset-0 bg-black/50' onClick={handleClose} />
          <div className='relative bg-white rounded-lg shadow-xl w-full max-w-md mx-4 p-6'>
            <h2 className='text-lg font-semibold text-gray-900 mb-4'>
              {modal.type === 'create' ? 'Add vehicle' : 'Edit vehicle'}
            </h2>
            {apiError && (
              <div className='mb-4 p-3 bg-red-50 border border-red-200 rounded-md'>
                <p className='text-xs text-red-700'>{apiError}</p>
              </div>
            )}
            <VehicleForm
              vehicle={modal.type === 'edit' ? modal.vehicle : undefined}
              onSubmit={modal.type === 'create' ? handleCreate : handleUpdate}
              onCancel={handleClose}
              isLoading={createVehicle.isPending || updateVehicle.isPending}
            />
          </div>
        </div>
      )}

      {modal?.type === 'delete' && (
        <VehicleDeleteDialog
          vehicleName={modal.vehicle.name}
          onConfirm={handleDelete}
          onCancel={handleClose}
          isLoading={deleteVehicle.isPending}
        />
      )}
    </div>
  )
}
