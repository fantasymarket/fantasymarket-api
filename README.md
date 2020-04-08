# fantasymarket-api

![CI](https://github.com/explodingcamera/fantasymarket-api/workflows/CI/badge.svg?branch=develop)

## Table of Contents

- [Project Structure](#structure)
- [Installation](#installation)
- [Development](#development)
  - [Recommended Tools](#recommended-tools)
  - [Testing](#testing)

## Structure

<big><pre>
**fantasymarket-app**
├── [api](api/)&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; _# rest api service_
├── [database](database/) &nbsp; _# database service_
├── [game](game/) &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; _# game service_
└── [utils](utils/) &nbsp;&nbsp;&nbsp;&nbsp; _# utility functions_</pre></big>

## Installation

### 1. Install Requirements

- [go >=1.13](https://golang.org/dl/)
- [Task](https://taskfile.dev/#/installation)
- Revive (optional) To install, you can also run `$ task install-linter` after installing task. 

### 2. Clone Repo

```bash
$ git clone https://github.com/explodingcamera/fantasymarket-app.git
$ cd fantasymarket-app
```

## Development

### Recommended Tools

VSCode with the official GO extension or goland 

### Testing

```bash
$ task test # run tests
$ task lint # lint code
```
