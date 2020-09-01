# fantasymarket-api ![CI](https://github.com/fantasymarket/fantasymarket-api/workflows/CI/badge.svg) [![Maintainability](https://api.codeclimate.com/v1/badges/0702b9a5e11f3a0b7629/maintainability)](https://codeclimate.com/github/fantasymarket/fantasymarket-api/maintainability) [![codecov](https://codecov.io/gh/fantasymarket/fantasymarket-api/branch/develop/graph/badge.svg)](https://codecov.io/gh/fantasymarket/fantasymarket-api)


## Table of Contents

- [Introduction](#introduction)
- [Project Structure](#structure)
- [Installation](#installation)
- [Development](#development)
  - [Recommended Tools](#recommended-tools)
  - [Testing](#testing)
  - [Beta impressions](#beta-impressions)
  - [Running beta](#running-beta)
- [Our Current Situation](#our-current-situation)

# Introduction

We are creating a working stock market simulation, where the user can invest in-game currency in a simulated market completely separated from the real world. By providing custom events that pop up randomly and alter the course of our stocks, coupled with a market that never sleeps, we are presenting a fun and fast-paced introduction into the stock market. Our beta is currently in its final phase and we plan to release it fully by the end of august.

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
- [Task](https://taskfile.dev/#/installation) (optional)
- Revive (optional) To install, you can also run `$ task install-linter` after installing task.
- go-bindata (required for building a binary) - To install, you can also run `$ task install-bindata` after installing task.
- Windows:
	- [gcc](https://sourceforge.net/projects/tdm-gcc/)
- OSX:
	- You might need to install [sqlite](https://github.com/mattn/go-sqlite3#mac-osx) (should be installed already)
- Linux:
	- You need to have the [development tools package for your distro](https://github.com/mattn/go-sqlite3#linux) installed

NOTE: Due to some changes in gcc 10, you might see some warnings because of the sqlite bindings we use. This won't cause any issues. 

### 2. Clone Repo

```bash
$ git clone https://github.com/fantasymarket/fantasymarket-app.git
$ cd fantasymarket-app
```

## Development

### Recommended Tools

VSCode with the official GO extension or goland 

### Testing

With `Task` installed:

```bash
$ task test # run tests
$ task lint # lint code
```

Alternative:
```bash
$ go test ./...
```

## Beta impressions

| Landing Page      | Chart View      | Trading View |
|------------|-------------|-------------|
| <img src=https://i.imgur.com/lK5VEYK.jpg> | <img src=https://i.imgur.com/duyFyZS.jpg> | <img src=https://i.imgur.com/AncR4ah.png> |

### Running Beta:

After installing all the requirements, start the program by running `$ go run main.go`.\
The console should output information about the status of the program:

```bash
INF successfully connected to the database
INF successfully started the game loop
INF successfully started the http server address=localhost:5000
```

It will then load the last tick from the Database before printing it and the ingame time:

```bash
DBG running tick date="2020-01-01 00:00:00 +0000 UTC" tick=2
```

After this setup, the program prints the indices from the two currently implemented stocks **Google** and **Apple**. 
This is done every 10 seconds, like this:

```bash
DBG updated stock index=60032 name=GOOG
DBG updated stock index=60024 name=APPL
```
## Our Current Situation

Currently, the project needs some final refactoring, some last bits of database logic and ultimately testing. 
Even though some tests are already completed, more are needed to ensure a well tested project.
Finally, we have also started to link the frontend to the backend. If all goes smoothly, we will be able to release the project fully by the end of september.
