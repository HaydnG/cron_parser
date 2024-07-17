package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func TestMainFunction(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "Missing arguments",
			args:           []string{"cron_parser"},
			expectedOutput: "",
			expectedError:  "Missing arguments",
		},
		{
			name:           "Insufficient cron arguments",
			args:           []string{"cron_parser", "*/15 0 1,15 * 1-5"},
			expectedOutput: "",
			expectedError:  "insufficient args provided for the cron command",
		},
		// {
		// 	name:           "Invalid cron command (invalid minute)",
		// 	args:           []string{"cron_parser", "*/61 0 1,15 * 1-5 /usr/bin/find"},
		// 	expectedOutput: "",
		// 	expectedError:  "minute field parsing failed: incorrect cron time field: '*/61'",
		// },
		{
			name:           "Invalid cron command (invalid hour)",
			args:           []string{"cron_parser", "*/15 24 1,15 * 1-5 /usr/bin/find"},
			expectedOutput: "",
			expectedError:  "hour field parsing failed err: outside expected range [0, 23], field: 24",
		},
		{
			name:           "Invalid cron command (invalid day of month)",
			args:           []string{"cron_parser", "*/15 0 32 * 1-5 /usr/bin/find"},
			expectedOutput: "",
			expectedError:  "day of month field parsing failed err: outside expected range [1, 31], field: 32",
		},
		{
			name:           "Invalid cron command (invalid month)",
			args:           []string{"cron_parser", "*/15 0 1,15 13 * /usr/bin/find"},
			expectedOutput: "",
			expectedError:  "month field parsing failed err: outside expected range [1, 12], field: 13",
		},
		{
			name: "Valid cron command with zero week",
			args: []string{"cron_parser", "*/15 0 1,15 * 0 /usr/bin/find"},
			expectedOutput: `minute         0 15 30 45
hour           0
day of month   1 15
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    0
command        /usr/bin/find`,
			expectedError: "",
		},
		{
			name: "Valid cron command",
			args: []string{"cron_parser", "*/15 0 1,15 * 1-5 /usr/bin/find"},
			expectedOutput: `minute         0 15 30 45
hour           0
day of month   1 15
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    1 2 3 4 5
command        /usr/bin/find`,
			expectedError: "",
		},
		{
			name: "Valid cron command with every minute",
			args: []string{"cron_parser", "* * * * * /usr/bin/every-minute"},
			expectedOutput: `minute         0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59
hour           0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month   1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    0 1 2 3 4 5 6 7
command        /usr/bin/every-minute`,
			expectedError: "",
		},
		{
			name: "Valid cron command with every hour",
			args: []string{"cron_parser", "0 * * * * /usr/bin/every-hour"},
			expectedOutput: `minute         0
hour           0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month   1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    0 1 2 3 4 5 6 7
command        /usr/bin/every-hour`,
			expectedError: "",
		},
		{
			name: "Valid cron command with every day at midnight",
			args: []string{"cron_parser", "0 0 * * * /usr/bin/once-a-day"},
			expectedOutput: `minute         0
hour           0
day of month   1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    0 1 2 3 4 5 6 7
command        /usr/bin/once-a-day`,
			expectedError: "",
		},
		{
			name: "Valid cron command with every day at midnight - Additional args",
			args: []string{"cron_parser", "0 0 * * * /usr/bin/once-a-day -ls"},
			expectedOutput: `minute         0
hour           0
day of month   1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    0 1 2 3 4 5 6 7
command        /usr/bin/once-a-day -ls`,
			expectedError: "",
		},
		{
			name: "Valid cron command wraparound",
			args: []string{"cron_parser", "*/15 0 1,15 * 5-1 /usr/bin/find"},
			expectedOutput: `minute         0 15 30 45
hour           0
day of month   1 15
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    5 6 7 0 1
command        /usr/bin/find`,
			expectedError: "",
		},
		{
			name: "Valid cron command Year param",
			args: []string{"cron_parser", "*/15 0 1,15 * 5-1 * /usr/bin/find"},
			expectedOutput: `minute         0 15 30 45
hour           0
day of month   1 15
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    5 6 7 0 1
year		   2024 2025 2026 2027 2028 2029 2030 2031 2032 2033 2034 2035 2036 2037 2038 2039 2040 2041 2042 2043 2044
command        /usr/bin/find`,
			expectedError: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.args

			// Capture log output
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)

			// Capture fmt output
			old := os.Stdout // keep original stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			main()

			// Redirect output back to stdout when done
			w.Close()
			os.Stdout = old
			var buf bytes.Buffer
			io.Copy(&buf, r)

			out := buf.String()
			if test.expectedError != "" {
				if !bytes.Contains([]byte(logBuf.String()), []byte(test.expectedError)) {
					t.Errorf("Command: %s, \n\rExpected error message '%s', got '%s', output: \n\r%s", test.args, test.expectedError, logBuf.String(), out)
				}
			} else if logBuf.String() != "" && test.expectedError == "" {
				t.Errorf("Command: %s, \n\rUnexpected error message: '%s'", test.args, logBuf.String())
			} else {
				if out != test.expectedOutput {
					t.Errorf("Command: %s, \n\rExpected output: \n\r'\n%s\n\r'\n\rgot: \n\r'\n%s\n\r'", test.args, test.expectedOutput, out)
				}
			}
		})
	}
}
