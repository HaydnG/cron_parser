package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"cronparser/parser"
)

const (
	expectedArgs     = 2
	expectedCronArgs = 6
	textPadding      = 14
)

func main() {

	var err error

	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%s\n%s", rec, debug.Stack())
		}

		if err != nil {
			log.Printf("cron_parser - Error: %v,\n\rArgs: %v,\n\rUsage: cronparser '*/15 0 1,15 * 1-5 /usr/bin/find'", err, os.Args)
			return
		}
	}()

	if len(os.Args) < expectedArgs || len(os.Args[1]) == 0 {
		err = fmt.Errorf("Missing arguments")
		return
	}

	// separate our argument into its separate values
	cronValues := strings.Split(os.Args[1], " ")
	if len(cronValues) < expectedCronArgs {
		err = fmt.Errorf("insufficient args provided for the cron command. Expected: %d, got: %d", expectedCronArgs, len(cronValues))
		return
	}

	// build our parsed cron values, passing in the acceptable ranged for each time definition
	output, err := parseCronTimeField(0, 59, cronValues[0], "minute")
	if err != nil {
		return
	}
	fmt.Println(output)

	output, err = parseCronTimeField(0, 23, cronValues[1], "hour")
	if err != nil {
		return
	}
	fmt.Println(output)

	output, err = parseCronTimeField(1, 31, cronValues[2], "day of month")
	if err != nil {
		return
	}
	fmt.Println(output)

	output, err = parseCronTimeField(1, 12, cronValues[3], "month")
	if err != nil {
		return
	}
	fmt.Println(output)

	output, err = parseCronTimeField(0, 7, cronValues[4], "day of week")
	if err != nil {
		return
	}
	fmt.Println(output)
	fmt.Printf("%-*s %s", textPadding, "command", cronValues[5])
}

// parseCronTimeField accepts valid start and end ranges, then validates and parsed the cron value
func parseCronTimeField(start, end int, field, name string) (string, error) {

	var intervals []string
	var err error

	// execute our parser ruleset, IF a comparitor is found
	for i := range parser.Rules {
		if parser.Rules[i].Comparitor(field) {
			intervals, err = parser.Rules[i].Parse(start, end, field)
			if err != nil {
				return "", fmt.Errorf("%s field parsing failed: %v", name, err)
			}
			break
		}
	}

	// fallback to standalone integers that require `simple` parsing
	if len(intervals) == 0 {
		value, err := strconv.Atoi(field)
		if err != nil {
			return "", fmt.Errorf("%s field parsing failed err: %v, field: %s", name, err, field)
		}
		if value < start || value > end {
			return "", fmt.Errorf("%s field parsing failed err: outside expected range [%d, %d], field: %s", name, start, end, field)
		}

		intervals = append(intervals, field)
	}

	// format our time intervals with the correct padding
	return fmt.Sprintf("%-*s %s", textPadding, name, strings.Join(intervals, " ")), nil
}
