# resurfaceio-logger-go
Easily log API requests and responses to your own [system of record](https://resurface.io).

pkg.go.dev [documentation](https://pkg.go.dev/github.com/resurfaceio/logger-go/v3) for this Go module.

## Contents

<ul>
  <li><a href="#installation">Installation</a></li>
	<li><a href="#resurface_setup">Setup the Resurface app</a></li>
  <li><a href="#logging_from_mux">Logging from gorilla/mux</a></li>
  <li><a href="#privacy">Protecting User Privacy</a></li>
</ul>

<a name="installation"/>

## Installation

In the same directory as your project's `go.mod` and `go.sum` files.

```
go get github.com/resurfaceio/logger-go/v3
```

<a name="resurface_setup"/>

## Setup the Resurface app

If you don't already have Docker installed, you can do so by following these [instructions](https://docs.docker.com/get-docker/).

Register [here](https://resurface.io/installation) to get access to the Resurface private container registry.

From the terminal, run this Docker command to start up the Resurface app using docker.

```
docker run -v resurface:/db -d --name resurface -p 7700:7700 -p 7701:7701 --cpus=4 --memory=8g docker.resurface.io/release/resurface:3.0.34
```

Point your browser at `http://localhost:7700`.


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
		Rules:   "include debug\n",
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
