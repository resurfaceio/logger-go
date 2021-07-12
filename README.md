# resurfaceio-logger-go
Easily log API requests and responses to your own [system of record](https://resurface.io).

## Contents

<ul>
  <li><a href="#installation">Installation</a></li>
  <li><a href="#logging_from_mux">Logging from gorilla/mux</a></li>
  <li><a href="#privacy">Protecting User Privacy</a></li>
</ul>

<a name="installation"/>

## Installation

In the same directory as your project's `go.mod` and `go.sum` files.

```
go get github.com/resurfaceio/logger-go
```

<a name="logging_from_mux"/>

## Logging from gorilla/mux

```golang
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/resurfaceio/logger-go/src/logger"
)


func main() {
	router := mux.NewRouter()
  
	options := logger.Options{
		Rules:   "allow_http_url\nskip_compression\n ...",
		Url:     "http://localhost:4001/message",
		Enabled: true,
		Queue:   nil,
	}

	httpLoggerForMux, err := logger.NewHttpLoggerForMuxOptions(options)

	if err != nil {
		log.Fatal(err)
	}

	app.Router.Use(httpLoggerForMux.StartResponse)

	log.Fatal(http.ListenAndServe(":5000", router))
}
```

<a name="privacy"/>

## Protecting User Privacy

Loggers always have an active set of <a href="https://resurface.io/rules.html">rules</a> that control what data is logged
and how sensitive data is masked. All of the examples above apply a predefined set of rules (`include debug`),
but logging rules are easily customized to meet the needs of any application.

<a href="https://resurface.io/rules.html">Logging rules documentation</a>
