[![CircleCI](https://circleci.com/gh/philwinder/gocoverage.svg?style=svg)](https://circleci.com/gh/philwinder/gocoverage) [![Coverage Status](https://coveralls.io/repos/github/philwinder/gocoverage/badge.svg?branch=master)](https://coveralls.io/github/philwinder/gocoverage?branch=master)
# gocoverage
A simple flexible tool to generate a unified coverage file for all your Go code.

## Continuous Integration with [coveralls.io](https://coveralls.io/)

Use with [mattn/goveralls](https://github.com/mattn/goveralls) to send metrics
to [coveralls.io](https://coveralls.io/):

```sh
go get github.com/philwinder/gocoverage
go get github.com/mattn/goveralls
gocoverage
goveralls -coverprofile=profile.cov -service=circle-ci -repotoken=${COVERALLS_TOKEN}
```

### Options

```sh
Usage of gocoverage:
  -dir string
        Directory to start recursing for tests (default ".")
  -ignore string
        RegEx that ignores files and folders. Default ignores hidden folders and vendor folder. (default "\\/(vendor|\\.\\w+)")
  -output string
        Filename for the output coverage file. (default "profile.cov")

```
