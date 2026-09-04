package scenario

import "fleet/simulator/internal/core/domain"

// All returns the list of all available simulation scenarios.
func All() []domain.Scenario {
	return []domain.Scenario{Normal, Stopped, Moving, Duplicate, Invalid}
}

// ByName returns the scenario with the given name and whether it was found.
func ByName(name string) (domain.Scenario, bool) {
	for _, s := range All() {
		if s.Name == name {
			return s, true
		}
	}
	return domain.Scenario{}, false
}

// Normal is a typical city driving route through Bogotá.
var Normal = domain.Scenario{
	Name:        "normal",
	Description: "Typical city driving route with 20 waypoints.",
	Points: []domain.Coordinate{
		{Latitude: 4.7110, Longitude: -74.0721}, {Latitude: 4.7125, Longitude: -74.0710},
		{Latitude: 4.7140, Longitude: -74.0698}, {Latitude: 4.7155, Longitude: -74.0685},
		{Latitude: 4.7170, Longitude: -74.0672}, {Latitude: 4.7185, Longitude: -74.0659},
		{Latitude: 4.7200, Longitude: -74.0645}, {Latitude: 4.7215, Longitude: -74.0632},
		{Latitude: 4.7230, Longitude: -74.0618}, {Latitude: 4.7245, Longitude: -74.0605},
		{Latitude: 4.7260, Longitude: -74.0592}, {Latitude: 4.7245, Longitude: -74.0578},
		{Latitude: 4.7230, Longitude: -74.0565}, {Latitude: 4.7215, Longitude: -74.0551},
		{Latitude: 4.7200, Longitude: -74.0538}, {Latitude: 4.7185, Longitude: -74.0524},
		{Latitude: 4.7170, Longitude: -74.0511}, {Latitude: 4.7155, Longitude: -74.0497},
		{Latitude: 4.7140, Longitude: -74.0484}, {Latitude: 4.7125, Longitude: -74.0470},
	},
}

// Stopped emits the same coordinate for 8 consecutive points to trigger VEHICLE_STOPPED alert.
var Stopped = domain.Scenario{
	Name:        "stopped",
	Description: "Vehicle stops at a fixed point for over 1 minute, triggering a stopped alert.",
	Points: []domain.Coordinate{
		{Latitude: 7.1193, Longitude: -73.1227}, {Latitude: 7.1193, Longitude: -73.1227},
		{Latitude: 7.1193, Longitude: -73.1227}, {Latitude: 7.1193, Longitude: -73.1227},
		{Latitude: 7.1193, Longitude: -73.1227}, {Latitude: 7.1193, Longitude: -73.1227},
		{Latitude: 7.1193, Longitude: -73.1227}, {Latitude: 7.1193, Longitude: -73.1227},
	},
}

// Moving is a fast-moving vehicle that resolves a STOPPED alert.
var Moving = domain.Scenario{
	Name:        "moving",
	Description: "Fast-moving vehicle that resolves a stopped alert by changing position.",
	Points: []domain.Coordinate{
		{Latitude: 7.1193, Longitude: -73.1227}, {Latitude: 7.1250, Longitude: -73.1180},
		{Latitude: 7.1310, Longitude: -73.1130}, {Latitude: 7.1370, Longitude: -73.1080},
		{Latitude: 7.1430, Longitude: -73.1030}, {Latitude: 7.1490, Longitude: -73.0980},
		{Latitude: 7.1550, Longitude: -73.0930}, {Latitude: 7.1610, Longitude: -73.0880},
		{Latitude: 7.1670, Longitude: -73.0830}, {Latitude: 7.1730, Longitude: -73.0780},
	},
}

// Duplicate emits some points twice to test deduplication logic.
var Duplicate = domain.Scenario{
	Name:        "duplicate",
	Description: "Sends duplicate telemetry points to test deduplication.",
	Points: []domain.Coordinate{
		{Latitude: 4.6500, Longitude: -74.1000}, {Latitude: 4.6520, Longitude: -74.0990},
		{Latitude: 4.6540, Longitude: -74.0980}, {Latitude: 4.6560, Longitude: -74.0970},
		{Latitude: 4.6580, Longitude: -74.0960}, {Latitude: 4.6600, Longitude: -74.0950},
		{Latitude: 4.6620, Longitude: -74.0940}, {Latitude: 4.6640, Longitude: -74.0930},
	},
}

// Invalid emits some invalid coordinates to test validation.
var Invalid = domain.Scenario{
	Name:        "invalid",
	Description: "Sends invalid coordinates to test server-side validation.",
	Points: []domain.Coordinate{
		{Latitude: 4.7000, Longitude: -74.0700}, {Latitude: 200.0, Longitude: -74.0700},
		{Latitude: 4.7020, Longitude: -74.0690}, {Latitude: 4.7040, Longitude: 250.0},
		{Latitude: 4.7060, Longitude: -74.0670}, {Latitude: 4.7080, Longitude: -74.0660},
		{Latitude: 4.7100, Longitude: -74.0650}, {Latitude: 4.7120, Longitude: -74.0640},
	},
}
