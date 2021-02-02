package logger

import "net/http"

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

	formattedURL := formatURL(request)
	if formattedURL != "" {
		message = append(message, []string{"request_url", formattedURL})
	}

	message = append(message, []string{"response_code", string(resp.StatusCode)})

}

func appendResponseHeaders(resp *http.Response) {

}

func appendRequestParams(req *http.Request) string {

}

func appendRequestHeaders(req *http.Request) string {

}

func formatURL(req *http.Request) string {

}
