package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

/*
* Submits request and response through logger.
 */
func sendNetHttpClientMessage(logger *HttpLogger, resp *http.Response, now int64, interval float64) {
	request := resp.Request

	if !logger.isEnabled() {
		return
	}

	// copy details from request & response
	message := buildNetHttpClientMessage(resp)

	copySessionField := logger.rules.copySessionField

	// copy data from session if configured
	if len(copySessionField) != 0 {
		sessionCookies := request.Cookies()
		if len(sessionCookies) != 0 {
			for _, r := range copySessionField {
				for _, cookie := range sessionCookies {
					name := cookie.Name
					matched, err := regexp.MatchString(r.param1, name)
					if err == nil && matched == true {
						cookieVal := cookie.Value
						message = append(message,
							[]string{"session_field:" + name, cookieVal})
					}
				}
			}
		}
	}

	// add timing details
	if now == 0 {
		timeNow := time.Now()
		unixNano := timeNow.UnixNano()
		umillisec := unixNano / int64(time.Millisecond)

		now = umillisec
	}
	message = append(message, []string{"now", string(now)})

	if interval != 0 {
		message = append(message, []string{"interval", fmt.Sprint(interval)})
	}

	logger.submitIfPassing(message)

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

	appendRequestHeaders(&message, request)
	appendRequestParams(&message, request)
	appendResponseHeaders(&message, resp)

	reqBodyBytes, err := ioutil.ReadAll(request.Body)
	reqBody := string(reqBodyBytes)
	if err != nil && reqBody != "" {
		message = append(message, []string{"request_body", reqBody})
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	respBody := string(respBodyBytes)
	if err != nil && respBody != "" {
		message = append(message, []string{"response_body", respBody})
	} else {
		message = append(message, []string{"response_body", "ISO-8859-1"})
	}

	return message

}

/*
* Adds response headers to message.
 */
func appendResponseHeaders(message *[][]string, resp *http.Response) {
	respHeader := resp.Header
	for headerName, headerValues := range respHeader {
		name := "response_header:" + strings.ToLower(headerName)
		for _, value := range headerValues {
			*message = append(*message, []string{name, value})
		}
	}
}

/*
* Adds request params to message.
 */
func appendRequestParams(message *[][]string, req *http.Request) {
	reqParams := req.Form
	for paramName, params := range reqParams {
		name := "request_param:" + strings.ToLower(paramName)
		for _, param := range params {
			*message = append(*message, []string{name, param})
		}
	}
}

/*
* Adds request headers to message.
 */
func appendRequestHeaders(message *[][]string, req *http.Request) {
	reqHeaders := req.Header
	for headerName, headerValues := range reqHeaders {
		name := "request_header:" + strings.ToLower(headerName)
		for _, value := range headerValues {
			*message = append(*message, []string{name, value})
		}
	}
}
