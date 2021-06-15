package logger

import (
	"log"

	"github.com/gorilla/mux"
)

type HttpLoggerForMux struct {
	logger HttpLogger
	router mux.Router
}

func NewHttpLoggerForMux(r mux.Router, options Options) {

	lgr, err := NewHttpLogger(options)

	if err != nil {
		log.Fatal(err)
	}

	logger := HttpLoggerForMux{
		logger: *lgr,
		router: r,
	}

}
