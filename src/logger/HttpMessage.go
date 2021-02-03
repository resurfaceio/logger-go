package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
* Builds list of key/value pairs for HTTP request and response.
 */
func buildNetHttpClientMessage(logger *NetHttpClientLogger, resp *http.Response) {
	request := resp.Request

	if !logger.isEnabled {
		return
	}

	var message [][]string

	method := resp.Request.Method
	if method != "" {
		message = append(message, []string{"request_method", method})
	}

	message = append(message, []string{"request_url", request.URL.RequestURI()})
	message = append(message, []string{"response_code", fmt.Sprint(resp.StatusCode)})

	appendRequestHeaders(message, request)
	appendRequestParams(message, request)
	appendResponseHeaders(message, resp)

	reqBodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		message = append(message, []string{"request_body", string(reqBodyBytes)})
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		message = append(message, []string{"response_body", string(respBodyBytes)})
	}

}

/*
* Adds response headers to message.
 */
func appendResponseHeaders(message [][]string, resp *http.Response) {

}

/*
* Adds request params to message.
 */
func appendRequestParams(message [][]string, req *http.Request) {

}

/*
* Adds request headers to message.
 */
func appendRequestHeaders(message [][]string, req *http.Request) {

}
