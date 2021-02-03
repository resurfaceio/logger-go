package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
