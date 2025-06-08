// Copyright (C) 2025 Nguyen Nhat Tung
//
// Hybrid Profile Generator is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"text/template"

	"github.com/lukaz17/hybrid-profile-generator-go/hevc"
	"github.com/lukaz17/hybrid-profile-generator-go/video"
	"github.com/tforce-io/tf-golib/diag"
	"github.com/tforce-io/tf-golib/opx"
	"github.com/tforce-io/tf-golib/stdx/mathxt"
)

var logger = diag.DefaultLogger{}

// EncodeParams holds the parameters for encoding profiles.
type EncodeParams struct {
	Name          string
	Width         uint16
	Height        uint16
	FrameRate     float64
	ThreadCount   uint8
	RateFactor    float64
	RateFactorMax float64
	HEVCLevel     float64
	HEVCTier      string
	RefFrame      uint8
	MeRange       uint8
	BFrame        uint8
	KeyInterval   uint16
	RCLookahead   uint16
	AQStrength    float64
}

func main() {
	profiles := []*hevc.EncodeProfile{}
	// Generic profiles
	resolutions := []*video.Resolution{
		{Width: 960, Height: 720},
		{Width: 1280, Height: 720},
		{Width: 1280, Height: 960},
		{Width: 1440, Height: 1080},
		{Width: 1920, Height: 816},
		{Width: 1920, Height: 1080},
		{Width: 1920, Height: 1440},
		{Width: 2560, Height: 1440},
		{Width: 3840, Height: 1600},
		{Width: 3840, Height: 2160},
	}
	framerates := []float64{25, 30, 50, 60}
	qualities := []hevc.RateFactor{
		hevc.LowQuality,
		hevc.MediumQuality,
		hevc.HighQuality,
	}
	for _, resolution := range resolutions {
		for _, framerate := range framerates {
			for _, quality := range qualities {
				profile := &hevc.EncodeProfile{
					Width:       resolution.Width,
					Height:      resolution.Height,
					FrameRate:   framerate,
					RateFactor:  quality,
					ThreadCount: 16,
				}
				profiles = append(profiles, profile)
			}
		}
	}

	defaultProfile, err := ioutil.ReadFile("./presets/x264.xml")
	if err != nil {
		logger.Error(err, "failed to read template file")
	}
	template, err := template.New("x265").Parse(string(defaultProfile))
	if err != nil {
		logger.Error(err, "failed to parse profile")
	}
	for _, profile := range profiles {
		params := createSetting(profile)
		saveSetting(template, params)
	}
	logger.Info("x265 profiles generated successfully.")
}

// Create EncodeParams based on EncodeProfile.
func createSetting(profile *hevc.EncodeProfile) *EncodeParams {
	quality := "L"
	if float64(profile.RateFactor) <= float64(19) {
		quality = "H"
	} else if float64(profile.RateFactor) <= float64(24) {
		quality = "M"
	}
	params := &EncodeParams{
		Name:        opx.Ternary(profile.Name != "", profile.Name, fmt.Sprintf("%dx%d@%4.2f-%s", profile.Width, profile.Height, profile.FrameRate, quality)),
		Width:       profile.Width,
		Height:      profile.Height,
		FrameRate:   profile.FrameRate,
		RateFactor:  float64(profile.RateFactor),
		ThreadCount: profile.ThreadCount,
	}
	level := hevc.MinLevel(profile.Width, profile.Height, profile.FrameRate)
	meRange, aqStrength := factorsByResolution(profile.Width)
	refFrame, bFrame, aqStrengthModifier := factorsByRateFactor(profile.RateFactor, profile.FrameRate)

	params.RateFactorMax = float64(profile.RateFactor) - 5
	params.HEVCLevel = float64(level) / 10
	params.HEVCTier = opx.Ternary(level >= 40, "High", "Main")
	params.RefFrame = refFrame
	params.MeRange = meRange
	params.BFrame = bFrame
	params.KeyInterval = uint16(math.Ceil(profile.FrameRate) * 10)
	params.RCLookahead = mathxt.MinUint16(params.KeyInterval, 250)
	params.AQStrength = aqStrength + aqStrengthModifier
	return params
}

// Save the EnodeParms to disk.
func saveSetting(template *template.Template, params *EncodeParams) {
	fileName := fmt.Sprintf("x265 %s.xml", params.Name)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logger.Error(err, "cannot write to file", fileName)
		return
	}
	defer file.Close()

	err = template.Execute(file, params)
	if err != nil {
		logger.Error(err, fileName)
	}
}

// Determine the motion estimation range and AQ strength based on the video width.
func factorsByResolution(width uint16) (meRange uint8, aqStrength float64) {
	meRange = uint8(24)
	aqStrength = float64(1)

	if width >= (3840 * 15 / 16) {
		meRange = uint8(57)
		aqStrength = float64(0.5)
	} else if width >= (2560 * 15 / 16) {
		meRange = uint8(57)
		aqStrength = float64(0.6)
	} else if width >= (1920 * 7 / 8) {
		meRange = uint8(57)
		aqStrength = float64(0.7)
	} else if width >= (1280 * 7 / 8) {
		meRange = uint8(48)
		aqStrength = float64(0.9)
	} else {
		meRange = uint8(32)
		aqStrength = float64(0.9)
	}

	return meRange, aqStrength
}

// Determine the reference frame count, B-frame count, and AQ strength modifier based on the rate factor and frame rate.
func factorsByRateFactor(quality hevc.RateFactor, frameRate float64) (refFrame, bFrame uint8, aqStrengthModifier float64) {
	refFrame = opx.Ternary(frameRate >= 32, uint8(4), uint8(3))
	bFrame = uint8(7)
	aqStrengthModifier = float64(0.15)

	if float64(quality) <= float64(17) {
		refFrame += 2
		bFrame = uint8(16)
		aqStrengthModifier = float64(0)
	} else if float64(quality) <= float64(22) {
		refFrame += 1
		bFrame = uint8(12)
		aqStrengthModifier = float64(0.05)
	} else {
		refFrame += 0
		bFrame = uint8(7)
		aqStrengthModifier = float64(0.1)
	}

	return refFrame, bFrame, aqStrengthModifier
}
