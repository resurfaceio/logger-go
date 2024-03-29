// © 2016-2024 Graylog, Inc.

package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var helperOnce sync.Once

type helper struct {
	demoURL         string
	demoURL1        string
	mockAgent       string
	mockHTML        string
	mockHTML2       string
	mockHTML3       string
	mockHTML4       string
	mockHTML5       string
	mockJSON        string
	mockJSONescaped string
	mockNow         int64
	mockQueryString string
	mockURL         string
	mockURLSdenied  []string
	mockURLSinvalid []string
	mockFormData    string
}

var testHelper *helper

// MockGetRequest covers a get request to compare against logging.
// https://appdividend.com/2019/12/02/golang-http-example-get-post-http-requests-in-golang/
func (testHelper *helper) MockGetRequest() (request http.Request, err error) {
	resp, err := http.Get(testHelper.demoURL)
	if err != nil {
		log.Println(err)
		return
	}
	request = *resp.Request
	return
}

// func MockDoRequest() http.Request {
// 	helper := newTestHelper()
// 	resp, err := http.Get(helper.demoURL)
// 	resp, err = http.Do(resp)
// 	request := resp.Request
// 	return *request
// }

func (testHelper *helper) MockHeadRequest() (request http.Request, err error) {
	resp, err := http.Head(testHelper.demoURL)
	if err != nil {
		log.Println(err)
		return
	}
	request = *resp.Request
	return
}

func (testHelper *helper) MockPostRequest() (request http.Request, err error) {
	resp, err := http.Post(testHelper.demoURL, "html", bytes.NewBuffer([]byte(testHelper.mockJSON)))
	if err != nil {
		log.Println(err)
		return
	}
	request = *resp.Request
	return
}

func (testHelper *helper) MockPostFormRequest() (request http.Request, err error) {
	form := url.Values{}
	form.Add("username", "resurfaceio")
	resp, err := http.PostForm(testHelper.demoURL, form)
	if err != nil {
		log.Println(err)
		return
	}
	request = *resp.Request
	return
}

func (h *helper) MockRequestWithJson() (request *http.Request) {
	// request, _ := http.NewRequest("POST", h.mockURL, requestBody)
	// request.Header.Add("Content-Type", "Application/JSON")
	// request.PostForm.Add("message", "{ \"hello\" : \"world\" }")
	requestBody, _ := json.Marshal(h.mockJSON)
	request, _ = http.NewRequest("POST", h.mockURL, bytes.NewBuffer(requestBody))
	request.Header.Add("content-type", "Application/JSON")
	return
}

func (h *helper) MockRequestWithJson2() (request *http.Request) {
	request = h.MockRequestWithJson()
	request.Header.Add("ABC", "123")
	request.Header.Add("A", "1")
	request.Header.Add("A", "2")
	request.Header.Add("ABC", "123")
	request.Header.Add("ABC", "234")
	return
}

func (h *helper) MockResponse() (response *http.Response) {
	response = &http.Response{
		Body:   io.NopCloser(bytes.NewBufferString(h.mockHTML)),
		Header: map[string][]string{},
	}
	return
}

func (h *helper) MockResponseWithHtml() (response *http.Response) {
	response = h.MockResponse()
	response.StatusCode = 200
	response.Header.Add("content-type", "text/html; charset=utf-8")
	response.Request = h.MockRequestWithJson2()
	return
}

// https://golang.org/pkg/encoding/json/#example_Unmarshal
func parseable(msg string) bool {
	if msg == "" || !strings.HasPrefix(msg, "[") || !strings.HasSuffix(msg, "]") || strings.Contains(msg, "[]") || strings.Contains(msg, ",,") {
		return false
	}
	/* json.Valid won't work with our custom json formatted
	we need to test our custom string over the wire first
	then see if we can use Go native methods like json.Marshal
	and json.Unmarshal
	*/
	return json.Valid([]byte(msg))
}

func newTestHelper() *helper {
	helperOnce.Do(func() {
		testHelper = &helper{
			demoURL: "https://demo.resurface.io/ping",

			demoURL1: "https://demo.resurface.io",

			mockAgent: "helper.go",

			mockHTML: "<html>Hello World!</html>",

			mockHTML2: "<html>Hola Mundo!</html>",

			mockHTML3: "<html>1 World 2 World Red World Blue World!</html>",

			mockHTML4: "<html>1 World\n2 World\nRed World \nBlue World!\n</html>",

			mockHTML5: "<html>\n<input type=\"hidden\">SENSITIVE</input>\n<input class='foo' type=\"hidden\">\nSENSITIVE\n</input>\n</html>",

			mockJSON: `{ \"hello\" : \"world\" }`,

			mockJSONescaped: "{ \\'hello\\' : \\'world\\' }",

			mockNow: 1455908640173,

			mockQueryString: "foo=bar",

			mockURL: "http://something.com:3000/index.html",

			mockURLSdenied: []string{
				"https://demo.resurface.io/ping/noway3is5this1valid2",
				"https://www.noway3is5this1valid2.com/",
			},

			mockURLSinvalid: []string{"",
				"noway3is5this1valid2",
				"ftp:\\www.noway3is5this1valid2.com/",
				"urn:ISSN:1535–3613",
			},

			mockFormData: "\"username\": { \" ResurfaceIO \" ",
		}
	})

	return testHelper
}
