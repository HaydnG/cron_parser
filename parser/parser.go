package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type comparitor func(field string) bool
type parser func(start, end int, field string) ([]string, error)

type parseRule struct {
	Comparitor comparitor
	Parse      parser
}

// Rules provides functions for parsing different formats of the cronjob time formats
var Rules = []parseRule{
	// parse wildcards '*' '*/15' '*/30'
	{
		Comparitor: func(field string) bool {
			return strings.Contains(field, "*")
		},
		Parse: parseWildcards,
	},
	// parse ranges '1-15'
	{
		Comparitor: func(field string) bool {
			return strings.Contains(field, "-")
		},
		Parse: parseRanges,
	},
	// parse list '1,10,15'
	{
		Comparitor: func(field string) bool {
			return strings.Contains(field, ",")
		},
		Parse: parseList,
	},
}

func parseWildcards(start, end int, field string) ([]string, error) {
	var err error

	// always split, if no '/' we'd still just end up with '*'
	// ['*']
	// ['*', '15']
	parts := strings.Split(field, "/")
	if parts[0] != "*" || len(parts) > 2 {
		return nil, fmt.Errorf("incorrect cron time field: '%s'", field)
	}

	// default to a step of 1
	timeStep := 1

	// if we have a custom timeStep, parse and set it
	if len(parts) == 2 {
		timeStep, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("err: %v, incorrect cron time field: '%s'", err, field)
		}
	}

	if timeStep > end || timeStep < 1 {
		return nil, fmt.Errorf("incorrect cron time field: '%s'", field)
	}

	// generate our execution intervals
	var intervals []string
	for i := start; i <= end; i += timeStep {
		intervals = append(intervals, strconv.Itoa(i))
	}
	return intervals, nil
}

func parseRanges(start, end int, field string) ([]string, error) {
	var err error

	parts := strings.Split(field, "-")
	if len(parts) > 2 {
		return nil, fmt.Errorf("incorrect cron time field: '%s'", field)
	}

	startRange, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("err: %v, incorrect cron time range field: '%s'", err, field)
	}

	endRange, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("err: %v, incorrect cron time range field: '%s'", err, field)
	}
	if endRange > end || startRange < start {
		return nil, fmt.Errorf("cron time range specified exceeds limit for this time type, field: '%s',\n\rminStartRange: %d, maxEndRange: %d", field, start, end)
	}

	wrapAround := false
	if startRange > endRange {
		wrapAround = true
	}

	// generate our execution intervals
	var intervals []string
	for i := startRange; i <= end; i++ {

		intervals = append(intervals, strconv.Itoa(i))
		if i == endRange {
			break
		}

		if wrapAround && i == end {
			i = start - 1
		}
	}
	return intervals, nil
}

func parseList(start, end int, field string) ([]string, error) {

	parts := strings.Split(field, ",")

	var intervals []string
	for i := range parts {
		value, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("field parsing failed err: %v, value: %s, field: %s", err, parts[i], field)
		}
		if value < start || value > end {
			return nil, fmt.Errorf("field parsing failed err: outside expected range [%d, %d],value: %s, field: %s", start, end, parts[i], field)
		}

		intervals = append(intervals, parts[i])
	}

	return intervals, nil
}

// validInt checks if a string contains only valid numeric ASCII characters
func validInt(s string) bool {
	if s == "" || len(s) == 0 {
		return false
	}

	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
