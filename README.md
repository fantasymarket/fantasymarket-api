# fantasymarket-api

![CI](https://github.com/explodingcamera/fantasymarket-api/workflows/CI/badge.svg?branch=develop)

## Table of Contents

- [Project Structure](#structure)
- [Installation](#installation)
- [Development](#development)
  - [Recommended Tools](#recommended-tools)
  - [Testing](#testing)
  - [Running](#running-beta)
 - [Situation](#situation)

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
- go-bindata - To install, you can also run `$ task install-bindata` after installing task.
- Windows:
	- [gcc](https://sourceforge.net/projects/tdm-gcc/)
- OSX:
	- You might need to install [sqlite](https://github.com/mattn/go-sqlite3#mac-osx) (should be installed already)
- Linux:
	- You need to have the [development tools package for your distro](https://github.com/mattn/go-sqlite3#linux) installed

### 2. Clone Repo

```bash
$ git clone https://github.com/explodingcamera/fantasymarket-app.git
$ cd fantasymarket-app
```

## Development

### Recommended Tools

VSCode with the official GO extension or goland 

### Testing

Linux:

```bash
$ task test # run tests
$ task lint # lint code
```

Windows:

```bash
$ go test ./...
```

### Running Beta:

After installing all the requirements, start the program by running ```go run main.go```.
The console should output infomation about the status of the program:

      INF successfully connected to the database
      INF successfully started the game loop
      INF successfully started the http server address=localhost:5000

It will then load the last tick from the Database before printing it and the ingame time:

    DBG running tick date="2020-01-01 00:00:00 +0000 UTC" tick=2

After this setup, the program prints the indices from the two currently implemented stocks **Google** and **Apple**. 
This is done every 10 seconds, like this:

    DBG updated stock index=60032 name=GOOG
    DBG updated stock index=60024 name=APPL

## Situation:

Currently, the project needs some final refactoring, some last bits of database logic and ultimatly testing. 
Even though some tests are already completed, more are needed to ensure a well tested project.
We are also looking to start linking the frontend to the backend to complete. 
If all that goes smoothly, we will be able to release the project fully.
