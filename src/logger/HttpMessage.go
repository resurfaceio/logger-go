package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
* Submits request and response through logger.
 */
func sendNetHttpClientRequestResponseMessage(logger *HttpLogger, resp *http.Response, start int64) /* maybe return error */ {
	request := resp.Request
	if !logger.Enabled() {
		return
	}

	// copy details from request & response
	message := buildNetHttpClientMessage(request, resp)
	copySessionField := logger.rules.CopySessionField()

	// copy data from session if configured
	if len(copySessionField) != 0 {
		sessionCookies := request.Cookies()
		if sessionCookies != nil {
			for _, r := range copySessionField {
				for _, cookie := range sessionCookies {
					name := strings.ToLower(cookie.Name)
					matched := r.param1.(*regexp.Regexp).MatchString(name)
					if matched {
						cookieVal := cookie.Value
						message = append(message,
							[]string{"session_field:" + name, cookieVal})
					}
				}
			}
		}
	}
	// append time of logging
	now := time.Now().UnixNano() / int64(time.Millisecond)
	message = append(message, []string{"now", strconv.FormatInt(now, 10)})

	// append interval noting the time it took to log
	interval := now - start
	message = append(message, []string{"interval", strconv.FormatInt(interval, 10)})

	logger.submitIfPassing(message)
}

/*
* Builds list of key/value pairs for HTTP request and response.
 */
func buildNetHttpClientMessage(req *http.Request, resp *http.Response) [][]string {
	var message [][]string

	method := resp.Request.Method
	if method != "" {
		message = append(message, []string{"request_method", method})
	}

	message = append(message, []string{"request_url", req.URL.String()})
	message = append(message, []string{"response_code", fmt.Sprint(resp.StatusCode)})

	appendRequestHeaders(&message, req)
	appendRequestParams(&message, req)
	appendResponseHeaders(&message, resp)

	if req.Body != nil {
		reqBodyBytes, err := ioutil.ReadAll(req.Body)
		reqBody := string(reqBodyBytes)
		if err == nil && reqBody != "" {
			message = append(message, []string{"request_body", reqBody})
		}
	}

	if resp.Body != nil {
		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		respBody := string(respBodyBytes)
		if err == nil && respBody != "" {
			message = append(message, []string{"response_body", respBody})
		}
	}

	return message

}

// create Http message for a server side logger
func buildHttpMessage(req *http.Request, resp *http.Response) [][]string {
	var message [][]string

	method := req.Method
	if method != "" {
		message = append(message, []string{"request_method", method})
	}

	message = append(message, []string{"request_url", req.URL.String()})
	message = append(message, []string{"response_code", fmt.Sprint(resp.StatusCode)})

	appendRequestHeaders(&message, req)
	appendRequestParams(&message, req)
	appendResponseHeaders(&message, resp)

	if req.Body != nil {
		reqBodyBytes, err := ioutil.ReadAll(req.Body)
		reqBody := string(reqBodyBytes)
		if err == nil && reqBody != "" {
			message = append(message, []string{"request_body", reqBody})
		}
	}

	if resp.Body != nil {
		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		respBody := string(respBodyBytes)
		if err == nil && respBody != "" {
			message = append(message, []string{"response_body", respBody})
		}
	}

	return message

}

func sendHttpMessage(logger *HttpLogger, resp *http.Response, start int64) {

	request := resp.Request
	if !logger.Enabled() {
		return
	}

	// copy details from request & response
	message := buildHttpMessage(request, resp)
	copySessionField := logger.rules.CopySessionField()

	// copy data from session if configured
	if len(copySessionField) != 0 {
		sessionCookies := request.Cookies()
		if sessionCookies != nil {
			for _, r := range copySessionField {
				for _, cookie := range sessionCookies {
					name := strings.ToLower(cookie.Name)
					matched := r.param1.(*regexp.Regexp).MatchString(name)
					if matched {
						cookieVal := cookie.Value
						message = append(message,
							[]string{"session_field:" + name, cookieVal})
					}
				}
			}
		}
	}
	// append time of logging
	now := time.Now().UnixNano() / int64(time.Millisecond)
	message = append(message, []string{"now", strconv.FormatInt(now, 10)})

	// append interval noting the time it took to log
	interval := now - start
	message = append(message, []string{"interval", strconv.FormatInt(interval, 10)})

	logger.submitIfPassing(message)
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
