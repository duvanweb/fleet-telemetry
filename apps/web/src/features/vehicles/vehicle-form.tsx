import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

import type { Vehicle } from './types'

const schema = z.object({
  external_id: z.string().min(1, 'External ID is required'),
  plate: z.string().min(1, 'Plate is required').max(50, 'Plate must be 50 characters or less'),
  name: z.string().min(1, 'Name is required').max(255, 'Name must be 255 characters or less'),
})

type FormValues = z.infer<typeof schema>

interface VehicleFormProps {
  vehicle?: Vehicle
  onSubmit: (values: FormValues) => void
  onCancel: () => void
  isLoading?: boolean
}

export function VehicleForm({ vehicle, onSubmit, onCancel, isLoading }: VehicleFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: vehicle
      ? { external_id: vehicle.external_id, plate: vehicle.plate, name: vehicle.name }
      : undefined,
  })

  const isEdit = Boolean(vehicle)

  return (
    <form onSubmit={handleSubmit(onSubmit)} className='space-y-4'>
      {!isEdit && (
        <div>
          <label className='block text-sm font-medium text-gray-700 mb-1'>External ID</label>
          <input
            {...register('external_id')}
            className='w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent'
            placeholder='e.g. VEH-001'
          />
          {errors.external_id && (
            <p className='mt-1 text-xs text-red-600'>{errors.external_id.message}</p>
          )}
        </div>
      )}

      <div>
        <label className='block text-sm font-medium text-gray-700 mb-1'>Plate</label>
        <input
          {...register('plate')}
          className='w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent'
          placeholder='e.g. ABC-123'
        />
        {errors.plate && <p className='mt-1 text-xs text-red-600'>{errors.plate.message}</p>}
      </div>

      <div>
        <label className='block text-sm font-medium text-gray-700 mb-1'>Name</label>
        <input
          {...register('name')}
          className='w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent'
          placeholder='e.g. Delivery Van 1'
        />
        {errors.name && <p className='mt-1 text-xs text-red-600'>{errors.name.message}</p>}
      </div>

      <div className='flex justify-end gap-3 pt-2'>
        <button
          type='button'
          onClick={onCancel}
          className='px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 transition-colors'
        >
          Cancel
        </button>
        <button
          type='submit'
          disabled={isLoading}
          className='px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors'
        >
          {isLoading ? 'Saving…' : isEdit ? 'Update vehicle' : 'Create vehicle'}
        </button>
      </div>
    </form>
  )
}
