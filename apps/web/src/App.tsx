import { lazy, Suspense } from 'react'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'

const VehicleList = lazy(() => import('./features/vehicles/vehicle-list'))

function LoadingFallback() {
  return (
    <div className='min-h-screen bg-gray-50 flex items-center justify-center'>
      <div className='w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full animate-spin' />
    </div>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingFallback />}>
        <Routes>
          <Route path='/vehicles' element={<VehicleList />} />
          <Route path='*' element={<Navigate to='/vehicles' replace />} />
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
