// Copyright (C) 2025 Nguyen Nhat Tung
//
// Hybrid Profile Generator is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package avc

var profiles []*AVCProfile

// RateFactor represents the Constant Rate Factor (CRF) for encoding profiles.
type RateFactor float64

const (
	NormalQuality RateFactor = 24
	HighQuality   RateFactor = 20
	UltraQuality  RateFactor = 16
)

// EncodeProfile contains minimum parameters for encoding video in AVC.
type EncodeProfile struct {
	Name        string
	Width       uint16
	Height      uint16
	FrameRate   float64
	RateFactor  RateFactor
	ThreadCount uint8
}

// AVCProfile contains all constraints of an AVC Level.
type AVCProfile struct {
	Level         uint8
	MacroBlockMax uint32
	BitRateKBMax  uint32
	RefFrameMax   uint8
}

// Return minimum AVC level for specified resolution and framerate.
func MinLevel(width, height uint16, framerate float64) uint8 {
	if width == 0 || height == 0 || framerate <= 0 {
		return 0
	}
	requiredMacroBlocks := uint32(float64(width) * float64(height) * framerate / float64(256))

	level := uint8(0)
	for _, profile := range profiles {
		if profile.MacroBlockMax >= requiredMacroBlocks {
			level = profile.Level
			break
		}
	}
	return level
}

// Return full AVCProfile by specified level.
// Return nil if profile is not found.
func ProfileByLevel(level uint8) *AVCProfile {
	if level == 0 {
		return nil
	}
	for _, profile := range profiles {
		if profile.Level == level {
			return &AVCProfile{
				Level:         profile.Level,
				MacroBlockMax: profile.MacroBlockMax,
				BitRateKBMax:  profile.BitRateKBMax,
				RefFrameMax:   profile.RefFrameMax,
			}
		}
	}
	return nil
}
