package main

import (
	"flag"
	"fmt"
	"os"
)

type SampleSubcommand struct {
	replace bool
	numRows int
	seed    int
}

func (sub *SampleSubcommand) Name() string {
	return "sample"
}
func (sub *SampleSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SampleSubcommand) Description() string {
	return "Sample rows."
}
func (sub *SampleSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&sub.replace, "replace", false, "Sample with replacement")
	fs.IntVar(&sub.numRows, "n", 0, "Number of rows to sample")
	fs.IntVar(&sub.seed, "seed", 0, "Seed for random number generation")
}

func (sub *SampleSubcommand) Run(args []string) {
	if sub.numRows < 1 {
		fmt.Fprintln(os.Stderr, "Invalid required argument -n")
		os.Exit(1)
	}

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	Sample(inputCsvs[0], sub.numRows, sub.replace, sub.seed)
}

func Sample(inputCsv *InputCsv, numRows int, replace bool, seed int) {

	imc := NewInMemoryCsvFromInputCsv(inputCsv)

	if numRows > imc.NumRows() && !replace {
		fmt.Fprintln(os.Stderr, "Cannot sample more rows than exist")
		os.Exit(1)
	}

	rowIndices := imc.SampleRowIndices(numRows, replace, seed)

	outputCsv := NewOutputCsvFromInputCsv(inputCsv)

	// Write header.
	outputCsv.Write(imc.header)

	for _, rowIndex := range rowIndices {
		outputCsv.Write(imc.rows[rowIndex])
	}

}
