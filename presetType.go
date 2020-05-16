package main

import "strings"

// PresetType is a set of predefined constants
type PresetType int

const (
	// Standard lame preset
	Standard PresetType = 1 << iota // 1
	// Extreme lame preset
	Extreme PresetType = 1 << iota // 2
	// Insane lame preset
	Insane PresetType = 1 << iota // 4
)

func (t PresetType) String() string {
	switch t {
	case Standard:
		return "standard"
	case Extreme:
		return "extreme"
	case Insane:
		return "insane"
	default:
		return "unknown"
	}
}

// presetTypeFromString returns a `PresetType` for a given
// string. Defaults to Standard
func presetTypeFromString(s string) PresetType {
	temp := strings.ToLower(strings.TrimSpace(s))
	switch temp {
	case Extreme.String():
		return Extreme
	case Insane.String():
		return Insane
	default:
		return Standard
	}
}
