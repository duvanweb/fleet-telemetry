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
		{4.7110, -74.0721}, {4.7125, -74.0710}, {4.7140, -74.0698}, {4.7155, -74.0685},
		{4.7170, -74.0672}, {4.7185, -74.0659}, {4.7200, -74.0645}, {4.7215, -74.0632},
		{4.7230, -74.0618}, {4.7245, -74.0605}, {4.7260, -74.0592}, {4.7245, -74.0578},
		{4.7230, -74.0565}, {4.7215, -74.0551}, {4.7200, -74.0538}, {4.7185, -74.0524},
		{4.7170, -74.0511}, {4.7155, -74.0497}, {4.7140, -74.0484}, {4.7125, -74.0470},
	},
}

// Stopped emits the same coordinate for 8 consecutive points to trigger VEHICLE_STOPPED alert.
var Stopped = domain.Scenario{
	Name:        "stopped",
	Description: "Vehicle stops at a fixed point for over 1 minute, triggering a stopped alert.",
	Points: []domain.Coordinate{
		{7.1193, -73.1227}, {7.1193, -73.1227}, {7.1193, -73.1227}, {7.1193, -73.1227},
		{7.1193, -73.1227}, {7.1193, -73.1227}, {7.1193, -73.1227}, {7.1193, -73.1227},
	},
}

// Moving is a fast-moving vehicle that resolves a STOPPED alert.
var Moving = domain.Scenario{
	Name:        "moving",
	Description: "Fast-moving vehicle that resolves a stopped alert by changing position.",
	Points: []domain.Coordinate{
		{7.1193, -73.1227}, {7.1250, -73.1180}, {7.1310, -73.1130}, {7.1370, -73.1080},
		{7.1430, -73.1030}, {7.1490, -73.0980}, {7.1550, -73.0930}, {7.1610, -73.0880},
		{7.1670, -73.0830}, {7.1730, -73.0780},
	},
}

// Duplicate emits some points twice to test deduplication logic.
var Duplicate = domain.Scenario{
	Name:        "duplicate",
	Description: "Sends duplicate telemetry points to test deduplication.",
	Points: []domain.Coordinate{
		{4.6500, -74.1000}, {4.6520, -74.0990}, {4.6540, -74.0980}, {4.6560, -74.0970},
		{4.6580, -74.0960}, {4.6600, -74.0950}, {4.6620, -74.0940}, {4.6640, -74.0930},
	},
}

// Invalid emits some invalid coordinates to test validation.
var Invalid = domain.Scenario{
	Name:        "invalid",
	Description: "Sends invalid coordinates to test server-side validation.",
	Points: []domain.Coordinate{
		{4.7000, -74.0700}, {200.0, -74.0700}, {4.7020, -74.0690}, {4.7040, 250.0},
		{4.7060, -74.0670}, {4.7080, -74.0660}, {4.7100, -74.0650}, {4.7120, -74.0640},
	},
}
