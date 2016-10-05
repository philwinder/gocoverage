package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const test_dir = "testing"

var test_profile_file = filepath.Join(test_dir, cov_tmp_file)

func TestGoCoverage(t *testing.T) {
	// Test file is created
	testDir, err := os.Open(test_dir)
	if err != nil {
		t.Fatal(err)
	}
	stat, err := testDir.Stat()
	if err != nil {
		t.Fatal(err)
	}
	PerformCoverage(testDir.Name(), stat, nil)
	if covStat, err := os.Stat(test_profile_file); os.IsNotExist(err) || covStat.Size() == 0 {
		t.Fatal(test_profile_file + " does not exist or is empty")
	}

	// Test file is collated
	testCovFile, err := os.Open(test_profile_file)
	if err != nil {
		t.Fatal(err)
	}
	stat, err = testCovFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	CreateOutputFile()
	CollateCoverage(testCovFile.Name(), stat, nil)
	if covStat, err := os.Stat(output_filename); os.IsNotExist(err) || covStat.Size() == 0 {
		t.Fatal(output_filename + " does not exist or is empty")
	}

	// Check that resultant file is valid
	outputBytes, err := ioutil.ReadFile(output_filename)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(outputBytes), "mode: count") {
		t.Fatal("Does not contain mode: count")
	}

	// Test delete files
	DeleteFiles(testCovFile.Name(), stat, nil)
	if _, err := os.Stat(test_profile_file); !os.IsNotExist(err) {
		t.Fatal(test_profile_file + " still exists, and it should not")
	}

	// Remove output file
	os.Remove(output_filename)
}

func TestMainFunction(t *testing.T) {
	dir = test_dir
	main()
	// Remove output file
	os.Remove(filepath.Join(dir, output_filename))
}

func TestPanic(t *testing.T) {
	defer func() { recover() }()
	check(errors.New("This is an error"))
	t.Fatal("Should not get to this line")
}
