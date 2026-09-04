import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet'
import { divIcon } from 'leaflet'
import 'leaflet/dist/leaflet.css'

import type { Vehicle } from '../vehicles/types'
import type { VehicleState, VehicleStatus } from './types'

interface Props {
  vehicles: Vehicle[]
  vehicleStates: Map<string, VehicleState>
}

const MARKER_COLORS: Record<VehicleStatus, string> = {
  MOVING: '#22c55e',
  STOPPED: '#eab308',
  ALERT: '#ef4444',
  UNKNOWN: '#9ca3af',
}

function vehicleMarkerIcon(status: VehicleStatus) {
  const color = MARKER_COLORS[status]
  return divIcon({
    html: `<div style="width:14px;height:14px;border-radius:50%;background:${color};border:2px solid #fff;box-shadow:0 1px 4px rgba(0,0,0,.3)"></div>`,
    className: '',
    iconSize: [14, 14],
    iconAnchor: [7, 7],
  })
}

const DEFAULT_CENTER: [number, number] = [4.711, -74.0721]

export function FleetMap({ vehicles, vehicleStates }: Props) {
  const active = vehicles.filter((v) => !v.deleted_at)

  const positioned = active.filter((v) => {
    const s = vehicleStates.get(v.id)
    return s?.lastLatitude != null && s?.lastLongitude != null
  })

  return (
    <div className='rounded-xl border border-gray-100 overflow-hidden shadow-sm' style={{ height: 400 }}>
      <MapContainer center={DEFAULT_CENTER} zoom={12} style={{ height: '100%', width: '100%' }}>
        <TileLayer
          url='https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
          attribution='&copy; OpenStreetMap contributors'
        />
        {positioned.map((v) => {
          const state = vehicleStates.get(v.id)!
          const status = state.status
          const pos: [number, number] = [state.lastLatitude!, state.lastLongitude!]
          return (
            <Marker key={v.id} position={pos} icon={vehicleMarkerIcon(status)}>
              <Popup>
                <p className='font-semibold'>{v.name}</p>
                <p className='text-xs text-gray-500'>{v.plate}</p>
                <p className='text-xs mt-1 capitalize'>{status.toLowerCase()}</p>
              </Popup>
            </Marker>
          )
        })}
      </MapContainer>
    </div>
  )
}
