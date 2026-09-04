import { MapContainer, Polyline, TileLayer, CircleMarker, Tooltip } from 'react-leaflet'
import 'leaflet/dist/leaflet.css'
import type { TelemetryPoint } from './telemetry.service'

interface Props {
  points: TelemetryPoint[]
}

export function RouteHistory({ points }: Props) {
  if (points.length === 0) {
    return (
      <div className='rounded-xl border border-gray-100 bg-white p-6 text-center text-sm text-gray-400'>
        No route history available.
      </div>
    )
  }

  const positions = points.map((p) => [p.latitude, p.longitude] as [number, number])
  const first = positions[0]
  const last = positions[positions.length - 1]
  const center: [number, number] = [
    positions.reduce((s, p) => s + p[0], 0) / positions.length,
    positions.reduce((s, p) => s + p[1], 0) / positions.length,
  ]

  return (
    <div className='rounded-xl border border-gray-100 overflow-hidden shadow-sm' style={{ height: 360 }}>
      <MapContainer center={center} zoom={13} style={{ height: '100%', width: '100%' }}>
        <TileLayer
          url='https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
          attribution='&copy; OpenStreetMap contributors'
        />
        <Polyline positions={positions} color='#3b82f6' weight={3} opacity={0.8} />
        <CircleMarker center={first} radius={7} color='#22c55e' fillColor='#22c55e' fillOpacity={1}>
          <Tooltip permanent>Start</Tooltip>
        </CircleMarker>
        <CircleMarker center={last} radius={7} color='#ef4444' fillColor='#ef4444' fillOpacity={1}>
          <Tooltip permanent>Last</Tooltip>
        </CircleMarker>
      </MapContainer>
    </div>
  )
}
