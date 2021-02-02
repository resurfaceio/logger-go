package logger

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type NetHttpClientLogger struct {
	http.Client
	LOGFLAG   bool
	isEnabled bool
}

func (bcl *NetHttpClientLogger) Get(url string) (resp *http.Response, err error) {
	// capture the response or error
	getResp, getErr := bcl.Client.Get(url)

	ioWriter := os.Stdout

	// create or open a file to log to
	if bcl.LOGFLAG {
		ioWriter, err = os.OpenFile("./get.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	// create a log.Logger to direct output of logging
	logger := log.New(ioWriter, "", log.LstdFlags)

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

func newLogger() NetHttpClientLogger {
	return NetHttpClientLogger{
		LOGFLAG:   true,
		isEnabled: true,
	}
}

func (bcl *NetHttpClientLogger) SetLogFlag(flag bool) {
	bcl.LOGFLAG = flag
}
