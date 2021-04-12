package logger

import (
	"bytes"
	"encoding/json"
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

//This is a rough outline.
//I am thinking we need to cover the requests rather than the responses.
//MockGetRequest covers a get request to compare against loggin.
//https://appdividend.com/2019/12/02/golang-http-example-get-post-http-requests-in-golang/
func (testHelper *helper) MockGetRequest() http.Request {
	resp, err := http.Get(testHelper.demoURL)
	if err != nil {
		log.Fatal(err)
	}
	request := resp.Request
	return *request
}

// func MockDoRequest() http.Request {
// 	helper := GetTestHelper()
// 	resp, err := http.Get(helper.demoURL)
// 	resp, err = http.Do(resp)
// 	request := resp.Request
// 	return *request
// }

func (testHelper *helper) MockHeadRequest() http.Request {
	resp, err := http.Head(testHelper.demoURL)
	if err != nil {
		log.Fatal(err)
	}
	request := resp.Request
	return *request
}

func (testHelper *helper) MockPostRequest() http.Request {
	resp, err := http.Post(testHelper.demoURL, "html", bytes.NewBuffer([]byte(testHelper.mockJSON)))
	if err != nil {
		log.Fatal(err)
	}
	request := resp.Request
	return *request
}

func (testHelper *helper) MockPostFormRequest() http.Request {
	form := url.Values{}
	form.Add("username", "resurfaceio")
	resp, err := http.PostForm(testHelper.demoURL, form)
	if err != nil {
		log.Fatal(err)
	}
	request := resp.Request
	return *request
}

func (h *helper) MockRequestWithJson() *http.Request {
	// request, _ := http.NewRequest("POST", h.mockURL, requestBody)
	// request.Header.Add("Content-Type", "Application/JSON")
	// request.PostForm.Add("message", "{ \"hello\" : \"world\" }")
	request, _ := http.NewRequest("POST", "www.legitwebsite.com", strings.NewReader("stuff"))
	return request
}

func (h *helper) MockRequestWithJson2() *http.Request {
	request, _ := http.NewRequest("POST", "www.legitwebsite.com", strings.NewReader("stuff"))
	return request
}

func (h *helper) MockResponseWithHtml() *http.Response {
	response := http.Response{}
	return &response
}

//https://golang.org/pkg/encoding/json/#example_Unmarshal
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

func GetTestHelper() *helper {
	helperOnce.Do(func() {
		testHelper = &helper{
			demoURL: "https://demo.resurface.io/ping",

			demoURL1: "https://demo.resurface.io",

			mockAgent: "helper.go",

			mockHTML: "<html>Hello World!</html>",

			mockHTML2: "<html>Hola Mundo!</html>",

			mockHTML3: "<html>1 World 2 World Red World Blue World!</html>",

			mockHTML4: "<html>1 World\n2 World\nRed World \nBlue World!\n</html>",

			mockHTML5: `<html>\n
			<input type=\"hidden\">SENSITIVE</input>\n
			<input class='foo' type=\"hidden\">\n
			SENSITIVE\n
			</input>\n
			</html>`,

			mockJSON: "{ \"hello\" : \"world\" }",

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
