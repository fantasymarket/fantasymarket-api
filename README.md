# fantasymarket-api ![CI](https://github.com/fantasymarket/fantasymarket-api/workflows/CI/badge.svg) [![Maintainability](https://api.codeclimate.com/v1/badges/0702b9a5e11f3a0b7629/maintainability)](https://codeclimate.com/github/fantasymarket/fantasymarket-api/maintainability) [![codecov](https://codecov.io/gh/fantasymarket/fantasymarket-api/branch/develop/graph/badge.svg)](https://codecov.io/gh/fantasymarket/fantasymarket-api)


## Table of Contents

- [Introduction](#introduction)
- [Project Structure](#structure)
- [Installation](#installation)
- [Development](#development)
  - [Recommended Tools](#recommended-tools)
  - [Testing](#testing)
  - [Impressions](#beta-impressions)
  - [Running the app](#running-beta)

# Introduction

We are creating a working stock market simulation, where the user can invest in-game currency in a simulated market completely separated from the real world. By providing custom events that pop up randomly and alter the course of our stocks, coupled with a market that never sleeps, we are presenting a fun and fast-paced introduction into the stock market. This is the repository for the backend, the corresponding frontend can be found [here](https://github.com/fantasymarket/fantasymarket-app).

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

## Impressions

| Landing Page      | Chart View      | Trading View |
|------------|-------------|-------------|
| <img src=https://i.imgur.com/DNQ0Xw8.jpg> | <img src=https://i.imgur.com/w5ikPh4.jpg> | <img src=https://i.imgur.com/A6Ga5k0.jpg> |

# Impressssions
## Landing Page
![Landing Page](https://i.imgur.com/DNQ0Xw8.jpg)

## Chart View
![Chart View](https://i.imgur.com/w5ikPh4.jpg)

## Trading View
![Trading View](https://i.imgur.com/A6Ga5k0.jpg)

### Running the app:

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

The corresponding instructions for running the frontend can be found [here](https://github.com/fantasymarket/fantasymarket-app).
