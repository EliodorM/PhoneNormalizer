package main

import "testing"

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"123456789", "123456789"},
		{"(123)4294", "1234294"},
		{"2131 12134", "213112134"},
		{"312-121-145", "312121145"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.want {
				t.Errorf("got %s; want %s", actual, tc.want)
			}

		})

	}

}
