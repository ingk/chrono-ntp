package main

import "testing"

const testDurationMs = 100

var testSineWaveTable = makeSineWaveTable(testDurationMs)

func TestMakeSineWaveTable_Length(t *testing.T) {
	expectedSamples := 8820 // sampleRate * testDurationMs / 1000 (44100 Hz * 100 ms / 1000)
	if len(testSineWaveTable) != expectedSamples {
		t.Errorf("expected length %d, got %d", expectedSamples, len(testSineWaveTable))
	}
}

func TestMakeSineWaveTable_Range(t *testing.T) {
	for i := 0; i < len(testSineWaveTable); i += 2 {
		v := int16(testSineWaveTable[i]) | int16(testSineWaveTable[i+1])<<8
		if v < -maxInt16 || v > maxInt16 {
			t.Errorf("sample out of range: %d", v)
		}
	}
}
