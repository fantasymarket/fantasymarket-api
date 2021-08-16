[![Website](https://img.shields.io/website?label=staging&url=http%3A%2F%2Fdevelop--fantasymarket.netlify.app%2F)](https://develop--fantasymarket.netlify.app/) [![Website](https://img.shields.io/website?label=production&url=http%3A%2F%2Ffantasymarket.netlify.app%2F)](https://fantasymarket.netlify.app/)
[![netlify](https://img.shields.io/netlify/306db36d-47d1-40d3-9f52-c52a5b7633e5?style=flat)](https://app.netlify.com/sites/fantasymarket/overview)
[![codecov](https://codecov.io/gh/fantasymarket/fantasymarket-app/branch/develop/graph/badge.svg)](https://codecov.io/gh/fantasymarket/fantasymarket-app)
[![Security Headers](https://img.shields.io/security-headers?url=http%3A%2F%2Fdevelop--fantasymarket.netlify.app%2F)](https://securityheaders.com/?q=http%3A%2F%2Fdevelop--fantasymarket.netlify.app%2F&followRedirects=on)
[![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade/develop--fantasymarket.netlify.app?publish)](https://observatory.mozilla.org/analyze/develop--fantasymarket.netlify.app)
[![Maintainability](https://api.codeclimate.com/v1/badges/0b67777ccab5a08e0546/maintainability)](https://codeclimate.com/github/fantasymarket/fantasymarket-app/maintainability)

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

We are creating a stock market simulation, where the user can invest in-game currency in a simulated market completely separated from the real world. By providing custom events that pop up randomly and alter the course of our stocks, coupled with a market that never sleeps, we are presenting a fun and fast-paced introduction into the stock market. This is the repository for the backend, the corresponding frontend can be found [here](https://github.com/fantasymarket/fantasymarket-app).

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

# Impressions
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
