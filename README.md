# resurfaceio-logger-go
Easily log API requests and responses to your own [system of record](https://resurface.io).

[![Go project version](https://badge.fury.io/go/github.com%2Fresurfaceio%2Flogger-go.svg)](https://badge.fury.io/go/github.com%2Fresurfaceio%2Flogger-go)
[![CodeFactor](https://www.codefactor.io/repository/github/resurfaceio/logger-go/badge)](https://www.codefactor.io/repository/github/resurfaceio/logger-go)
[![License](https://img.shields.io/github/license/resurfaceio/logger-go)](https://github.com/resurfaceio/logger-go/blob/master/LICENSE)
[![Contributing](https://img.shields.io/badge/contributions-welcome-green.svg)](https://github.com/resurfaceio/logger-go/blob/master/CONTRIBUTING.md)

## Contents

<ul>
<li><a href="#dependencies">Dependencies</a></li>
<li><a href="#installation">Installation</a></li>
<li><a href="#logging_from_mux">Logging from gorilla/mux</a></li>
<li><a href="#privacy">Protecting User Privacy</a></li>
</ul>

<a name="dependencies"/>

## Dependencies

Requires go 1.15 or later.

<a name="installation"/>

## Installation

Run this command in the same directory as your project's `go.mod` and `go.sum` files:

```
go get github.com/resurfaceio/logger-go/v3
```

<a name="logging_from_mux"/>

## Logging from gorilla/mux

```golang
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/resurfaceio/logger-go/v3" //<----- 1
)


func main() {
	router := mux.NewRouter()
  
	options := logger.Options{ //<----- 2
		Rules:   "include_debug\n",
		Url:     "http://localhost:7701/message",
		Enabled: true,
		Queue:   nil,
	}

	httpLoggerForMux, err := logger.NewHttpLoggerForMuxOptions(options) //<----- 3

	if err != nil {
		log.Fatal(err)
	}

	router.Use(httpLoggerForMux.LogData) //<----- 4

	log.Fatal(http.ListenAndServe(":5000", router))
}
```

<a name="privacy"/>

## Protecting User Privacy

Loggers always have an active set of <a href="https://resurface.io/rules.html">rules</a> that control what data is logged
and how sensitive data is masked. All of the examples above apply a predefined set of rules, `include_debug`,
but logging rules are easily customized to meet the needs of any application.

<a href="https://resurface.io/rules.html">Logging rules documentation</a>

---
<small>&copy; 2016-2022 <a href="https://resurface.io">Resurface Labs Inc.</a></small>
