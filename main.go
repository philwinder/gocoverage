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
	cov_tmp_file = "profile.cov.tmp"
	cover_header = `mode: count
`
)

var (
	dir             string
	ignore          string
	output_filename string
	fileIgnore      *regexp.Regexp
)

func init() {
	flag.StringVar(&dir, "dir", ".", "Directory to start recursing for tests")
	flag.StringVar(&ignore, "ignore", `(vendor|\.\w+)`, "RegEx that ignores files and folders. Default ignores hidden folders and vendor folder.")
	flag.StringVar(&output_filename, "output", "profile.cov", "Filename for the output coverage file.")
	flag.Parse()
	fileIgnore = regexp.MustCompile(ignore)
}

func main() {
	filepath.Walk(dir, PerformCoverage)
	CreateOutputFile()
	filepath.Walk(dir, CollateCoverage)
	filepath.Walk(dir, DeleteFiles)
}

func PerformCoverage(path string, info os.FileInfo, err error) error {
	if err == nil && info.IsDir() && hasGoFile(path) && !fileIgnore.MatchString(path) {
		path = "./" + path
		log.Println(path)
		exec.Command("go", "test", "-covermode=count", "-coverprofile="+path+"/"+cov_tmp_file, path).Output()
	}
	return nil
}

func CreateOutputFile() {
	err := ioutil.WriteFile(output_filename, []byte(cover_header), 0644)
	check(err)
}

func CollateCoverage(path string, info os.FileInfo, err error) error {
	if err == nil && !info.IsDir() && strings.Contains(path, cov_tmp_file) {
		contentsB, err := ioutil.ReadFile(path)
		check(err)
		contents := strings.Replace(string(contentsB), "mode: count\n", "", 1)
		f, err := os.OpenFile(output_filename, os.O_APPEND|os.O_WRONLY, 0600)
		check(err)
		_, err = f.WriteString(contents)
		check(err)
		f.Close()
	}
	return nil
}

func DeleteFiles(path string, info os.FileInfo, err error) error {
	if err == nil && !info.IsDir() && strings.Contains(path, cov_tmp_file) {
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
