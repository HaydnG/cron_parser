# Cron Parser

This Go application parses and validates cron job time fields using a set of rules and generates execution intervals based on the parsed values.

## Overview

The application parses cron expressions for minute, hour, day of month, month, and day of week fields. It handles wildcards (\*), ranges (1-15), and lists (1,10,15) using modular parser functions.

Does NOT support special time strings, such as `@yearly`

## Features

- **Parsing Wildcards (\*, \*/15, \*/30)**
- **Parsing Ranges (1-15)**
- **Parsing Lists (1,10,15)**

## Usage

### Prerequisites

- Go (version 1.17 or later)

### Installation

1. Unzipping the given .zip file

### Running the Application
##### Running the application (via Go):

```sh
go run cronparser.go '*/15 0 1,15 * 1-5 /usr/bin/find'
```
Replace `*/15 0 1,15 * 1-5 /usr/bin/find` with your cron expression.

##### Building and running the application

An already built executable has been provided
```sh
go build cronparser.go
./cronparser '*/15 0 1,15 * 1-5 /usr/bin/find'
```

<details><summary>Example Output</summary>

```sh
$ ./cronparser '*/15 0 1,15 * 1-5 /usr/bin/find'
minute          0 15 30 45
hour            0
day of month    1 15
month           1 2 3 4 5 6 7 8 9 10 11 12
day of week     1 2 3 4 5
command         /usr/bin/find
```
</details>

## Rules
The application uses a set of rules to determine how to parse different formats of cron time fields:

- Wildcards Rule: Parses expressions like *, */15, */30.
- Ranges Rule: Parses expressions like 1-15.
- List Rule: Parses expressions like 1,10,15.


## Example Output
Upon successful parsing of each cron time field, the application outputs the parsed intervals:

```
minute          0 15 30 45
hour            0
day of month    1 15
month           1 2 3 4 5 6 7 8 9 10 11 12
day of week     1 2 3 4 5
command         /usr/bin/find
```

## Unit tests / benchmarking

### Running Tests
Tests can be ran with the following:
```sh
go test ./...
```

### Benchmarking
Since being a commandline tool, i feel benchmarking is not necessary

## Error Handling
The application provides error handling for incorrect cron time field formats and values outside expected ranges.

## Adding additional rules
This process should be simple. Within the parser package, you should see the ruleset defined.

For example:
```go
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
```

- The `Comparitor` function controls when the rule should be ran
- The `Parse` function does the actual parsing. And returns a slice of time intervals.