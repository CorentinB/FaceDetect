package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
)

var arguments = struct {
	Input       string
	Output      string
	Concurrency int
	Recursive   bool
}{
	// Default arguments
	Concurrency: 4,
}

func argumentParsing(args []string) {
	// Create new parser object
	parser := argparse.NewParser("facedetect", "Detect images with faces")

	input := parser.String("i", "input", &argparse.Options{
		Required: true,
		Help:     "Input folder"})
	output := parser.String("o", "output", &argparse.Options{
		Required: true,
		Help:     "Output folder"})
	concurrency := parser.Int("c", "concurrency", &argparse.Options{
		Required: false,
		Help:     "Concurrent images to process"})
	recursive := parser.Flag("r", "recursive", &argparse.Options{
		Required: false,
		Help:     "Process input directory recursively",
		Default:  false})

	// Parse input
	err := parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	// Convert path parameters to absolute paths
	inputDir, _ := filepath.Abs(*input)
	outputDir, _ := filepath.Abs(*output)

	// Create output dir
	os.MkdirAll(outputDir, 0755)

	// Finally save the collected flags
	arguments.Input = inputDir
	arguments.Output = outputDir
	arguments.Recursive = *recursive
	arguments.Concurrency = *concurrency
}
