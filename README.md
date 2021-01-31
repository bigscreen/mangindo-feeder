# Mangindo-Feeder

[![Build Status](https://travis-ci.org/bigscreen/mangindo-feeder.svg?branch=master)](https://travis-ci.org/bigscreen/mangindo-feeder)
[![Coverage Status](https://coveralls.io/repos/github/bigscreen/mangindo-feeder/badge.svg?branch=master)](https://coveralls.io/github/bigscreen/mangindo-feeder?branch=master)

This repository contains the source code of Mangindo-Feeder. 
Mangindo-Feeder is a service which has a responsibility to provide friendly responses that are consumable by Mangindo mobile apps, which the responses are produced by mapping original responses from [mangacan](http://mangacanblog.com) web API.

## Project Setup
Clone this repo inside `$GOPATH/src/github.com/bigscreen/`, then setup the project by running the following commands:
```
$ make copy-config
$ make setup
$ make build
```

## Running Tests
Run the following command:
```
$ make test
```
or run the following command to show coverage per package:
```
$ make test-cov
```

## Running Service
Run the following command to start this service *(ensure the setup commands have been ran)*:
```
make build
./out/mangindo-feeder start
```
