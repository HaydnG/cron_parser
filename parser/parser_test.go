package parser

import (
	"fmt"
	"testing"
)

func TestParseWildcards(t *testing.T) {
	tests := []struct {
		start          int
		end            int
		field          string
		expectedOutput []string
		expectedError  string
	}{
		{
			start:          0,
			end:            59,
			field:          "*/15",
			expectedOutput: []string{"0", "15", "30", "45"},
			expectedError:  "",
		},
		{
			start: 0,
			end:   59,
			field: "*",
			expectedOutput: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
				"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23",
				"24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36",
				"37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
				"50", "51", "52", "53", "54", "55", "56", "57", "58", "59"},
			expectedError: "",
		},
		{
			start:          0,
			end:            59,
			field:          "*/60",
			expectedOutput: nil,
			expectedError:  "incorrect cron time field: '*/60'",
		},
		{
			start:          0,
			end:            59,
			field:          "*/invalid",
			expectedOutput: nil,
			expectedError:  "err: strconv.Atoi: parsing \"invalid\": invalid syntax, incorrect cron time field: '*/invalid'",
		},
		{
			start:          0,
			end:            23,
			field:          "*/2",
			expectedOutput: []string{"0", "2", "4", "6", "8", "10", "12", "14", "16", "18", "20", "22"},
			expectedError:  "",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Wildcards_%s", test.field), func(t *testing.T) {
			output, err := parseWildcards(test.start, test.end, test.field)
			if test.expectedError != "" {
				if err == nil || err.Error() != test.expectedError {
					t.Errorf("Expected error '%s', but got '%v'", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !equalSlices(output, test.expectedOutput) {
					t.Errorf("Expected output '%v', but got '%v'", test.expectedOutput, output)
				}
			}
		})
	}
}

func TestParseRanges(t *testing.T) {
	tests := []struct {
		start          int
		end            int
		field          string
		expectedOutput []string
		expectedError  string
	}{
		{
			start:          1,
			end:            15,
			field:          "1-10",
			expectedOutput: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			expectedError:  "",
		},
		{
			start:          0,
			end:            23,
			field:          "12-18",
			expectedOutput: []string{"12", "13", "14", "15", "16", "17", "18"},
			expectedError:  "",
		},
		{
			start:          0,
			end:            59,
			field:          "30-60",
			expectedOutput: nil,
			expectedError:  "cron time range specified exceeds limit for this time type, field: '30-60',\n\rminStartRange: 0, maxEndRange: 59",
		},
		{
			start:          0,
			end:            59,
			field:          "invalid-10",
			expectedOutput: nil,
			expectedError:  "err: strconv.Atoi: parsing \"invalid\": invalid syntax, incorrect cron time range field: 'invalid-10'",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Ranges_%s", test.field), func(t *testing.T) {
			output, err := parseRanges(test.start, test.end, test.field)
			if test.expectedError != "" {
				if err == nil || err.Error() != test.expectedError {
					t.Errorf("Expected error '%s', but got '%v'", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !equalSlices(output, test.expectedOutput) {
					t.Errorf("Expected output '%v', but got '%v'", test.expectedOutput, output)
				}
			}
		})
	}
}

func TestParseList(t *testing.T) {
	tests := []struct {
		start          int
		end            int
		field          string
		expectedOutput []string
		expectedError  string
	}{
		{
			start:          0,
			end:            59,
			field:          "1,15,30",
			expectedOutput: []string{"1", "15", "30"},
			expectedError:  "",
		},
		{
			start:          0,
			end:            23,
			field:          "0,12,18",
			expectedOutput: []string{"0", "12", "18"},
			expectedError:  "",
		},
		{
			start:          0,
			end:            23,
			field:          "0,12,25",
			expectedOutput: nil,
			expectedError:  "field parsing failed err: outside expected range [0, 23],value: 25, field: 0,12,25",
		},
		{
			start:          5,
			end:            27,
			field:          "0,12,25",
			expectedOutput: nil,
			expectedError:  "field parsing failed err: outside expected range [5, 27],value: 0, field: 0,12,25",
		},
		{
			start:          0,
			end:            59,
			field:          "invalid,15,30",
			expectedOutput: nil,
			expectedError:  "field parsing failed err: strconv.Atoi: parsing \"invalid\": invalid syntax, value: invalid, field: invalid,15,30",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("List_%s", test.field), func(t *testing.T) {
			output, err := parseList(test.start, test.end, test.field)
			if test.expectedError != "" {
				if err == nil || err.Error() != test.expectedError {
					t.Errorf("Expected error '%s', but got '%v'", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !equalSlices(output, test.expectedOutput) {
					t.Errorf("Expected output '%v', but got '%v'", test.expectedOutput, output)
				}
			}
		})
	}
}

func TestValidInt(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput bool
	}{
		{"123", true},
		{"456", true},
		{"7890", true},
		{"456abc", false},
		{"12.34", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ValidInt_%s", test.input), func(t *testing.T) {
			output := validInt(test.input)
			if output != test.expectedOutput {
				t.Errorf("Expected %v for input '%s', but got %v", test.expectedOutput, test.input, output)
			}
		})
	}
}

// Helper function to compare string slices
func equalSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
