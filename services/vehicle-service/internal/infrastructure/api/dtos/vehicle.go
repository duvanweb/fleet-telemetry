package dtos

// CreateVehicleRequest is the request body for creating a vehicle.
type CreateVehicleRequest struct {
	ExternalID string `json:"external_id"`
	Plate      string `json:"plate"`
	Name       string `json:"name"`
}

// UpdateVehicleRequest is the request body for updating a vehicle.
type UpdateVehicleRequest struct {
	Plate string `json:"plate"`
	Name  string `json:"name"`
}

// VehicleResponse is the JSON representation of a vehicle.
type VehicleResponse struct {
	ID         string  `json:"id"`
	ExternalID string  `json:"external_id"`
	Plate      string  `json:"plate"`
	Name       string  `json:"name"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at,omitempty"`
}
