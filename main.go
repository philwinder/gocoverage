package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	covTmpFile  = "profile.cov.tmp"
	coverHeader = `mode: count
`
)

var (
	dir            string
	ignore         string
	outputFilename string
	fileIgnore     *regexp.Regexp
)

func init() {
	flag.StringVar(&dir, "dir", ".", "Directory to start recursing for tests")
	flag.StringVar(&ignore, "ignore", `(vendor|\.\w+)`, "RegEx that ignores files and folders. Default ignores hidden folders and vendor folder.")
	flag.StringVar(&outputFilename, "output", "profile.cov", "Filename for the output coverage file.")
	flag.Parse()
	fileIgnore = regexp.MustCompile(ignore)
}

func main() {
	filepath.Walk(dir, performCoverage)
	createOutputFile()
	filepath.Walk(dir, collateCoverage)
	filepath.Walk(dir, deleteFiles)
}

// A filepath.Walk function to use `go test` to generate all the coverage reports
func performCoverage(path string, info os.FileInfo, err error) error {
	if err == nil && info.IsDir() && hasGoFile(path) && !fileIgnore.MatchString(path) {
		path = "./" + path
		log.Println(path)
		exec.Command("go", "test", "-covermode=count", "-coverprofile="+path+"/"+covTmpFile, path).Output()
	}
	return nil
}

// Creates the final output file
func createOutputFile() {
	err := ioutil.WriteFile(outputFilename, []byte(coverHeader), 0644)
	check(err)
}

// A filepath.Walk function to collate the coverage reports together and saves into output file
func collateCoverage(path string, info os.FileInfo, err error) error {
	if err == nil && !info.IsDir() && strings.Contains(path, covTmpFile) {
		contentsB, err := ioutil.ReadFile(path)
		check(err)
		contents := strings.Replace(string(contentsB), "mode: count\n", "", 1)
		f, err := os.OpenFile(outputFilename, os.O_APPEND|os.O_WRONLY, 0600)
		check(err)
		_, err = f.WriteString(contents)
		check(err)
		f.Close()
	}
	return nil
}

// A filepath.Walk function to delete all the temporary files
func deleteFiles(path string, info os.FileInfo, err error) error {
	if err == nil && !info.IsDir() && strings.Contains(path, covTmpFile) {
		os.Remove(path)
	}
	return nil
}

func hasGoFile(path string) bool {
	matches, _ := filepath.Glob(path + "/*.go")
	return (matches != nil && len(matches) > 0)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
