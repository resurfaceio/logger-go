package logger

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type BaseClientLogger struct {
	http.Client
	LOG_FLAG bool
}

func (bcl *BaseClientLogger) Get(url string) (resp *http.Response, err error) {
	// capture the response or error
	getResp, getErr := bcl.Client.Get(url)

	ioWriter := os.Stdout

	// create or open a file to log to
	if bcl.LOG_FLAG {
		ioWriter, err = os.OpenFile("./get.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	// create a log.Logger to direct output of logging
	logger := log.New(ioWriter,"", log.LstdFlags)

	// logging conditions
	if err != nil {
		logger.Println("FAILURE: GET Request")
		req, _ := httputil.DumpRequest(getResp.Request, true)
		logger.Println("URL: " + string(req))
		logger.Println("ERROR: " + err.Error())
		logger.Println("STATUS: " + getResp.Status)
		//logger.Println("STATUS CODE: " + getResp.StatusCode)
	} else {
		bodyBytes, _ := ioutil.ReadAll(getResp.Body)
		logger.Println("SUCCESS: GET Request")
		logger.Println("URL: " + url)
		logger.Println("RESPONSE: " + string(bodyBytes))
	}

	return getResp, getErr
}

func NewLogger() BaseClientLogger {
	return BaseClientLogger{
		LOG_FLAG: true,
	}
}

func (bcl *BaseClientLogger) SetLogFlag(flag bool) {
	bcl.LOG_FLAG = flag
}
