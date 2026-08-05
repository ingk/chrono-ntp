package audio

import (
	"testing"
)

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

func TestReadEndOfFile(t *testing.T) {
	reader := &beepReader{data: testSineWaveTable}
	buf := make([]byte, len(testSineWaveTable)+10) // buffer larger than data
	n, err := reader.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != len(testSineWaveTable) {
		t.Errorf("expected to read %d bytes, got %d", len(testSineWaveTable), n)
	}

	// Read again to hit EOF
	n, err = reader.Read(buf)
	if err != nil && err.Error() != "EOF" {
		t.Errorf("expected EOF error, got: %v", err)
	}
	if n != 0 {
		t.Errorf("expected to read 0 bytes at EOF, got %d", n)
	}
}
