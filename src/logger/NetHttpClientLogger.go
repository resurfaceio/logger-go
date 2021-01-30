package logger

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

<<<<<<< HEAD:src/logger/NetHttpClientLogger.go
type netHttpClientLogger struct {
=======
type baseClientLogger struct {
>>>>>>> 7c6f55721f57030eaff4cd21234ba3cece9a0c17:src/logger/BaseLogger.go
	http.Client
	LOGFLAG bool
}

<<<<<<< HEAD:src/logger/NetHttpClientLogger.go
func (bcl *netHttpClientLogger) Get(url string) (resp *http.Response, err error) {
=======
func (bcl *baseClientLogger) Get(url string) (resp *http.Response, err error) {
>>>>>>> 7c6f55721f57030eaff4cd21234ba3cece9a0c17:src/logger/BaseLogger.go
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

<<<<<<< HEAD:src/logger/NetHttpClientLogger.go
func newLogger() netHttpClientLogger {
	return netHttpClientLogger{
=======
func newLogger() baseClientLogger {
	return baseClientLogger{
>>>>>>> 7c6f55721f57030eaff4cd21234ba3cece9a0c17:src/logger/BaseLogger.go
		LOGFLAG: true,
	}
}

<<<<<<< HEAD:src/logger/NetHttpClientLogger.go
func (bcl *netHttpClientLogger) SetLogFlag(flag bool) {
=======
func (bcl *baseClientLogger) SetLogFlag(flag bool) {
>>>>>>> 7c6f55721f57030eaff4cd21234ba3cece9a0c17:src/logger/BaseLogger.go
	bcl.LOGFLAG = flag
}
