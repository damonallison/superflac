package main

// PresetType is a set of predefined lame present constants
type PresetType string

const (
	// PresetTypeStandard is the standard lame preset
	//
	// This preset should generally be transparent to most people on most music
	// and is already quite high in quality
	PresetTypeStandard = PresetType("standard")

	// PresetTypeExtreme is the extreme lame preset
	//
	// If you have extremely good hearing and similar equipment, this preset
	// will generally provide slightly higher quality than the "standard" mode.
	PresetTypeExtreme = PresetType("extreme")

	// PresetTypeInsane is the insane  lame preset
	//
	// This preset will usually be overkill for most people and most situations,
	// but if you must have the absolute highest quality with no regard to
	// filesize, this is the way to go.
	PresetTypeInsane = PresetType("insane")
)
