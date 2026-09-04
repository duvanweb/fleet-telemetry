export interface Vehicle {
  id: string
  external_id: string
  plate: string
  name: string
  created_at: string
  updated_at: string
  deleted_at?: string | null
}

export interface CreateVehicleRequest {
  external_id: string
  plate: string
  name: string
}

export interface UpdateVehicleRequest {
  plate: string
  name: string
}
