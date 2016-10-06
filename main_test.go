package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testDir = "testing"

var testProfileFile = filepath.Join(testDir, covTmpFile)

func TestGoCoverage(t *testing.T) {
	// Test file is created
	testDir, err := os.Open(testDir)
	if err != nil {
		t.Fatal(err)
	}
	stat, err := testDir.Stat()
	if err != nil {
		t.Fatal(err)
	}
	performCoverage(testDir.Name(), stat, nil)
	if covStat, err := os.Stat(testProfileFile); os.IsNotExist(err) || covStat.Size() == 0 {
		t.Fatal(testProfileFile + " does not exist or is empty")
	}

	// Test file is collated
	testCovFile, err := os.Open(testProfileFile)
	if err != nil {
		t.Fatal(err)
	}
	stat, err = testCovFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	createOutputFile()
	collateCoverage(testCovFile.Name(), stat, nil)
	if covStat, err := os.Stat(outputFilename); os.IsNotExist(err) || covStat.Size() == 0 {
		t.Fatal(outputFilename + " does not exist or is empty")
	}

	// Check that resultant file is valid
	outputBytes, err := ioutil.ReadFile(outputFilename)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(outputBytes), "mode: count") {
		t.Fatal("Does not contain mode: count")
	}

	// Test delete files
	deleteFiles(testCovFile.Name(), stat, nil)
	if _, err := os.Stat(testProfileFile); !os.IsNotExist(err) {
		t.Fatal(testProfileFile + " still exists, and it should not")
	}

	// Remove output file
	os.Remove(outputFilename)
}

func TestMainFunction(t *testing.T) {
	dir = testDir
	main()
	// Remove output file
	os.Remove(filepath.Join(dir, outputFilename))
}

func TestPanic(t *testing.T) {
	defer func() { recover() }()
	check(errors.New("This is an error"))
	t.Fatal("Should not get to this line")
}
