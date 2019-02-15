# Licensei

[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.11.4-orange.svg?style=flat-square)](https://github.com/goph/licensei)
[![Build Status](https://travis-ci.com/goph/licensei.svg?branch=master)](https://travis-ci.com/goph/licensei)
[![Go Report Card](https://goreportcard.com/badge/github.com/goph/licensei?style=flat-square)](https://goreportcard.com/report/github.com/goph/licensei)
[![GolangCI](https://golangci.com/badges/github.com/goph/licensei.svg)](https://golangci.com)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/goph/licensei)

**Library and various tools for working with project licenses.**


## Installation

Install the latest pre-built binary from the [Releases](https://github.com/goph/licensei/releases) page.

Alternatively, place the following code in your `Makefile`:

```makefile
LICENSEI_VERSION = 0.1.0
bin/licensei: bin/licensei-${LICENSEI_VERSION}
	@ln -sf licensei-${LICENSEI_VERSION} bin/licensei
bin/licensei-${LICENSEI_VERSION}:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@
```


## Usage

```
Usage:
  licensei [command]

Available Commands:
  cache       Cache licenses of dependencies in the project
  check       Check licenses of dependencies in the project
  help        Help about any command
  list        List licenses of dependencies in the project

Flags:
      --config string   config file (default is $PWD/.licensei.yaml)
  -h, --help            help for licensei

Use "licensei [command] --help" for more information about a command.
```

See an example integration into [CircleCI](http://circleci.com/) [here](https://github.com/banzaicloud/pipeline/blob/master/.circleci/config.yml#L56-L80).


## Configuration example

Place the following content into `.licensei.toml`.

```toml
approved = [
  "mit",
  "apache-2.0",
  "bsd-3-clause",
  "bsd-2-clause",
  "mpl-2.0",
]

ignored = [
  "github.com/aliyun/aliyun-oss-go-sdk",
  "github.com/ghodss/yaml",
  "github.com/gogo/protobuf",
  "github.com/golang/protobuf",
  "github.com/stretchr/testify",

  "github.com/davecgh/go-spew", # ISC license
  "github.com/howeyc/gopass", # ISC license
  "github.com/oracle/oci-go-sdk", # UPL-1.0
]
```


## Development

The project uses Go Modules, so the minimum Go version is 1.11.4.

When all coding and testing is done, please run the test suite:

``` bash
$ make check
```


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
