package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Field struct {
	name string
	min  int
	max  int
}

func parseField(field string, min, max int) ([]int, error) {
	var result []int
	parts := strings.Split(field, ",")

	for _, part := range parts {
		step := 1
		rangeStr := part

		if strings.Contains(part, "/") {
			split := strings.Split(part, "/")
			if len(split) != 2 {
				return nil, fmt.Errorf("invalid step format: %s", part)
			}
			rangeStr = split[0]
			stepStr := split[1]
			var err error
			step, err = strconv.Atoi(stepStr)
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step value: %s", stepStr)
			}
		}

		var fieldMin, fieldMax int
		if rangeStr == "*" {
			fieldMin = min
			fieldMax = max
		} else if strings.Contains(rangeStr, "-") {
			split := strings.Split(rangeStr, "-")
			if len(split) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", rangeStr)
			}
			var err error
			fieldMin, err = strconv.Atoi(split[0])
			if err != nil {
				return nil, fmt.Errorf("invalid min value: %s", split[0])
			}
			fieldMax, err = strconv.Atoi(split[1])
			if err != nil {
				return nil, fmt.Errorf("invalid max value: %s", split[1])
			}
			if fieldMin > fieldMax {
				return nil, fmt.Errorf("min > max in range: %s", rangeStr)
			}
		} else {
			var err error
			fieldMin, err = strconv.Atoi(rangeStr)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", rangeStr)
			}
			fieldMax = fieldMin
		}

		for i := fieldMin; i <= fieldMax && i <= max; i += step {
			if i >= min {
				result = append(result, i)
			}
		}
	}

	sort.Ints(result)
	unique := []int{}
	prev := -1
	for _, val := range result {
		if val != prev {
			unique = append(unique, val)
			prev = val
		}
	}
	return unique, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cron_parser \"cron-string\"")
		fmt.Println("Example: cron_parser \"*/15 0 1,15 * 1-5 /usr/bin/find\"")
		os.Exit(1)
	}

	input := os.Args[1]
	parts := strings.Fields(input)
	if len(parts) < 6 {
		fmt.Println("Invalid input: cron string must have 5 time fields and a command")
		os.Exit(1)
	}

	timeFields := parts[:5]
	command := strings.Join(parts[5:], " ")

	fields := []Field{
		{"minute", 0, 59},
		{"hour", 0, 23},
		{"day of month", 1, 31},
		{"month", 1, 12},
		{"day of week", 1, 7},
	}

	for i, field := range fields {
		values, err := parseField(timeFields[i], field.min, field.max)
		if err != nil {
			fmt.Printf("Error parsing field %s: %v\n", field.name, err)
			os.Exit(1)
		}

		strValues := make([]string, len(values))
		for j, v := range values {
			strValues[j] = strconv.Itoa(v)
		}

		fmt.Printf("%-14s %s\n", field.name, strings.Join(strValues, " "))
	}

	fmt.Printf("%-14s %s\n", "command", command)
}
