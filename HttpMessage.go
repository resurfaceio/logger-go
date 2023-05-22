// © 2016-2023 Resurface Labs Inc.

package logger

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/brotli/go/cbrotli"
)

// helper function to read body bytes
func readBody(rBody io.ReadCloser, encoding string) (string, error) {
	const bodyLimit = 1024 * 1024
	var reader io.Reader
	var err error
	defer rBody.Close()

	switch encoding {
	case "gzip", "x-gzip":
		reader, err = gzip.NewReader(rBody)
		if err != nil {
			return "", err
		}
	case "deflate", "zlib", "deflated":
		reader, err = zlib.NewReader(rBody)
		if err != nil {
			return "", err
		}
	case "br":
		reader = cbrotli.NewReader(rBody)
	case "", "identity":
		reader = rBody
	default:
		return "", io.ErrNoProgress
	}
	reader = io.LimitReader(reader, bodyLimit)

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", nil
	}

	return string(bodyBytes), err
}

// create Http message for any logger
func buildHttpMessage(req *http.Request, resp *http.Response) [][]string {
	var message [][]string

	method := req.Method
	if method != "" {
		message = append(message, []string{"request_method", method})
	}

	//protocol is not currently used
	// message = append(message, []string{"request_protocol", req.Proto})

	var fullUrl string
	if req.URL.IsAbs() {
		fullUrl = req.RequestURI
	} else {
		//Not sure of a better way to do this at the moment - 6/24/21
		//check for other tls proto
		if req.TLS == nil {
			fullUrl = "http://" + req.Host + req.URL.Path
		} else {
			fullUrl = "https://" + req.Host + req.URL.Path
		}
		// ---
	}

	message = append(message, []string{"request_url", fullUrl})
	message = append(message, []string{"response_code", fmt.Sprint(resp.StatusCode)})

	if req.Body != nil {
		var contentEncoding string
		if encodings, encoded := req.Header["Content-Encoding"]; encoded {
			contentEncoding = encodings[0]
		}
		requestBody, err := readBody(req.Body, contentEncoding)
		if err != nil {
			log.Println(err)
		}
		message = append(message, []string{"request_body", requestBody})

		// Unescaped semicolons in querystring make ParseForm return a non-nil error
		req.URL.RawQuery = strings.ReplaceAll(req.URL.RawQuery, ";", "%3B")
		err = req.ParseForm()
		if err != nil {
			log.Println(err)
		}
	}

	appendRequestHeaders(&message, req)
	appendRequestParams(&message, req)
	appendResponseHeaders(&message, resp)

	if resp.Body != nil {
		var contentEncoding string
		if encodings, encoded := resp.Header["Content-Encoding"]; encoded {
			contentEncoding = encodings[0]
		}
		responseBody, err := readBody(resp.Body, contentEncoding)
		if err != nil {
			log.Println(err)
		}
		message = append(message, []string{"response_body", responseBody})
	}

	return message

}

// SendHttpMessage(l *HttpLogger, resp *http.Response, req *http.Request, now int64, interval int64) Uses logger l to send a log of the given resp and req to the loggers url
// here, now refers to the time at which the request was received and interval corresponds to the time between request and response. customFields are used to pass custom information fields through the logger to Resurface.
func SendHttpMessage(logger *HttpLogger, resp *http.Response, req *http.Request, now int64, interval int64, customFields map[string]string) {

	if !logger.Enabled() {
		return
	}

	// copy details from request & response
	message := buildHttpMessage(req, resp)
	copySessionField := logger.rules.CopySessionField()

	// copy data from session if configured
	if len(copySessionField) != 0 {
		sessionCookies := req.Cookies()
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
	// append request time, if given. If not, append logging time
	if now == 0 {
		now = time.Now().UnixNano() / int64(time.Millisecond)
	}
	message = append(message, []string{"now", strconv.FormatInt(now, 10)})

	// append interval noting the time between request and response
	if interval != 0 {
		message = append(message, []string{"interval", strconv.FormatInt(interval, 10)})
	} else {
		message = append(message, []string{"interval", strconv.FormatInt(1, 10)})
	}

	logger.submitIfPassing(message, customFields)
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

	if req.RemoteAddr != "" {
		addr := req.RemoteAddr
		if idx := strings.Index(addr, ":"); idx != -1 {
			addr = addr[:idx]
		}
		*message = append(*message, []string{"request_header:x-forwarded-for", addr})
	}
}
