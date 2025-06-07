// Copyright (C) 2025 Nguyen Nhat Tung
//
// Hybrid Profile Generator is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package hevc

var profiles []*HEVCProfile

// RateFactor represents the Constant Rate Factor (CRF) for encoding profiles.
type RateFactor float64

const (
	LowQuality    RateFactor = 27
	MediumQuality RateFactor = 22
	HighQuality   RateFactor = 17
)

// EncodeProfile contains minimum parameters for encoding video in HEVC.
type EncodeProfile struct {
	Name        string
	Width       uint16
	Height      uint16
	FrameRate   float64
	RateFactor  RateFactor
	ThreadCount uint8
}

// HEVCProfile contains all constraints of an HEVC Level.
type HEVCProfile struct {
	Level             uint8
	LumaSampleRateMax uint32
	BitRateKBMax      uint32
}

// Return minimum AVC level for specified resolution and framerate.
func MinLevel(width, height uint16, framerate float64) uint8 {
	if width == 0 || height == 0 || framerate <= 0 {
		return 0
	}
	requiredLumaSample := uint32(float64(width) * float64(height) * framerate)

	level := uint8(0)
	for _, profile := range profiles {
		if profile.LumaSampleRateMax >= requiredLumaSample {
			level = profile.Level
			break
		}
	}
	return level
}

// Return full HEVCProfile by specified level.
// Return nil if profile is not found.
func ProfileByLevel(level uint8) *HEVCProfile {
	if level == 0 {
		return nil
	}
	for _, profile := range profiles {
		if profile.Level == level {
			return &HEVCProfile{
				Level:             profile.Level,
				LumaSampleRateMax: profile.LumaSampleRateMax,
				BitRateKBMax:      profile.BitRateKBMax,
			}
		}
	}
	return nil
}
