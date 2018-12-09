package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type HeadSubcommand struct {
	numRowsStr string
}

func (sub *HeadSubcommand) Name() string {
	return "head"
}
func (sub *HeadSubcommand) Aliases() []string {
	return []string{}
}
func (sub *HeadSubcommand) Description() string {
	return "Extract the first N rows from a CSV."
}
func (sub *HeadSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.numRowsStr, "n", "10", "Number of rows to include")
}

func (sub *HeadSubcommand) Run(args []string) {
	numRowsRegex := regexp.MustCompile("^\\+?\\d+$")
	if !numRowsRegex.MatchString(sub.numRowsStr) {
		fmt.Fprintln(os.Stderr, "Invalid argument to -n")
		os.Exit(1)
		return
	}

	inputCsvs := GetInputCsvsOrPanic(args, 1)

	if strings.HasPrefix(sub.numRowsStr, "+") {
		sub.numRowsStr = strings.TrimPrefix(sub.numRowsStr, "+")
		numRows, err := strconv.Atoi(sub.numRowsStr)
		if err != nil {
			ExitWithError(err)
		}
		HeadFromBottom(inputCsvs[0], numRows)
	} else {
		numRows, err := strconv.Atoi(sub.numRowsStr)
		if err != nil {
			ExitWithError(err)
		}
		HeadFromTop(inputCsvs[0], numRows)
	}
}

func HeadFromBottom(inputCsv *InputCsv, numRows int) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

	rows, err := inputCsv.ReadAll()
	if err != nil {
		ExitWithError(err)
	}

	// Write header.
	outputCsv.Write(rows[0])

	// Write rows up to last `numRows` rows.
	maxRow := len(rows) - numRows
	if maxRow < 1 {
		return
	}
	for i := 1; i < maxRow; i++ {
		outputCsv.Write(rows[i])
	}
}

func HeadFromTop(inputCsv *InputCsv, numRows int) {
	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	outputCsv.Write(header)

	// Write first `numRows` rows.
	curRow := 0
	for {
		if curRow == numRows {
			break
		}
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		curRow++
		outputCsv.Write(row)
	}
}
