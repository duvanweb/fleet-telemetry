package vehicle_client

// Configuration holds the vehicle-service HTTP client settings.
type Configuration struct {
	VehicleServiceURL string `env:"VEHICLE_SERVICE_URL" envDefault:"http://localhost:8081"`
}
