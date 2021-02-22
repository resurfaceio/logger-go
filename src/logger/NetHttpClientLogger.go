package logger

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type netHttpClientLogger struct {
	http.Client
	LOGFLAG      bool
	isEnableable bool
}

func (bcl *netHttpClientLogger) Get(url string) (resp *http.Response, err error) {
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

func (bcl *netHttpClientLogger) Post(url string) (resp *http.Response, err error) {
	// capture the response or error
	postResp, postErr := bcl.Client.Get(url)

	ioWriter := os.Stdout

	// create or open a file to log to
	if bcl.LOGFLAG {
		ioWriter, err = os.OpenFile("./post.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	// create a log.Logger to direct output of logging
	logger := log.New(ioWriter, "", log.LstdFlags)

	// logging conditions
	if err != nil {
		logger.Println("FAILURE: POST Request")
		req, _ := httputil.DumpRequest(postResp.Request, true)
		logger.Println("URL: " + string(req))
		logger.Println("ERROR: " + err.Error())
		logger.Println("STATUS: " + postResp.Status)
	} else {
		bodyBytes, _ := ioutil.ReadAll(postResp.Body)
		logger.Println("SUCCESS: POST Request")
		logger.Println("URL: " + url)
		logger.Println("RESPONSE: " + string(bodyBytes))
	}

	return postResp, postErr
}

func (bcl *netHttpClientLogger) Delete(url string) (resp *http.Response, err error) {
	// capture the response or error
	deleteResp, deleteErr := bcl.Client.Get(url)

	ioWriter := os.Stdout

	// create or open a file to log to
	if bcl.LOGFLAG {
		ioWriter, err = os.OpenFile("./delete.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	// create a log.Logger to direct output of logging
	logger := log.New(ioWriter, "", log.LstdFlags)

	// logging conditions
	if err != nil {
		logger.Println("FAILURE: DELETE Request")
		req, _ := httputil.DumpRequest(deleteResp.Request, true)
		logger.Println("URL: " + string(req))
		logger.Println("ERROR: " + err.Error())
		logger.Println("STATUS: " + deleteResp.Status)
	} else {
		bodyBytes, _ := ioutil.ReadAll(deleteResp.Body)
		logger.Println("SUCCESS: DELETE Request")
		logger.Println("URL: " + url)
		logger.Println("RESPONSE: " + string(bodyBytes))
	}

	return deleteResp, deleteErr
}

func (bcl *netHttpClientLogger) Put(url string) (resp *http.Response, err error) {
	// capture the response or error
	putResp, putErr := bcl.Client.Get(url)

	ioWriter := os.Stdout

	// create or open a file to log to
	if bcl.LOGFLAG {
		ioWriter, err = os.OpenFile("./put.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	// create a log.Logger to direct output of logging
	logger := log.New(ioWriter, "", log.LstdFlags)

	// logging conditions
	if err != nil {
		logger.Println("FAILURE: PUT Request")
		req, _ := httputil.DumpRequest(putResp.Request, true)
		logger.Println("URL: " + string(req))
		logger.Println("ERROR: " + err.Error())
		logger.Println("STATUS: " + putResp.Status)
	} else {
		bodyBytes, _ := ioutil.ReadAll(putResp.Body)
		logger.Println("SUCCESS: PUT Request")
		logger.Println("URL: " + url)
		logger.Println("RESPONSE: " + string(bodyBytes))
	}

	return putResp, putErr
}

func (bcl *netHttpClientLogger) Patch(url string) (resp *http.Response, err error) {
	// capture the response or error
	patchResp, patchErr := bcl.Client.Get(url)

	ioWriter := os.Stdout

	// create or open a file to log to
	if bcl.LOGFLAG {
		ioWriter, err = os.OpenFile("./patch.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	// create a log.Logger to direct output of logging
	logger := log.New(ioWriter, "", log.LstdFlags)

	// logging conditions
	if err != nil {
		logger.Println("FAILURE: PATCH Request")
		req, _ := httputil.DumpRequest(patchResp.Request, true)
		logger.Println("URL: " + string(req))
		logger.Println("ERROR: " + err.Error())
		logger.Println("STATUS: " + patchResp.Status)
	} else {
		bodyBytes, _ := ioutil.ReadAll(patchResp.Body)
		logger.Println("SUCCESS: PATCH Request")
		logger.Println("URL: " + url)
		logger.Println("RESPONSE: " + string(bodyBytes))
	}

	return patchResp, patchErr
}

func newLogger() netHttpClientLogger {
	return netHttpClientLogger{
		LOGFLAG: true,
	}
}

func (bcl *netHttpClientLogger) SetLogFlag(flag bool) {
	bcl.LOGFLAG = flag
}
