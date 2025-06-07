// Copyright (C) 2025 Nguyen Nhat Tung
//
// Hybrid Profile Generator is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package hevc

// init hevc package internal variables
func init() {
	profiles = []*HEVCProfile{
		{Level: 10, LumaSampleRateMax: 552960, BitRateKBMax: 128},
		{Level: 20, LumaSampleRateMax: 3686400, BitRateKBMax: 1500},
		{Level: 21, LumaSampleRateMax: 7372800, BitRateKBMax: 3000},
		{Level: 30, LumaSampleRateMax: 16588800, BitRateKBMax: 6000},
		{Level: 31, LumaSampleRateMax: 33177600, BitRateKBMax: 10000},
		{Level: 40, LumaSampleRateMax: 66846720, BitRateKBMax: 12000},
		{Level: 41, LumaSampleRateMax: 133693440, BitRateKBMax: 20000},
		{Level: 50, LumaSampleRateMax: 267386880, BitRateKBMax: 25000},
		{Level: 51, LumaSampleRateMax: 534773760, BitRateKBMax: 40000},
		{Level: 52, LumaSampleRateMax: 1069547520, BitRateKBMax: 60000},
		{Level: 60, LumaSampleRateMax: 1069547520, BitRateKBMax: 60000},
		{Level: 61, LumaSampleRateMax: 2139095040, BitRateKBMax: 120000},
		{Level: 62, LumaSampleRateMax: 4278190080, BitRateKBMax: 240000},
	}
}
