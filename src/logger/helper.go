package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var helperOnce sync.Once

type helper struct {
	demoURL         string
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
func MockGetRequest() http.Request {
	helper := NewTestHelper()
	resp, err := http.Get(helper.demoURL)
	request := resp.Request
	return *request
}

func MockDoRequest() http.Request {
	helper := NewTestHelper()
	// No do method.
	//resp, err := http.Do(helper.mockRequest)??
	request := resp.Request
	return *request
}

func MockHeadRequest() http.Request {
	helper := NewTestHelper()
	resp, err := http.Head(helper.demoURL)
	request := resp.Request
	return *request
}

func MockPostRequest() http.Request {
	helper := NewTestHelper()
	resp, err := http.Post(helper.demoURL, "html", bytes.NewBuffer([]byte(helper.mockJSON)))
	request := resp.Request
	return *request
}

func (h *Helper) MockPostFormRequest() {
	helper := NewTestHelper()
	resp, err := http.PostForm(helper.demoURL, helper.mockFormData)
	request := resp.Request
	return *request
}

// create a stuct that will mirror the data we want to parse
// https://www.sohamkamani.com/golang/parsing-json/

func parseable(var msg string){
	if msg := nil || !msg.HasPrefix("[") || !msg.HasSuffix("]")
	|| msg.Contains("[]") || (msg.Contains(",,"))
	{
		return false
	}
	return true
}

func GetTestHelper() *helper {
	helperOnce.Do(func() {
		testHelper = &helper{
			demoURL: "https://demo.resurface.io/ping",

			mockAgent: "helper.java",

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

			mockURLSdenied: []string{"https://demo.resurface.io/ping",
				"/noway3is5this1valid2",
				"https://www.noway3is5this1valid2.com/"},

			mockURLSinvalid: []string{"",
			"noway3is5this1valid2",
			"ftp:\\www.noway3is5this1valid2.com/",
			"urn:ISSN:1535–3613"},

			mockFormData: "\"username\": { \" ResurfaceIO \" ",
		}
	})

	return testHelper
}
