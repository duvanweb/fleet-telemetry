import { lazy, Suspense } from 'react'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'

const Dashboard = lazy(() => import('./features/dashboard/dashboard'))
const VehicleDetail = lazy(() => import('./features/vehicles/vehicle-detail'))
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
          <Route path='/dashboard' element={<Dashboard />} />
          <Route path='/vehicles' element={<VehicleList />} />
          <Route path='/vehicles/:id' element={<VehicleDetail />} />
          <Route path='*' element={<Navigate to='/dashboard' replace />} />
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
