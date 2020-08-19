package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/CorentinB/FaceDetect/face"
	pigo "github.com/esimov/pigo/core"
	"github.com/remeh/sizedwaitgroup"

	_ "image/jpeg"

	"github.com/labstack/gommon/color"
)

func listFiles(path string, recursive bool) []string {
	var files []string
	var currentDirFiles []string
	var subDirFiles []string

	// Read all files in path
	items, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Unable to read directory " + path)
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	// Go through all files, and process subdirectories if needed
	for _, item := range items {
		if item.IsDir() {
			if recursive {
				subPath, err := filepath.Abs(path + "/" + item.Name())
				if err != nil {
					fmt.Println("Unable to get absolute path of " + item.Name())
					fmt.Println("Error: " + err.Error())
					os.Exit(1)
				}
				subDirFiles = append([]string{}, append(listFiles(subPath, recursive), subDirFiles...)...)
			}
		} else {
			absPath, err := filepath.Abs(path + "/" + item.Name())
			if err != nil {
				fmt.Println("Unable to get absolute path of " + item.Name())
				fmt.Println("Error: " + err.Error())
				os.Exit(1)
			}
			currentDirFiles = append(currentDirFiles, absPath)
		}
	}

	files = append([]string{}, append(currentDirFiles, subDirFiles...)...)

	return files
}

func processFile(path, output string, worker *sizedwaitgroup.SizedWaitGroup) {
	defer worker.Done()

	// Detect faces on the image
	src, err := pigo.GetImage(path)
	if err != nil {
		logError("Error reading the source file " + path + ": " + err.Error())
		return
	}

	facesCount := face.Detect(src, fileNameWithoutExtension(filepath.Base(path)), arguments.Output)
	logSuccess(color.Green(strconv.Itoa(facesCount)) + color.Yellow(" faces found in "+filepath.Base(path)))
}

func main() {
	start := time.Now()
	argumentParsing(os.Args)

	var worker = sizedwaitgroup.New(arguments.Concurrency)

	// Get a list of all files
	files := listFiles(arguments.Input, arguments.Recursive)
	logSuccess(color.Green(strconv.Itoa(len(files))) + color.Yellow(" pictures ready for processing"))

	// Go through all files and detect if there are faces
	for _, path := range files {
		worker.Add()
		go processFile(path, arguments.Output, &worker)
	}

	worker.Wait()

	logSuccess(color.Green(strconv.Itoa(len(files))) + color.Yellow(" pictures sorted in ") + color.Green(time.Since(start)))
}
