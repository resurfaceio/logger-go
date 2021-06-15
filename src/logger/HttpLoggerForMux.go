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

	loggerMux := HttpLoggerForMux{
		logger: *lgr,
		router: r,
	}

	log.Println("Testing", loggerMux.logger)
	return nil, nil
}

func NewHttpLoggerForMuxOptions(options Options) (*HttpLoggerForMux, error) {

	lgr, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	loggerMux := HttpLoggerForMux{
		logger: *lgr,
		// router: r,
	}

	log.Println("Testing", loggerMux.logger.host)
	return nil, nil
}
