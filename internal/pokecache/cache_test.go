package pokecache

import (
	"bytes"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) { // we dont do deepequal cache only because of the mutex
	// Define the map[string]struct
	tests := map[string]struct {
		interval time.Duration
		expected time.Duration // Only checking the interval against expected
	}{
		"simple": {
			interval: 5 * time.Second,
			expected: 5 * time.Second,
		},
		"millisecond": {
			interval: time.Millisecond,
			expected: time.Millisecond,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Call the constructor (assuming it handles its own mutex initialization)
			actual := NewCache(test.interval)

			// 1. Verify the interval matches the expected value
			if actual.interval != test.expected {
				t.Errorf("expected interval %v, got %v", test.expected, actual.interval)
			}

			// 2. Verify the map was properly allocated and isn't nil
			if actual.myCache == nil {
				t.Errorf("expected myCache map to be initialized, but it was nil")
			}
		})
	}
}

func TestAdd(t *testing.T) {
	// Define the map[string]struct
	tests := map[string]struct {
		keyIn         string
		valIn         []byte
		expectedValue []byte
	}{
		"simple": {
			keyIn:         "atest",
			valIn:         []byte("Hello"),
			expectedValue: []byte("Hello"),
		},
		"no bytes": {
			keyIn:         "nobytes",
			valIn:         []byte{},
			expectedValue: []byte{},
		},
		"bytesHex": {
			keyIn:         "Hex",
			valIn:         []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expectedValue: []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualCache := NewCache(5 * time.Second)
			actualCache.Add(test.keyIn, test.valIn)

			actualVal, ok := actualCache.Get(test.keyIn)
			if !ok || !bytes.Equal(actualVal, test.expectedValue) {
				t.Errorf("expected %v, got %v", test.expectedValue, actualVal)
			}

		})
	}
}

// do test get and reaploop
