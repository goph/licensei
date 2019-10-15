# Licensei

![Build Status](https://github.com/goph/licensei/workflows/CI/badge.svg?style=flat-square)
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
	curl -sfL https://git.io/licensei | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@
```


## Usage

```
Usage:
  licensei [command]

Available Commands:
  cache       Cache licenses of dependencies in the project
  check       Check licenses of dependencies in the project
  header      Check license header of files
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

[header]
ignorePaths = ["vendor", "client", ".gen"]
ignoreFiles = ["mock_*.go", "*_gen.go"]
template = """// Copyright © :YEAR: OWNER
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License."""
```


## GitHub Authentication

Licensei uses Github to lookup licenses for most of the packages.
The used parts of GitHub API does not require authentication, but it's easy to get rate limited,
so it is recommended to configure a personal access token in your environment:

```bash
export GITHUB_TOKEN=xyz
```

## Development

The project uses Go Modules, so the minimum Go version is 1.11.4.

When all coding and testing is done, please run the test suite:

``` bash
$ make check
```


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
