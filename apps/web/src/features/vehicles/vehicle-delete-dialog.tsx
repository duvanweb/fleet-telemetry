interface VehicleDeleteDialogProps {
  vehicleName: string
  onConfirm: () => void
  onCancel: () => void
  isLoading?: boolean
}

export function VehicleDeleteDialog({
  vehicleName,
  onConfirm,
  onCancel,
  isLoading,
}: VehicleDeleteDialogProps) {
  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center'>
      <div className='absolute inset-0 bg-black/50' onClick={onCancel} />
      <div className='relative bg-white rounded-lg shadow-xl w-full max-w-md mx-4 p-6'>
        <h2 className='text-lg font-semibold text-gray-900 mb-2'>Delete vehicle</h2>
        <p className='text-sm text-gray-600 mb-6'>
          Are you sure you want to delete{' '}
          <span className='font-medium text-gray-900'>{vehicleName}</span>? This action cannot be
          undone.
        </p>
        <div className='flex justify-end gap-3'>
          <button
            onClick={onCancel}
            className='px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 transition-colors'
          >
            Cancel
          </button>
          <button
            onClick={onConfirm}
            disabled={isLoading}
            className='px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors'
          >
            {isLoading ? 'Deleting…' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  )
}
