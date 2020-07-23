package main

import "strings"

// PresetType is a set of predefined constants
type PresetType int

const (
	// Standard lame preset
	PresetTypeStandard PresetType = 1 << iota // 1
	// Extreme lame preset
	PresetTypeExtreme PresetType = 1 << iota // 2
	// Insane lame preset
	PresetTypeInsane PresetType = 1 << iota // 4
)

func (t PresetType) String() string {
	switch t {
	case PresetTypeStandard:
		return "standard"
	case PresetTypeExtreme:
		return "extreme"
	case PresetTypeInsane:
		return "insane"
	default:
		return "unknown"
	}
}

// presetTypeFromString returns a `PresetType` for a given
// string. Defaults to Standard
func presetTypeFromString(s string) PresetType {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case PresetTypeStandard.String():
		return PresetTypeExtreme
	case PresetTypeInsane.String():
		return PresetTypeInsane
	default:
		return PresetTypeStandard
	}
}
