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

	"github.com/lukaz17/hybrid-profile-generator-go/avc"
	"github.com/lukaz17/hybrid-profile-generator-go/video"
	"github.com/tforce-io/tf-golib/diag"
	"github.com/tforce-io/tf-golib/opx"
	"github.com/tforce-io/tf-golib/stdx/mathxt"
)

var logger = diag.DefaultLogger{}

// EncodeParams holds the parameters for encoding profiles.
type EncodeParams struct {
	Name           string
	Width          uint16
	Height         uint16
	FrameRate      float64
	ThreadCount    uint8
	RateFactor     float64
	AVCLevel       float64
	RefFrame       uint8
	MeRange        uint8
	BFrame         uint8
	KeyInterval    uint16
	InputLookahead uint8
	RCLookahead    uint16
	AQStrength     float64
}

func main() {
	profiles := []*avc.EncodeProfile{
		{Name: "NTSC DVD", Width: 640, Height: 480, FrameRate: 30, RateFactor: avc.MediumQuality, ThreadCount: 16},
		{Name: "PAL DVD", Width: 768, Height: 576, FrameRate: 25, RateFactor: avc.MediumQuality, ThreadCount: 16},
		{Name: "NTSC-WIDE DVD", Width: 864, Height: 480, FrameRate: 30, RateFactor: avc.MediumQuality, ThreadCount: 16},
		{Name: "PAL-WIDE DVD", Width: 1024, Height: 576, FrameRate: 25, RateFactor: avc.MediumQuality, ThreadCount: 16},
	}
	// Generic profiles
	resolutions := []*video.Resolution{
		{Width: 640, Height: 360},
		{Width: 640, Height: 480},
		{Width: 960, Height: 540},
		{Width: 960, Height: 720},
		{Width: 1280, Height: 720},
		{Width: 1280, Height: 960},
		{Width: 1440, Height: 1080},
		{Width: 1920, Height: 816},
		{Width: 1920, Height: 1080},
		{Width: 1920, Height: 1440},
	}
	framerates := []float64{25, 30, 50, 60}
	qualities := []avc.RateFactor{
		avc.LowQuality,
		avc.MediumQuality,
		avc.HighQuality,
	}
	for _, resolution := range resolutions {
		for _, framerate := range framerates {
			for _, quality := range qualities {
				profile := &avc.EncodeProfile{
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
	template, err := template.New("x264").Parse(string(defaultProfile))
	if err != nil {
		logger.Error(err, "failed to parse profile")
	}
	for _, profile := range profiles {
		params := createSetting(profile)
		saveSetting(template, params)
	}
	logger.Info("x264 profiles generated successfully.")
}

// Create EncodeParams based on EncodeProfile.
func createSetting(profile *avc.EncodeProfile) *EncodeParams {
	quality := "L"
	if float64(profile.RateFactor) <= float64(17) {
		quality = "H"
	} else if float64(profile.RateFactor) <= float64(22) {
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
	level := avc.MinLevel(profile.Width, profile.Height, profile.FrameRate)
	x264Profile := avc.ProfileByLevel(level)
	meRange, aqStrength := factorsByResolution(profile.Width)
	refFrame, bFrame, aqStrengthModifier := factorsByRateFactor(profile.RateFactor, profile.FrameRate)

	params.AVCLevel = float64(level) / 10
	params.RefFrame = mathxt.MinUint8(x264Profile.RefFrameMax, refFrame)
	params.MeRange = meRange
	params.BFrame = bFrame
	params.KeyInterval = uint16(math.Ceil(profile.FrameRate) * 10)
	params.InputLookahead = mathxt.MaxUint8(params.ThreadCount*5, 30)
	params.RCLookahead = uint16(math.Ceil(profile.FrameRate) * 2)
	params.AQStrength = aqStrength + aqStrengthModifier
	return params
}

// Save the EnodeParms to disk.
func saveSetting(template *template.Template, params *EncodeParams) {
	fileName := fmt.Sprintf("x264 %s.xml", params.Name)

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
		meRange = uint8(64)
		aqStrength = float64(0.7)
	} else if width >= (2560 * 15 / 16) {
		meRange = uint8(48)
		aqStrength = float64(0.75)
	} else if width >= (1920 * 7 / 8) {
		meRange = uint8(32)
		aqStrength = float64(0.9)
	} else if width >= (1280 * 7 / 8) {
		meRange = uint8(32)
		aqStrength = float64(1)
	} else {
		meRange = uint8(24)
		aqStrength = float64(1.1)
	}

	return meRange, aqStrength
}

// Determine the reference frame count, B-frame count, and AQ strength modifier based on the rate factor and frame rate.
func factorsByRateFactor(quality avc.RateFactor, frameRate float64) (refFrame, bFrame uint8, aqStrengthModifier float64) {
	refFrame = opx.Ternary(frameRate >= 32, uint8(5), uint8(3))
	bFrame = uint8(7)
	aqStrengthModifier = float64(0.15)

	if float64(quality) <= float64(17) {
		refFrame += 2
		bFrame = uint8(16)
		aqStrengthModifier = float64(0.05)
	} else if float64(quality) <= float64(22) {
		refFrame += 1
		bFrame = uint8(12)
		aqStrengthModifier = float64(0.1)
	} else {
		refFrame += 0
		bFrame = uint8(7)
		aqStrengthModifier = float64(0.15)
	}

	return refFrame, bFrame, aqStrengthModifier
}
