package main

import "strings"

// PresetType is a set of predefined constants
type PresetType int

const (
	// Medium lame preset
	Medium PresetType = 1 << iota // 1
	// Standard lame preset
	Standard PresetType = 1 << iota // 2
	// Extreme lame preset
	Extreme PresetType = 1 << iota // 4
	// Insane lame preset
	Insane PresetType = 1 << iota // 8
)

func (t PresetType) String() string {
	if t == Medium {
		return "medium"
	}
	if t == Standard {
		return "standard"
	}
	if t == Extreme {
		return "extreme"
	}
	if t == Insane {
		return "insane"
	}
	return "unknown"
}

// presetTypeFromString converts PresetType into a
// string value to send into `flac`.
// Defaults to Standard
func presetTypeFromString(s string) PresetType {
	temp := strings.ToLower(strings.Trim(s, " "))
	if temp == Medium.String() {
		return Medium
	}
	if temp == Extreme.String() {
		return Extreme
	}
	if temp == Insane.String() {
		return Insane
	}
	return Standard
}
