package logger

import (
	"log"

	"github.com/gorilla/mux"
)

type HttpLoggerForMux struct {
	logger HttpLogger
	router mux.Router
}

func NewHttpLoggerForMux(r mux.Router) (*HttpLoggerForMux, error) {

	options := Options{}
	lgr, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	logger := HttpLoggerForMux{
		logger: *lgr,
		router: r,
	}

}

func NewHttpLoggerForMuxOptions(r mux.Router, options Options) (*HttpLoggerForMux, error) {

	lgr, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	logger := HttpLoggerForMux{
		logger: *lgr,
		router: r,
	}

	log.Println("Testing")
}
