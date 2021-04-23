[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-infra/jenkins-version)](https://goreportcard.com/report/github.com/jenkins-infra/jenkins-version)
[![Downloads](https://img.shields.io/github/downloads/jenkins-infra/jenkins-version/total.svg)]()

# jenkins-version

The goal of this tool is to provide a small, simple CLI that can be used to determine the latest Jenkins version, whether that be in the stable or weekly release train, from maven metadata.

This is designed to used in scripts / automation and even github actions to determine that a new Jenkins release is available.

This replaces a python script (getJenkinsVersion.py) that was used to similar but was difficult to distribute.

## To Install

```
brew tap jenkins-infra/tap
brew install jv
```

This can be used a docker container with the following:

```
docker run -it jenkins-infra/jv:main
```

## Usage

To get the latest weekly release:

```
jv get [--username <username> --password <password>]
```

To get the latest LTS release:

```
jv get --version-identifier lts [--username <username> --password <password>]
```

To get the latest LTS for a particular release train:

```
jv get --version-identifier 2.249 [--username <username> --password <password>]
```

## Documentation

More indepth documentation can be found [here](./docs/jv.md)

## Development

To build the application:

```
make build
```

To test:

```
make test
```
