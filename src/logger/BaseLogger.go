package logger

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type BaseClientLogger struct {
	http.Client
}

func (cl *BaseClientLogger) Get(url string) (resp *http.Response, err error) {
	// capture the response or error
	getResp, getErr := cl.Client.Get(url)

	// create or open a file to log to
	f, err := os.OpenFile("get.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// create a log.Logger to direct output of logging
	logger := log.New(f,"", log.LstdFlags)

	// logging conditions
	if err != nil {
		logger.Println("FAILURE: GET Request")
		req, _ := httputil.DumpRequest(getResp.Request, true)
		logger.Println("URL: " + string(req))
		logger.Println("ERROR: " + err.Error())
		logger.Println("STATUS: " + getResp.Status)
		//logger.Println("STATUS CODE: " + getResp.StatusCode)
	} else {
		//bodyBytes, _ := ioutil.ReadAll(getResp.Body)
		logger.Println("SUCCESS: GET Request")
		logger.Println("URL: " + url)
		//logger.Println("RESPONSE: " + string(bodyBytes))
	}

	return getResp, getErr
}

// takes
func NewLogger() BaseClientLogger {
	return BaseClientLogger{

	}
}