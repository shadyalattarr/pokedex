package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	tests := map[string]struct {
		textInput string
		expected  []string
	}{
		"simple":                   {textInput: "bananas are cool!", expected: []string{"bananas", "are", "cool!"}},
		"no space":                 {textInput: "NOSPACEBOIlol", expected: []string{"nospaceboilol"}},
		"spaces before":            {textInput: "  space Before", expected: []string{"space", "before"}},
		"spaces after":             {textInput: "space after   ", expected: []string{"space", "after"}},
		"spaces before and after ": {textInput: "   Space before & after   ", expected: []string{"space", "before", "&", "after"}},
		"too many spaces":          {textInput: "a   LOT  of      spaces!", expected: []string{"a", "lot", "of", "spaces!"}},
	}

	for testName, test := range tests {
		// t.Run creates a subtest for each map entry
		t.Run(testName, func(t *testing.T) {
			actual := cleanInput(test.textInput)
			fmt.Println(actual)

			if len(actual) != len(test.expected) {
				// This now stops ONLY this subtest, and the loop continues
				fmt.Println("Debugging:")
				fmt.Printf("Expected: %#v\n", test.expected)
				fmt.Printf("got: %#v\n", actual)
				t.Fatalf("lengths don't match: %v != %v", len(actual), len(test.expected))
			}

			for i, word := range actual {
				expectedWord := test.expected[i]
				if word != expectedWord {
					// This also stops ONLY this subtest, and the loop continues
					fmt.Println("Debugging:")
					fmt.Printf("Expected: %#v\n", test.expected)
					fmt.Printf("got: %#v\n", actual)
					t.Fatalf("%v is not same as %v in index %d.", word, expectedWord, i)

				}
			}
		})
	}
}
