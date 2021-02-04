package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// just for testing
type HttpLogger struct {
	enabled bool
}

/*
* Submits request and response through logger.
 */
func sendNetHttpClientMessage(logger *HttpLogger, resp *http.Response, now int64, interval float64) {

	if !logger.enabled {
		return
	}

	// copy details from request & response
	message := buildNetHttpClientMessage(resp)

	// copy data from session if configured
	if len(logger.rules.copySessionField) != 0 {

	}

}

/*
* Builds list of key/value pairs for HTTP request and response.
 */
func buildNetHttpClientMessage(resp *http.Response) [][]string {
	request := resp.Request

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
	} else {
		message = append(message, []string{"response_body", "ISO-8859-1"})
	}

	return message

}

/*
* Adds response headers to message.
 */
func appendResponseHeaders(message [][]string, resp *http.Response) {
	respHeader := resp.Header
	for headerName, headerValues := range respHeader {
		name := "response_header:" + strings.ToLower(headerName)
		for _, value := range headerValues {
			message = append(message, []string{name, value})
		}
	}
}

/*
* Adds request params to message.
 */
func appendRequestParams(message [][]string, req *http.Request) {
	reqParams := req.Form
	for paramName, params := range reqParams {
		name := "request_param:" + strings.ToLower(paramName)
		for _, param := range params {
			message = append(message, []string{name, param})
		}
	}
}

/*
* Adds request headers to message.
 */
func appendRequestHeaders(message [][]string, req *http.Request) {
	reqHeaders := req.Header
	for headerName, headerValues := range reqHeaders {
		name := "request_header:" + strings.ToLower(headerName)
		for _, value := range headerValues {
			message = append(message, []string{name, value})
		}
	}
}
