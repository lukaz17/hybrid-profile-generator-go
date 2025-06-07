// Copyright (C) 2025 Nguyen Nhat Tung
//
// Hybrid Profile Generator is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package avc

// init avc package internal variables
func init() {
	profiles = []*AVCProfile{
		{Level: 10, MacroBlockMax: 1485, BitRateKBMax: 64, RefFrameMax: 2},
		{Level: 11, MacroBlockMax: 3000, BitRateKBMax: 192, RefFrameMax: 2},
		{Level: 12, MacroBlockMax: 6000, BitRateKBMax: 384, RefFrameMax: 2},
		{Level: 13, MacroBlockMax: 11880, BitRateKBMax: 768, RefFrameMax: 2},
		{Level: 20, MacroBlockMax: 11880, BitRateKBMax: 2000, RefFrameMax: 2},
		{Level: 21, MacroBlockMax: 19800, BitRateKBMax: 4000, RefFrameMax: 2},
		{Level: 22, MacroBlockMax: 20250, BitRateKBMax: 4000, RefFrameMax: 2},
		{Level: 30, MacroBlockMax: 40500, BitRateKBMax: 10000, RefFrameMax: 2},
		{Level: 31, MacroBlockMax: 108000, BitRateKBMax: 14000, RefFrameMax: 3},
		{Level: 32, MacroBlockMax: 216000, BitRateKBMax: 20000, RefFrameMax: 4},
		{Level: 40, MacroBlockMax: 245760, BitRateKBMax: 20000, RefFrameMax: 6},
		{Level: 41, MacroBlockMax: 245760, BitRateKBMax: 50000, RefFrameMax: 6},
		{Level: 42, MacroBlockMax: 522240, BitRateKBMax: 50000, RefFrameMax: 7},
		{Level: 50, MacroBlockMax: 589824, BitRateKBMax: 135000, RefFrameMax: 16},
		{Level: 51, MacroBlockMax: 983040, BitRateKBMax: 240000, RefFrameMax: 16},
		{Level: 52, MacroBlockMax: 2073600, BitRateKBMax: 240000, RefFrameMax: 16},
		{Level: 60, MacroBlockMax: 4177920, BitRateKBMax: 240000, RefFrameMax: 16},
		{Level: 61, MacroBlockMax: 8355840, BitRateKBMax: 480000, RefFrameMax: 16},
		{Level: 62, MacroBlockMax: 16711680, BitRateKBMax: 800000, RefFrameMax: 16},
	}
}
