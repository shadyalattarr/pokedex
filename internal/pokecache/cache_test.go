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

func TestInvalidKeyGet(t *testing.T) {
	// Define the map[string]struct
	tests := map[string]struct {
		keyIn         string
		valIn         []byte
		expectedValue []byte
	}{
		"simple": {
			keyIn:         "atest",
			valIn:         []byte("Hello"),
			expectedValue: nil,
		},
		"no bytes": {
			keyIn:         "nobytes",
			valIn:         []byte{},
			expectedValue: nil,
		},
		"bytesHex": {
			keyIn:         "Hex",
			valIn:         []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expectedValue: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualCache := NewCache(5 * time.Second)
			actualCache.Add(test.keyIn, test.valIn)

			actualVal, ok := actualCache.Get("invalid Key")
			if ok { // it shouldn't be ok
				t.Errorf("expected %v, got %v", test.expectedValue, actualVal)
			}

		})
	}
}

func TestReapLoop(t *testing.T) {
	// Define the map[string]struct
	tests := map[string]struct {
		interval        time.Duration
		within_interval time.Duration
	}{
		"2 s // 1 s": {
			interval:        2 * time.Second,
			within_interval: 1 * time.Second,
		},
		"1 s // 250 ms": {
			interval:        1 * time.Second,
			within_interval: 500 * time.Millisecond,
		},
		"5 seconds // 4 s": {
			interval:        5 * time.Second,
			within_interval: 4 * time.Second,
		},
		"5 seconds // 2 s": {
			interval:        5 * time.Second,
			within_interval: 2 * time.Second,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel() // great addition
			actualCache := NewCache(test.interval)
			actualCache.Add("reaploop_test", []byte("Hello"))

			time.Sleep(test.within_interval)

			_, ok := actualCache.Get("reaploop_test")
			if !ok {
				t.Errorf("expected item to exist after %v, but it did not exist", test.within_interval)
			}

			whatsLeft := test.interval - test.within_interval
			time.Sleep(whatsLeft + test.interval/10) // guard time

			actualValue, ok := actualCache.Get("reaploop_test")
			if ok {
				t.Errorf("expected item to be removed after %v, but it does exist, \n expected %v, got %v", test.interval, nil, actualValue)
			}

		})
	}
}
